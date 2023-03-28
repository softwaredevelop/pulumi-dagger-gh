package lint

import "dagger.io/dagger"

func Yl(c *dagger.Client, id dagger.ContainerID) *dagger.Container {
	return yl(c, id)
}

func yl(c *dagger.Client, id dagger.ContainerID) *dagger.Container {
	install := c.
		Container().
		From("python:alpine").
		WithExec([]string{"pip", "install", "yamllint"})

	return c.Container(dagger.ContainerOpts{ID: id}).
		WithFile("/usr/bin", install.File("/usr/local/bin/yamllint"))
}
