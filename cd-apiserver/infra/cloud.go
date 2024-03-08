package infra

import (
	"fmt"
	"pulumigcp/infra/azure"
	"pulumigcp/infra/gcp"
	"pulumigcp/server"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type CloudInformation struct {
	Name         string // cloud provider name :GCP or Azure
	ClientConfig *server.ClientConfig
}

// cloud provider毎に分ける
func NewCloud(provider *CloudInformation) (Cloud, []string, error) {
	switch provider.Name {
	case "GCP":
		return gcp.NewConfig(provider.ClientConfig)
	case "Azure":
		cloud := &azure.Provider{}
		return cloud, nil, nil
	}
	return nil, nil, nil
}

// クラウドの操作を抽象化
// pulumi以外に操作しないのでここはpulumiContextのままで
type Cloud interface {
	NetworkFunc(ctx *pulumi.Context, network_key string) error
	NetworkSubnetFunc(*pulumi.Context, string) error
	NetworkPeeringFunc(*pulumi.Context, string) error
	VPNGatewayFunc(*pulumi.Context, string) error
	VPNTunnelFunc(*pulumi.Context, string) error
	VMInstanceFunc(*pulumi.Context, string) error
}

// 設定を元に順番にStackする操作群を作成する
func CreateNetworksForStacking(cloud Cloud, keys []string) pulumi.RunFunc {
	return func(ctx *pulumi.Context) error {
		fmt.Println("Creating Upload Stack")
		var err error
		for _, key := range keys {
			err = cloud.NetworkFunc(ctx, key)
			if err != nil {
				fmt.Printf("Failed create Network Func: %v", err)
				return err
			}
			err = cloud.NetworkSubnetFunc(ctx, key)
			if err != nil {
				fmt.Printf("Failed create NetworkSubnet Func: %v", err)
				return err
			}
			err = cloud.VPNGatewayFunc(ctx, key)
			if err != nil {
				fmt.Printf("Failed create VpnGw Func: %v", err)
				return err
			}
			err = cloud.VPNTunnelFunc(ctx, key)
			if err != nil {
				fmt.Printf("Failed create VpnTunnnel Func: %v", err)
				return err
			}
			err = cloud.VMInstanceFunc(ctx, key)
			if err != nil {
				fmt.Printf("Failed create VMInstance Func: %v", err)
				return err
			}

		}
		// ここはなんとか並列処理等で解決したい
		for _, key := range keys {
			err = cloud.NetworkPeeringFunc(ctx, key)
			if err != nil {
				fmt.Printf("Failed create NetworkPeering Func: %v", err)
				return err
			}
		}
		return nil
	}
}
