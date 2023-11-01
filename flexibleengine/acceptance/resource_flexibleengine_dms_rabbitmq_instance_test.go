package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/dms/v2/rabbitmq/instances"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getDmsRabbitMqInstanceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.DmsV2Client(OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating FlexibleEngine DMS client(V2): %s", err)
	}
	return instances.Get(client, state.Primary.ID).Extract()
}

func getKafkaInstanceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.DmsV2Client(OS_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating FlexibleEngine DMS client(V2): %s", err)
	}
	return instances.Get(client, state.Primary.ID).Extract()
}

func TestAccDmsRabbitmqInstances_basic(t *testing.T) {
	var instance instances.Instance
	rName := acceptance.RandomAccResourceNameWithDash()
	updateName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "flexibleengine_dms_rabbitmq_instance.test"
	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getDmsRabbitMqInstanceFunc,
	)

	// DMS instances use the tenant-level shared lock, the instances cannot be created or modified in parallel.
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDmsRabbitmqInstance_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "engine", "rabbitmq"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
				),
			},
			{
				Config: testAccDmsRabbitmqInstance_update(rName, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "product_id", "${data.flexibleengine_dms_product.test2.id}"),
					resource.TestCheckResourceAttr(resourceName, "description", "rabbitmq test update"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform_update"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password", "used_storage_space",
				},
			},
		},
	})
}

func TestAccDmsRabbitmqInstances_newFormat_cluster(t *testing.T) {
	var instance instances.Instance
	rName := acceptance.RandomAccResourceNameWithDash()
	updateName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "flexibleengine_dms_rabbitmq_instance.test"
	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getDmsRabbitMqInstanceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDmsRabbitmqInstance_newFormat_cluster(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "rabbitmq test"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
				),
			},
			{
				Config: testAccDmsRabbitmqInstance_newFormat_cluster_update(rName, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "description", "rabbitmq test update"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform_update"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password", "used_storage_space",
				},
			},
		},
	})
}

func TestAccDmsRabbitmqInstances_newFormat_single(t *testing.T) {
	var instance instances.Instance
	rName := acceptance.RandomAccResourceNameWithDash()
	updateName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "flexibleengine_dms_rabbitmq_instance.test"
	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getDmsRabbitMqInstanceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDmsRabbitmqInstance_newFormat_single(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "rabbitmq test"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
				),
			},
			{
				Config: testAccDmsRabbitmqInstance_newFormat_single_update(rName, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "description", "rabbitmq test update"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform_update"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password", "used_storage_space",
				},
			},
		},
	})
}

func TestAccDmsRabbitmqInstances_prePaid(t *testing.T) {
	var instance instances.Instance
	rName := acceptance.RandomAccResourceNameWithDash()
	updateName := rName + "update"
	resourceName := "flexibleengine_dms_rabbitmq_instance.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getKafkaInstanceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDmsRabbitmqInstance_prePaid(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "rabbitmq test"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "prePaid"),
					resource.TestCheckResourceAttr(resourceName, "broker_num", "3"),
				),
			},
			{
				Config: testAccDmsRabbitmqInstance_prePaid_update(rName, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "description", "rabbitmq test update"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform_update"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "prePaid"),
					resource.TestCheckResourceAttr(resourceName, "broker_num", "3"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
					"auto_renew",
					"period",
					"period_unit",
				},
			},
		},
	})
}

func TestAccDmsRabbitmqInstances_withEpsId(t *testing.T) {
	var instance instances.Instance
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "flexibleengine_dms_rabbitmq_instance.test"
	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getDmsRabbitMqInstanceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckEpsID(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDmsRabbitmqInstance_withEpsId(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "engine", "rabbitmq"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", OS_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func TestAccDmsRabbitmqInstances_compatible(t *testing.T) {
	var instance instances.Instance
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "flexibleengine_dms_rabbitmq_instance.test"
	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getDmsRabbitMqInstanceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDmsRabbitmqInstance_compatible(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "engine", "rabbitmq"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password", "used_storage_space",
				},
			},
		},
	})
}

func TestAccDmsRabbitmqInstances_single(t *testing.T) {
	var instance instances.Instance
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "flexibleengine_dms_rabbitmq_instance.test"
	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getDmsRabbitMqInstanceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDmsRabbitmqInstance_single(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
		},
	})
}

