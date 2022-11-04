package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"dagger.io/dagger"
)

func main() {
	if len(os.Args) < 2 {
		log.Println("Must pass the workdir")
		os.Exit(1)
	}

	if err := build(context.Background(), os.Args[1]); err != nil {
		fmt.Println(err)
	}
}

func build(ctx context.Context, workdir string) error {
	fmt.Printf("Building example project")

	// Initialize Dagger client
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout), dagger.WithWorkdir(workdir))
	if err != nil {
		return err
	}

	defer client.Close()

	gc := client.Container().
		From("golang:latest").
		WithMountedDirectory("/host", client.Host().Workdir()).
		WithWorkdir("/host")

	// Define build command
	path := "build/"
	gc = gc.Exec(dagger.ContainerExecOpts{
		Args: []string{"go", "build", "-o", path},
	})

	// Get reference to build output dir
	output := gc.Directory(path)

	// Write contents of container build/ directory to host
	if _, err = output.Export(ctx, path); err != nil {
		return err
	}

	return nil
}
