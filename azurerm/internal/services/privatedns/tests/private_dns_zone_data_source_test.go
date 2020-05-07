package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMPrivateDNSZone_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_private_dns_zone", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPrivateDnsZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePrivateDNSZone_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMPrivateDNSZone_tags(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_private_dns_zone", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPrivateDnsZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePrivateDNSZone_tags(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.hello", "world"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMPrivateDNSZone_withoutResourceGroupName(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_private_dns_zone", "test")
	resourceGroupName := fmt.Sprintf("acctestRG-%d", data.RandomInteger)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPrivateDnsZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePrivateDNSZone_onlyNamePrep(data, resourceGroupName),
			},
			{
				Config: testAccDataSourcePrivateDNSZone_onlyName(data, resourceGroupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "resource_group_name", resourceGroupName),
				),
			},
		},
	})
}

func testAccDataSourcePrivateDNSZone_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_private_dns_zone" "test" {
  name                = "acctestzone%d.internal"
  resource_group_name = azurerm_resource_group.test.name
}

data "azurerm_private_dns_zone" "test" {
  name                = azurerm_private_dns_zone.test.name
  resource_group_name = azurerm_private_dns_zone.test.resource_group_name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccDataSourcePrivateDNSZone_tags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_private_dns_zone" "test" {
  name                = "acctestzone%d.internal"
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    hello = "world"
  }
}

data "azurerm_private_dns_zone" "test" {
  name                = azurerm_private_dns_zone.test.name
  resource_group_name = azurerm_private_dns_zone.test.resource_group_name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccDataSourcePrivateDNSZone_onlyNamePrep(data acceptance.TestData, resourceGroupName string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "%s"
  location = "%s"
}

resource "azurerm_private_dns_zone" "test" {
  name                = "acctestzone%d.internal"
  resource_group_name = azurerm_resource_group.test.name
}
`, resourceGroupName, data.Locations.Primary, data.RandomInteger)
}

func testAccDataSourcePrivateDNSZone_onlyName(data acceptance.TestData, resourceGroupName string) string {
	template := testAccDataSourcePrivateDNSZone_onlyNamePrep(data, resourceGroupName)
	return fmt.Sprintf(`
%s

data "azurerm_private_dns_zone" "test" {
  name = azurerm_private_dns_zone.test.name
}
`, template)
}
