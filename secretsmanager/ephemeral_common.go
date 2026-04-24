package secretsmanager

import (
	"context"
	"encoding/base64"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/keeper-security/secrets-manager-go/core"
)

// fileRefEphemeralAttribute returns the file_ref as a computed list nested attribute.
func fileRefEphemeralAttribute() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Computed:    true,
		Sensitive:   true,
		Description: "The secret file references.",
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"uid": schema.StringAttribute{
					Computed:    true,
					Description: "The file ref UID.",
				},
				"title": schema.StringAttribute{
					Computed:    true,
					Description: "The file title.",
				},
				"name": schema.StringAttribute{
					Computed:    true,
					Description: "The file name.",
				},
				"type": schema.StringAttribute{
					Computed:    true,
					Description: "The file type.",
				},
				"size": schema.Int64Attribute{
					Computed:    true,
					Description: "The file size.",
				},
				"last_modified": schema.StringAttribute{
					Computed:    true,
					Description: "The file last modified date.",
				},
				"content_base64": schema.StringAttribute{
					Computed:    true,
					Description: "The file content (base64).",
				},
			},
		},
	}
}

// hostEphemeralAttribute returns the host as a computed list nested attribute.
func hostEphemeralAttribute() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Computed:    true,
		Description: "Hostname and port.",
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"host_name": schema.StringAttribute{
					Computed:    true,
					Description: "The hostname.",
				},
				"port": schema.StringAttribute{
					Computed:    true,
					Description: "The port.",
				},
			},
		},
	}
}

var fileRefObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"uid":            types.StringType,
		"title":          types.StringType,
		"name":           types.StringType,
		"type":           types.StringType,
		"size":           types.Int64Type,
		"last_modified":  types.StringType,
		"content_base64": types.StringType,
	},
}

var hostObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"host_name": types.StringType,
		"port":      types.StringType,
	},
}

var totpObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"url":   types.StringType,
		"token": types.StringType,
		"ttl":   types.Int64Type,
	},
}

var keyPairObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"public_key":  types.StringType,
		"private_key": types.StringType,
	},
}

// fileItemsToListValue converts KSM file items to a Framework types.List.
func fileItemsToListValue(ctx context.Context, files []*core.KeeperFile) (types.List, diag.Diagnostics) {
	if len(files) == 0 {
		return types.ListValueMust(fileRefObjectType, []attr.Value{}), nil
	}

	items := make([]attr.Value, 0, len(files))
	for _, f := range files {
		content := ""
		if f.FileData != nil {
			content = base64.StdEncoding.EncodeToString(f.FileData)
		}
		lastModified := ""
		if f.LastModified > 0 {
			lastModified = time.Unix(int64(f.LastModified/1000), 0).Format(time.RFC3339)
		}

		obj, diags := types.ObjectValue(fileRefObjectType.AttrTypes, map[string]attr.Value{
			"uid":            types.StringValue(f.Uid),
			"title":          types.StringValue(f.Title),
			"name":           types.StringValue(f.Name),
			"type":           types.StringValue(f.Type),
			"size":           types.Int64Value(int64(f.Size)),
			"last_modified":  types.StringValue(lastModified),
			"content_base64": types.StringValue(content),
		})
		if diags.HasError() {
			return types.ListNull(fileRefObjectType), diags
		}
		items = append(items, obj)
	}

	return types.ListValue(fileRefObjectType, items)
}

// hostToListValue converts KSM host field data to a Framework types.List.
func hostToListValue(ctx context.Context, secret *core.Record) (types.List, diag.Diagnostics) {
	fields := secret.GetFieldsByType("host")
	if len(fields) == 0 {
		return types.ListValueMust(hostObjectType, []attr.Value{}), nil
	}

	hostName := ""
	port := ""
	if values, ok := fields[0]["value"].([]interface{}); ok && len(values) > 0 {
		if vmap, ok := values[0].(map[string]interface{}); ok {
			if val, ok := vmap["hostName"].(string); ok {
				hostName = val
			}
			if val, ok := vmap["port"].(string); ok {
				port = val
			}
		}
	}

	obj, diags := types.ObjectValue(hostObjectType.AttrTypes, map[string]attr.Value{
		"host_name": types.StringValue(hostName),
		"port":      types.StringValue(port),
	})
	if diags.HasError() {
		return types.ListNull(hostObjectType), diags
	}

	return types.ListValue(hostObjectType, []attr.Value{obj})
}

