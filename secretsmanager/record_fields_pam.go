package secretsmanager

import (
	"encoding/json"
	"fmt"

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
				"type": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Field type.",
				},
				"label": {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    true,
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
				"type": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Field type.",
				},
				"label": {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    true,
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
				"type": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Field type.",
				},
				"label": {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    true,
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

// suppressEquivalentJSON compares two JSON strings semantically, ignoring field order and whitespace.
// Returns true if the JSON is semantically equivalent (suppresses diff), false otherwise.
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

	// Re-marshal to normalized JSON (consistent ordering)
	oldNormalized, err1 := json.Marshal(oldJSON)
	newNormalized, err2 := json.Marshal(newJSON)
	if err1 != nil || err2 != nil {
		return false
	}

	// Compare normalized JSON
	return string(oldNormalized) == string(newNormalized)
}

// schemaPamSettingsField returns the schema for PAM Settings field.
// This field contains protocol-specific connection configuration stored as JSON.
// The structure varies significantly by protocol (RDP, SSH, MySQL, PostgreSQL, etc.).
//
// NOTE: Using JSON string instead of typed struct because:
// 1. Go SDK's PamSettings struct is incomplete (missing 15+ fields)
// 2. Field structure varies drastically by protocol (RDP has 9 fields, SSH has 9 different fields, database has 6 fields)
// 3. Backend stores as encrypted JSON blob without validation
// 4. Prevents data loss on round-trip operations
// 5. Forward-compatible with new protocols and fields
//
// See PAM_SCHEMA_ANALYSIS.md for complete field definitions by protocol.
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

// schemaPamResourcesField returns the schema for PAM Resources field.
// NOTE: This field is for PAM Configuration records, NOT for pamUser/pamMachine/pamDatabase/pamDirectory.
//
// This function is reserved for future PAM Configuration record implementation (RT_PAM_CONFIGURATION).
func schemaPamResourcesField() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "PAM Resources field data.",
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
					Computed:    true,
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
							// NOTE: The following allowed_* fields are stored in DAG (access control system)
							// and are NOT accessible via KSM API. They cannot be read or written through
							// this Terraform provider. Removed to avoid confusion.
							// - allowed_connections
							// - allowed_port_forwards
							// - allowed_rotation
							// - allowed_session_recording
							// - allowed_typescript_recording
						},
					},
				},
			},
		},
	}
}

// schemaDatabaseTypeField returns the schema for Database Type field.
// Currently used in resource_pam_database.go.
//
// Supported values are based on Keeper's PAM connection protocols:
// - postgresql: PostgreSQL (port 5432)
// - mysql: MySQL (port 3306)
// - mariadb: MariaDB (port 3306)
// - mariadb-flexible: Azure MariaDB Flexible Server (port 3306)
// - mssql: Microsoft SQL Server (port 1433)
// - oracle: Oracle Database (port 1521)
// - mongodb: MongoDB (port 27017)
func schemaDatabaseTypeField() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		ValidateFunc: validation.StringInSlice([]string{
			"postgresql",
			"mysql",
			"mariadb",
			"mariadb-flexible",
			"mssql",
			"oracle",
			"mongodb",
		}, false),
		Description: "Database type. Must be one of: postgresql, mysql, mariadb, mariadb-flexible, " +
			"mssql, oracle, mongodb. Invalid values will render the connection unusable.",
	}
}

// schemaDirectoryTypeField returns the schema for Directory Type field.
// NOTE: Currently not used in any PAM resource but mapped in provider.go.
// May be needed for pamUser or pamDirectory record types in the future.
//
// Supported values are based on Keeper's PAM directory protocols:
// - Active Directory: Microsoft Active Directory (port 636 LDAPS required)
// - OpenLDAP: OpenLDAP directory service (port 389 or 636)
func schemaDirectoryTypeField() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		ValidateFunc: validation.StringInSlice([]string{
			"Active Directory",
			"OpenLDAP",
		}, false),
		Description: "Directory type. Must be one of: 'Active Directory', 'OpenLDAP'. " +
			"Invalid values will render the connection unusable.",
	}
}

// schemaScheduleField returns the schema for Schedule field (rotation schedules).
// NOTE: Currently not used in any PAM resource but mapped in provider.go.
// May be needed for rotation schedule configuration in PAM records.
//
// Go SDK Schedule structure (secrets-manager-go/core/record_data.go):
// - Type: Schedule type (e.g., "WEEKLY", "MONTHLY", "DAILY")
// - Cron: Cron expression for complex schedules
// - Time: Time of day (e.g., "02:00")
// - Tz: Timezone (e.g., "America/New_York")
// - Weekday: Day of week for weekly schedules (e.g., "Sunday")
// - IntervalCount: Interval count for recurring schedules
func schemaScheduleField() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "Schedule field data for rotation schedules.",
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
					Computed:    true,
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
					Description: "Schedule configuration.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"type": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "Schedule type (e.g., WEEKLY, MONTHLY, DAILY).",
							},
							"cron": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "Cron expression for complex schedules.",
							},
							"time": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "Time of day (e.g., 02:00).",
							},
							"tz": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "Timezone (e.g., America/New_York).",
							},
							"weekday": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "Day of week for weekly schedules (e.g., Sunday).",
							},
							"interval_count": {
								Type:        schema.TypeInt,
								Optional:    true,
								Description: "Interval count for recurring schedules.",
							},
						},
					},
				},
			},
		},
	}
}

// createPamSettingsFieldFromJSON creates a pamSettings field from a JSON string.
// Uses json.RawMessage to preserve exact JSON bytes, bypassing the incomplete PamSetting struct.
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
	// This bypasses the typed PamSetting struct which drops unknown fields
	field := &struct {
		core.KeeperRecordField
		Value json.RawMessage `json:"value"`
	}{
		KeeperRecordField: core.KeeperRecordField{Type: "pamSettings"},
		Value:             json.RawMessage(jsonStr),
	}

	return field, nil
}

// pamSettingsFieldToJSON converts a pamSettings field to a JSON string.
// Returns the JSON string representation of the field's value, or empty string if no value.
// Handles both *core.PamSettings objects and raw field maps from GetFieldsByType().
// The returned JSON is compact (no whitespace) for consistent comparison.
func pamSettingsFieldToJSON(field interface{}) (string, error) {
	if field == nil {
		return "", nil
	}

	var value interface{}

	// Try as *core.PamSettings first
	if pamSettings, ok := field.(*core.PamSettings); ok {
		if len(pamSettings.Value) == 0 {
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

	// Serialize to compact JSON (no whitespace, consistent ordering)
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		return "", fmt.Errorf("failed to serialize pam_settings to JSON: %w", err)
	}

	return string(jsonBytes), nil
}
