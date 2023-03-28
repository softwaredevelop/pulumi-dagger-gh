package util

import (
	"context"

	"dagger.io/dagger"
)

func ConnectToDagger() (*dagger.Client, error) {
	return connectToDagger()
}

func connectToDagger() (*dagger.Client, error) {
	ctx := context.Background()
	c, err := dagger.Connect(ctx)
	if err != nil {
		return nil, err
	}
	return c, nil
}
