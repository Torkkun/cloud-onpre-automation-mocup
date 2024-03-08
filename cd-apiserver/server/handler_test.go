package server_test

import (
	"context"
	"fmt"
	"os"
	"pulumigcp/automation"
	"pulumigcp/infra"
	"pulumigcp/server"
	"strconv"
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
)

/* func TestCreateNetwork(t *testing.T) {
	ctx := context.Background()
	stackName := "test_stack_minimum"
	projectName := "testNetworkproject"

	// NetworkTest()
	testNwStack := automation.CreateOrSelectStack(ctx, projectName, stackName, NetworkTest, "GCP")

	stdoutStreamer := optup.ProgressStreams(os.Stdout)

	// run update to deploy
	_, err := testNwStack.Up(ctx, stdoutStreamer)
	if err != nil {
		fmt.Printf("Failed to update stack: %v\n\n", err)
		os.Exit(1)
	}
	fmt.Println("TestNetwork stack update succeeded!")
}

func NetworkTest(ctx *pulumi.Context) error {
	for i := 0; i < 1; i++ {
		_, err := compute.NewNetwork(ctx, fmt.Sprintf("Just test network-%d", i), &compute.NetworkArgs{
			Name:                  pulumi.String(fmt.Sprintf("minitestnetwork-%d", i)),
			AutoCreateSubnetworks: pulumi.Bool(false),
		})
		if err != nil {
			return err
		}
	}
	return nil
} */

func TestNetworkStackCreate(t *testing.T) {
	var err error
	ctx := context.Background()
	stackName := "test_stack"
	projectName := "testNetworkproject"

	var nwconfigs []server.NetworkConfig
	var peercon []server.NetworkPeeringConfig

	for i := 0; i < 2; i++ {
		var subnetclc []string
		var isVpnGw bool
		var vpnclc server.VPNConfig
		if i == 1 {
			subnetclc = []string{"172.26.0.0/16"}
			isVpnGw = false // trueでないとVPN関連のものが作成されない
			vpnclc = server.VPNConfig{
				PeerIp:       "203.178.146.10",
				SharedSecret: "ZnVq4v8J",
				DestRange:    "172.26.0.0/24",
			}
		} else {
			subnetclc = []string{fmt.Sprintf("192.168.%d.0/24", i)}
		}
		nwcon := &server.NetworkConfig{
			ID:                   "testid" + strconv.Itoa(i),
			AutoCreateSubnetwork: false,
			AddressRange:         subnetclc,
			AutoCreateVPN:        isVpnGw,
			VPN:                  vpnclc,
			CreateVm:             false,
			Vm: server.VmConfig{
				OS:     "ubuntu",
				Amount: 1,
			},
		}
		nwconfigs = append(nwconfigs, *nwcon)

	}
	peercon = append(peercon, server.NetworkPeeringConfig{
		Peer1_nw_id: "testid0",
		Peer2_nw_id: "testid1",
	})

	clientConfig := &server.ClientConfig{
		NetworkC: nwconfigs,
		Peering:  peercon,
	}

	provider := &infra.CloudInformation{
		Name:         "GCP",
		ClientConfig: clientConfig,
	}

	// provider別のCloudに送るデータを設定
	// Configの設定、Functionの設定
	cloudgcp, keys, err := infra.NewCloud(provider)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create New Network Stack Test
	fmt.Println("Starting TestNetwork stack update")

	stackFunc := infra.CreateNetworksForStacking(cloudgcp, keys)

	fmt.Println("Create Network stack")
	testNwStack := automation.CreateOrSelectStack(ctx, projectName, stackName, stackFunc, provider.Name)
	fmt.Println("Create or Select stack success!")
	// wire up our update to stream progress to stdout
	stdoutStreamer := optup.ProgressStreams(os.Stdout)
	// run update to deploy
	_, err = testNwStack.Up(ctx, stdoutStreamer)
	if err != nil {
		fmt.Printf("Failed to update testNwStack.Up: %v\n\n", err)
		os.Exit(1)
	}
	fmt.Println("TestNetwork stack update succeeded!")
}
