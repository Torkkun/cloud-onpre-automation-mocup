package gcp

import (
	"fmt"
	"pulumigcp/server"
	"testing"
)

func TestConfig(t *testing.T) {
	clc := &server.ClientConfig{
		NetworkC: []server.NetworkConfig{
			{
				ID:                   "test-peer1",
				AutoCreateSubnetwork: false,
			},
			{
				ID: "test-peer2",
			},
			{ID: "test-peer3"},
		},
		Peering: []server.NetworkPeeringConfig{
			{
				Peer1_nw_id: "test-peer1",
				Peer2_nw_id: "test-peer2",
			},
		},
	}
	_, _, keypeerIDs, err := newNetwork(clc.NetworkC)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//compute := newCompute(keys)
	peering, err := newNetworkPeering(clc.Peering, keypeerIDs)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("IDList:%v\npeering:%v", keypeerIDs, peering)
}

func TestNewConfig(t *testing.T) {
	_, keys, err := NewConfig(&server.ClientConfig{
		NetworkC: []server.NetworkConfig{
			{
				ID:                   "test-peer1",
				AutoCreateSubnetwork: false,
			},
			//{},
		},
		Peering: []server.NetworkPeeringConfig{
			{
				Peer1_nw_id: "test-peer1",
				Peer2_nw_id: "test-peer2",
			},
		},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(keys)
}

func TestNewNwConfig(t *testing.T) {
	tests := []struct {
		name     string
		nwConfig []server.NetworkConfig
	}{
		{
			name: "Normal form test-1",
			nwConfig: []server.NetworkConfig{
				{
					ID:                   "test-peer1",
					AutoCreateSubnetwork: false,
					AddressRange:         []string{"172.24.1.0/24", "172.23.23.23/32"},
				},
			},
		},
		{
			name:     "Abnormal form test-1",
			nwConfig: []server.NetworkConfig{},
		},
	}

	for _, tt := range tests {
		//tt := tt
		t.Run(tt.name, func(t *testing.T) {
			//t.Parallel()
			configmap, keys, id_nwkeymap, err := newNetwork(tt.nwConfig)
			if err != nil {
				fmt.Println("Error:" + err.Error())
				return
			}
			for _, key := range keys {
				config := configmap[key]
				fmt.Printf("NetworkArgs_Name:%s\n", config.Nw.Args.Name)
				fmt.Printf("AutoCreateSubnet:%v\n", config.Nw.Args.AutoCreateSubnetworks)
				for _, v := range config.Subnets {
					fmt.Printf("SubnetArgs_IPCidrRange%v\n", v.Args.IpCidrRange)
				}
			}
			fmt.Printf("ID_keymap:%v\n", id_nwkeymap)
		})
	}
}

func TestNewNwPeeringConfig(t *testing.T) {
	tests := []struct {
		name        string
		peerConfig  []server.NetworkPeeringConfig
		id_nwkeymap map[string]string
	}{
		{
			name: "Normal form test-1",
			peerConfig: []server.NetworkPeeringConfig{
				{Peer1_nw_id: "test-peer1", Peer2_nw_id: "test-peer2"},
			},
			id_nwkeymap: map[string]string{"test-peer1": "network-1", "test-peer2": "network-2"},
		},
		{
			name:        "Normal form test-2",
			peerConfig:  []server.NetworkPeeringConfig{},
			id_nwkeymap: map[string]string{},
		},
		{
			name: "Abnormal form test-1",
			peerConfig: []server.NetworkPeeringConfig{
				{Peer1_nw_id: "test-peer1", Peer2_nw_id: "test-peer2"},
			},
			id_nwkeymap: map[string]string{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			peermaping, err := newNetworkPeering(tt.peerConfig, tt.id_nwkeymap)
			if err != nil {
				fmt.Println("Error:" + err.Error())
				return
			}
			fmt.Printf("test-name:%s\nOutput:%v\n", tt.name, peermaping)
		})
	}
}
