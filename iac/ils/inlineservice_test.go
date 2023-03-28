package main

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/pulumi/pulumi-command/sdk/go/command/local"
	"github.com/pulumi/pulumi-github/sdk/v5/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/debug"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optdestroy"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optpreview"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optrefresh"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optremove"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	"github.com/pulumi/pulumi/sdk/v3/go/common/workspace"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/require"
)

func TestUpsertStackInlineSourceRefresh(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	pon := os.Getenv("PULUMI_ORG_NAME")
	orgName := pon
	projectName := "testproject"
	stackNameA := "testUpsertStackInlineSourceRefreshA"
	desc := "A inline source Go Pulumi program Test"
	ws, err := auto.NewLocalWorkspace(ctx, auto.Project(workspace.Project{
		Name:        tokens.PackageName(projectName),
		Runtime:     workspace.NewProjectRuntimeInfo("go", nil),
		Description: &desc,
	}))
	require.NoError(t, err)
	require.NotNil(t, ws)

	prj, err := ws.ProjectSettings(ctx)
	require.NoError(t, err)
	require.NotNil(t, prj)

	s, err := auto.UpsertStackInlineSource(ctx, stackNameA, prj.Name.String(), func(ctx *pulumi.Context) error {
		return nil
	})

	require.NoError(t, err)
	require.NotNil(t, s)

	defer func() {
		dr, err := s.Destroy(ctx, optdestroy.Message("Successfully destroyed stack :"+stackNameA))
		require.NoError(t, err)
		log.Println(dr.Summary.Kind + " " + dr.Summary.Message)
		err = s.Workspace().RemoveStack(ctx, s.Name(), optremove.Force())
		require.NoError(t, err)
	}()

	err = s.SetAllConfig(ctx, auto.ConfigMap{
		"bar:token": auto.ConfigValue{
			Value:  "def",
			Secret: true,
		},
		"buzz:owner": auto.ConfigValue{
			Value:  "xyz",
			Secret: true,
		},
	})
	require.NoError(t, err)

	rr, err := s.Refresh(ctx, optrefresh.Message("Refresh stack "+stackNameA))
	require.NoError(t, err)
	require.NotNil(t, rr)
	log.Println(rr.Summary.Kind + " " + rr.Summary.Message)

	values, err := s.GetAllConfig(ctx)
	require.NoError(t, err)

	require.Equal(t, "def", values["bar:token"].Value)
	require.Equal(t, "xyz", values["buzz:owner"].Value)

	s.Workspace().SetProgram(func(pCtx *pulumi.Context) error {

		hello, err := local.NewCommand(pCtx, "hello", &local.CommandArgs{
			Create: pulumi.String("echo \"Hello Pulumi\""),
		})
		if err != nil {
			return err
		}

		pCtx.Export("hello", hello.Stdout)

		return nil
	})

	prev, err := s.Preview(ctx, optpreview.DebugLogging(debug.LoggingOptions{
		Debug: true,
	}))
	require.NoError(t, err)
	log.Println(prev.StdOut)

	up, err := s.Up(ctx, optup.DebugLogging(debug.LoggingOptions{
		Debug: true,
	}))
	require.NoError(t, err)
	log.Println(up.StdOut)

	stackNameB := "testUpsertStackInlineSourceRefreshB"

	require.NoError(t, err)
	require.NotNil(t, ws)

	prj, err = ws.ProjectSettings(ctx)
	require.NoError(t, err)
	require.NotNil(t, prj)

	ss, err := auto.UpsertStackInlineSource(ctx, stackNameB, prj.Name.String(), func(ctx *pulumi.Context) error {
		return nil
	})
	require.NoError(t, err)
	require.NotNil(t, ss)

	defer func() {
		dr, err := ss.Destroy(ctx, optdestroy.Message("Successfully destroyed stack :"+stackNameB))
		require.NoError(t, err)
		log.Println(dr.Summary.Kind + " " + dr.Summary.Message)
		err = ss.Workspace().RemoveStack(ctx, ss.Name(), optremove.Force())
		require.NoError(t, err)
	}()

	ss.Workspace().SetProgram(func(pCtx *pulumi.Context) error {

		qualifiedStackName := auto.FullyQualifiedStackName(orgName, prj.Name.String(), s.Name())
		require.NotNil(t, qualifiedStackName)
		require.Equal(t, orgName+"/"+projectName+"/"+s.Name(), qualifiedStackName)

		stackReff, err := pulumi.NewStackReference(pCtx, qualifiedStackName, nil)
		require.NoError(t, err)
		require.NotNil(t, stackReff)

		outputValue := stackReff.GetOutput(pulumi.String("hello"))
		require.NotNil(t, outputValue)
		log.Println(outputValue)

		return nil

	})

	prev, err = ss.Preview(ctx, optpreview.DebugLogging(debug.LoggingOptions{
		Debug: true,
	}))
	require.NoError(t, err)
	log.Println(prev.StdOut)

	up, err = ss.Up(ctx, optup.DebugLogging(debug.LoggingOptions{
		Debug: true,
	}))
	require.NoError(t, err)
	log.Println(up.StdOut)
}

