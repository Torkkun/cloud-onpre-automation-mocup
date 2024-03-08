package gcp

import (
	"fmt"
	"pulumigcp/server"

	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/compute"
)

type Provider struct {
	Config            map[string]*Config     // need user or automation setting
	Compute           map[string]*GcpCompute // gcp original data
	NwPeering_keyList map[string][]string
}

// GCP用のコンフィグ諸々
// 1NW毎にコンフィグ毎
type Config struct {
	Nw      *VpcNetworkConfig
	Subnets []*VpcSubnetworkConfig
	VpnGw   *VpnGwConfig
	VpnT    *VpnTunnelConfig
	Vm      *VmInstanceConfig
	Peering *NwPeeringConfig
}

// GCPのComputeデータ
type GcpCompute struct {
	network *compute.Network
	subnet  []*compute.Subnetwork
	//peering   *compute.NetworkPeering
	vpnGw     *compute.VPNGateway
	vpnTunnel *compute.VPNTunnel
	//instance  *compute.Instance
}

func NewConfig(clc *server.ClientConfig) (*Provider, []string, error) {
	// config設定をする
	config, keys, keypeerIDs, err := newNetwork(clc.NetworkC)
	if err != nil {
		return nil, nil, err
	}
	compute := newCompute(keys)
	peering, err := newNetworkPeering(clc.Peering, keypeerIDs)
	if err != nil {
		return nil, nil, err
	}
	return &Provider{Config: config, Compute: compute, NwPeering_keyList: peering}, keys, nil
}

func newNetwork(clcNw []server.NetworkConfig) (map[string]*Config, []string, map[string]string, error) {
	if len(clcNw) == 0 {
		return nil, nil, nil, fmt.Errorf("network settings do not exist")
	}

	configmap := map[string]*Config{}
	var nwkeyL []string
	// clientConfigで与えられているIDと割り振ったIDの紐付け
	peering_nwkeyM := map[string]string{}

	for i, v := range clcNw {
		key_name := fmt.Sprintf("network-%d", i)
		// peeringがあれば追加
		if v.ID != "" {
			peering_nwkeyM[v.ID] = key_name
		}
		// create network
		networkcon := NewNetworkConfig(key_name, v.AutoCreateSubnetwork)

		// create subnetwork configs
		subnetcon, err := NewSubnetworkConfig(key_name, v.AddressRange, v.Subnet)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("NewsubnetConfig Error:%v", err)
		}

		// create vpn config
		var vpngwcon *VpnGwConfig
		var vpntcon *VpnTunnelConfig
		if v.AutoCreateVPN {
			// 1to1 network and vpn
			vpngwcon = NewVPNGwConfig(key_name)
			vpntcon = NewVPNTunnelConfig(key_name, v.VPN)
		}

		var vmcon *VmInstanceConfig
		if v.CreateVm {
			vmcon = NewVMInstanceConfig(key_name, v.Vm)
		}

		// Note:現在必要なし
		//peering := NewNetworkPeeringConfig()

		// VPC network setting
		config := &Config{
			Nw:      networkcon,
			Subnets: subnetcon,
			VpnGw:   vpngwcon,
			VpnT:    vpntcon,
			Vm:      vmcon,
			//Peering: peering,
		}
		configmap[key_name] = config

		nwkeyL = append(nwkeyL, key_name)
	}

	return configmap, nwkeyL, peering_nwkeyM, nil
}

// もっといい方法があるか？
func newCompute(keys []string) map[string]*GcpCompute {
	computemap := make(map[string]*GcpCompute)
	for _, v := range keys {
		computemap[v] = new(GcpCompute)
	}

	return computemap
}

func newNetworkPeering(clcPeer []server.NetworkPeeringConfig, Id_keymap map[string]string) (map[string][]string, error) {
	if len(clcPeer) == 0 {
		return nil, nil
	}
	if len(Id_keymap) == 0 {
		return nil, fmt.Errorf("failed to set peering. Need an ID for the appropriate network")
	}

	peerListmap := map[string][]string{}
	for _, keypeer := range clcPeer {
		nw_key1 := Id_keymap[keypeer.Peer1_nw_id]
		nw_key2 := Id_keymap[keypeer.Peer2_nw_id]
		if nw_key2 == "" {
			return nil, fmt.Errorf("insufficient peering settings")
		}
		peerListmap[nw_key1] = append(peerListmap[nw_key1], nw_key2)
		peerListmap[nw_key2] = append(peerListmap[nw_key2], nw_key1)
	}

	return peerListmap, nil
}
