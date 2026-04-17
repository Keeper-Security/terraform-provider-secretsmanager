package secretsmanager

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/keeper-security/secrets-manager-go/core"
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

// TestCustomFieldsFromSchemaTypeNormalization verifies that custom field types
// are matched case-insensitively to the canonical vault API type strings.
//
// Bug: case-sensitive switch caused type="paymentcard" to fall through to default
// handler and be stored as core.Text instead of core.PaymentCards — silent data
// corruption. Also, type="" was silently stored instead of erroring.
//
// Fix: normalize user input to lowercase, look up canonical type string in a map,
// and return error for empty types.
func TestCustomFieldsFromSchemaTypeNormalization(t *testing.T) {
	tests := []struct {
		name             string
		fieldType        string
		fieldValue       string
		expectError      bool
		expectStructType string // "Text", "PaymentCards", "BankAccounts", etc.
	}{
		// Canonical cases (should work before and after fix)
		{
			name:             "canonical_text",
			fieldType:        "text",
			fieldValue:       "test_value",
			expectError:      false,
			expectStructType: "*core.Text",
		},
		{
			name:             "canonical_paymentCard",
			fieldType:        "paymentCard",
			fieldValue:       `{"cardNumber":"4111","cardExpirationDate":"12/25","cardSecurityCode":"123"}`,
			expectError:      false,
			expectStructType: "*core.PaymentCards",
		},
		{
			name:             "canonical_bankAccount",
			fieldType:        "bankAccount",
			fieldValue:       `{"accountType":"Checking","routingNumber":"123","accountNumber":"456"}`,
			expectError:      false,
			expectStructType: "*core.BankAccounts",
		},

		// Case variants (should work after fix)
		{
			name:             "uppercase_Text",
			fieldType:        "Text",
			fieldValue:       "test_value",
			expectError:      false,
			expectStructType: "*core.Text",
		},
		{
			name:             "lowercase_paymentcard",
			fieldType:        "paymentcard",
			fieldValue:       `{"cardNumber":"4111","cardExpirationDate":"12/25","cardSecurityCode":"123"}`,
			expectError:      false,
			expectStructType: "*core.PaymentCards",
		},
		{
			name:             "lowercase_bankaccount",
			fieldType:        "bankaccount",
			fieldValue:       `{"accountType":"Checking","routingNumber":"123","accountNumber":"456"}`,
			expectError:      false,
			expectStructType: "*core.BankAccounts",
		},
		{
			name:             "uppercase_PaymentCard",
			fieldType:        "PaymentCard",
			fieldValue:       `{"cardNumber":"4111","cardExpirationDate":"12/25","cardSecurityCode":"123"}`,
			expectError:      false,
			expectStructType: "*core.PaymentCards",
		},

		// Edge cases (should error after fix)
		{
			name:             "empty_type",
			fieldType:        "",
			fieldValue:       "test_value",
			expectError:      true,
			expectStructType: "",
		},
		{
			name:             "whitespace_only",
			fieldType:        "   ",
			fieldValue:       "test_value",
			expectError:      true,
			expectStructType: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items := []interface{}{
				map[string]interface{}{
					"type":     tt.fieldType,
					"label":    "TestField",
					"value":    tt.fieldValue,
					"required": false,
				},
			}
			result, err := customFieldsFromSchema(items)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error for type=%q, got nil", tt.fieldType)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error for type=%q: %v", tt.fieldType, err)
			}

			if len(result) == 0 {
				t.Errorf("expected result to contain field, got empty")
				return
			}

			// Type-assert to check the struct type
			field := result[0]
			actualType := fmt.Sprintf("%T", field)
			if actualType != tt.expectStructType {
				t.Errorf("type=%q: expected struct type %s, got %s", tt.fieldType, tt.expectStructType, actualType)
			}

			// For known types, verify base.Type is canonical
			if f, ok := field.(*core.Text); ok {
				if f.Type == "" {
					t.Errorf("type=%q: expected base.Type to be set (not empty)", tt.fieldType)
				}
			}
		})
	}
}

// TestCustomFieldTypeDiffSuppressFunc verifies that the DiffSuppressFunc on the
// custom field type attribute suppresses case-only differences, preventing
// perpetual diffs when users write non-canonical casing (e.g. "paymentcard").
//
// KSM-908 follow-up: the write path normalizes type to canonical casing before
// writing to vault, so the vault stores "paymentCard". The Read path returns the
// canonical form into state. Without DiffSuppressFunc, a config that still has
// "paymentcard" mismatches state on every subsequent plan.
func TestCustomFieldTypeDiffSuppressFunc(t *testing.T) {
	customSchema := schemaCustomField()
	typeElem := customSchema.Elem.(*schema.Resource).Schema["type"]

	if typeElem.DiffSuppressFunc == nil {
		t.Fatal("DiffSuppressFunc is nil — type attribute does not suppress case-only diffs")
	}

	fn := typeElem.DiffSuppressFunc
	tests := []struct {
		old  string
		new  string
		want bool
	}{
		{"paymentCard", "paymentcard", true},
		{"paymentCard", "PAYMENTCARD", true},
		{"text", "Text", true},
		{"text", "text", true},
		{"text", "secret", false},
		{"paymentCard", "bankAccount", false},
	}

	for _, tt := range tests {
		got := fn("custom.0.type", tt.old, tt.new, nil)
		if got != tt.want {
			t.Errorf("DiffSuppressFunc(%q, %q) = %v, want %v", tt.old, tt.new, got, tt.want)
		}
	}
}
