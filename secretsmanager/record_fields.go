package secretsmanager

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func schemaAccountNumberField() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "Account number field data.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": { // accountNumber
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Field type.",
				},
				"label": {
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
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Field value.",
				},
			},
		},
	}
}

func schemaAddressRefField() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "AddressRef field data.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": { // addressRef
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Field type.",
				},
				"label": {
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
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Field value.",
				},
			},
		},
	}
}

func schemaAddressField() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "Address field data.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": { // address
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Field type.",
				},
				"label": {
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
					Type:     schema.TypeList,
					Optional: true,
					// MaxItems:    1,
					Description: "Field value.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"street1": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "Street line 1.",
							},
							"street2": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "Street line 2.",
							},
							"city": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "City.",
							},
							"state": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "State.",
							},
							"country": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "Country.",
							},
							"zip": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "ZIP code.",
							},
						},
					},
				},
			},
		},
	}
}

func schemaBankAccountField() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "Bank account field data.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": { // bankAccount
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Field type.",
				},
				"label": {
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
					Type:     schema.TypeList,
					Optional: true,
					// MaxItems:    1,
					Description: "Field value.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"account_type": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "Account type.",
							},
							"routing_number": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "Routing number.",
							},
							"account_number": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "Account number.",
							},
							"other_type": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "Other type info.",
							},
						},
					},
				},
			},
		},
	}
}

func schemaBirthDateField() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "Birth date field data.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": { // birthDate
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Field type.",
				},
				"label": {
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
					Type:        schema.TypeInt,
					Optional:    true,
					Description: "Field value.",
				},
			},
		},
	}
}

func schemaCardRefField() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "CardRef field data.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": { // cardRef
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Field type.",
				},
				"label": {
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
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Field value.",
				},
			},
		},
	}
}

func schemaDateField() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "Date field data.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": { // date
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Field type.",
				},
				"label": {
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
					Type:        schema.TypeInt,
					Optional:    true,
					Description: "Field value.",
				},
			},
		},
	}
}

func schemaEmailField() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "Email field data.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": { // email
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Field type.",
				},
				"label": {
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
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Field value.",
				},
			},
		},
	}
}

func schemaExpirationDateField() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "Expiration date field data.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": { // expirationDate
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Field type.",
				},
				"label": {
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
					Type:        schema.TypeInt,
					Optional:    true,
					Description: "Field value.",
				},
			},
		},
	}
}

func schemaFileRefField() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "FileRef field data.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": { // fileRef
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
					Computed:    true,
					Description: "Required flag.",
				},
				"value": {
					Type:        schema.TypeList,
					Optional:    true,
					Computed:    true,
					Description: "Field value (File UID list).",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							// "path": { // TODO: Enable with file upload = abs. file path
							// 	Type:        schema.TypeString,
							// 	Computed:    true,
							// 	Optional:    true,
							// 	ForceNew:    true,
							// 	Description: "Absolute filepath including the filename.",
							// 	ValidateDiagFunc: func(i interface{}, p cty.Path) diag.Diagnostics {
							// 		var diags diag.Diagnostics
							// 		filePath := i.(string)
							// 		if _, err := os.Stat(filePath); err != nil {
							// 			errMessage := "is not accessible"
							// 			if os.IsNotExist(err) {
							// 				errMessage = "does not exist"
							// 			}
							// 			diag := diag.Diagnostic{
							// 				Severity:      diag.Error,
							// 				Summary:       "wrong value",
							// 				Detail:        fmt.Sprintf("Bad file path: %q %q", filePath, errMessage),
							// 				AttributePath: p,
							// 			}
							// 			diags = append(diags, diag)
							// 		}
							// 		return diags
							// 	},
							// },
							"uid": {
								Type:        schema.TypeString,
								Optional:    true,
								Computed:    true,
								Description: "The file ref UID.",
								// ConflictsWith: []string{"file_ref.value.file.path"},
								ValidateDiagFunc: func(i interface{}, p cty.Path) diag.Diagnostics {
									var diags diag.Diagnostics
									fuid := i.(string)
									if validUid := validateUid(fuid); !validUid {
										diag := diag.Diagnostic{
											Severity:      diag.Error,
											Summary:       "invalid fileRef UID format - expected unpadded base64url encoded value (RFC 4648)",
											Detail:        fmt.Sprintf("Invalid fileRef UID: %q", fuid),
											AttributePath: p,
										}
										diags = append(diags, diag)
									}
									return diags
								},
							},
							"title": {
								Type:        schema.TypeString,
								Computed:    true,
								Description: "The file title.",
							},
							"name": {
								Type:        schema.TypeString,
								Computed:    true,
								Description: "The file name.",
							},
							"type": {
								Type:        schema.TypeString,
								Computed:    true,
								Description: "The file type.",
							},
							"size": {
								Type:        schema.TypeInt,
								Computed:    true,
								Description: "The file size.",
							},
							"last_modified": {
								Type:        schema.TypeString,
								Computed:    true,
								Description: "The file last modified date.",
							},
							"content_base64": {
								Type:        schema.TypeString,
								Computed:    true,
								Description: "The file content (base64).",
							},
						},
					},
				},
			},
		},
	}
}

func schemaHostField() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "Host field data.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": { // host
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Field type.",
				},
				"label": {
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
					Type:     schema.TypeList,
					Optional: true,
					// MaxItems:    1,
					Description: "Field value.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"host_name": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "Hostname.",
							},
							"port": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "Port.",
							},
						},
					},
				},
			},
		},
	}
}

