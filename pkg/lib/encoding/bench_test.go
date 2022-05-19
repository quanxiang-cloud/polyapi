package encoding

import (
	"fmt"
	"testing"
)

var benchData = MustFromJSON(`
{
  "req1": {
    "action": "DescribeInstancesResponse",
    "instance_set": [
      {
        "logic_volumes": [
          {
            "device": "vda1",
            "volume_id": "os"
          }
        ],
        "vxnets": [
          {
            "ipv6_address": "",
            "vxnet_type": 2,
            "vxnet_id": "vxnet-r2wedwp",
            "vxnet_name": "vxnet-0",
            "role": 1,
            "private_ip": "10.180.11.2",
            "security_group": {
              "is_default": 1,
              "security_group_name": "default security group",
              "security_group_id": "sg-c6ys6akl"
            },
            "nic_id": "52:54:9e:20:33:79",
            "security_groups": [
              {
                "is_default": 1,
                "security_group_name": "default security group",
                "security_group_id": "sg-c6ys6akl"
              }
            ]
          }
        ],
        "memory_current": 1024,
        "fence": null,
        "extra": {
          "mem_max": 0,
          "nic_mqueue": 0,
          "read_throughput": 0,
          "ivshmem": [],
          "gpu_pci_nums": "",
          "cpu_max": 0,
          "cpu_model": "",
          "bandwidth": 200,
          "iops": 660,
          "throughput": 39936,
          "read_iops": 0,
          "hypervisor": "kvm",
          "os_disk_encryption": 0,
          "gpu": 0,
          "os_disk_size": 20,
          "gpu_class": 0,
          "features": 4
        },
        "image": {
          "ui_type": "tui",
          "processor_type": "64bit",
          "platform": "linux",
          "features_supported": {
            "set_keypair": 1,
            "disk_hot_plug": 1,
            "root_fs_rw_online": 1,
            "user_data": 1,
            "set_pwd": 1,
            "root_fs_rw_offline": 1,
            "ipv6_supported": 1,
            "nic_hot_plug": 1,
            "join_multiple_managed_vxnets": 0,
            "reset_fstab": 1
          },
          "image_size": 20,
          "image_name": "Ubuntu Server 20.04.1 LTS 64bit",
          "image_id": "focal1x64",
          "os_family": "ubuntu",
          "provider": "system",
          "features": 64
        },
        "graphics_passwd": "xxx",
        "dns_aliases": [],
        "alarm_status": "",
        "owner": "usr-msnWeHKp",
        "security_groups": [
            {
              "is_default": 1,
              "security_group_name": "default security group",
              "security_group_id": "sg-c6ys6akl"
            }
          ],
        "keypair_ids": [
            "kp-jm02koxi"
          ],
        "vcpus_current": 1,
        "instance_id": "i-pc8k55du",
        "sub_code": 0,
        "graphics_protocol": "vnc",
        "platform": "linux",
        "instance_class": 101,
        "status_time": "2021-05-29T17:11:26Z",
        "status": "running",
        "description": null,
        "cpu_topology": "",
        "tags": [],
        "transition_status": "",
        "eips": [],
        "repl": "rpp-00000002",
        "volume_ids": [],
        "zone_id": "pek3d",
        "lastest_snapshot_time": null,
        "instance_group": null,
        "instance_name": "allyvpc1",
        "instance_type": "s1.small.r1",
        "create_time": "2021-05-27T13:20:42Z",
        "volumes": [],
        "security_group": {
                    "is_default": 1,
                    "security_group_name": "default security group",
                    "security_group_id": "sg-c6ys6akl"
                  },
        "resource_project_info": []
      }
    ],
    "total_count": 1,
    "ret_code": 0
  },
  "req2": {
    "action": "StopInstancesResponse",
    "job_id": "j-1v211t6p04l",
    "ret_code": 0
  },
  "merged": {
    "InstNummmmmmmmmmmm": 1,
    "CodeReq11111111111111": 0,
    "action": "StopInstancesResponse",
    "job_id": "j-1v211t6p04l",
    "ret_code": 0
  }
} 
`)

func Benchmark_Json(b *testing.B) {
	s := MustToJSON(benchData, true)
	MustFromJSON(s)
}

func Benchmark_Yaml(b *testing.B) {
	s := MustToYAML(benchData, true)
	MustFromYAML(s)
}

func Benchmark_Xml(b *testing.B) {
	s := MustToXML(benchData, true)
	MustFromXML(s)
}

func TestEncoding2(t *testing.T) {
	jsn := MustToJSON(benchData, true)
	yml := MustToYAML(benchData, true)
	xml := MustToXML(benchData, true)
	fmt.Printf("data: \n%#v\n", benchData)
	fmt.Printf("json: \n%s\n", jsn)
	fmt.Printf("yaml: \n%s\n", yml)
	fmt.Printf("xml : \n%s\n", xml)
}