// keyPairToListValue converts KSM key pair field data to a Framework types.List.
func keyPairToListValue(ctx context.Context, secret *core.Record) (types.List, diag.Diagnostics) {
	fields := secret.GetFieldsByType("keyPair")
	if len(fields) == 0 {
		return types.ListValueMust(keyPairObjectType, []attr.Value{}), nil
	}

	publicKey := ""
	privateKey := ""
	if values, ok := fields[0]["value"].([]interface{}); ok && len(values) > 0 {
		if vmap, ok := values[0].(map[string]interface{}); ok {
			if val, ok := vmap["publicKey"].(string); ok {
				publicKey = val
			}
			if val, ok := vmap["privateKey"].(string); ok {
				privateKey = val
			}
		}
	}

	obj, diags := types.ObjectValue(keyPairObjectType.AttrTypes, map[string]attr.Value{
		"public_key":  types.StringValue(publicKey),
		"private_key": types.StringValue(privateKey),
	})
	if diags.HasError() {
		return types.ListNull(keyPairObjectType), diags
	}

	return types.ListValue(keyPairObjectType, []attr.Value{obj})
}

// totpToListValue converts a TOTP URL to a Framework types.List with generated token.
func totpToListValue(ctx context.Context, totpUrl string) (types.List, diag.Diagnostics) {
	if totpUrl == "" {
		return types.ListValueMust(totpObjectType, []attr.Value{}), nil
	}

	code, seconds, err := getTotpCode(totpUrl)
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError("Error generating TOTP code", err.Error())
		return types.ListNull(totpObjectType), diags
	}

	obj, diags := types.ObjectValue(totpObjectType.AttrTypes, map[string]attr.Value{
		"url":   types.StringValue(totpUrl),
		"token": types.StringValue(code),
		"ttl":   types.Int64Value(int64(seconds)),
	})
	if diags.HasError() {
		return types.ListNull(totpObjectType), diags
	}

	return types.ListValue(totpObjectType, []attr.Value{obj})
}

// dateFieldToString converts a KSM date field value (unix millis) to YYYY-MM-DD string.
func dateFieldToString(dateValue string) string {
	dateValue = strings.TrimSpace(dateValue)
	if dateValue == "" {
		return ""
	}
	if unixTime, err := strconv.ParseInt(dateValue, 10, 64); err == nil {
		return time.Unix(unixTime/1000, 0).UTC().Format("2006-01-02")
	}
	return dateValue
}

// nameEphemeralAttribute returns the name as a computed list nested attribute.
func nameEphemeralAttribute() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Computed:    true,
		Description: "Full name.",
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"first": schema.StringAttribute{
					Computed:    true,
					Description: "First name.",
				},
				"middle": schema.StringAttribute{
					Computed:    true,
					Description: "Middle name.",
				},
				"last": schema.StringAttribute{
					Computed:    true,
					Description: "Last name.",
				},
			},
		},
	}
}

// addressEphemeralAttribute returns the address as a computed list nested attribute.
func addressEphemeralAttribute() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Computed:    true,
		Description: "The address information.",
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"street1": schema.StringAttribute{
					Computed:    true,
					Description: "Street line one.",
				},
				"street2": schema.StringAttribute{
					Computed:    true,
					Description: "Street line two.",
				},
				"city": schema.StringAttribute{
					Computed:    true,
					Description: "City.",
				},
				"state": schema.StringAttribute{
					Computed:    true,
					Description: "State.",
				},
				"zip": schema.StringAttribute{
					Computed:    true,
					Description: "ZIP code.",
				},
				"country": schema.StringAttribute{
					Computed:    true,
					Description: "Country.",
				},
			},
		},
	}
}

