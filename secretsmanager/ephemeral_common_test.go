package secretsmanager

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/keeper-security/secrets-manager-go/core"
)

// Regression test for dateFieldToString: must use UTC and ParseInt (not Atoi) to avoid
// wrong-day output in negative-offset timezones and int32 overflow on 32-bit platforms.
func TestDateFieldToString(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		// 2024-06-15 00:00:00 UTC in epoch ms — in UTC-4 this is 2024-06-14 20:00 local;
		// without .UTC() the result would be "2024-06-14", not "2024-06-15".
		{"UTC date regression", "1718409600000", "2024-06-15"},
		// Empty input returns empty string.
		{"empty", "", ""},
		// Whitespace-only returns empty string.
		{"whitespace", "   ", ""},
		// Non-numeric passthrough.
		{"non-numeric", "not-a-date", "not-a-date"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := dateFieldToString(tc.input)
			if got != tc.want {
				t.Errorf("dateFieldToString(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

// Regression test for KSM-884: pamHostnameToListValue read vmap["hostname"] (all-lowercase)
// instead of vmap["hostName"] (camelCase), so host_name was always empty.
// The KSM API wire format uses "hostName". This test fails on the buggy code and passes after fix.
func TestPamHostnameToListValue_ReadsHostNameKey(t *testing.T) {
	record := &core.Record{
		RecordDict: map[string]interface{}{
			"type": "pamMachine",
			"fields": []interface{}{
				map[string]interface{}{
					"type": "pamHostname",
					"value": []interface{}{
						map[string]interface{}{
							"hostName": "myhost.example.com",
							"port":     "8022",
						},
					},
				},
			},
		},
	}

	list, diags := pamHostnameToListValue(context.Background(), record)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	elements := list.Elements()
	if len(elements) != 1 {
		t.Fatalf("expected 1 element, got %d", len(elements))
	}

	obj, ok := elements[0].(types.Object)
	if !ok {
		t.Fatalf("expected types.Object, got %T", elements[0])
	}
	attrs := obj.Attributes()

	hostName := attrs["host_name"].(types.String).ValueString()
	if hostName != "myhost.example.com" {
		t.Errorf("host_name = %q, want %q", hostName, "myhost.example.com")
	}

	port := attrs["port"].(types.String).ValueString()
	if port != "8022" {
		t.Errorf("port = %q, want %q", port, "8022")
	}
}

// Verify empty-field path still returns an empty list (not a crash).
func TestPamHostnameToListValue_EmptyField(t *testing.T) {
	record := &core.Record{
		RecordDict: map[string]interface{}{
			"type":   "pamMachine",
			"fields": []interface{}{},
		},
	}

	list, diags := pamHostnameToListValue(context.Background(), record)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if len(list.Elements()) != 0 {
		t.Errorf("expected empty list, got %d elements", len(list.Elements()))
	}
}
