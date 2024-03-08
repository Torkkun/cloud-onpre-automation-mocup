package server

// graphqlに移行する

type ClientConfig struct {
	CloudProvider []string //GCPやAzureなど
	NetworkC      []NetworkConfig
	Peering       []NetworkPeeringConfig
}

type NetworkConfig struct {
	ID                   string   `json:"nw_id"` // peerのコンフィグと紐づける用
	AutoCreateSubnetwork bool     `json:"auto_create_subnw"`
	AddressRange         []string `json:"ip_cidr"`
	Subnet               []SubnetworkConfig
	AutoCreateVPN        bool `json:"auto_create_vpn"`
	VPN                  VPNConfig
	CreateVm             bool
	Vm                   VmConfig
}

type NetworkPeeringConfig struct {
	Peer1_nw_id string `json:"peer1_nw_id"`
	Peer2_nw_id string `json:"peer2_nw_id"`
}

type SubnetworkConfig struct {
	AddressRange string
	Ip_cidr      []string
}

type VPNConfig struct {
	PeerIp       string
	SharedSecret string
	DestRange    string
}

type VmConfig struct {
	OS     string
	Amount int
}