func testAccDmsRabbitmqInstance_Base(rName string) string {
	return fmt.Sprintf(`
%s

data "flexibleengine_availability_zones" "test" {}

data "flexibleengine_dms_product" "test1" {
  engine        = "kafka"
  bandwidth     = "100MB"
}

data "flexibleengine_dms_product" "test2" {
  engine        = "kafka"
  bandwidth     = "100MB"
}
`, testBaseNetwork(rName))
}

func testAccDmsRabbitmqInstance_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_dms_rabbitmq_instance" "test" {
  name        = "%s"
  description = "rabbitmq test"
  
  vpc_id             = flexibleengine_vpc_v1.test.id
  network_id         = flexibleengine_vpc_subnet_v1.test.id
  security_group_id  = flexibleengine_networking_secgroup_v2.test.id
  availability_zones = [
    data.flexibleengine_availability_zones.test.names[0]
  ]

  product_id        = data.flexibleengine_dms_product.test1.id
  storage_spec_code = "dms.physical.storage.high"

  access_user = "user"
  password    = "Rabbitmqtest@123"

  tags = {
    key   = "value"
    owner = "terraform"
  }
}
`, testAccDmsRabbitmqInstance_Base(rName), rName)
}

func testAccDmsRabbitmqInstance_update(rName, updateName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_dms_rabbitmq_instance" "test" {
  name        = "%s"
  description = "rabbitmq test update"

  vpc_id             = flexibleengine_vpc_v1.test.id
  network_id         = flexibleengine_vpc_subnet_v1.test.id
  security_group_id  = flexibleengine_networking_secgroup_v2.test.id
  availability_zones = [
    data.flexibleengine_availability_zones.test.names[0]
  ]

  product_id        = data.flexibleengine_dms_product.test2.id
  engine_version    = data.flexibleengine_dms_product.test2.version

  access_user = "user"
  password    = "Rabbitmqtest@123"

  tags = {
    key1  = "value"
    owner = "terraform_update"
  }
}
`, testAccDmsRabbitmqInstance_Base(rName), updateName)
}

func testAccDmsRabbitmqInstance_withEpsId(rName string) string {
	return fmt.Sprintf(`
%s

resource "flexibleengine_dms_rabbitmq_instance" "test" {
  name                  = "%s"
  description           = "rabbitmq test"
  enterprise_project_id = "%s"

  vpc_id             = flexibleengine_vpc_v1.test.id
  network_id         = flexibleengine_vpc_subnet_v1.test.id
  security_group_id  = flexibleengine_networking_secgroup_v2.test.id
  availability_zones = [
    data.flexibleengine_availability_zones.test.names[0],
    data.flexibleengine_availability_zones.test.names[1]
  ]
  
  product_id        = data.flexibleengine_dms_product.test1.id

  access_user = "user"
  password    = "Rabbitmqtest@123"

  tags = {
    key   = "value"
    owner = "terraform"
  }
}
`, testAccDmsRabbitmqInstance_Base(rName), rName, OS_ENTERPRISE_PROJECT_ID_TEST)
}

// After the 1.31.1 version, arguments storage_space and available_zones are deprecated.
func testAccDmsRabbitmqInstance_compatible(rName string) string {
	return fmt.Sprintf(`
%s

data "flexibleengine_dms_az" "test" {}
resource "flexibleengine_dms_rabbitmq_instance" "test" {
  name        = "%s"
  description = "rabbitmq test"
  
  vpc_id            = flexibleengine_vpc_v1.test.id
  network_id        = flexibleengine_vpc_subnet_v1.test.id
  security_group_id = flexibleengine_networking_secgroup_v2.test.id
  available_zones   = [data.flexibleengine_dms_az.test.id]

  product_id        = data.flexibleengine_dms_product.test1.id
  storage_space     = data.flexibleengine_dms_product.test1.storage


  access_user = "user"
  password    = "Rabbitmqtest@123"

  tags = {
    key   = "value"
    owner = "terraform"
  }
}
`, testAccDmsRabbitmqInstance_Base(rName), rName)
}

func testAccDmsRabbitmqInstance_single(rName string) string {
	randPwd := fmt.Sprintf("%s!#%d", acctest.RandString(5), acctest.RandIntRange(0, 999))
	return fmt.Sprintf(`
%[1]s

data "flexibleengine_dms_product" "single" {
  engine           = "rabbitmq"
  instance_type    = "single"
  version          = "3.8.35"
  node_num         = 1
}

resource "flexibleengine_dms_rabbitmq_instance" "test" {
  availability_zones = [
    data.flexibleengine_availability_zones.test.names[0],
  ]

  name              = "%[2]s"
  vpc_id            = flexibleengine_vpc_v1.test.id
  network_id        = flexibleengine_vpc_subnet_v1.test.id
  security_group_id = flexibleengine_networking_secgroup_v2.test.id

  product_id        = data.flexibleengine_dms_product.single.id
  engine_version    = data.flexibleengine_dms_product.single.version
  storage_space     = data.flexibleengine_dms_product.single.storage

  access_user = "root"
  password    = "%[3]s"
}
`, testAccDmsRabbitmqInstance_Base(rName), rName, randPwd)
}

