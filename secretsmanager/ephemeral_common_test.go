package secretsmanager

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/keeper-security/secrets-manager-go/core"
)

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
