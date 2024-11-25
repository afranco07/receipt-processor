package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/afranco07/receipt-processor/database"
	"github.com/afranco07/receipt-processor/receipt"
)

func TestReceiptHandler_GetPointsForID(t *testing.T) {
	type fields struct {
		store store
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &ReceiptHandler{
				store: tt.fields.store,
			}
			h.GetPointsForID(tt.args.w, tt.args.r)
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
				store: db,
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
