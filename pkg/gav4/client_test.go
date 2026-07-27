package gav4

import "testing"

func TestMetadataResourceName(t *testing.T) {
	cases := []struct {
		name       string
		propertyID string
		want       string
	}{
		{"bare numeric id", "390931478", "properties/390931478/metadata"},
		// WebPropertyID is stored as the full resource name everywhere else
		// (account summaries, RunReport) — getMetadata must not double the
		// "properties/" prefix when it's already present.
		{"full resource name", "properties/390931478", "properties/390931478/metadata"},
		{"empty falls back to 0", "", "properties/0/metadata"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := metadataResourceName(tc.propertyID); got != tc.want {
				t.Errorf("metadataResourceName(%q) = %q, want %q", tc.propertyID, got, tc.want)
			}
		})
	}
}