// phoneEphemeralAttribute returns the phone as a computed list nested attribute.
func phoneEphemeralAttribute() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Computed:    true,
		Description: "Phone number.",
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"region": schema.StringAttribute{
					Computed:    true,
					Description: "Region.",
				},
				"number": schema.StringAttribute{
					Computed:    true,
					Description: "Phone number.",
				},
				"ext": schema.StringAttribute{
					Computed:    true,
					Description: "Extension.",
				},
				"type": schema.StringAttribute{
					Computed:    true,
					Description: "Type - Mobile, Home or Work.",
				},
			},
		},
	}
}

// paymentCardEphemeralAttribute returns the payment card as a computed list nested attribute.
func paymentCardEphemeralAttribute() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Computed:    true,
		Description: "The payment card information.",
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"card_number": schema.StringAttribute{
					Computed:    true,
					Sensitive:   true,
					Description: "The card number.",
				},
				"card_expiration_date": schema.StringAttribute{
					Computed:    true,
					Sensitive:   true,
					Description: "The card expiration date.",
				},
				"card_security_code": schema.StringAttribute{
					Computed:    true,
					Sensitive:   true,
					Description: "The card security code.",
				},
			},
		},
	}
}

// bankAccountEphemeralAttribute returns the bank account as a computed list nested attribute.
func bankAccountEphemeralAttribute() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Computed:    true,
		Description: "The bank account information.",
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"account_type": schema.StringAttribute{
					Computed:    true,
					Description: "The account type.",
				},
				"other_type": schema.StringAttribute{
					Computed:    true,
					Description: "The other type.",
				},
				"routing_number": schema.StringAttribute{
					Computed:    true,
					Sensitive:   true,
					Description: "The routing number.",
				},
				"account_number": schema.StringAttribute{
					Computed:    true,
					Sensitive:   true,
					Description: "The account number.",
				},
			},
		},
	}
}

// addressRefEphemeralAttribute returns the address_ref as a computed list nested attribute.
func addressRefEphemeralAttribute() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Computed:    true,
		Description: "The referenced address record.",
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"uid": schema.StringAttribute{
					Computed:    true,
					Description: "The address ref UID.",
				},
				"street1": schema.StringAttribute{
					Computed:    true,
					Description: "Street line one.",
				},
				"street2": schema.StringAttribute{
					Computed:    true,
					Description: "Street line two.",
				},
				"city": schema.StringAttribute{
					Computed:    true,
					Description: "City.",
				},
				"state": schema.StringAttribute{
					Computed:    true,
					Description: "State.",
				},
				"zip": schema.StringAttribute{
					Computed:    true,
					Description: "ZIP code.",
				},
				"country": schema.StringAttribute{
					Computed:    true,
					Description: "Country.",
				},
			},
		},
	}
}

// cardRefEphemeralAttribute returns the card_ref as a computed list nested attribute.
func cardRefEphemeralAttribute() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Computed:    true,
		Description: "The referenced card record.",
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"uid": schema.StringAttribute{
					Computed:    true,
					Description: "The card ref UID.",
				},
				"payment_card": paymentCardEphemeralAttribute(),
				"cardholder_name": schema.StringAttribute{
					Computed:    true,
					Description: "The cardholder name.",
				},
				"pin_code": schema.StringAttribute{
					Computed:    true,
					Sensitive:   true,
					Description: "The PIN code.",
				},
			},
		},
	}
}

var nameObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"first":  types.StringType,
		"middle": types.StringType,
		"last":   types.StringType,
	},
}

var addressObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"street1": types.StringType,
		"street2": types.StringType,
		"city":    types.StringType,
		"state":   types.StringType,
		"zip":     types.StringType,
		"country": types.StringType,
	},
}

var phoneObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"region": types.StringType,
		"number": types.StringType,
		"ext":    types.StringType,
		"type":   types.StringType,
	},
}

var paymentCardObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"card_number":          types.StringType,
		"card_expiration_date": types.StringType,
		"card_security_code":   types.StringType,
	},
}

var bankAccountObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"account_type":   types.StringType,
		"other_type":     types.StringType,
		"routing_number": types.StringType,
		"account_number": types.StringType,
	},
}

var addressRefObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"uid":     types.StringType,
		"street1": types.StringType,
		"street2": types.StringType,
		"city":    types.StringType,
		"state":   types.StringType,
		"zip":     types.StringType,
		"country": types.StringType,
	},
}

var cardRefObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"uid":             types.StringType,
		"payment_card":    types.ListType{ElemType: paymentCardObjectType},
		"cardholder_name": types.StringType,
		"pin_code":        types.StringType,
	},
}

// nameToListValue converts KSM name field data to a Framework types.List.
func nameToListValue(ctx context.Context, secret *core.Record) (types.List, diag.Diagnostics) {
	fields := secret.GetFieldsByType("name")
	if len(fields) == 0 {
		return types.ListValueMust(nameObjectType, []attr.Value{}), nil
	}

	first, middle, last := "", "", ""
	if values, ok := fields[0]["value"].([]interface{}); ok && len(values) > 0 {
		if vmap, ok := values[0].(map[string]interface{}); ok {
			if val, ok := vmap["first"].(string); ok {
				first = val
			}
			if val, ok := vmap["middle"].(string); ok {
				middle = val
			}
			if val, ok := vmap["last"].(string); ok {
				last = val
			}
		}
	}

	obj, diags := types.ObjectValue(nameObjectType.AttrTypes, map[string]attr.Value{
		"first":  types.StringValue(first),
		"middle": types.StringValue(middle),
		"last":   types.StringValue(last),
	})
	if diags.HasError() {
		return types.ListNull(nameObjectType), diags
	}

	return types.ListValue(nameObjectType, []attr.Value{obj})
}

// addressToListValue converts KSM address field data to a Framework types.List.
func addressToListValue(ctx context.Context, secret *core.Record) (types.List, diag.Diagnostics) {
	fields := secret.GetFieldsByType("address")
	if len(fields) == 0 {
		return types.ListValueMust(addressObjectType, []attr.Value{}), nil
	}

	street1, street2, city, state, zip, country := "", "", "", "", "", ""
	if values, ok := fields[0]["value"].([]interface{}); ok && len(values) > 0 {
		if vmap, ok := values[0].(map[string]interface{}); ok {
			if val, ok := vmap["street1"].(string); ok {
				street1 = val
			}
			if val, ok := vmap["street2"].(string); ok {
				street2 = val
			}
			if val, ok := vmap["city"].(string); ok {
				city = val
			}
			if val, ok := vmap["state"].(string); ok {
				state = val
			}
			if val, ok := vmap["zip"].(string); ok {
				zip = val
			}
			if val, ok := vmap["country"].(string); ok {
				country = val
			}
		}
	}

	obj, diags := types.ObjectValue(addressObjectType.AttrTypes, map[string]attr.Value{
		"street1": types.StringValue(street1),
		"street2": types.StringValue(street2),
		"city":    types.StringValue(city),
		"state":   types.StringValue(state),
		"zip":     types.StringValue(zip),
		"country": types.StringValue(country),
	})
	if diags.HasError() {
		return types.ListNull(addressObjectType), diags
	}

	return types.ListValue(addressObjectType, []attr.Value{obj})
}

// phoneToListValue converts KSM phone field data to a Framework types.List.
func phoneToListValue(ctx context.Context, secret *core.Record) (types.List, diag.Diagnostics) {
	fields := secret.GetFieldsByType("phone")
	if len(fields) == 0 {
		return types.ListValueMust(phoneObjectType, []attr.Value{}), nil
	}

	region, number, ext, phoneType := "", "", "", ""
	if values, ok := fields[0]["value"].([]interface{}); ok && len(values) > 0 {
		if vmap, ok := values[0].(map[string]interface{}); ok {
			if val, ok := vmap["region"].(string); ok {
				region = val
			}
			if val, ok := vmap["number"].(string); ok {
				number = val
			}
			if val, ok := vmap["ext"].(string); ok {
				ext = val
			}
			if val, ok := vmap["type"].(string); ok {
				phoneType = val
			}
		}
	}

	obj, diags := types.ObjectValue(phoneObjectType.AttrTypes, map[string]attr.Value{
		"region": types.StringValue(region),
		"number": types.StringValue(number),
		"ext":    types.StringValue(ext),
		"type":   types.StringValue(phoneType),
	})
	if diags.HasError() {
		return types.ListNull(phoneObjectType), diags
	}

	return types.ListValue(phoneObjectType, []attr.Value{obj})
}

