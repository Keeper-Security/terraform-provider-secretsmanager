package secretsmanager

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	core "github.com/keeper-security/secrets-manager-go/core"
)

// PAM-specific field schema functions

func schemaCheckboxField() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "Checkbox field data.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"field_type": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Field type.",
				},
				"field_label": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Field label.",
				},
				"required": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "Required flag.",
				},
				"value": {
					Type:        schema.TypeList,
					Optional:    true,
					MaxItems:    1,
					Description: "Field value.",
					Elem:        &schema.Schema{Type: schema.TypeBool},
				},
			},
		},
	}
}

func schemaScriptField() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Script field data.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"field_type": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Field type.",
				},
				"field_label": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Field label.",
				},
				"required": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "Required flag.",
				},
				"privacy_screen": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "Privacy screen flag.",
				},
				"value": {
					Type:        schema.TypeList,
					Optional:    true,
					Description: "Script values.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"file_ref": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "File reference UID.",
							},
							"command": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "Script command.",
							},
							"record_ref": {
								Type:        schema.TypeList,
								Optional:    true,
								Description: "Record reference UIDs.",
								Elem:        &schema.Schema{Type: schema.TypeString},
							},
						},
					},
				},
			},
		},
	}
}

func schemaPamHostnameField() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "PAM Hostname field data.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"field_type": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Field type.",
				},
				"field_label": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Field label.",
				},
				"required": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "Required flag.",
				},
				"privacy_screen": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "Privacy screen flag.",
				},
				"value": {
					Type:        schema.TypeList,
					Optional:    true,
					Description: "Field value.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"hostname": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "Hostname or IP address.",
							},
							"port": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "Port number.",
							},
						},
					},
				},
			},
		},
	}
}

// suppressEquivalentJSON compares two JSON strings for structural equality
func suppressEquivalentJSON(k, oldValue, newValue string, d *schema.ResourceData) bool {
	if oldValue == newValue {
		return true
	}

	// Parse both JSON strings
	var oldJSON, newJSON interface{}
	if err := json.Unmarshal([]byte(oldValue), &oldJSON); err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(newValue), &newJSON); err != nil {
		return false
	}

	return reflect.DeepEqual(oldJSON, newJSON)
}

// schemaPamSettingsField returns the schema for PAM Settings field.
// This field contains protocol-specific connection configuration stored as JSON.
// The structure varies significantly by protocol (RDP, SSH, MySQL, PostgreSQL, etc.).
func schemaPamSettingsField() *schema.Schema {
	return &schema.Schema{
		Type:             schema.TypeString,
		Optional:         true,
		ValidateFunc:     validation.StringIsJSON,
		DiffSuppressFunc: suppressEquivalentJSON,
		Description: "PAM connection settings as JSON string. Structure varies by protocol:\n" +
			"- RDP: protocol, port, recordingIncludeKeys, security, ignoreCert, resizeMethod, enableFullWindowDrag, enableWallpaper, sftp\n" +
			"- SSH: protocol, port, recordingIncludeKeys, colorScheme, allowSupplyUser, hostKey, command, fontSize, sftp\n" +
			"- Database: protocol, port, recordingIncludeKeys, allowSupplyUser, database, allowSupplyHost\n" +
			"All protocols support portForward sub-object with port and reusePort fields.",
	}
}

// createPamSettingsFieldFromJSON creates a PAM settings field from JSON string
func createPamSettingsFieldFromJSON(jsonStr string) (interface{}, error) {
	if jsonStr == "" {
		return nil, nil
	}

	// Validate JSON
	var test interface{}
	if err := json.Unmarshal([]byte(jsonStr), &test); err != nil {
		return nil, fmt.Errorf("failed to parse pam_settings JSON: %w", err)
	}

	// Create a field struct that uses RawMessage to preserve exact JSON
	field := &struct {
		core.KeeperRecordField
		Value json.RawMessage `json:"value"`
	}{
		KeeperRecordField: core.KeeperRecordField{Type: "pamSettings"},
		Value:             json.RawMessage(jsonStr),
	}

	return field, nil
}

