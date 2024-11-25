package main

import (
	"github.com/afranco07/receipt-processor/database"
	"github.com/afranco07/receipt-processor/handler"
	"net/http"
)

func main() {
	db := database.NewInMemoryDatabase()
	receiptHandler := handler.New(db)

	http.HandleFunc("GET /receipts/{id}/points", receiptHandler.GetPointsForID)
	http.HandleFunc("POST /receipts/process", receiptHandler.ProcessReceipt)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
