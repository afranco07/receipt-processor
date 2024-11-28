package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/afranco07/receipt-processor/database"
	"github.com/afranco07/receipt-processor/receipt"
	"github.com/go-playground/validator/v10"
)

type store interface {
	Get(string) (int, error)
	Insert(receipt.Receipt, int) (string, error)
}

type errorMessage struct {
	Message string `json:"message"`
}

type ReceiptHandler struct {
	store     store
	validator *validator.Validate
}

func New(store store) ReceiptHandler {

	return ReceiptHandler{
		store:     store,
		validator: validator.New(validator.WithRequiredStructEnabled()),
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
		var timeErr *time.ParseError
		if errors.As(err, &timeErr) {
			w.WriteHeader(http.StatusBadRequest)
			_ = enc.Encode(errorMessage{Message: fmt.Sprintf("%s is not a valid value", timeErr.Value)})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		_ = enc.Encode(errorMessage{Message: "something went wrong"})
		return
	}

	validationErrors, err := rcpt.ValidateReceipt(h.validator)
	if err != nil {
		log.Printf("error validating receipt: %v", validationErrors)
		w.WriteHeader(http.StatusBadRequest)
		_ = enc.Encode(errorMessage{Message: err.Error()})
		return
	}

	score, err := rcpt.GetScore()
	if err != nil {
		log.Printf("error getting receipt score: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = enc.Encode(errorMessage{Message: "something went wrong"})
		return
	}

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

func (h *ReceiptHandler) validateReceipt(rcpt receipt.Receipt) (validator.ValidationErrors, error) {
	err := h.validator.Struct(rcpt)
	if err == nil {
		return nil, nil
	}

	var validationErrors validator.ValidationErrors
	e := errors.As(err, &validationErrors)
	if !e {
		return nil, nil
	}

	var fields string
	for _, v := range validationErrors {
		fields += v.Field() + ", "
	}

	return validationErrors, errors.New(fmt.Sprintf("validation errors for the following fields: %s", fields))
}