// pamSettingsFieldToJSON converts a PAM settings field to JSON string
func pamSettingsFieldToJSON(field interface{}) (string, error) {
	if field == nil {
		return "", nil
	}

	var value interface{}

	// Try as *core.PamSettings first
	if pamSettings, ok := field.(*core.PamSettings); ok {
		if pamSettings.Value == nil || len(pamSettings.Value) == 0 {
			return "", nil
		}
		value = pamSettings.Value
	} else if fieldMap, ok := field.(map[string]interface{}); ok {
		// Handle raw field map from GetFieldsByType()
		if v, found := fieldMap["value"]; found && v != nil {
			value = v
		} else {
			return "", nil
		}
	} else {
		return "", fmt.Errorf("field is not a PamSettings type or field map")
	}

	// Serialize to compact JSON
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		return "", fmt.Errorf("failed to serialize pam_settings to JSON: %w", err)
	}

	return string(jsonBytes), nil
}

func schemaPamResourcesField() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "PAM Resources field data.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"field_type": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Field type.",
				},
				"field_label": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Field label.",
				},
				"required": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "Required flag.",
				},
				"value": {
					Type:        schema.TypeList,
					Optional:    true,
					Description: "PAM resource values.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"controller_uid": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "Controller UID.",
							},
							"folder_uid": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "Folder UID.",
							},
							"resource_ref": {
								Type:        schema.TypeList,
								Optional:    true,
								Description: "Resource reference UIDs.",
								Elem:        &schema.Schema{Type: schema.TypeString},
							},
							"allowed_connections": {
								Type:        schema.TypeBool,
								Optional:    true,
								Description: "Allow connections.",
							},
							"allowed_port_forwards": {
								Type:        schema.TypeBool,
								Optional:    true,
								Description: "Allow port forwards.",
							},
							"allowed_rotation": {
								Type:        schema.TypeBool,
								Optional:    true,
								Description: "Allow rotation.",
							},
							"allowed_session_recording": {
								Type:        schema.TypeBool,
								Optional:    true,
								Description: "Allow session recording.",
							},
							"allowed_typescript_recording": {
								Type:        schema.TypeBool,
								Optional:    true,
								Description: "Allow typescript recording.",
							},
						},
					},
				},
			},
		},
	}
}

func schemaDatabaseTypeField() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "Database type field data.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"field_type": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Field type.",
				},
				"field_label": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Field label.",
				},
				"required": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "Required flag.",
				},
				"value": {
					Type:        schema.TypeList,
					Optional:    true,
					MaxItems:    1,
					Description: "Database type value.",
					Elem:        &schema.Schema{Type: schema.TypeString},
				},
			},
		},
	}
}

func schemaDirectoryTypeField() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "Directory type field data.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"field_type": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Field type.",
				},
				"field_label": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Field label.",
				},
				"required": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "Required flag.",
				},
				"value": {
					Type:        schema.TypeList,
					Optional:    true,
					MaxItems:    1,
					Description: "Directory type value.",
					Elem:        &schema.Schema{Type: schema.TypeString},
				},
			},
		},
	}
}

func schemaScheduleField() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "Schedule field data.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"field_type": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Field type.",
				},
				"field_label": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Field label.",
				},
				"required": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "Required flag.",
				},
				"value": {
					Type:        schema.TypeList,
					Optional:    true,
					MaxItems:    1,
					Description: "Schedule value.",
					Elem:        &schema.Schema{Type: schema.TypeString},
				},
			},
		},
	}
}

