//revive:disable:package-comments,exported
package main

import (
	"github.com/pulumi/pulumi-github/sdk/v5/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {

	pulumi.Run(func(ctx *pulumi.Context) error {

		repositoryName := "pulumiloc-dagger-gh"
		repository, err := github.NewRepository(ctx, "newRepository", &github.RepositoryArgs{
			DeleteBranchOnMerge: pulumi.Bool(true),
			Description:         pulumi.String("This is a test repository for Pulumi repository creation with Dagger CI/CD"),
			HasIssues:           pulumi.Bool(true),
			HasProjects:         pulumi.Bool(true),
			Name:                pulumi.String(repositoryName),
			Topics:              pulumi.StringArray{pulumi.String("pulumi"), pulumi.String("dagger"), pulumi.String("github"), pulumi.String("test")},
			Visibility:          pulumi.String("public"),
		})
		if err != nil {
			return err
		}

		_, err = github.NewBranchProtection(ctx, "branchProtection", &github.BranchProtectionArgs{
			RepositoryId:          repository.NodeId,
			Pattern:               pulumi.String("main"),
			RequiredLinearHistory: pulumi.Bool(true),
		})
		if err != nil {
			return err
		}

		_, err = github.NewIssueLabel(ctx, "newIssueLabelGhActions", &github.IssueLabelArgs{
			Color:       pulumi.String("E66E01"),
			Description: pulumi.String("This issue is related to github-actions dependencies"),
			Name:        pulumi.String("github-actions dependencies"),
			Repository:  repository.Name,
		})
		if err != nil {
			return err
		}

		ctx.Export("repository", repository.Name)
		ctx.Export("repositoryUrl", repository.HtmlUrl)

		return nil

	})

}