func testAccDmsRabbitmqInstance_newFormat_cluster(rName string) string {
	return fmt.Sprintf(`
%s

data "flexibleengine_availability_zones" "test" {}

data "flexibleengine_dms_rabbitmq_flavors" "test" {
  type = "cluster"
}

locals {
  query_results = data.flexibleengine_dms_rabbitmq_flavors.test
  flavor        = data.flexibleengine_dms_rabbitmq_flavors.test.flavors[0]
}

resource "flexibleengine_dms_rabbitmq_instance" "test" {
  name        = "%s"
  description = "rabbitmq test"
  
  vpc_id            = flexibleengine_vpc_v1.test.id
  network_id        = flexibleengine_vpc_subnet_v1.test.id
  security_group_id = flexibleengine_networking_secgroup_v2.test.id

  availability_zones = [
    data.flexibleengine_availability_zones.test.names[0]
  ]

  flavor_id         = local.flavor.id
  engine_version    = element(local.query_results.versions, length(local.query_results.versions)-1)
  storage_space     = local.flavor.properties[0].min_broker * local.flavor.properties[0].min_storage_per_node
  broker_num        = 3
  access_user       = "user"
  password          = "Rabbitmqtest@123"

  tags = {
    key   = "value"
    owner = "terraform"
  }
}`, testBaseNetwork(rName), rName)
}

func testAccDmsRabbitmqInstance_newFormat_cluster_update(rName, updateName string) string {
	return fmt.Sprintf(`
%s

data "flexibleengine_availability_zones" "test" {}

data "flexibleengine_dms_rabbitmq_flavors" "test" {
  type = "cluster"
}

locals {
  query_results = data.flexibleengine_dms_rabbitmq_flavors.test
  flavor        = data.flexibleengine_dms_rabbitmq_flavors.test.flavors[0]
  newFlavor     = data.flexibleengine_dms_rabbitmq_flavors.test.flavors[1]
}

resource "flexibleengine_dms_rabbitmq_instance" "test" {
  name        = "%s"
  description = "rabbitmq test update"
  
  vpc_id            = flexibleengine_vpc_v1.test.id
  network_id        = flexibleengine_vpc_subnet_v1.test.id
  security_group_id = flexibleengine_networking_secgroup_v2.test.id

  availability_zones = [
    data.flexibleengine_availability_zones.test.names[0]
  ]

  flavor_id         = local.newFlavor.id
  engine_version    = element(local.query_results.versions, length(local.query_results.versions)-1)
  storage_space     = 1000
  storage_spec_code = local.flavor.ios[0].storage_spec_code
  broker_num        = 5
  access_user       = "user"
  password          = "Rabbitmqtest@123"

  tags = {
    key1  = "value"
    owner = "terraform_update"
  }
}`, testBaseNetwork(rName), updateName)
}

func testAccDmsRabbitmqInstance_newFormat_single(rName string) string {
	return fmt.Sprintf(`
%s

data "flexibleengine_availability_zones" "test" {}

data "flexibleengine_dms_rabbitmq_flavors" "test" {
  type = "single"
}

locals {
  query_results = data.flexibleengine_dms_rabbitmq_flavors.test
  flavor        = data.flexibleengine_dms_rabbitmq_flavors.test.flavors[0]
}

resource "flexibleengine_dms_rabbitmq_instance" "test" {
  name        = "%s"
  description = "rabbitmq test"
  
  vpc_id            = flexibleengine_vpc_v1.test.id
  network_id        = flexibleengine_vpc_subnet_v1.test.id
  security_group_id = flexibleengine_networking_secgroup_v2.test.id

  availability_zones = [
    data.flexibleengine_availability_zones.test.names[0]
  ]

  flavor_id         = local.flavor.id
  engine_version    = element(local.query_results.versions, length(local.query_results.versions)-1)
  storage_space     = local.flavor.properties[0].min_broker * local.flavor.properties[0].min_storage_per_node
  storage_spec_code = local.flavor.ios[0].storage_spec_code
  access_user       = "user"
  password          = "Rabbitmqtest@123"

  tags = {
    key   = "value"
    owner = "terraform"
  }
}`, testBaseNetwork(rName), rName)
}