func schemaKeyPairField() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "Key pair field data.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": { // keyPair
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Field type.",
				},
				"label": {
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
					Type:     schema.TypeList,
					Optional: true,
					// MaxItems:    1,
					Description: "Field value.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"public_key": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "Public key.",
							},
							"private_key": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "Private key.",
							},
						},
					},
				},
			},
		},
	}
}

func schemaLicenseNumberField() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "License number field data.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": { // licenseNumber
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Field type.",
				},
				"label": {
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
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Field value.",
				},
			},
		},
	}
}

func schemaLoginField() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		// Computed:    true, // only when used as datasource
		Description: "Login field data.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": { // login
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Field type.",
				},
				"label": {
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
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Field value.",
				},
			},
		},
	}
}

/*
func schemaMultilineField() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "Multiline field data.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": { // multiline
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Field type.",
				},
				"label": {
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
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Field value.",
				},
			},
		},
	}
}
*/

func schemaNameField() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "Name field data.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": { // name
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Field type.",
				},
				"label": {
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
					MaxItems:    1,
					Description: "Field value.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"first": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "First name.",
							},
							"middle": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "MIddle name.",
							},
							"last": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "Last name.",
							},
						},
					},
				},
			},
		},
	}
}

func schemaOneTimeCodeField() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "TOTP field data.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": { // oneTimeCode
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Field type.",
				},
				"label": {
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
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Field value.",
				},
			},
		},
	}
}

func schemaPasswordField(attributeName string) *schema.Schema {
	attributeName = strings.TrimSpace(attributeName)
	if attributeName == "" {
		attributeName = "password"
	}
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "Password field data.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": { // password
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
					Type:     schema.TypeString,
					Optional: true,
					// ConflictsWith: []string{attributeName + ".0.value"},
					Description: "Flag to force password generation (when set to 'yes' or 'true').",
					ValidateDiagFunc: func(i interface{}, p cty.Path) diag.Diagnostics {
						var diags diag.Diagnostics
						valid := []string{"true", "yes"}
						v := i.(string)
						for _, str := range valid {
							if v == str {
								return diags
							}
						}
						diag := diag.Diagnostic{
							Severity:      diag.Error,
							Summary:       fmt.Sprintf("invalid generate = %s", v),
							Detail:        fmt.Sprintf("expected 'generate' to be one of %v, got %s", valid, v),
							AttributePath: p,
						}
						diags = append(diags, diag)
						return diags
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
				"enforce_generation": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "Enforce generation flag.",
				},
				"complexity": {
					Type:        schema.TypeList,
					Optional:    true,
					MaxItems:    1,
					Description: "Password complexity.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"length": {
								Type:        schema.TypeInt,
								Optional:    true,
								Description: "Password length.",
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
					Type:          schema.TypeString,
					Computed:      true,
					Optional:      true,
					Sensitive:     true,
					ConflictsWith: []string{attributeName + ".0.generate"},
					Description:   "Field value.",
				},
			},
		},
	}
}

func schemaPaymentCardField() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "Payment card field data.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": { // paymentCard
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Field type.",
				},
				"label": {
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
					Type:     schema.TypeList,
					Optional: true,
					// MaxItems:    1,
					Description: "Field value.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"card_number": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "Card number.",
							},
							"card_expiration_date": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "Card expiration date.",
							},
							"card_security_code": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "Card security code.",
							},
						},
					},
				},
			},
		},
	}
}

func schemaPhoneField() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "Phone field data.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": { // phone
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Field type.",
				},
				"label": {
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
					Type:     schema.TypeList,
					Optional: true,
					// MaxItems:    1,
					Description: "Field value.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"region": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "Region code - ex. US",
							},
							"number": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "Phone number - ex. 510-222-5555",
							},
							"ext": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "Extension number - ex. 9987",
							},
							"type": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "Phone number type - ex. Home, Work, or Mobile",
								// ValidateFunc: validation.StringInSlice([]string{}, false), // deprecated
								ValidateDiagFunc: func(i interface{}, p cty.Path) diag.Diagnostics {
									var diags diag.Diagnostics
									values := map[string]struct{}{
										"Home":   {},
										"Mobile": {},
										"Work":   {},
									}

									value := i.(string)
									if _, found := values[value]; !found {
										diag := diag.Diagnostic{
											Severity:      diag.Error,
											Summary:       "wrong value",
											Detail:        fmt.Sprintf("%q is not in %q", value, values),
											AttributePath: p,
										}
										diags = append(diags, diag)
									}
									return diags
								},
							},
						},
					},
				},
			},
		},
	}
}

func schemaPinCodeField() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "PinCode field data.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": { // pinCode
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Field type.",
				},
				"label": {
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
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Field value.",
				},
			},
		},
	}
}

/*
func schemaSecretField() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "Secret field data.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": { // secret
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Field type.",
				},
				"label": {
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
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Field value.",
				},
			},
		},
	}
}
*/

func schemaSecureNoteField() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "Secure note field data.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": { // note
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Field type.",
				},
				"label": {
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
					Type:        schema.TypeString,
					Optional:    true,
					Sensitive:   true,
					Description: "Field value.",
				},
			},
		},
	}
}

/*
func schemaSecurityQuestionField() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "Security question field data.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": { // securityQuestion
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Field type.",
				},
				"label": {
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
					Type:     schema.TypeList,
					Optional: true,
					// MaxItems:    1,
					Description: "Field value.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"question": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "Security question.",
							},
							"answer": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "Answer to the security question.",
							},
						},
					},
				},
			},
		},
	}
}
*/

func schemaTextField() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "Text field data.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": { // text
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Field type.",
				},
				"label": {
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
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Field value.",
				},
			},
		},
	}
}

func schemaUrlField() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "URL field data.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": { // url
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Field type.",
				},
				"label": {
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
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Field value.",
				},
			},
		},
	}
}
