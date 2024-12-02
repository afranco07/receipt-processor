package main

import (
	"log"
	"net/http"

	"github.com/afranco07/receipt-processor/database"
	"github.com/afranco07/receipt-processor/handler"
)

func main() {
	db := database.NewInMemoryDatabase()
	receiptHandler := handler.New(db)

	http.HandleFunc("GET /receipts/{id}/points", receiptHandler.GetPointsForID)
	http.HandleFunc("POST /receipts/process", receiptHandler.ProcessReceipt)

	log.Println("Starting server on port :8080…")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
