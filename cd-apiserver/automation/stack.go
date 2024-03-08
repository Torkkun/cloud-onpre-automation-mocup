package automation

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// stack作成
// 1Networkで一つのStackを作成
func CreateOrSelectStack(ctx context.Context, projectName, stackName string, deployFunc pulumi.RunFunc, providerName string) auto.Stack {
	var pName, version string
	var cfg auto.ConfigMap
	switch providerName {
	case "GCP":
		pName = gcp
		version = gcp_sdk_ver
		// stack config setting
		cfg = auto.ConfigMap{
			tag_project: auto.ConfigValue{Value: gcp_project},
			tag_zone:    auto.ConfigValue{Value: gcp_zone},
		}
		// deploy function毎のコンフィグも設定するようにしたい
	}

	s, err := auto.UpsertStackInlineSource(ctx, stackName, projectName, deployFunc)
	if err != nil {
		log.Printf("Failed to create stack: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created stack %q\n", stackName)

	w := s.Workspace()
	fmt.Printf("Installing the %s plugin\n", pName)
	err = w.InstallPlugin(ctx, pName, version)
	if err != nil {
		fmt.Printf("Failed to install program plugins %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Successfully installed %s plugin\n", pName)

	// stackのコンフィグ設定
	err = s.SetAllConfig(ctx, cfg)
	if err != nil {
		fmt.Printf("Failed to set Config %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully setting %s config\n", stackName)

	_, err = s.Refresh(ctx)
	if err != nil {
		fmt.Printf("Failed to refresh stack: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Refresh succeeded!\n")

	return s

}
