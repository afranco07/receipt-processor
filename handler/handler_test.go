package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/afranco07/receipt-processor/database"
	"github.com/afranco07/receipt-processor/receipt"
	"github.com/go-playground/validator/v10"
)

func TestReceiptHandler_GetPointsForID(t *testing.T) {
	db := database.NewInMemoryDatabase()
	// adding this file to test already submitted receipt
	testFile, err := os.ReadFile("../examples/test-receipt.json")
	if err != nil {
		t.Fatal(err)
	}

	var rcpt receipt.Receipt
	err = json.Unmarshal(testFile, &rcpt)
	if err != nil {
		t.Fatal(err)
	}

	score := 10
	id, err := db.Insert(rcpt, score)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name             string
		id               string
		score            int
		expectScore      int
		expectStatusCode int
		setupDB          func() store
	}{
		{
			name:             "successful get",
			id:               id,
			score:            score,
			expectScore:      10,
			expectStatusCode: http.StatusOK,
		},
		{
			name:             "id not found",
			id:               "does-not-exist",
			score:            score,
			expectScore:      -1,
			expectStatusCode: http.StatusNotFound,
		},
		{
			name:             "id path value missing",
			id:               "",
			score:            score,
			expectScore:      -1,
			expectStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &ReceiptHandler{
				store: db,
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/receipts/%s/points", tt.id), nil)
			r.SetPathValue("id", tt.id)

			h.GetPointsForID(w, r)

			if w.Result().StatusCode != tt.expectStatusCode {
				t.Errorf("the response status code did not match. Got %d, want %d", w.Result().StatusCode, tt.expectStatusCode)
			}

			err = json.NewEncoder(w).Encode(rcpt)
			if err != nil {
				t.Fatal(err)
			}

			if tt.expectScore >= 0 && tt.expectScore != tt.score {
				t.Errorf("the response score did not match. Got %d, want %d", w.Result().Body, tt.score)
			}
		})
	}
}

func TestReceiptHandler_ProcessReceipt(t *testing.T) {
	db := database.NewInMemoryDatabase()
	// adding this file to test already submitted receipt
	testFile, err := os.ReadFile("../examples/test-receipt.json")
	if err != nil {
		t.Fatal(err)
	}

	var rcpt receipt.Receipt
	err = json.Unmarshal(testFile, &rcpt)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Insert(rcpt, 10)
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		body string
	}
	tests := []struct {
		name             string
		args             args
		filepath         string
		expectResponse   string
		expectStatusCode int
	}{
		{
			name:             "test example #1",
			filepath:         "../examples/morning-receipt.json",
			expectStatusCode: http.StatusCreated,
		},
		{
			name:             "test example #2",
			filepath:         "../examples/simple-receipt.json",
			expectStatusCode: http.StatusCreated,
		},
		{
			name:             "attempt to add a receipt that has already been inserted previously",
			filepath:         "../examples/test-receipt.json",
			expectStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := os.Open(tt.filepath)
			if err != nil {
				t.Fatal(err)
			}

			h := &ReceiptHandler{
				store:     db,
				validator: validator.New(validator.WithRequiredStructEnabled()),
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/receipts/process", f)
			h.ProcessReceipt(w, r)

			if w.Result().StatusCode != tt.expectStatusCode {
				t.Errorf("the response status code did not match. Got %d, want %d", w.Result().StatusCode, tt.expectStatusCode)
			}
		})
	}
}
