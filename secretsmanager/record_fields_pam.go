package secretsmanager

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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