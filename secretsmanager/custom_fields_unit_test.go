package secretsmanager

import (
	"encoding/json"
	"testing"
)

// TestEpochMsToDateUTC verifies that date conversion always uses UTC.
//
// The bug this guards against: using local time instead of UTC shifts
// midnight-UTC dates back by one day on machines in negative-offset timezones
// (e.g. UTC-8 US/Pacific), causing a perpetual Terraform diff on every plan.
//
// These test cases use midnight UTC so the expected value is unambiguous —
// any local timezone behind UTC would return the previous day without the fix.
func TestEpochMsToDateUTC(t *testing.T) {
	tests := []struct {
		name  string
		ms    float64
		want  string
	}{
		{
			// 2025-01-01T00:00:00Z — midnight UTC, New Year's Day.
			// UTC-1 through UTC-12 would return "2024-12-31" without .UTC().
			name: "midnight_utc_jan_1",
			ms:   1735689600000,
			want: "2025-01-01",
		},
		{
			// 2024-03-01T00:00:00Z — midnight UTC on first day of March in a leap year.
			// UTC-1 through UTC-12 would return "2024-02-29" without .UTC().
			name: "midnight_utc_march_1_leap_year",
			ms:   1709251200000,
			want: "2024-03-01",
		},
		{
			// 2025-07-15T00:00:00Z — mid-year date, positive sanity check.
			name: "midnight_utc_july_15",
			ms:   1752537600000,
			want: "2025-07-15",
		},
		{
			// 2023-12-31T23:59:59.999Z — just before midnight UTC, should stay in 2023.
			name: "just_before_midnight_utc",
			ms:   1704067199999,
			want: "2023-12-31",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := epochMsToDateUTC(tt.ms)
			if got != tt.want {
				t.Errorf("epochMsToDateUTC(%v) = %q, want %q", tt.ms, got, tt.want)
			}
		})
	}
}

// TestParseJSONItems verifies that parseJSONItems correctly distinguishes a single
// JSON object from a JSON array, returning one item in both cases but with the
// right structure so callers can unmarshal each entry uniformly.
func TestParseJSONItems(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantCount int
		wantErr   bool
	}{
		{
			name:      "single_object",
			input:     `{"region":"US","number":"555-1234","type":"Work"}`,
			wantCount: 1,
		},
		{
			name:      "array_one_entry",
			input:     `[{"region":"US","number":"555-1234","type":"Work"}]`,
			wantCount: 1,
		},
		{
			name:      "array_two_entries",
			input:     `[{"region":"US","number":"555-1234","type":"Work"},{"region":"US","number":"555-5678","type":"Mobile"}]`,
			wantCount: 2,
		},
		{
			name:      "empty_array",
			input:     `[]`,
			wantCount: 0,
		},
		{
			// Invalid JSON that starts with '[' is caught immediately by json.Unmarshal.
			// Invalid JSON that does NOT start with '[' is passed through as a raw message
			// and caught later by the caller's per-entry Unmarshal — not parseJSONItems itself.
			name:    "invalid_array_json",
			input:   `[not valid json]`,
			wantErr: true,
		},
		{
			name:    "empty_string",
			input:   ``,
			wantErr: true,
		},
		{
			name:    "whitespace_only",
			input:   `   `,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items, err := parseJSONItems(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("parseJSONItems(%q): expected error, got nil", tt.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("parseJSONItems(%q): unexpected error: %v", tt.input, err)
			}
			if len(items) != tt.wantCount {
				t.Errorf("parseJSONItems(%q): got %d items, want %d", tt.input, len(items), tt.wantCount)
			}
			// Verify each item is valid JSON
			for i, item := range items {
				if !json.Valid(item) {
					t.Errorf("item[%d] is not valid JSON: %s", i, item)
				}
			}
		})
	}
}
