package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/afranco07/receipt-processor/database"
	"github.com/afranco07/receipt-processor/receipt"
)

type store interface {
	Get(string) (int, error)
	Insert(receipt.Receipt, int) (string, error)
}

type errorMessage struct {
	Message string `json:"message"`
}

type ReceiptHandler struct {
	store store
}

func New(store store) ReceiptHandler {
	return ReceiptHandler{
		store: store,
	}
}

type getPointsResponse struct {
	Points int `json:"points"`
}

func (h *ReceiptHandler) GetPointsForID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)

	id := r.PathValue("id")
	if id == "" {
		log.Println("Missing id parameter")
		w.WriteHeader(http.StatusBadRequest)
		_ = enc.Encode(errorMessage{Message: "id is required"})
		return
	}

	score, err := h.store.Get(id)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			log.Printf("receipt with ID '%s' not found", id)
			w.WriteHeader(http.StatusNotFound)
			_ = enc.Encode(errorMessage{Message: fmt.Sprintf("receipt with ID '%s' not found", id)})
			return
		}

		log.Printf("error getting receipt with id %s: %v", id, err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = enc.Encode(errorMessage{Message: "something went wrong"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = enc.Encode(getPointsResponse{Points: score})
	return
}

type processReceiptResponse struct {
	Id string `json:"id"`
}

func (h *ReceiptHandler) ProcessReceipt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	dec := json.NewDecoder(r.Body)

	var rcpt receipt.Receipt
	if err := dec.Decode(&rcpt); err != nil {
		log.Printf("error marshalling receipt: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = enc.Encode(errorMessage{Message: "something went wrong"})
		return
	}

	score := rcpt.GetScore()

	id, err := h.store.Insert(rcpt, score)
	if err != nil {
		log.Printf("error inserting receipt into database: %v", err)
		if errors.Is(err, database.ErrReceiptAlreadyExists) {
			w.WriteHeader(http.StatusBadRequest)
			_ = enc.Encode(errorMessage{Message: "receipt has already been submitted"})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		_ = enc.Encode(errorMessage{Message: "something went wrong"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = enc.Encode(processReceiptResponse{Id: id})
}
