---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_server"
subcategory: "Compute"
description: |-
  Manages a SakuraCloud Server.
---

# sakuracloud_server

Manages a SakuraCloud Server.

## Example Usage

```hcl
resource "sakuracloud_server" "foobar" {
  name        = "foobar"
  disks       = [sakuracloud_disk.foobar.id]
  description = "description"
  tags        = ["tag1", "tag2"]

  network_interface {
    upstream         = "shared"
    packet_filter_id = data.sakuracloud_packet_filter.foobar.id
  }

  disk_edit_parameter {
    hostname        = "hostname"
    password        = "password"
    disable_pw_auth = true

    # ssh_keys    = ["ssh-rsa xxxxx"]
    # ssh_key_ids = ["<ID>", "<ID>"]
    # note {
    #  id         = "<ID>"
    #  api_key_id = "<ID>"
    #  variables = {
    #    foo = "bar"
    #  }
    # }
  }

  # If you use cloud-init instead of disk_edit_parameter
  
  # user_data = join("\n", [
  #   "#cloud-config",
  #   yamlencode({
  #     hostname: "hostname",
  #     password: "password",
  #     chpasswd: {
  #       expire: false,
  #     }
  #     ssh_pwauth: false,
  #     ssh_authorized_keys: [
  #       file("~/.ssh/id_rsa.pub"),
  #     ],
  #   }),
  # ])
}

data "sakuracloud_packet_filter" "foobar" {
  filter {
    names = ["foobar"]
  }
}

data "sakuracloud_archive" "ubuntu" {
  os_type = "ubuntu"
}

resource "sakuracloud_disk" "foobar" {
  name              = "foobar"
  source_archive_id = data.sakuracloud_archive.ubuntu.id
}
```

## Argument Reference

* `name` - (Required) The name of the Server. The length of this value must be in the range [`1`-`64`].
* `cdrom_id` - (Optional) The id of the CD-ROM to attach to the Server.
* `force_shutdown` - (Optional) The flag to use force shutdown when need to reboot/shutdown while applying.

#### Spec

* `commitment` - (Optional) The policy of how to allocate virtual CPUs to the server. This must be one of [`standard`/`dedicatedcpu`]. Default:`standard`.
* `core` - (Optional) The number of virtual CPUs. Default:`1`.
* `memory` - (Optional) The size of memory in GiB. Default:`1`.
* `cpu_model` - (Optional) The model of CPU.
* `gpu` - (Optional) The number of GPUs.
* `gpu_model` - (Optional) The model of GPU.
* `network_interface` - (Optional) One or more `network_interface` blocks as defined below.
* `interface_driver` - (Optional) The driver name of network interface. This must be one of [`virtio`/`e1000`]. Default:`virtio`.
* `private_host_id` - (Optional) The id of the PrivateHost which the Server is assigned.

The values that can be specified for `commitment`, `core`, `memory`, `cpu_model`, `gpu`, and `gpu_model` can be found with the following command.

```bash
usacloud iaas server-plan list --zone is1a
```

Note: This command requires usacloud v1.17 or later.

---

A `network_interface` block supports the following:

* `upstream` - (Required) The upstream type or upstream switch id. This must be one of [`shared`/`disconnect`/`<switch id>`].
* `packet_filter_id` - (Optional) The id of the packet filter to attach to the network interface.
* `user_ip_address` - (Optional) The IP address for only display. This value doesn't affect actual NIC settings.



#### Disks

* `disk_edit_parameter` - (Optional) A `disk_edit_parameter` block as defined below. This parameter conflicts with [`user_data`].
* `user_data` - (Optional) A string representing the user data used by cloud-init. This parameter conflicts with [`disk_edit_parameter`].
* `disks` - (Optional) A list of disk id connected to the server.

---

A `disk_edit_parameter` block supports the following:

* `change_partition_uuid` - (Optional) The flag to change partition uuid.
* `disable_pw_auth` - (Optional) The flag to disable password authentication.
* `enable_dhcp` - (Optional) The flag to enable DHCP client.
* `gateway` - (Optional) The gateway address used by the Server.
* `hostname` - (Optional) The hostname of the Server. The length of this value must be in the range [`1`-`64`].
* `ip_address` - (Optional) The IP address to assign to the Server.
* `netmask` - (Optional) The bit length of the subnet to assign to the Server.
* `note` - (Optional) A list of the `note` block as defined below.
* `note_ids` - (Optional/Deprecated) A list of the Note id.  
Note: **The `note_ids` will be removed in a future version. Please use the `note` instead**
* `password` - (Optional) The password of default user. The length of this value must be in the range [`8`-`64`].
* `ssh_key_ids` - (Optional) A list of the SSHKey id.
* `ssh_keys` - (Optional) A list of the SSHKey text.

---

A `note` block supports the following:

* `id` - (Required) The id of the Note/StartupScript.
* `api_key_id` - (Optional) The id of the API key to be injected into the Note/StartupScript when editing the disk.
* `variables` - (Optional) The value of the variable that be injected into the Note/StartupScript when editing the disk.

#### Common Arguments

* `description` - (Optional) The description of the Server. The length of this value must be in the range [`1`-`512`].
* `icon_id` - (Optional) The icon id to attach to the Server.
* `tags` - (Optional) Any tags to assign to the Server.
* `zone` - (Optional) The name of zone that the Server will be created. (e.g. `is1a`, `tk1a`). Changing this forces a new resource to be created.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the Server
* `update` - (Defaults to 5 minutes) Used when updating the Server
* `delete` - (Defaults to 20 minutes) Used when deleting Server

## Attribute Reference

* `id` - The id of the Server.
* `dns_servers` - A list of IP address of DNS server in the zone.
* `gateway` - The IP address of the gateway used by Server.
* `hostname` - The hostname of the Server.
* `ip_address` - The IP address assigned to the Server.
* `netmask` - The bit length of the subnet assigned to the Server.
* `network_address` - The network address which the `ip_address` belongs.
* `private_host_name` - The id of the PrivateHost which the Server is assigned.

---

A `network_interface` block exports the following:

* `mac_address` - The MAC address.

