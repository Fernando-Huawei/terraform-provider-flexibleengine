package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccEvsVolumesDataSource_basic(t *testing.T) {
	dataSourceName := "data.flexibleengine_evs_volumes.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)
	rName := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEvsVolumesDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "volumes.#", "5"),
				),
			},
		},
	})
}

func testAccEvsVolumesDataSource_base(rName string) string {
	return fmt.Sprintf(`
variable "volume_configuration" {
  type = list(object({
    suffix      = string
    size        = number
    device_type = string
    multiattach = bool
  }))
  default = [
    {suffix = "vbd_normal_volume", size = 100, device_type = "VBD", multiattach = false},
    {suffix = "vbd_share_volume", size = 100, device_type = "VBD", multiattach = true},
    {suffix = "scsi_normal_volume", size = 100, device_type = "SCSI", multiattach = false},
    {suffix = "scsi_share_volume", size = 100, device_type = "SCSI", multiattach = true},
  ]
}

%[1]s

resource "flexibleengine_compute_instance_v2" "test" {
  availability_zone = data.flexibleengine_availability_zones.test.names[0]
  name              = "%[2]s"
  image_id          = data.flexibleengine_images_image.test.id
  flavor_id         = data.flexibleengine_compute_flavors_v2.test.flavors[0]

  network {
    uuid = flexibleengine_vpc_subnet_v1.test.id
  }
}

resource "flexibleengine_evs_volume" "test" {
  count = length(var.volume_configuration)
  
  availability_zone = data.flexibleengine_availability_zones.test.names[0]
  volume_type       = "SSD"
  name              = "%[2]s_${var.volume_configuration[count.index].suffix}"
  size              = var.volume_configuration[count.index].size
  device_type       = var.volume_configuration[count.index].device_type
  multiattach       = var.volume_configuration[count.index].multiattach

  tags = {
    index = tostring(count.index)
  }
}

resource "flexibleengine_compute_volume_attach_v2" "test" {
  count = length(flexibleengine_evs_volume.test)

  instance_id = flexibleengine_compute_instance_v2.test.id
  volume_id   = flexibleengine_evs_volume.test[count.index].id
}
`, testBaseComputeResources(rName), rName)
}

func testAccEvsVolumesDataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "flexibleengine_evs_volumes" "test" {
  depends_on = [flexibleengine_compute_volume_attach_v2.test]

  availability_zone = data.flexibleengine_availability_zones.test.names[0]
  server_id         = flexibleengine_compute_instance_v2.test.id
  status            = "in-use"
}
`, testAccEvsVolumesDataSource_base(rName))
}
