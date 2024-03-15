package main

import (
	"context"
	"log"
	"os"

	"dagger.io/dagger"
)

func main() {
	verbose := true
	var c *dagger.Client
	var err error
	ctx := context.Background()
	if verbose {
		c, err = dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout), dagger.WithLogOutput(os.Stderr))
	} else {
		c, err = dagger.Connect(ctx)
	}
	if err != nil {
		log.Fatal(err)
	}

	src := c.Host().Directory(".")
	//src := pipeline.Host().Directory(".")
	ctr := c.Container().From("amazoncorretto:8").WithExec([]string{"amazon-linux-extras", "install", "redis6"})

	service := c.Container().From("docker.io/stobias123/redis-cluster:latest").
		WithExposedPort(30000).
		WithExposedPort(30001).
		WithExposedPort(30002).
		WithExposedPort(30003).
		WithExposedPort(30004).
		WithExposedPort(30005).
		AsService()

	ctr.WithMountedDirectory("/app", src).
		WithWorkdir("/app").
		WithServiceBinding("redis-cluster", service).
		//WithEntrypoint([]string{"/bin/bash"}).
		WithExec([]string{"redis-cli", "-h", "redis-cluster", "-p", "30000", "ping"}).
		WithExec([]string{"./gradlew", ":test"}).Sync(ctx)
}
