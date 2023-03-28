//revive:disable:package-comments,exported
package test

import "dagger.io/dagger"

func GoTest(c *dagger.Client, id dagger.ContainerID, version string) *dagger.Container {
	return goTest(c, id, version)
}

func goTest(c *dagger.Client, id dagger.ContainerID, version string) *dagger.Container {
	download := c.
		Container().
		From("busybox:glibc").
		WithExec([]string{"wget", "https://go.dev/dl/go" + version + ".linux-amd64.tar.gz"}).
		WithExec([]string{"mkdir", "-p", "/usr/local/go/bin"}).
		WithExec([]string{"tar", "-xzf", "go1.20.2.linux-amd64.tar.gz", "--strip-components=2", "go/bin/go", "-C", "/usr/local/go/bin/"})

	return c.Container(dagger.ContainerOpts{ID: id}).
		WithExec([]string{"mkdir", "-p", "/usr/local/go/bin"}).
		WithFile("/usr/local/go/bin", download.File("/usr/local/go/bin/go"))
}
