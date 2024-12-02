package receipt

import "testing"

func TestItem_scoreDescription(t *testing.T) {
	type fields struct {
		ShortDescription string
		Price            string
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "test example #1 mountain dew",
			fields: fields{
				ShortDescription: "Mountain Dew 12PK",
				Price:            "6.49",
			},
			want: 0,
		},
		{
			name: "test example #1 emils cheese pizza",
			fields: fields{
				ShortDescription: "Emils Cheese Pizza",
				Price:            "12.25",
			},
			want: 3,
		},
		{
			name: "test example #1 Klarbrunn",
			fields: fields{
				ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ",
				Price:            "12.00",
			},
			want: 3,
		},
		{
			name: "test example #2 gatorade",
			fields: fields{
				ShortDescription: "Gatorade",
				Price:            "2.25",
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := Item{
				ShortDescription: tt.fields.ShortDescription,
				Price:            tt.fields.Price,
			}
			if got := i.scoreDescription(); got != tt.want {
				t.Errorf("scoreDescription() = %v, want %v", got, tt.want)
			}
		})
	}
}