// schemaPrivatePemKeyField defines the schema for the "Private PEM Key" field
// on PAM record types (pamUser, pamMachine). This is a secret-type standard field
// that stores a PEM-encoded private key. Supports SSH key generation.
func schemaPrivatePemKeyField() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "Private PEM Key field data. Stored as a secret field labeled 'Private PEM Key'.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Field type.",
				},
				"label": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Field label.",
				},
				"generate": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Flag to force SSH key generation (when set to 'yes' or 'true').",
					ValidateDiagFunc: func(i interface{}, p cty.Path) diag.Diagnostics {
						v := i.(string)
						if v == "" || v == "true" || v == "yes" {
							return nil
						}
						return diag.Diagnostics{diag.Diagnostic{
							Severity:      diag.Error,
							Summary:       fmt.Sprintf("invalid generate = %s", v),
							Detail:        fmt.Sprintf("expected 'generate' to be one of ['true', 'yes', ''], got %s", v),
							AttributePath: p,
						}}
					},
				},
				"key_type": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "ssh-ed25519",
					Description: "SSH key type. One of: ssh-ed25519 (default), ssh-rsa, ecdsa-sha2-nistp256, ecdsa-sha2-nistp384, ecdsa-sha2-nistp521.",
					ValidateDiagFunc: func(i interface{}, p cty.Path) diag.Diagnostics {
						v := i.(string)
						valid := []string{"ssh-ed25519", "ssh-rsa", "ecdsa-sha2-nistp256", "ecdsa-sha2-nistp384", "ecdsa-sha2-nistp521"}
						for _, s := range valid {
							if v == s {
								return nil
							}
						}
						return diag.Diagnostics{diag.Diagnostic{
							Severity:      diag.Error,
							Summary:       fmt.Sprintf("invalid key_type = %s", v),
							Detail:        fmt.Sprintf("expected 'key_type' to be one of %v, got %s", valid, v),
							AttributePath: p,
						}}
					},
				},
				"key_bits": {
					Type:        schema.TypeInt,
					Optional:    true,
					Default:     4096,
					Description: "Key size in bits. Only used for ssh-rsa. Valid: 2048, 3072, 4096.",
					ValidateDiagFunc: func(i interface{}, p cty.Path) diag.Diagnostics {
						v := i.(int)
						if v == 2048 || v == 3072 || v == 4096 {
							return nil
						}
						return diag.Diagnostics{diag.Diagnostic{
							Severity:      diag.Error,
							Summary:       fmt.Sprintf("invalid key_bits = %d", v),
							Detail:        fmt.Sprintf("expected 'key_bits' to be one of [2048, 3072, 4096], got %d", v),
							AttributePath: p,
						}}
					},
				},
				"required": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "Required flag.",
				},
				"privacy_screen": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "Privacy screen flag.",
				},
				"value": {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    true,
					Sensitive:   true,
					Description: "Private key in PEM format. Computed when generate is set.",
				},
				"public_key": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Public key in OpenSSH format. Computed when generate is set.",
				},
			},
		},
	}
}

// schemaPrivateKeyPassphraseField defines the schema for the "Private Key Passphrase"
// custom field on PAM record types. This is stored as a custom secret field.
func schemaPrivateKeyPassphraseField() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "Private Key Passphrase. Stored as a custom field labeled 'Private Key Passphrase'.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Field type.",
				},
				"generate": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Flag to force passphrase generation (when set to 'yes' or 'true').",
					ValidateDiagFunc: func(i interface{}, p cty.Path) diag.Diagnostics {
						v := i.(string)
						if v == "" || v == "true" || v == "yes" {
							return nil
						}
						return diag.Diagnostics{diag.Diagnostic{
							Severity:      diag.Error,
							Summary:       fmt.Sprintf("invalid generate = %s", v),
							Detail:        fmt.Sprintf("expected 'generate' to be one of ['true', 'yes', ''], got %s", v),
							AttributePath: p,
						}}
					},
				},
				"complexity": {
					Type:        schema.TypeList,
					Optional:    true,
					MaxItems:    1,
					Description: "Passphrase complexity.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"length": {
								Type:        schema.TypeInt,
								Optional:    true,
								Description: "Passphrase length.",
							},
							"caps": {
								Type:        schema.TypeInt,
								Optional:    true,
								Description: "Number of uppercase characters.",
							},
							"lowercase": {
								Type:        schema.TypeInt,
								Optional:    true,
								Description: "Number of lowercase characters.",
							},
							"digits": {
								Type:        schema.TypeInt,
								Optional:    true,
								Description: "Number of digits.",
							},
							"special": {
								Type:        schema.TypeInt,
								Optional:    true,
								Description: "Number of special characters.",
							},
						},
					},
				},
				"value": {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    true,
					Sensitive:   true,
					Description: "Passphrase value. Computed when generate is set.",
				},
			},
		},
	}
}
