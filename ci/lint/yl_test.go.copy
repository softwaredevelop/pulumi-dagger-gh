package lint

import (
	"context"
	"testing"

	"dagger.io/dagger"
	"github.com/stretchr/testify/require"
)

func TestYl(t *testing.T) {
	t.Parallel()

	// c, err := util.ConnectToDagger()
	ctx := context.Background()

	c, err := dagger.Connect(ctx)
	require.NoError(t, err)
	require.NotNil(t, c)
	defer c.Close()

	c = c.Pipeline("yamllint")
	require.NotNil(t, c)

	container := c.Container().From("busybox:uclibc")
	require.NotNil(t, container)

	id, err := container.ID(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	container = Yl(c, id)
	require.NotNil(t, container)

	out, err := container.
		WithExec([]string{"ls", "/usr/bin/yamllint"}).
		Stdout(ctx)
	require.NoError(t, err)
	require.Equal(t, "/usr/bin/yamllint\n", out)

	// out, err = container.
	// 	WithExec([]string{"yamllint"}).
	// 	Stdout(ctx)
	// require.NoError(t, err)
	// require.Equal(t, "", out)
	// fmt.Println(err)

}
