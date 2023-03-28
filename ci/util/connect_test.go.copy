package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConnectToDaggerPipeline(t *testing.T) {
	t.Parallel()
	c, err := ConnectToDagger()
	require.NoError(t, err)
	require.NotNil(t, c)
	defer c.Close()

	p1 := c.Pipeline("p1")
	require.NotNil(t, p1)
	p2 := c.Pipeline("p2")
	require.NotNil(t, p2)

	c1 := p1.Container().
		From("alpine")
	require.NotNil(t, c1)

	c2 := p2.Container().
		From("busybox")
	require.NotNil(t, c2)
}

func TestConnectToDagger(t *testing.T) {
	c, err := ConnectToDagger()
	require.NoError(t, err)
	require.NotNil(t, c)
	defer c.Close()
}
