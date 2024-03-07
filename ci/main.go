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
	ctr := c.Container().From("amazoncorretto:8")
	ctr = BindDockerCompose(c, ctr, "docker-compose.ci.yml")
	ctr.WithMountedDirectory("/app", src).
        WithExec([]string{"amazon-linux-extras", "install", "redis6"}).
		WithWorkdir("/app").
		//WithEntrypoint([]string{"/bin/bash"}).
		WithExec([]string{"./gradlew", ":test", "--tests", "ActiveRideCacheTest"}).Sync(ctx)

}
