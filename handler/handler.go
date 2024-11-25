package handler

import (
	"encoding/json"
	"errors"
	"github.com/afranco07/receipt-processor/database"
	"github.com/afranco07/receipt-processor/receipt"
	"log"
	"net/http"
)

type store interface {
	Get(string) (int, error)
	Insert(receipt.Receipt, int) (string, error)
}

type ReceiptHandler struct {
	store store
}

func New(store store) ReceiptHandler {
	return ReceiptHandler{
		store: store,
	}
}

func (h *ReceiptHandler) GetPointsForID(w http.ResponseWriter, r *http.Request) {}

func (h *ReceiptHandler) ProcessReceipt(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)

	var rcpt receipt.Receipt
	if err := dec.Decode(&rcpt); err != nil {
		log.Printf("error marshalling receipt: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{ "message": "something went wrong" }`))
		return
	}

	score := rcpt.GetScore()

	id, err := h.store.Insert(rcpt, score)
	if err != nil {
		log.Printf("error inserting receipt into database: %v", err)
		if errors.Is(err, database.ErrReceiptAlreadyExists) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{ "message": "receipt has already been submitted" }`))
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{ "message": "something went wrong" }`))
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(`{ "id": "` + id + `" }`))
}
