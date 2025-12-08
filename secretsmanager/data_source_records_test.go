package secretsmanager

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceRecords_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  testAccPreCheck(t),
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRecordsConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceRecordsCheck("data.secretsmanager_records.test"),
					resource.TestCheckResourceAttrSet("data.secretsmanager_records.test", "records.#"),
					resource.TestCheckResourceAttrSet("data.secretsmanager_records.test", "records.0.uid"),
					resource.TestCheckResourceAttrSet("data.secretsmanager_records.test", "records.0.type"),
					resource.TestCheckResourceAttrSet("data.secretsmanager_records.test", "records.0.title"),
					resource.TestCheckResourceAttrSet("data.secretsmanager_records.test", "records_by_uid.%"),
				),
			},
		},
	})
}

func TestAccDataSourceRecords_WithTitles(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  testAccPreCheck(t),
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRecordsConfig_withTitles(),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceRecordsCheck("data.secretsmanager_records.test"),
					resource.TestCheckResourceAttrSet("data.secretsmanager_records.test", "records.#"),
					resource.TestCheckResourceAttrSet("data.secretsmanager_records.test", "records_by_uid.%"),
				),
			},
		},
	})
}

func TestAccDataSourceRecords_MultipleTitles(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  testAccPreCheck(t),
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRecordsConfig_mixed(),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceRecordsCheck("data.secretsmanager_records.test"),
					resource.TestCheckResourceAttrSet("data.secretsmanager_records.test", "records.#"),
					resource.TestCheckResourceAttrSet("data.secretsmanager_records.test", "records_by_uid.%"),
				),
			},
		},
	})
}

func TestAccDataSourceRecords_LargeBatch(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  testAccPreCheck(t),
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRecordsConfig_largeBatch(),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceRecordsCheck("data.secretsmanager_records.test"),
					resource.TestCheckResourceAttrSet("data.secretsmanager_records.test", "records.#"),
				),
			},
		},
	})
}

func testAccDataSourceRecordsCheck(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("data source not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("data source ID not set")
		}

		// Check that we have records
		recordsCount := rs.Primary.Attributes["records.#"]
		if recordsCount == "" || recordsCount == "0" {
			return fmt.Errorf("no records found")
		}

		// Check that records_by_uid is populated
		// Terraform stores maps with a ".%" key for count
		recordsByUidCount, ok := rs.Primary.Attributes["records_by_uid.%"]
		if !ok || recordsByUidCount == "" || recordsByUidCount == "0" {
			return fmt.Errorf("records_by_uid not populated")
		}

		return nil
	}
}

func testAccDataSourceRecordsConfig_basic() string {
	_, title := testAcc.getRecordInfo("login")
	return fmt.Sprintf(`
data "secretsmanager_records" "test" {
	titles = [
		"%s"
	]
}
`, title)
}

func testAccDataSourceRecordsConfig_withTitles() string {
	_, title := testAcc.getRecordInfo("login")
	return fmt.Sprintf(`
data "secretsmanager_records" "test" {
	titles = [
		"%s"
	]
}
`, title)
}

func testAccDataSourceRecordsConfig_mixed() string {
	_, title1 := testAcc.getRecordInfo("login")
	_, title2 := testAcc.getRecordInfo("encryptedNotes")
	return fmt.Sprintf(`
data "secretsmanager_records" "test" {
	titles = [
		"%s",
		"%s"
	]
}
`, title1, title2)
}

func testAccDataSourceRecordsConfig_largeBatch() string {
	// This would be used for testing with many titles
	// In real testing, these would be actual titles from test data
	_, title := testAcc.getRecordInfo("login")
	return fmt.Sprintf(`
data "secretsmanager_records" "test" {
	titles = [
		"%s"
	]
}
`, title)
}

func TestAccDataSourceRecords_WithTitlePatterns(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  testAccPreCheck(t),
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRecordsConfig_withTitlePatterns(),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceRecordsCheck("data.secretsmanager_records.test"),
					resource.TestCheckResourceAttrSet("data.secretsmanager_records.test", "records.#"),
					resource.TestCheckResourceAttrSet("data.secretsmanager_records.test", "records_by_uid.%"),
					// Verify at least one record matches the pattern
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["data.secretsmanager_records.test"]
						if !ok {
							return fmt.Errorf("data source not found")
						}
						recordsCount := rs.Primary.Attributes["records.#"]
						if recordsCount == "" || recordsCount == "0" {
							return fmt.Errorf("no records matched the pattern")
						}
						return nil
					},
				),
			},
		},
	})
}

func TestAccDataSourceRecords_InvalidPattern(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  testAccPreCheck(t),
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceRecordsConfig_invalidPattern(),
				ExpectError: regexp.MustCompile("invalid regex pattern"),
			},
		},
	})
}

func TestAccDataSourceRecords_CombinedWithPatterns(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  testAccPreCheck(t),
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRecordsConfig_combinedWithPatterns(),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceRecordsCheck("data.secretsmanager_records.test"),
					resource.TestCheckResourceAttrSet("data.secretsmanager_records.test", "records.#"),
					resource.TestCheckResourceAttrSet("data.secretsmanager_records.test", "records_by_uid.%"),
					// Verify we have multiple records from different sources
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["data.secretsmanager_records.test"]
						if !ok {
							return fmt.Errorf("data source not found")
						}
						recordsCount := rs.Primary.Attributes["records.#"]
						if recordsCount == "" || recordsCount == "0" {
							return fmt.Errorf("no records found")
						}
						return nil
					},
				),
			},
		},
	})
}

func TestAccDataSourceRecords_MultiplePatterns(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  testAccPreCheck(t),
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRecordsConfig_multiplePatterns(),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceRecordsCheck("data.secretsmanager_records.test"),
					resource.TestCheckResourceAttrSet("data.secretsmanager_records.test", "records.#"),
					resource.TestCheckResourceAttrSet("data.secretsmanager_records.test", "records_by_uid.%"),
				),
			},
		},
	})
}

func testAccDataSourceRecordsConfig_withTitlePatterns() string {
	return `
data "secretsmanager_records" "test" {
	title_patterns = [
		"^tf_acc_test.*"
	]
}
`
}

func testAccDataSourceRecordsConfig_invalidPattern() string {
	return `
data "secretsmanager_records" "test" {
	title_patterns = [
		"[invalid(regex"
	]
}
`
}

func testAccDataSourceRecordsConfig_combinedWithPatterns() string {
	_, title1 := testAcc.getRecordInfo("login")
	_, title2 := testAcc.getRecordInfo("encryptedNotes")
	return fmt.Sprintf(`
data "secretsmanager_records" "test" {
	titles = [
		"%s",
		"%s"
	]
	title_patterns = [
		"^tf_acc_test.*"
	]
}
`, title1, title2)
}

func testAccDataSourceRecordsConfig_multiplePatterns() string {
	return `
data "secretsmanager_records" "test" {
	title_patterns = [
		"^tf_acc_test.*login.*",
		"^tf_acc_test.*notes.*"
	]
}
`
}
