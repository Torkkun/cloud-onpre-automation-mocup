package gcp

import (
	"pulumigcp/server"
	"strconv"

	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/compute"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type VmInstanceConfig struct {
	Args   compute.InstanceArgs
	Amount int
}

// Zones are currently hardcoded with fixed values
// MachinType e2-medium fixed value
func NewVMInstanceConfig(key_name string, vmclc server.VmConfig) *VmInstanceConfig {
	var image string
	switch vmclc.OS {
	case "ubuntu":
		image = "debian-cloud/debian-11"
	default:
		image = "debian-cloud/debian-11"
	}
	return &VmInstanceConfig{
		Args: compute.InstanceArgs{
			MachineType: pulumi.String("e2-medium"),
			Zone:        pulumi.String("asia-northeast1-b"),
			BootDisk: &compute.InstanceBootDiskArgs{
				InitializeParams: &compute.InstanceBootDiskInitializeParamsArgs{
					Image: pulumi.String(image),
				},
			},
		},
		Amount: vmclc.Amount,
	}
}

func (p *Provider) VMInstanceFunc(ctx *pulumi.Context, key string) error {
	config := *p.Config[key]
	if config.Vm == nil {
		// skip create vm instance
		return nil
	}
	// 指定した数だけ作成
	for i := 0; i < config.Vm.Amount; i++ {
		thisCompute := *p.Compute[key]
		// vm に割り当てるサブネットのコンフィグについて考える
		subnet := thisCompute.subnet[0]

		config.Vm.Args.Name = pulumi.String(key + "-vm-" + strconv.Itoa(i))

		config.Vm.Args.NetworkInterfaces = compute.InstanceNetworkInterfaceArray{
			&compute.InstanceNetworkInterfaceArgs{
				Subnetwork: subnet.Name,
			},
		}
		_, err := compute.NewInstance(ctx, key+"-vm-"+strconv.Itoa(i), &config.Vm.Args)
		if err != nil {
			return err
		}
		// ピアリングと同様の理由につきコメントアウトしている
		//thisCompute.instance = instance
	}
	return nil
}
