package receipt

import (
	"testing"
	"time"
)

func Test_purchaseDate_UnmarshalJSON(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name     string
		d        purchaseDate
		args     args
		wantErr  bool
		wantTime time.Time
	}{
		{
			name: "test example #1 parses and marshals correctly",
			d:    purchaseDate{},
			args: args{
				b: []byte(`"2022-01-01"`),
			},
			wantErr:  false,
			wantTime: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "test example #2 parses and marshals correctly",
			d:    purchaseDate{},
			args: args{
				b: []byte(`"2022-03-20"`),
			},
			wantErr:  false,
			wantTime: time.Date(2022, 3, 20, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.d.UnmarshalJSON(tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantTime.Equal(time.Time(tt.d)) {
				t.Errorf("incorrect time parsed, want = %v, got = %v", tt.wantTime, time.Time(tt.d))
			}
		})
	}
}

func Test_purchaseDate_scoreDay(t *testing.T) {
	tests := []struct {
		name string
		d    purchaseDate
		want int
	}{
		{
			name: "even day",
			d:    purchaseDate(time.Date(2022, 3, 20, 0, 0, 0, 0, time.UTC)),
			want: 0,
		},
		{
			name: "odd day",
			d:    purchaseDate(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)),
			want: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.scoreDay(); got != tt.want {
				t.Errorf("scoreDay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_purchaseTime_UnmarshalJSON(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name     string
		pt       purchaseTime
		args     args
		wantErr  bool
		wantTime time.Time
	}{
		{
			name: "test example #1 parses and marshals correctly",
			pt:   purchaseTime{},
			args: args{
				b: []byte(`"13:01"`),
			},
			wantErr:  false,
			wantTime: time.Date(0, 1, 1, 13, 1, 0, 0, time.UTC),
		},
		{
			name: "test example #2 parses and marshals correctly",
			pt:   purchaseTime{},
			args: args{
				b: []byte(`"14:33"`),
			},
			wantErr:  false,
			wantTime: time.Date(0, 1, 1, 14, 33, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.pt.UnmarshalJSON(tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantTime.Equal(time.Time(tt.pt)) {
				t.Errorf("incorrect time parsed, want = %v, got = %v", tt.wantTime, time.Time(tt.pt))
			}
		})
	}
}

func Test_purchaseTime_scoreTime(t *testing.T) {
	tests := []struct {
		name string
		pt   purchaseTime
		want int
	}{
		{
			name: "time between 2 and 4",
			pt:   purchaseTime(time.Date(0, 1, 1, 14, 33, 0, 0, time.UTC)),
			want: 10,
		},
		{
			name: "time not between 2 and 4",
			pt:   purchaseTime(time.Date(0, 1, 1, 8, 0, 0, 0, time.UTC)),
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pt.scoreTime(); got != tt.want {
				t.Errorf("scoreTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
