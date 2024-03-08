package gcp

import (
	"fmt"
	"pulumigcp/server"

	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/compute"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type VpnGwConfig struct {
	vpnType string
	Args    compute.VPNGatewayArgs
}

func NewVPNGwConfig(key_name string) *VpnGwConfig {
	return &VpnGwConfig{
		Args: compute.VPNGatewayArgs{
			Name: pulumi.String("vpn-gw-" + key_name),
		},
		vpnType: "classic",
	}
}

func (p *Provider) VPNGatewayFunc(ctx *pulumi.Context, key string) error {
	config := *p.Config[key]
	if config.VpnGw == nil {
		// skip create vpn gateway
		return nil
	}

	thisCompute := *p.Compute[key]

	fmt.Printf("##VpnGwConfig:%v\n", *p.Config[key].VpnGw)
	fmt.Printf("##VpnGwCom:%v\n", *p.Compute[key])

	config.VpnGw.Args.Network = thisCompute.network.ID()

	switch config.VpnGw.vpnType {
	case "classic":
		classicVpnGw, err := compute.NewVPNGateway(ctx, key, &config.VpnGw.Args)
		if err != nil {
			return err
		}
		thisCompute.vpnGw = classicVpnGw
		*p.Compute[key] = thisCompute
		return nil

	case "HA":
		return fmt.Errorf("no plans for implementation, only location")
	default:
		return fmt.Errorf("the specified vpn type does not exist.\nclassic or HA")
	}
}

type VpnTunnelConfig struct {
	Args      compute.VPNTunnelArgs
	DestRange string
}

func NewVPNTunnelConfig(key_name string, vpncon server.VPNConfig) *VpnTunnelConfig {
	return &VpnTunnelConfig{
		Args: compute.VPNTunnelArgs{
			Name:   pulumi.String(key_name + "-tunnel"),
			PeerIp: pulumi.String(vpncon.PeerIp),
			// 暗号化する必要がある
			SharedSecret: pulumi.String(vpncon.SharedSecret),
			LocalTrafficSelectors: pulumi.StringArray{
				pulumi.String("0.0.0.0/0"),
			},
			RemoteTrafficSelectors: pulumi.StringArray{
				pulumi.String("0.0.0.0/0"),
			},
		},
		DestRange: vpncon.DestRange,
	}
}

func (p *Provider) VPNTunnelFunc(ctx *pulumi.Context, key string) error {
	config := *p.Config[key]
	if config.VpnT == nil {
		// skip create VPN Tunnel
		return nil
	}

	thisCompute := *p.Compute[key]
	config.VpnT.Args.TargetVpnGateway = thisCompute.vpnGw.ID()

	// forwarding rules
	frs, err := defaultVPNforwardingRule(ctx, p, key)
	if err != nil {
		return err
	}
	// vpn tunnel
	vpnTunnel, err := compute.NewVPNTunnel(ctx, key+"-vpn-Tunnel", &config.VpnT.Args, pulumi.DependsOn(frs))
	if err != nil {
		return err
	}
	thisCompute.vpnTunnel = vpnTunnel
	// vpn Route
	_, err = compute.NewRoute(ctx, key+"vpn-Tunnel-Route", &compute.RouteArgs{
		Name:             pulumi.String(key + "-vpn-tunnel-route"),
		Network:          thisCompute.network.Name,
		DestRange:        pulumi.String(config.VpnT.DestRange),
		Priority:         pulumi.Int(1000),
		NextHopVpnTunnel: thisCompute.vpnGw.ID(),
	})
	if err != nil {
		return err
	}
	return nil
}

// 通常コマンド等で作成するとデフォルトでこの３つのforwardingRuleが作成される
func defaultVPNforwardingRule(ctx *pulumi.Context, p *Provider, key string) ([]pulumi.Resource, error) {
	thisCompute := *p.Compute[key]
	// クラウド側のvpn static ipをセッティング
	vpnStaticIp, err := compute.NewAddress(ctx, key+"vpngw"+"-StaticIP", &compute.AddressArgs{
		Name: pulumi.String(key + "vpn-gw-static-ip"),
	})
	if err != nil {
		return nil, err
	}
	// frESP,frUDP500,frUDP4500を設定
	var frs []pulumi.Resource
	frEsp, err := compute.NewForwardingRule(ctx, "frEsp", &compute.ForwardingRuleArgs{
		Name:       pulumi.String("fr-esp"),
		IpProtocol: pulumi.String("ESP"),
		IpAddress:  vpnStaticIp.Address,
		Target:     thisCompute.vpnGw.ID(),
	})
	if err != nil {
		return nil, err
	}
	frUdp500, err := compute.NewForwardingRule(ctx, "frUdp500", &compute.ForwardingRuleArgs{
		Name:       pulumi.String("fr-udp500"),
		IpProtocol: pulumi.String("UDP"),
		PortRange:  pulumi.String("500"),
		IpAddress:  vpnStaticIp.Address,
		Target:     thisCompute.vpnGw.ID(),
	})
	if err != nil {
		return nil, err
	}
	frUdp4500, err := compute.NewForwardingRule(ctx, "frUdp4500", &compute.ForwardingRuleArgs{
		Name:       pulumi.String("fr-udp4500"),
		IpProtocol: pulumi.String("UDP"),
		PortRange:  pulumi.String("4500"),
		IpAddress:  vpnStaticIp.Address,
		Target:     thisCompute.vpnGw.ID(),
	})
	if err != nil {
		return nil, err
	}
	frs = append(frs, frEsp, frUdp500, frUdp4500)
	return frs, nil
}
