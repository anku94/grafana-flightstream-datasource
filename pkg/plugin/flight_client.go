package plugin

import (
	"context"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/apache/arrow/go/v17/arrow/flight"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GrpcDialOptions() ([]grpc.DialOption, error) {
	transport := grpc.WithTransportCredentials(insecure.NewCredentials())
	opts := []grpc.DialOption{
		transport,
	}

	return opts, nil
}

func NewFlightClient(ip_port string) (flight.Client, error) {
	dialOptions, err := GrpcDialOptions()
	if err != nil {
		return nil, fmt.Errorf("grpc dial options: %s", err)
	}

	client, err := flight.NewClientWithMiddleware(ip_port, nil, nil, dialOptions...)
	if err != nil {
		return nil, err
	}

	return client, nil
}

type FlightStreamClient struct {
	mu      sync.Mutex
	client  flight.Client
	tickets map[string]*flight.Ticket
}

func NewFlightStreamClient(ip_port string) (*FlightStreamClient, error) {
	client, err := NewFlightClient(ip_port)
	if err != nil {
		return nil, err
	}

	return &FlightStreamClient{
		client:  client,
		tickets: make(map[string]*flight.Ticket),
	}, nil
}

func (c *FlightStreamClient) ListFlights(ctx context.Context) ([]string, error) {
	stream, err := c.client.ListFlights(ctx, &flight.Criteria{})
	if err != nil {
		return nil, err
	}

	var flight_names []string
	for {
		info, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		flight_name := strings.Join(info.FlightDescriptor.Path, "/")
		flight_names = append(flight_names, flight_name)
	}

	return flight_names, nil
}

func (c *FlightStreamClient) GetFlightInfo(ctx context.Context, flight_name string) (*flight.FlightInfo, error) {
	desc := &flight.FlightDescriptor{
		Type: flight.DescriptorPATH,
		Path: []string{flight_name},
	}

	info, err := c.client.GetFlightInfo(ctx, desc)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (c *FlightStreamClient) GetFlightTicket(ctx context.Context, flight_name string) (*flight.Ticket, error) {
	// we cache flight tickets, return from there if exists
	log.DefaultLogger.Info("GetFlightTicket for flight_name: " + flight_name)

	c.mu.Lock()
	ticket, ok := c.tickets[flight_name]
	c.mu.Unlock()

	if ok {
		log.DefaultLogger.Info("GetFlightTicket for flight_name: " + flight_name + " found in cache")
		return ticket, nil
	}

	log.DefaultLogger.Info("GetFlightTicket for flight_name: " + flight_name + " not found in cache, fetching")

	info, err := c.GetFlightInfo(ctx, flight_name)
	if err != nil {
		log.DefaultLogger.Error("GetFlightTicket for flight_name: " + flight_name + " error: " + err.Error())
		return nil, err
	}

	if len(info.Endpoint) == 0 {
		log.DefaultLogger.Error("GetFlightTicket for flight_name: " + flight_name + " no endpoints found")
		return nil, fmt.Errorf("no endpoints found for flight %s", flight_name)
	}

	tkt := info.Endpoint[0].Ticket

	c.mu.Lock()
	c.tickets[flight_name] = tkt
	c.mu.Unlock()

	return tkt, nil
}

func (c *FlightStreamClient) GetStreamData(ctx context.Context, flight_name string) (*data.Frame, error) {
	ticket, err := c.GetFlightTicket(ctx, flight_name)
	if err != nil {
		return nil, err
	}

	stream, err := c.client.DoGet(ctx, ticket)
	if err != nil {
		return nil, err
	}

	reader, err := flight.NewRecordReader(stream)
	if err != nil {
		return nil, err
	}

	frame, err := frameForRecords(reader)
	if err != nil {
		return nil, err
	}

	if frame == nil {
		return nil, fmt.Errorf("no data found for flight %s", flight_name)
	}

	return frame, nil
}

func (c *FlightStreamClient) InvalidateTicket(flight_name string) {
	c.mu.Lock()
	delete(c.tickets, flight_name)
	c.mu.Unlock()
}
