---
subcategory: "Virtual Private Cloud (VPC)"
description: ""
page_title: "flexibleengine_vpc_route"
---

# flexibleengine_vpc_route

Manages a VPC route resource within Flexibleengine.

## Example Usage

### Add route to the default route table

```hcl
variable "nexthop" {}

resource "flexibleengine_vpc_v1" "example_vpc" {
  name = "example-vpc"
  cidr = "192.168.0.0/16"
}

resource "flexibleengine_vpc_route" "vpc_route" {
  vpc_id      = flexibleengine_vpc_v1.example_vpc.id
  destination = "192.168.0.0/16"
  type        = "peering"
  nexthop     = var.nexthop
}
```

### Add route to a custom route table

```hcl
variable "nexthop" {}

resource "flexibleengine_vpc_v1" "example_vpc" {
  name = "example-vpc"
  cidr = "192.168.0.0/16"
}

data "flexibleengine_vpc_route_table" "rtb" {
  vpc_id = flexibleengine_vpc_v1.example_vpc.id
  name   = "demo"
}

resource "flexibleengine_vpc_route" "vpc_route" {
  vpc_id         = flexibleengine_vpc_v1.example_vpc.id
  route_table_id = data.flexibleengine_vpc_route_table.rtb.id
  destination    = "172.16.8.0/24"
  type           = "ecs"
  nexthop        = var.nexthop
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the VPC route. If omitted, the provider-level
  region will be used. Changing this creates a new resource.

* `vpc_id` (Required, String, ForceNew) - Specifies the VPC for which a route is to be added. Changing this creates a
  new resource.

* `destination` (Required, String, ForceNew) - Specifies the destination address in the CIDR notation format,
  for example, 192.168.200.0/24. The destination of each route must be unique and cannot overlap with any
  subnet in the VPC. Changing this creates a new resource.

* `type` (Required, String) - Specifies the route type. Currently, the value can be:
  **ecs**, **eni**, **vip**, **nat**, **peering**, **vpn** and **dc**.

* `nexthop` (Required, String) - Specifies the next hop.
  + If the route type is **ecs**, the value is an ECS instance ID in the VPC.
  + If the route type is **eni**, the value is the extension NIC of an ECS in the VPC.
  + If the route type is **vip**, the value is a virtual IP address.
  + If the route type is **nat**, the value is a VPN gateway ID.
  + If the route type is **peering**, the value is a VPC peering connection ID.
  + If the route type is **vpn**, the value is a VPN gateway ID.
  + If the route type is **dc**, the value is a Direct Connect gateway ID.

* `description` (Optional, String) - Specifies the supplementary information about the route.
  The value is a string of no more than 255 characters and cannot contain angle brackets (< or >).

* `route_table_id` (Optional, String, ForceNew) - Specifies the route table ID for which a route is to be added.
  If the value is not set, the route will be added to the *default* route table.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The route ID, the format is `<route_table_id>/<destination>`

* `route_table_name` - The name of route table.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minute.
* `delete` - Default is 10 minute.

## Import

VPC routes can be imported using the route table ID and their `destination` separated by a slash, e.g.

```shell
terraform import flexibleengine_vpc_route.test <route_table_id>/<destination>
```