// paymentCardToListValue converts KSM paymentCard field data to a Framework types.List.
func paymentCardToListValue(ctx context.Context, secret *core.Record) (types.List, diag.Diagnostics) {
	fields := secret.GetFieldsByType("paymentCard")
	if len(fields) == 0 {
		return types.ListValueMust(paymentCardObjectType, []attr.Value{}), nil
	}

	cardNumber, cardExpDate, cardSecCode := "", "", ""
	if values, ok := fields[0]["value"].([]interface{}); ok && len(values) > 0 {
		if vmap, ok := values[0].(map[string]interface{}); ok {
			if val, ok := vmap["cardNumber"].(string); ok {
				cardNumber = val
			}
			if val, ok := vmap["cardExpirationDate"].(string); ok {
				cardExpDate = val
			}
			if val, ok := vmap["cardSecurityCode"].(string); ok {
				cardSecCode = val
			}
		}
	}

	obj, diags := types.ObjectValue(paymentCardObjectType.AttrTypes, map[string]attr.Value{
		"card_number":          types.StringValue(cardNumber),
		"card_expiration_date": types.StringValue(cardExpDate),
		"card_security_code":   types.StringValue(cardSecCode),
	})
	if diags.HasError() {
		return types.ListNull(paymentCardObjectType), diags
	}

	return types.ListValue(paymentCardObjectType, []attr.Value{obj})
}

// bankAccountToListValue converts KSM bankAccount field data to a Framework types.List.
func bankAccountToListValue(ctx context.Context, secret *core.Record) (types.List, diag.Diagnostics) {
	fields := secret.GetFieldsByType("bankAccount")
	if len(fields) == 0 {
		return types.ListValueMust(bankAccountObjectType, []attr.Value{}), nil
	}

	accountType, otherType, routingNumber, accountNumber := "", "", "", ""
	if values, ok := fields[0]["value"].([]interface{}); ok && len(values) > 0 {
		if vmap, ok := values[0].(map[string]interface{}); ok {
			if val, ok := vmap["accountType"].(string); ok {
				accountType = val
			}
			if val, ok := vmap["otherType"].(string); ok {
				otherType = val
			}
			if val, ok := vmap["routingNumber"].(string); ok {
				routingNumber = val
			}
			if val, ok := vmap["accountNumber"].(string); ok {
				accountNumber = val
			}
		}
	}

	obj, diags := types.ObjectValue(bankAccountObjectType.AttrTypes, map[string]attr.Value{
		"account_type":   types.StringValue(accountType),
		"other_type":     types.StringValue(otherType),
		"routing_number": types.StringValue(routingNumber),
		"account_number": types.StringValue(accountNumber),
	})
	if diags.HasError() {
		return types.ListNull(bankAccountObjectType), diags
	}

	return types.ListValue(bankAccountObjectType, []attr.Value{obj})
}

