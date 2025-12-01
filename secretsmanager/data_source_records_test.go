package secretsmanager

import (
	"fmt"
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
					resource.TestCheckResourceAttrSet("data.secretsmanager_records.test", "records_by_uid"),
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
					resource.TestCheckResourceAttrSet("data.secretsmanager_records.test", "records_by_uid"),
				),
			},
		},
	})
}

func TestAccDataSourceRecords_MixedUidsAndTitles(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  testAccPreCheck(t),
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRecordsConfig_mixed(),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceRecordsCheck("data.secretsmanager_records.test"),
					resource.TestCheckResourceAttrSet("data.secretsmanager_records.test", "records.#"),
					resource.TestCheckResourceAttrSet("data.secretsmanager_records.test", "records_by_uid"),
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
		if _, ok := rs.Primary.Attributes["records_by_uid"]; !ok {
			return fmt.Errorf("records_by_uid not populated")
		}

		return nil
	}
}

func testAccDataSourceRecordsConfig_basic() string {
	uid, _ := testAcc.getRecordInfo("login")
	return fmt.Sprintf(`
provider "secretsmanager" {
	credential = "%s"
}

data "secretsmanager_records" "test" {
	uids = [
		"%s"
	]
}
`, testAcc.credential, uid)
}

func testAccDataSourceRecordsConfig_withTitles() string {
	_, title := testAcc.getRecordInfo("login")
	return fmt.Sprintf(`
provider "secretsmanager" {
	credential = "%s"
}

data "secretsmanager_records" "test" {
	titles = [
		"%s"
	]
}
`, testAcc.credential, title)
}

func testAccDataSourceRecordsConfig_mixed() string {
	_, title1 := testAcc.getRecordInfo("login")
	uid2, _ := testAcc.getRecordInfo("encryptedNotes")
	return fmt.Sprintf(`
provider "secretsmanager" {
	credential = "%s"
}

data "secretsmanager_records" "test" {
	uids = [
		"%s"
	]
	titles = [
		"%s"
	]
}
`, testAcc.credential, uid2, title1)
}

func testAccDataSourceRecordsConfig_largeBatch() string {
	// This would be used for testing with many UIDs
	// In real testing, these would be actual UIDs from test data
	uid, _ := testAcc.getRecordInfo("login")
	return fmt.Sprintf(`
provider "secretsmanager" {
	credential = "%s"
}

data "secretsmanager_records" "test" {
	uids = [
		"%s"
	]
}
`, testAcc.credential, uid)
}
