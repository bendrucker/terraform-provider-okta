package okta

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAppGroupAssignment_crud(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := fmt.Sprintf("%s.test", appGroupAssignment)
	mgr := newFixtureManager(appGroupAssignment)
	config := mgr.GetFixtures("basic.tf", ri, t)
	updatedConfig := mgr.GetFixtures("updated.tf", ri, t)
	newUpdate := mgr.GetFixtures("force_new_update.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureAppGroupAssignmentExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "app_id"),
					resource.TestCheckResourceAttrSet(resourceName, "group_id"),
					resource.TestCheckResourceAttr(resourceName, "profile", "{}"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureAppGroupAssignmentExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "app_id"),
					resource.TestCheckResourceAttrSet(resourceName, "group_id"),
					resource.TestCheckResourceAttr(resourceName, "profile", "{}"),
				),
			},
			{
				Config: newUpdate,
				Check: resource.ComposeTestCheckFunc(
					ensureAppGroupAssignmentExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "app_id"),
					resource.TestCheckResourceAttrSet(resourceName, "group_id"),
					resource.TestCheckResourceAttr(resourceName, "profile", "{}"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("failed to find %s", resourceName)
					}

					appId := rs.Primary.Attributes["app_id"]
					groupId := rs.Primary.Attributes["group_id"]

					return fmt.Sprintf("%s/%s", appId, groupId), nil
				},
				ImportStateCheck: func(s []*terraform.InstanceState) error {
					if len(s) != 1 {
						return errors.New("Failed to import schema into state")
					}

					return nil
				},
			},
		},
	})
}

func ensureAppGroupAssignmentExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		missingErr := fmt.Errorf("resource not found: %s", name)
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return missingErr
		}

		appId := rs.Primary.Attributes["app_id"]
		groupId := rs.Primary.Attributes["group_id"]
		client := getOktaClientFromMetadata(testAccProvider.Meta())

		g, _, err := client.Application.GetApplicationGroupAssignment(appId, groupId, nil)
		if err != nil {
			return err
		} else if g == nil {
			return missingErr
		}

		return nil
	}
}