// addressRefToListValue fetches the referenced address record and converts it to a Framework types.List.
// If the UID is empty, returns an empty list. If the referenced record cannot be fetched, returns a
// partial object (uid only, fields empty) with a warning diagnostic so callers are aware.
func addressRefToListValue(ctx context.Context, secret *core.Record, client core.SecretsManager) (types.List, diag.Diagnostics) {
	uid := strings.TrimSpace(secret.GetFieldValueByType("addressRef"))
	if uid == "" {
		return types.ListValueMust(addressRefObjectType, []attr.Value{}), nil
	}

	attrs := map[string]attr.Value{
		"uid":     types.StringValue(uid),
		"street1": types.StringValue(""),
		"street2": types.StringValue(""),
		"city":    types.StringValue(""),
		"state":   types.StringValue(""),
		"zip":     types.StringValue(""),
		"country": types.StringValue(""),
	}

	var diags diag.Diagnostics
	refs, err := getSecrets(client, []string{uid})
	if err != nil || len(refs) == 0 {
		diags.AddWarning("Referenced Address Record Not Found",
			"Could not fetch addressRef record with UID '"+uid+"'. Address fields will be empty.")
	} else {
		if fields := refs[0].GetFieldsByType("address"); len(fields) > 0 {
			if values, ok := fields[0]["value"].([]interface{}); ok && len(values) > 0 {
				if vmap, ok := values[0].(map[string]interface{}); ok {
					if val, ok := vmap["street1"].(string); ok {
						attrs["street1"] = types.StringValue(val)
					}
					if val, ok := vmap["street2"].(string); ok {
						attrs["street2"] = types.StringValue(val)
					}
					if val, ok := vmap["city"].(string); ok {
						attrs["city"] = types.StringValue(val)
					}
					if val, ok := vmap["state"].(string); ok {
						attrs["state"] = types.StringValue(val)
					}
					if val, ok := vmap["zip"].(string); ok {
						attrs["zip"] = types.StringValue(val)
					}
					if val, ok := vmap["country"].(string); ok {
						attrs["country"] = types.StringValue(val)
					}
				}
			}
		}
	}

	obj, objDiags := types.ObjectValue(addressRefObjectType.AttrTypes, attrs)
	diags.Append(objDiags...)
	if diags.HasError() {
		return types.ListNull(addressRefObjectType), diags
	}

	list, listDiags := types.ListValue(addressRefObjectType, []attr.Value{obj})
	diags.Append(listDiags...)
	return list, diags
}

// cardRefToListValue fetches the referenced card record and converts it to a Framework types.List.
// If the UID is empty, returns an empty list. If the referenced record cannot be fetched, returns a
// partial object (uid only, fields empty) with a warning diagnostic so callers are aware.
func cardRefToListValue(ctx context.Context, secret *core.Record, client core.SecretsManager) (types.List, diag.Diagnostics) {
	uid := strings.TrimSpace(secret.GetFieldValueByType("cardRef"))
	if uid == "" {
		return types.ListValueMust(cardRefObjectType, []attr.Value{}), nil
	}

	// Build a default payment_card empty list.
	emptyPaymentCard, pcDiags := types.ListValue(paymentCardObjectType, []attr.Value{})
	if pcDiags.HasError() {
		return types.ListNull(cardRefObjectType), pcDiags
	}

	cardAttrs := map[string]attr.Value{
		"uid":             types.StringValue(uid),
		"payment_card":    emptyPaymentCard,
		"cardholder_name": types.StringValue(""),
		"pin_code":        types.StringValue(""),
	}

	var diags diag.Diagnostics
	refs, err := getSecrets(client, []string{uid})
	if err != nil || len(refs) == 0 {
		diags.AddWarning("Referenced Card Record Not Found",
			"Could not fetch cardRef record with UID '"+uid+"'. Card fields will be empty.")
	} else {
		ref := refs[0]

		cardNumber, cardExpDate, cardSecCode := "", "", ""
		if fields := ref.GetFieldsByType("paymentCard"); len(fields) > 0 {
			if values, ok := fields[0]["value"].([]interface{}); ok && len(values) > 0 {
				if vmap, ok := values[0].(map[string]interface{}); ok {
					if val, ok := vmap["cardNumber"].(string); ok {
						cardNumber = val
					}
					if val, ok := vmap["cardExpirationDate"].(string); ok {
						cardExpDate = val
					}
					if val, ok := vmap["cardSecurityCode"].(string); ok {
						cardSecCode = val
					}
				}
			}
		}

		pcObj, pcDiags := types.ObjectValue(paymentCardObjectType.AttrTypes, map[string]attr.Value{
			"card_number":          types.StringValue(cardNumber),
			"card_expiration_date": types.StringValue(cardExpDate),
			"card_security_code":   types.StringValue(cardSecCode),
		})
		if pcDiags.HasError() {
			return types.ListNull(cardRefObjectType), pcDiags
		}

		pcList, pcDiags := types.ListValue(paymentCardObjectType, []attr.Value{pcObj})
		if pcDiags.HasError() {
			return types.ListNull(cardRefObjectType), pcDiags
		}

		cardAttrs["payment_card"] = pcList
		cardAttrs["cardholder_name"] = types.StringValue(ref.GetFieldValueByType("text"))
		cardAttrs["pin_code"] = types.StringValue(ref.GetFieldValueByType("pinCode"))
	}

	obj, objDiags := types.ObjectValue(cardRefObjectType.AttrTypes, cardAttrs)
	diags.Append(objDiags...)
	if diags.HasError() {
		return types.ListNull(cardRefObjectType), diags
	}

	list, listDiags := types.ListValue(cardRefObjectType, []attr.Value{obj})
	diags.Append(listDiags...)
	return list, diags
}

