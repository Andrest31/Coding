package main

import (
	"encoding/json"
	"net/http"

	"st/coding"
	"st/norm"

	"golang.org/x/exp/rand"
)

type SegmentRequest struct {
	Segment string `json:"segment"`
}

type SegmentResponse struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

func SegmentHandler(w http.ResponseWriter, r *http.Request) {
	var req SegmentRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Невалидный JSON", http.StatusBadRequest)
		return
	}

	if req.Segment == "" {
		http.Error(w, "Пустой сегмент", http.StatusBadRequest)
		return
	}

	result, err := coding.ProcessMessage(req.Segment, rand.NewSource(123))
	if err != nil {
		resp := SegmentResponse{Error: err.Error()}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := SegmentResponse{Result: result}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
