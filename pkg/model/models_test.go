package model

import "testing"

func TestGetColumnType(t *testing.T) {
	numberTypes := []string{
		"TYPE_INTEGER", "TYPE_FLOAT", "TYPE_CURRENCY", "TYPE_MILLISECONDS", "TYPE_SECONDS",
		"CURRENCY", "INTEGER", "FLOAT", "PERCENT",
	}
	for _, ty := range numberTypes {
		if got := getColumnType(ty); got != ColumTypeNumber {
			t.Errorf("getColumnType(%q) = %q, want %q", ty, got, ColumTypeNumber)
		}
	}

	if got := getColumnType("TIME"); got != ColumTypeTime {
		t.Errorf("getColumnType(TIME) = %q, want %q", got, ColumTypeTime)
	}

	// Anything else (including GA string types) falls back to string.
	for _, ty := range []string{"STRING", "BOOLEAN", "UNKNOWN", ""} {
		if got := getColumnType(ty); got != ColumTypeString {
			t.Errorf("getColumnType(%q) = %q, want %q", ty, got, ColumTypeString)
		}
	}
}

func TestNewColumnDefinition(t *testing.T) {
	cd := NewColumnDefinition("activeUsers", 3, "TYPE_INTEGER")
	if cd.Header != "activeUsers" {
		t.Errorf("Header = %q, want %q", cd.Header, "activeUsers")
	}
	if cd.ColumnIndex != 3 {
		t.Errorf("ColumnIndex = %d, want 3", cd.ColumnIndex)
	}
	if cd.GetType() != ColumTypeNumber {
		t.Errorf("GetType() = %q, want %q", cd.GetType(), ColumTypeNumber)
	}
}
