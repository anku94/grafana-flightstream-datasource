package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/pdl/orcastream/pkg/models"
)

// Make sure Datasource implements required interfaces. This is important to do
// since otherwise we will only get a not implemented error response from plugin in
// runtime. In this example datasource instance implements backend.QueryDataHandler,
// backend.CheckHealthHandler interfaces. Plugin should not implement all these
// interfaces - only those which are required for a particular task.
var (
	_ backend.QueryDataHandler      = (*Datasource)(nil)
	_ backend.CheckHealthHandler    = (*Datasource)(nil)
	_ instancemgmt.InstanceDisposer = (*Datasource)(nil)
	_ backend.StreamHandler         = (*Datasource)(nil)
)

type DatasourceConfig struct {
	Addr string `json:"addr"`
}

// Datasource is an example datasource which can respond to data queries, reports
// its health and has streaming skills.
type Datasource struct {
	fsc *FlightStreamClient
}

// NewDatasource creates a new datasource instance.
func NewDatasource(ctx context.Context, instanceSettings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	var cfg DatasourceConfig
	err := json.Unmarshal(instanceSettings.JSONData, &cfg)
	if err != nil {
		return nil, err
	}

	// IP and port as string
	ip_port := "host.docker.internal:8815" // temporarily hardcode
	fsc, err := NewFlightStreamClient(ip_port)
	if err != nil {
		return nil, err
	}

	ds := &Datasource{
		fsc: fsc,
	}

	return ds, nil
}

// Dispose here tells plugin SDK that plugin wants to clean up resources when a new instance
// created. As soon as datasource settings change detected by SDK old datasource instance will
// be disposed and a new one will be created using NewSampleDatasource factory function.
func (d *Datasource) Dispose() {
	// Clean up datasource instance resources.
	d.fsc = nil
}

// QueryData handles multiple queries and returns multiple responses.
// req contains the queries []DataQuery (where each query contains RefID as a unique identifier).
// The QueryDataResponse contains a map of RefID to the response for each query, and each response
// contains Frames ([]*Frame).
func (d *Datasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	// create response struct
	response := backend.NewQueryDataResponse()

	// loop over queries and execute them individually.
	for _, q := range req.Queries {
		res := d.query(ctx, req.PluginContext, q)

		// save the response in a hashmap
		// based on with RefID as identifier
		response.Responses[q.RefID] = res
	}

	return response, nil
}

type queryModel struct{}

func (d *Datasource) query(_ context.Context, pCtx backend.PluginContext, query backend.DataQuery) backend.DataResponse {
	var response backend.DataResponse

	// Unmarshal the JSON into our queryModel.
	var qm queryModel

	err := json.Unmarshal(query.JSON, &qm)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("json unmarshal: %v", err.Error()))
	}

	// create data frame response.
	// For an overview on data frames and how grafana handles them:
	// https://grafana.com/developers/plugin-tools/introduction/data-frames
	frame := data.NewFrame("response")

	// add fields.
	frame.Fields = append(frame.Fields,
		data.NewField("time", nil, []time.Time{query.TimeRange.From, query.TimeRange.To}),
		data.NewField("values", nil, []int64{10, 20}),
	)

	// add the frames to the response.
	response.Frames = append(response.Frames, frame)

	return response
}

// CheckHealth handles health checks sent from Grafana to the plugin.
// The main use case for these health checks is the test button on the
// datasource configuration page which allows users to verify that
// a datasource is working as expected.
func (d *Datasource) CheckHealth(_ context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	res := &backend.CheckHealthResult{}
	config, err := models.LoadPluginSettings(*req.PluginContext.DataSourceInstanceSettings)

	if err != nil {
		res.Status = backend.HealthStatusError
		res.Message = "Unable to load settings"
		return res, nil
	}

	if config.Secrets.ApiKey == "" {
		res.Status = backend.HealthStatusError
		res.Message = "API key is missing"
		return res, nil
	}

	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: "Data source is working",
	}, nil
}

func (d *Datasource) SubscribeStream(ctx context.Context, req *backend.SubscribeStreamRequest) (*backend.SubscribeStreamResponse, error) {
	log.DefaultLogger.Info("SubscribeStream for path: %s", req.Path)

	// path: /ds/<uuid>/<flight_name>
	// path_parts := strings.Split(req.Path, "/")

	// if len(path_parts) != 4 {
	// 	return &backend.SubscribeStreamResponse{
	// 		Status: backend.SubscribeStreamStatusPermissionDenied,
	// 	}, nil
	// }

	// flight_name := path_parts[3]
	flight_name := req.Path
	log.DefaultLogger.Info("flight_name: " + flight_name)

	// flight_name := req.Path
	_, err := d.fsc.GetFlightTicket(ctx, flight_name)
	if err != nil {
		return &backend.SubscribeStreamResponse{
			Status: backend.SubscribeStreamStatusNotFound,
		}, nil
	}

	return &backend.SubscribeStreamResponse{
		Status: backend.SubscribeStreamStatusOK,
	}, nil
}

func (d *Datasource) PublishStream(ctx context.Context, req *backend.PublishStreamRequest) (*backend.PublishStreamResponse, error) {
	log.DefaultLogger.Info("PublishStream for path: " + req.Path)

	return &backend.PublishStreamResponse{
		Status: backend.PublishStreamStatusPermissionDenied,
	}, nil
}

func (d *Datasource) RunStream(ctx context.Context, req *backend.RunStreamRequest, sender *backend.StreamSender) error {
	log.DefaultLogger.Info("RunStream for path: " + req.Path)

	// path_parts := strings.Split(req.Path, "/")

	// if len(path_parts) != 4 {
	// 	return fmt.Errorf("invalid path: %s", req.Path)
	// }

	// flight_name := path_parts[3]
	// log.DefaultLogger.Info("flight_name: " + flight_name)
	flight_name := req.Path

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			log.DefaultLogger.Info("trying fetch flight_name: " + flight_name)
			frame, err := d.fsc.GetStreamData(ctx, flight_name)
			log.DefaultLogger.Info("fetched flight_name: " + flight_name + "row count: " + strconv.Itoa(frame.Rows()))
			if err != nil {
				return err
			}

			if frame.Rows() > 0 {
				if err = sender.SendFrame(frame, data.IncludeAll); err != nil {
					return err
				}
			}

			time.Sleep(1 * time.Second)
		}
	}
}