func testAccDmsRabbitmqInstance_newFormat_single_update(rName, updateName string) string {
	return fmt.Sprintf(`
%s

data "flexibleengine_availability_zones" "test" {}

data "flexibleengine_dms_rabbitmq_flavors" "test" {
  type = "single"
}

locals {
  query_results = data.flexibleengine_dms_rabbitmq_flavors.test
  newFlavor     = data.flexibleengine_dms_rabbitmq_flavors.test.flavors[1]
}

resource "flexibleengine_dms_rabbitmq_instance" "test" {
  name        = "%s"
  description = "rabbitmq test update"
  
  vpc_id            = flexibleengine_vpc_v1.test.id
  network_id        = flexibleengine_vpc_subnet_v1.test.id
  security_group_id = flexibleengine_networking_secgroup_v2.test.id

  availability_zones = [
    data.flexibleengine_availability_zones.test.names[0]
  ]

  flavor_id         = local.newFlavor.id
  engine_version    = element(local.query_results.versions, length(local.query_results.versions)-1)
  storage_space     = 600
  storage_spec_code = local.newFlavor.ios[0].storage_spec_code
  access_user       = "user"
  password          = "Rabbitmqtest@123"

  tags = {
    key1  = "value"
    owner = "terraform_update"
  }
}`, testBaseNetwork(rName), updateName)
}

func testAccDmsRabbitmqInstance_prePaid(rName string) string {
	return fmt.Sprintf(`
%s

data "flexibleengine_availability_zones" "test" {}

data "flexibleengine_dms_rabbitmq_flavors" "test" {
  type = "cluster"
}

locals {
  query_results = data.flexibleengine_dms_rabbitmq_flavors.test
  flavor        = data.flexibleengine_dms_rabbitmq_flavors.test.flavors[0]
}

resource "flexibleengine_dms_rabbitmq_instance" "test" {
  name        = "%s"
  description = "rabbitmq test"
  
  vpc_id            = flexibleengine_vpc_v1.test.id
  network_id        = flexibleengine_vpc_subnet_v1.test.id
  security_group_id = flexibleengine_networking_secgroup_v2.test.id

  availability_zones = [
    data.flexibleengine_availability_zones.test.names[0]
  ]

  flavor_id         = local.flavor.id
  engine_version    = element(local.query_results.versions, length(local.query_results.versions)-1)
  storage_space     = local.flavor.properties[0].min_broker * local.flavor.properties[0].min_storage_per_node
  storage_spec_code = local.flavor.ios[0].storage_spec_code
  broker_num        = 3

  access_user = "user"
  password    = "Rabbitmqtest@123"

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = true

  tags = {
    key   = "value"
    owner = "terraform"
  }
}`, testBaseNetwork(rName), rName)
}

func testAccDmsRabbitmqInstance_prePaid_update(rName, updateName string) string {
	return fmt.Sprintf(`
%s

data "flexibleengine_availability_zones" "test" {}

data "flexibleengine_dms_rabbitmq_flavors" "test" {
  type = "cluster"
}

locals {
  query_results = data.flexibleengine_dms_rabbitmq_flavors.test
  flavor        = data.flexibleengine_dms_rabbitmq_flavors.test.flavors[0]
}

resource "flexibleengine_dms_rabbitmq_instance" "test" {
  name        = "%s"
  description = "rabbitmq test update"
  
  vpc_id            = flexibleengine_vpc_v1.test.id
  network_id        = flexibleengine_vpc_subnet_v1.test.id
  security_group_id = flexibleengine_networking_secgroup_v2.test.id

  availability_zones = [
    data.flexibleengine_availability_zones.test.names[0]
  ]

  flavor_id         = local.flavor.id
  engine_version    = element(local.query_results.versions, length(local.query_results.versions)-1)
  storage_space     = local.flavor.properties[0].min_broker * local.flavor.properties[0].min_storage_per_node
  storage_spec_code = local.flavor.ios[0].storage_spec_code
  broker_num        = 3

  access_user = "user"
  password    = "Rabbitmqtest@123"

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = true

  tags = {
    key1  = "value"
    owner = "terraform_update"
  }
}`, testBaseNetwork(rName), updateName)
}
