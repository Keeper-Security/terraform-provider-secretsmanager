package secretsmanager

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/keeper-security/secrets-manager-go/core"
)

func dataSourceRecords() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRecordsRead,
		Schema: map[string]*schema.Schema{
			"uids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of record UIDs to fetch",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"titles": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of record titles to fetch (requires fetching all records first)",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"records": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of fetched records",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The record UID",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The secret type",
						},
						"title": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The secret title",
						},
						"notes": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The secret notes",
						},
						"fields": schemaGenericField(),
						"custom": schemaGenericField(),
						"file_ref": {
							Type:        schema.TypeList,
							Computed:    true,
							Sensitive:   true,
							Description: "The secret file references",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"uid": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The file ref UID",
									},
									"title": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The file title",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The file name",
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The file type",
									},
									"size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The file size",
									},
									"last_modified": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The file last modified date",
									},
									"content_base64": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The file content (base64)",
									},
								},
							},
						},
					},
				},
			},
			"records_by_uid": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Map of records keyed by UID for direct access (JSON-encoded)",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceRecordsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(providerMeta)
	client := *provider.client
	var diags diag.Diagnostics

	// Get UIDs and titles from config
	uidsRaw := d.Get("uids").([]interface{})
	titlesRaw := d.Get("titles").([]interface{})

	// Validate that at least one is provided
	if len(uidsRaw) == 0 && len(titlesRaw) == 0 {
		return diag.Errorf("at least one of 'uids' or 'titles' must be provided")
	}

	// Convert to string slices
	uids := make([]string, len(uidsRaw))
	for i, uid := range uidsRaw {
		uids[i] = strings.TrimSpace(uid.(string))
	}

	titles := make([]string, len(titlesRaw))
	for i, title := range titlesRaw {
		titles[i] = strings.TrimSpace(title.(string))
	}

	// Future enhancement: enforce batch size limit
	// const maxBatchSize = 500
	// totalRecords := len(uids) + len(titles)
	// if totalRecords > maxBatchSize {
	//     return diag.Errorf("batch size exceeds maximum of %d records (requested %d)", maxBatchSize, totalRecords)
	// }

	var secrets []*core.Record
	var err error

	// Optimization: If we have titles, we need to fetch all records anyway
	// So we can filter both UIDs and titles from the same result set
	if len(titles) > 0 {
		// Fetch all records once
		allSecrets, err := client.GetSecrets([]string{})
		if err != nil {
			return diag.Errorf("failed to fetch all records: %v", err)
		}

		// Create maps for efficient lookup
		uidMap := make(map[string]bool)
		for _, uid := range uids {
			uidMap[uid] = true
		}

		titleMap := make(map[string]bool)
		for _, title := range titles {
			titleMap[title] = true
		}

		// Filter records by UIDs and titles
		for _, record := range allSecrets {
			if uidMap[record.Uid] || titleMap[record.Title()] {
				secrets = append(secrets, record)
			}
		}

		// Validate that we found all requested records
		foundUids := make(map[string]bool)
		foundTitles := make(map[string]bool)
		for _, record := range secrets {
			foundUids[record.Uid] = true
			foundTitles[record.Title()] = true
		}

		// Check for missing UIDs
		for _, uid := range uids {
			if !foundUids[uid] {
				return diag.Errorf("record not found - UID: %s", uid)
			}
		}

		// Check for missing titles
		for _, title := range titles {
			if !foundTitles[title] {
				return diag.Errorf("record not found - title: %s", title)
			}
		}
	} else {
		// Only UIDs provided - efficient batch fetch
		secrets, err = client.GetSecrets(uids)
		if err != nil {
			return diag.Errorf("failed to fetch records: %v", err)
		}

		// Validate that we got all requested records
		if len(secrets) != len(uids) {
			// Find missing UIDs
			foundUids := make(map[string]bool)
			for _, record := range secrets {
				foundUids[record.Uid] = true
			}
			for _, uid := range uids {
				if !foundUids[uid] {
					return diag.Errorf("record not found - UID: %s", uid)
				}
			}
		}
	}

	// Convert records to Terraform schema format
	recordsList := make([]interface{}, len(secrets))
	recordsMap := make(map[string]interface{})
	
	for i, secret := range secrets {
		record := make(map[string]interface{})
		
		record["uid"] = secret.Uid
		record["type"] = secret.Type()
		record["title"] = secret.Title()
		record["notes"] = secret.Notes()

		// Process fields
		fieldItems := getFieldItemsData(secret.RecordDict, "fields")
		record["fields"] = fieldItems

		// Process custom fields
		customItems := getFieldItemsData(secret.RecordDict, "custom")
		record["custom"] = customItems

		// Process file references
		fileItems := getFileItemsData(secret.Files)
		record["file_ref"] = fileItems

		recordsList[i] = record
		
		// Store JSON-encoded record in map for UID-based access
		// This allows users to decode and access the full record structure
		if jsonData, err := json.Marshal(record); err == nil {
			recordsMap[secret.Uid] = string(jsonData)
		}
	}

	// Set the records list in the data source
	if err := d.Set("records", recordsList); err != nil {
		return diag.FromErr(err)
	}

	// Set the records map for UID-based access
	if err := d.Set("records_by_uid", recordsMap); err != nil {
		return diag.FromErr(err)
	}

	// Generate a consistent ID for this data source
	// Use a hash of sorted UIDs for consistency
	allIds := make([]string, 0, len(secrets))
	for _, record := range secrets {
		allIds = append(allIds, record.Uid)
	}
	sort.Strings(allIds)
	
	h := sha256.New()
	h.Write([]byte(strings.Join(allIds, ",")))
	d.SetId(fmt.Sprintf("%x", h.Sum(nil)))

	return diags
}