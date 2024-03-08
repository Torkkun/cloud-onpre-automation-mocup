package gcp

import (
	"fmt"
	"pulumigcp/server"

	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/compute"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type VpcNetworkConfig struct {
	Args compute.NetworkArgs
}

func NewNetworkConfig(key_name string, subnet_bool bool) *VpcNetworkConfig {
	return &VpcNetworkConfig{
		Args: compute.NetworkArgs{
			Name:                  pulumi.String("vpc-" + key_name),
			AutoCreateSubnetworks: pulumi.Bool(subnet_bool),
		}}
}

func (p *Provider) NetworkFunc(ctx *pulumi.Context, key string) error {
	config := *p.Config[key]
	if config.Nw == nil {
		return fmt.Errorf("no network, Please create a network")
	}

	network, err := compute.NewNetwork(ctx, key, &config.Nw.Args)
	if err != nil {
		return err
	}

	// init & create GcpCompute data
	*p.Compute[key] = GcpCompute{network: network}
	return nil
}

type VpcSubnetworkConfig struct {
	Args compute.SubnetworkArgs
}

func NewSubnetworkConfig(key_name string, subnetRanges []string, secondary []server.SubnetworkConfig) ([]*VpcSubnetworkConfig, error) {
	configs := []*VpcSubnetworkConfig{}

	if secondary != nil {
		return nil, fmt.Errorf("secondary subnetwork not yet implemented")
	}

	for i, subnetRange := range subnetRanges {
		config := VpcSubnetworkConfig{
			Args: compute.SubnetworkArgs{
				Name:        pulumi.String(fmt.Sprintf(key_name+"-subnet-%d", i)),
				IpCidrRange: pulumi.String(subnetRange),
				Region:      pulumi.String("asia-northeast1"),
			},
		}
		configs = append(configs, &config)
	}

	return configs, nil
}

func (p *Provider) NetworkSubnetFunc(ctx *pulumi.Context, key string) error {
	config := *p.Config[key]
	if config.Subnets == nil {
		// skip create subnetwork
		return nil
	}

	thisCompute := *p.Compute[key]
	for i, con := range config.Subnets {
		//fmt.Printf("%v\n", )

		con.Args.Network = thisCompute.network.ID()
		subnet, err := compute.NewSubnetwork(ctx, fmt.Sprintf("%s-subnet-%d", key, i), &con.Args)
		if err != nil {
			return err
		}
		thisCompute.subnet = append(thisCompute.subnet, subnet)
	}

	*p.Compute[key] = thisCompute
	return nil
}

type NwPeeringConfig struct {
	Args compute.NetworkPeeringArgs
}

// Not Used
func NewNetworkPeeringConfig() *NwPeeringConfig {
	return &NwPeeringConfig{Args: compute.NetworkPeeringArgs{}}
}

// gcpは相互にピアを各々で貼る必要がある
func (p *Provider) NetworkPeeringFunc(ctx *pulumi.Context, key string) error {
	fmt.Println("### peer start")
	if p.NwPeering_keyList[key] == nil {
		fmt.Println("### perkymap lis nil")
	}
	peerkeys := p.NwPeering_keyList[key]
	if peerkeys == nil {
		fmt.Println("### perkys lis nil")
	}
	fmt.Println("### peer start")
	//config := *p.Config[key] // not used
	var config compute.NetworkPeeringArgs
	thisCompute := *p.Compute[key]
	//config.Peering.Args.Network = thisCompute.network.SelfLink
	config.Network = thisCompute.network.SelfLink

	for _, peerkey := range peerkeys {
		peerCompute := *p.Compute[peerkey]

		//config.Peering.Args.PeerNetwork = peerCompute.network.SelfLink
		//config.Peering.Args.Name = pulumi.String(key + "-to-" + peerkey + "-peer")

		config.PeerNetwork = peerCompute.network.SelfLink
		config.Name = pulumi.String(key + "-to-" + peerkey + "-peer")
		fmt.Println("### peer create start")
		_, err := compute.NewNetworkPeering(ctx, key+"to"+peerkey+"peering", &config)
		if err != nil {
			return err
		}
		// 一個しか入れられないので、更新されていってしまうためmapにするか考える。そもそも使わないので消すかも
		//thisCompute.peering = peering
	}
	fmt.Println("### peer END")
	return nil
}
