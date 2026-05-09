package util

import (
	"testing"
	"time"
)

func TestParseAndTimezoneTime(t *testing.T) {
	utc, _ := time.LoadLocation("UTC")
	seoul, _ := time.LoadLocation("Asia/Seoul")

	tests := []struct {
		name    string
		input   string
		tz      *time.Location
		wantY   int
		wantM   time.Month
		wantD   int
		wantH   int
		wantErr bool
	}{
		{"date only (8 chars)", "20240912", utc, 2024, time.September, 12, 0, false},
		{"date + hour (10 chars)", "2024091215", utc, 2024, time.September, 12, 15, false},
		{"date + hour + minute (12 chars)", "202409121530", utc, 2024, time.September, 12, 15, false},
		{"parses into target timezone", "2024091215", seoul, 2024, time.September, 12, 15, false},
		{"unparseable (other) aggregate", "(other)", utc, 0, 0, 0, 0, true},
		{"unparseable non-numeric", "foo", utc, 0, 0, 0, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseAndTimezoneTime(tt.input, tt.tz)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr = %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if got != nil {
					t.Errorf("expected nil time on error, got %v", got)
				}
				return
			}
			if got.Year() != tt.wantY || got.Month() != tt.wantM || got.Day() != tt.wantD {
				t.Errorf("date = %d-%s-%d, want %d-%s-%d",
					got.Year(), got.Month(), got.Day(), tt.wantY, tt.wantM, tt.wantD)
			}
			if got.Hour() != tt.wantH {
				t.Errorf("hour = %d, want %d", got.Hour(), tt.wantH)
			}
			if got.Location().String() != tt.tz.String() {
				t.Errorf("location = %s, want %s", got.Location(), tt.tz)
			}
		})
	}
}

func TestTimeArithmetic(t *testing.T) {
	base := time.Date(2024, 9, 12, 10, 30, 0, 0, time.UTC)

	if got := AddOneMinute(base); !got.Equal(base.Add(time.Minute)) {
		t.Errorf("AddOneMinute mismatch: got %v", got)
	}
	if got := AddOneHour(base); !got.Equal(base.Add(time.Hour)) {
		t.Errorf("AddOneHour mismatch: got %v", got)
	}
	if got := AddOneDay(base); !got.Equal(base.Add(24 * time.Hour)) {
		t.Errorf("AddOneDay mismatch: got %v", got)
	}
	if got := SubOneMinute(base); !got.Equal(base.Add(-time.Minute)) {
		t.Errorf("SubOneMinute mismatch: got %v", got)
	}
	if got := SubOneHour(base); !got.Equal(base.Add(-time.Hour)) {
		t.Errorf("SubOneHour mismatch: got %v", got)
	}
	if got := SubOneDay(base); !got.Equal(base.Add(-24 * time.Hour)) {
		t.Errorf("SubOneDay mismatch: got %v", got)
	}
}

func TestFillArray(t *testing.T) {
	out := FillArray(make([]string, 3), "x")
	if len(out) != 3 {
		t.Fatalf("length = %d, want 3", len(out))
	}
	for i, v := range out {
		if v != "x" {
			t.Errorf("index %d = %q, want %q", i, v, "x")
		}
	}
}

func TestTypeConverter(t *testing.T) {
	type sample struct {
		A int    `json:"a"`
		B string `json:"b"`
	}

	got, err := TypeConverter[sample](map[string]any{"a": 7, "b": "ok"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil || got.A != 7 || got.B != "ok" {
		t.Errorf("converted = %+v, want {A:7 B:ok}", got)
	}
}
