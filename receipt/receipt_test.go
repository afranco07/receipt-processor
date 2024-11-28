package receipt

import (
	"testing"
	"time"
)

func TestReceipt_GetScore(t *testing.T) {
	targetPurchaseDate, _ := time.Parse(time.DateOnly, "2022-01-01")
	targetPurchaseTime, _ := time.Parse(timeOnly, "13:01")

	mmPurchaseDate, _ := time.Parse(time.DateOnly, "2022-03-20")
	mmPurchaseTime, _ := time.Parse(timeOnly, "14:33")

	type fields struct {
		Retailer     string
		PurchaseDate purchaseDate
		PurchaseTime purchaseTime
		Items        []Item
		Total        string
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "Target returns 28 points and wins",
			fields: fields{
				Retailer:     "Target",
				PurchaseDate: purchaseDate(targetPurchaseDate),
				PurchaseTime: purchaseTime(targetPurchaseTime),
				Items: []Item{
					{
						ShortDescription: "Mountain Dew 12PK",
						Price:            "6.49",
					},
					{
						ShortDescription: "Emils Cheese Pizza",
						Price:            "12.25",
					},
					{
						ShortDescription: "Knorr Creamy Chicken",
						Price:            "1.26",
					},
					{
						ShortDescription: "Doritos Nacho Cheese",
						Price:            "3.35",
					},
					{
						ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ",
						Price:            "12.00",
					},
				},
				Total: "35.35",
			},
			want: 28,
		},
		{
			name: "M&M Corner Market receipt and you miss",
			fields: fields{
				Retailer:     "M&M Corner Market",
				PurchaseDate: purchaseDate(mmPurchaseDate),
				PurchaseTime: purchaseTime(mmPurchaseTime),
				Items: []Item{
					{
						ShortDescription: "Gatorade",
						Price:            "2.25",
					},
					{
						ShortDescription: "Gatorade",
						Price:            "2.25",
					},
					{
						ShortDescription: "Gatorade",
						Price:            "2.25",
					},
					{
						ShortDescription: "Gatorade",
						Price:            "2.25",
					},
				},
				Total: "9.00",
			},
			want: 109,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Receipt{
				Retailer:     tt.fields.Retailer,
				PurchaseDate: tt.fields.PurchaseDate,
				PurchaseTime: tt.fields.PurchaseTime,
				Items:        tt.fields.Items,
				Total:        tt.fields.Total,
			}
			if got, _ := r.GetScore(); got != tt.want {
				t.Errorf("GetScore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReceipt_scoreItems(t *testing.T) {
	type fields struct {
		Items []Item
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "5 items (2 pairs @ 5 points each)",
			fields: fields{
				Items: []Item{
					{
						ShortDescription: "Mountain Dew 12PK",
						Price:            "6.49",
					},
					{
						ShortDescription: "Emils Cheese Pizza",
						Price:            "12.25",
					},
					{
						ShortDescription: "Knorr Creamy Chicken",
						Price:            "1.26",
					},
					{
						ShortDescription: "Doritos Nacho Cheese",
						Price:            "3.35",
					},
					{
						ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ",
						Price:            "12.00",
					},
				},
			},
			want: 10,
		},
		{
			name: "4 items (2 pairs @ 5 points each)",
			fields: fields{
				Items: []Item{
					{
						ShortDescription: "Gatorade",
						Price:            "2.25",
					},
					{
						ShortDescription: "Gatorade",
						Price:            "2.25",
					},
					{
						ShortDescription: "Gatorade",
						Price:            "2.25",
					},
					{
						ShortDescription: "Gatorade",
						Price:            "2.25",
					},
					{
						ShortDescription: "Gatorade",
						Price:            "2.25",
					},
				},
			},
			want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Receipt{
				Items: tt.fields.Items,
			}
			if got := r.scoreItems(); got != tt.want {
				t.Errorf("scoreItems() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReceipt_scoreRetailer(t *testing.T) {
	type fields struct {
		Retailer string
	}

	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "all alphanumeric retailer",
			fields: fields{
				Retailer: "Target",
			},
			want: 6,
		},
		{
			name: "all alphanumeric retailer",
			fields: fields{
				Retailer: "M&M Corner Market",
			},
			want: 14,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Receipt{
				Retailer: tt.fields.Retailer,
			}
			if got := r.scoreRetailer(); got != tt.want {
				t.Errorf("scoreRetailer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReceipt_scoreTotal(t *testing.T) {
	type fields struct {
		Total string
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "total not divisible by 0.25 and not round dollar amount",
			fields: fields{
				Total: "35.35",
			},
			want: 0,
		},
		{
			name: "total is divisible by 0.25 and round dollar amount",
			fields: fields{
				Total: "9.00",
			},
			want: 75,
		},
		{
			name: "total is divisible by 0.25 and not round dollar amount",
			fields: fields{
				Total: "2.50",
			},
			want: 25,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Receipt{
				Total: tt.fields.Total,
			}
			if got, _ := r.scoreTotal(); got != tt.want {
				t.Errorf("scoreTotal() = %v, want %v", got, tt.want)
			}
		})
	}
}
