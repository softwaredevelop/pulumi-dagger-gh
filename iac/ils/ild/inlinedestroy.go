//revive:disable:package-comments,exported
package main

import (
	"context"
	"log"
	"os"

	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optdestroy"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	"github.com/pulumi/pulumi/sdk/v3/go/common/workspace"
)

func main() {

	ctx := context.Background()

	pon := os.Getenv("PULUMI_ORG_NAME")
	orgName := pon
	projectName := "iac"
	stackNameA := "ilsA"
	stackNameB := "ilsB"
	desc := "A inline source Go Pulumi program Test"
	ws, err := auto.NewLocalWorkspace(ctx, auto.Project(workspace.Project{
		Name:        tokens.PackageName(projectName),
		Runtime:     workspace.NewProjectRuntimeInfo("go", nil),
		Description: &desc,
	}))
	if err != nil {
		panic(err)
	}

	ws.SetEnvVars(map[string]string{
		"PULUMI_SKIP_UPDATE_CHECK": "true",
		"PULUMI_CONFIG_PASSPHRASE": "",
		"PULUMI_ACCESS_TOKEN":      os.Getenv("PULUMI_ACCESS_TOKEN"),
	})

	prj, err := ws.ProjectSettings(ctx)
	if err != nil {
		panic(err)
	}

	qualifiedStackNameB := auto.FullyQualifiedStackName(orgName, prj.Name.String(), stackNameB)
	sB, err := auto.SelectStack(ctx, qualifiedStackNameB, ws)
	if err != nil {
		panic(err)
	}
	dr, err := sB.Destroy(ctx, optdestroy.Message("Successfully destroyed stack :"+qualifiedStackNameB))
	if err != nil {
		panic(err)
	}
	log.Println(dr.Summary.Kind + " " + dr.Summary.Message)

	qualifiedStackNameA := auto.FullyQualifiedStackName(orgName, prj.Name.String(), stackNameA)
	sA, err := auto.SelectStack(ctx, qualifiedStackNameA, ws)
	if err != nil {
		panic(err)
	}
	dr, err = sA.Destroy(ctx, optdestroy.Message("Successfully destroyed stack :"+qualifiedStackNameA))
	if err != nil {
		panic(err)
	}
	log.Println(dr.Summary.Kind + " " + dr.Summary.Message)

}
