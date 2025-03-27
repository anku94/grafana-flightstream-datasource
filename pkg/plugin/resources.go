package plugin

import (
	"context"
	"encoding/json"
	"net/http"
)

type StreamsResponse struct {
	Streams []string `json:"streams"`
}

func (d *Datasource) GetStreams(w http.ResponseWriter, r *http.Request) {
	streams, err := d.fsc.ListFlights(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(StreamsResponse{Streams: streams})
}