func TestNewStackInlineSourceActionsSecret(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	projectName := "testproject"
	stackNameA := "testInlineSourceActionsSecretA"
	desc := "A inline source Go Pulumi program Test"
	ws, err := auto.NewLocalWorkspace(ctx, auto.Project(workspace.Project{
		Name:        tokens.PackageName(projectName),
		Runtime:     workspace.NewProjectRuntimeInfo("go", nil),
		Description: &desc,
		Config: map[string]workspace.ProjectConfigType{
			"bar:token": {
				Value: "abc",
			},
		},
	}))
	require.NoError(t, err)
	require.NotNil(t, ws)

	prj, err := ws.ProjectSettings(ctx)
	require.NoError(t, err)
	require.NotNil(t, prj)

	s, err := auto.NewStackInlineSource(ctx, stackNameA, prj.Name.String(), func(ctx *pulumi.Context) error {
		return nil
	})
	if err != nil && auto.IsCreateStack409Error(err) {
		log.Println("stack " + stackNameA + " already exists")
		s, err = auto.UpsertStackInlineSource(ctx, stackNameA, prj.Name.String(), func(ctx *pulumi.Context) error {
			return nil
		})
	}
	require.NoError(t, err)
	require.NotNil(t, s)

	stackNameB := "testInlineSourceActionsSecretB"
	desc = "A inline source Go Pulumi program Test"

	prj, err = ws.ProjectSettings(ctx)
	require.NoError(t, err)
	require.NotNil(t, prj)

	ss, err := auto.NewStackInlineSource(ctx, stackNameB, prj.Name.String(), func(ctx *pulumi.Context) error {
		return nil
	})
	if err != nil && auto.IsCreateStack409Error(err) {
		log.Println("stack " + stackNameB + " already exists")
		ss, err = auto.UpsertStackInlineSource(ctx, stackNameB, prj.Name.String(), func(ctx *pulumi.Context) error {
			return nil
		})
	}
	require.NoError(t, err)
	require.NotNil(t, ss)

	defer func() {
		dr, err := s.Destroy(ctx, optdestroy.Message("Successfully destroyed stack :"+stackNameA))
		require.NoError(t, err)
		log.Println(dr.Summary.Kind + " " + dr.Summary.Message)
		err = s.Workspace().RemoveStack(ctx, s.Name(), optremove.Force())
		require.NoError(t, err)
	}()

	defer func() {
		dr, err := ss.Destroy(ctx, optdestroy.Message("Successfully destroyed stack :"+stackNameB))
		require.NoError(t, err)
		log.Println(dr.Summary.Kind + " " + dr.Summary.Message)
		err = ss.Workspace().RemoveStack(ctx, ss.Name(), optremove.Force())
		require.NoError(t, err)
	}()

	prj, err = s.Workspace().ProjectSettings(ctx)
	require.NoError(t, err)
	require.NotNil(t, prj)
	log.Println("project name: " + prj.Name.String())
	log.Println("project runtime: " + prj.Runtime.Name())

	values, err := s.GetAllConfig(ctx)
	require.NoError(t, err)

	for _, s := range values {
		log.Println("config: " + s.Value)
	}

	ght := os.Getenv("GITHUB_TOKEN")
	gho := os.Getenv("GITHUB_OWNER")
	err = s.SetAllConfig(ctx, auto.ConfigMap{
		"github:token": auto.ConfigValue{
			Value:  ght,
			Secret: true,
		},
		"github:owner": auto.ConfigValue{
			Value:  gho,
			Secret: true,
		},
	})
	require.NoError(t, err)

	values, err = s.GetAllConfig(ctx)
	require.NoError(t, err)

	require.Equal(t, ght, values["github:token"].Value)
	require.Equal(t, gho, values["github:owner"].Value)

	s.Workspace().SetProgram(func(pCtx *pulumi.Context) error {

		repositoryName := "pulumi-dagger-gh-inline-source-actions-secret"
		_, err := github.NewRepository(pCtx, "newRepository", &github.RepositoryArgs{
			DeleteBranchOnMerge: pulumi.Bool(true),
			Description:         pulumi.String("This is a test repository for Pulumi repository creation with Dagger CI/CD"),
			HasIssues:           pulumi.Bool(true),
			HasProjects:         pulumi.Bool(true),
			Name:                pulumi.String(repositoryName),
			Topics:              pulumi.StringArray{pulumi.String("pulumi"), pulumi.String("dagger"), pulumi.String("github"), pulumi.String("test")},
			Visibility:          pulumi.String("public"),
		})
		require.NoError(t, err)

		pCtx.Export("repositoryName", pulumi.String(repositoryName))

		return nil
	})
	require.NoError(t, err)

	prev, err := s.Preview(ctx, optpreview.DebugLogging(debug.LoggingOptions{
		Debug: true,
	}))
	require.NoError(t, err)
	log.Println(prev.StdOut)

	up, err := s.Up(ctx, optup.DebugLogging(debug.LoggingOptions{
		Debug: true,
	}))
	require.NoError(t, err)
	log.Println(up.StdOut)

	err = ss.SetAllConfig(ctx, auto.ConfigMap{
		"github:token": auto.ConfigValue{
			Value:  ght,
			Secret: true,
		},
		"github:owner": auto.ConfigValue{
			Value:  gho,
			Secret: true,
		},
	})
	require.NoError(t, err)

	values, err = ss.GetAllConfig(ctx)
	require.NoError(t, err)

	require.Equal(t, ght, values["github:token"].Value)
	require.Equal(t, gho, values["github:owner"].Value)

	ss.Workspace().SetProgram(func(pCtx *pulumi.Context) error {

		_, err = github.GetActionsPublicKey(pCtx, &github.GetActionsPublicKeyArgs{
			Repository: "pulumi-dagger-gh-inline-source-actions-secret",
		}, nil)
		require.NoError(t, err)

		_, err = github.NewActionsSecret(pCtx, "newActionsSecret", &github.ActionsSecretArgs{
			Repository: pulumi.String("pulumi-dagger-gh-inline-source-actions-secret"),
			SecretName: pulumi.String("TOKEN"),
		})
		require.NoError(t, err)

		return nil
	})

	prev, err = ss.Preview(ctx, optpreview.DebugLogging(debug.LoggingOptions{
		Debug: true,
	}))
	require.NoError(t, err)
	log.Println(prev.StdOut)

	up, err = ss.Up(ctx, optup.DebugLogging(debug.LoggingOptions{
		Debug: true,
	}))
	require.NoError(t, err)
	log.Println(up.StdOut)
}
