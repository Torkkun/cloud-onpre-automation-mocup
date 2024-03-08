package infra

import (
	"fmt"
	"pulumigcp/server"
	"testing"
)

func TestNewCloud(t *testing.T) {
	nwcon := server.NetworkConfig{
		AutoCreateSubnetwork: false,
		Subnet: []server.SubnetworkConfig{
			{
				AddressRange: "172.24.0.0/16",
				Ip_cidr:      []string{"172.24.1.0/24"},
			},
		},
	}
	clc := server.ClientConfig{
		NetworkC: []server.NetworkConfig{nwcon},
	}
	_, keys, err := NewCloud(&CloudInformation{
		Name:         "GCP",
		ClientConfig: &clc,
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(keys)
}