// pamHostnameToListValue converts the pamHostname field to a Framework types.List.
// pamHostname has the same hostname/port structure as the host field.
func pamHostnameToListValue(ctx context.Context, secret *core.Record) (types.List, diag.Diagnostics) {
	fields := secret.GetFieldsByType("pamHostname")
	if len(fields) == 0 {
		return types.ListValueMust(hostObjectType, []attr.Value{}), nil
	}

	hostName := ""
	port := ""
	if values, ok := fields[0]["value"].([]interface{}); ok && len(values) > 0 {
		if vmap, ok := values[0].(map[string]interface{}); ok {
			if val, ok := vmap["hostName"].(string); ok {
				hostName = val
			}
			if val, ok := vmap["port"].(string); ok {
				port = val
			}
		}
	}

	obj, diags := types.ObjectValue(hostObjectType.AttrTypes, map[string]attr.Value{
		"host_name": types.StringValue(hostName),
		"port":      types.StringValue(port),
	})
	if diags.HasError() {
		return types.ListNull(hostObjectType), diags
	}

	return types.ListValue(hostObjectType, []attr.Value{obj})
}

// pamSettingsToString converts a pamSettings field to a JSON string.
func pamSettingsToString(secret *core.Record) string {
	fields := secret.GetFieldsByType("pamSettings")
	if len(fields) == 0 {
		return ""
	}
	jsonStr, err := pamSettingsFieldToJSON(fields[0])
	if err != nil {
		return ""
	}
	return jsonStr
}

// pamFieldString extracts a simple string value from a PAM field by type and section.
func pamFieldString(fieldType, section string, secret *core.Record) string {
	flds := getFieldDicts(fieldType, section, secret.RecordDict)
	if len(flds) == 0 {
		return ""
	}
	if fmap, ok := flds[0].(map[string]interface{}); ok {
		if values, ok := fmap["value"].([]interface{}); ok && len(values) > 0 {
			if str, ok := values[0].(string); ok {
				return str
			}
		}
	}
	return ""
}

// pamFieldStringWithLabel extracts a simple string value from a PAM field by type, section, and label.
func pamFieldStringWithLabel(fieldType, section string, secret *core.Record, label string) string {
	flds := getFieldDicts(fieldType, section, secret.RecordDict)
	for _, fld := range flds {
		if fmap, ok := fld.(map[string]interface{}); ok {
			if lblValue, found := fmap["label"]; found && lblValue == label {
				if values, ok := fmap["value"].([]interface{}); ok && len(values) > 0 {
					if str, ok := values[0].(string); ok {
						return str
					}
				}
			}
		}
	}
	return ""
}

// pamFieldBoolWithLabel extracts a bool value from a PAM field by type, section, and label.
func pamFieldBoolWithLabel(fieldType, section string, secret *core.Record, label string) bool {
	flds := getFieldDicts(fieldType, section, secret.RecordDict)
	for _, fld := range flds {
		if fmap, ok := fld.(map[string]interface{}); ok {
			if lblValue, found := fmap["label"]; found && lblValue == label {
				if values, ok := fmap["value"].([]interface{}); ok && len(values) > 0 {
					if b, ok := values[0].(bool); ok {
						return b
					}
				}
			}
		}
	}
	return false
}

// pamFieldStringByType extracts a simple string value from a field (like databaseType/directoryType).
func pamFieldStringByType(fieldType string, secret *core.Record) string {
	fields := secret.GetFieldsByType(fieldType)
	if len(fields) == 0 {
		return ""
	}
	if values, ok := fields[0]["value"].([]interface{}); ok && len(values) > 0 {
		if str, ok := values[0].(string); ok {
			return str
		}
	}
	return ""
}
