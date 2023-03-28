package test

import (
	"context"
	"log"
	"os"
	"testing"

	"dagger.io/dagger"
	"github.com/stretchr/testify/require"
)

func TestGoTest(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	require.NoError(t, err)
	defer c.Close()

	version := "1.20.2"
	id, err := c.Container().
		From("busybox:glibc").
		WithExec([]string{"wget", "https://go.dev/dl/go" + version + ".linux-amd64.tar.gz"}).
		WithExec([]string{"mkdir", "-p", "/usr/local/go/bin"}).
		WithExec([]string{"tar", "-xzf", "go" + version + ".linux-amd64.tar.gz", "--strip-components=2", "go/bin/go", "-C", "/usr/local/go/bin/"}).
		ID(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	out, err := c.
		Container(dagger.ContainerOpts{ID: id}).
		WithExec([]string{"/usr/local/go/bin/go", "version"}).
		Stdout(ctx)
	require.NoError(t, err)
	log.Println(out)

	container := c.Container().From("busybox:glibc")
	require.NotNil(t, container)

	id, err = container.ID(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	container = GoTest(c, id, version)
	require.NotNil(t, container)
	out, err = container.
		WithExec([]string{"/usr/local/go/bin/go", "version"}).
		Stdout(ctx)
	require.NoError(t, err)
	require.Contains(t, out, "go version go"+version+" linux/amd64")
	log.Println(out)
}
