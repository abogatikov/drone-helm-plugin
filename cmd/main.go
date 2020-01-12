package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"

	"github.com/abogatikov/drone-helm-plugin/pkg"
)

func main() {
	if _, err := fmt.Fprintln(os.Stderr, "Drone Helm Plugin"); err != nil {
		panic(err)
	}

	if _, err := fmt.Fprintf(os.Stderr, "Release: %s \n", pkg.Release); err != nil {
		panic(err)
	}

	if _, err := fmt.Fprintf(os.Stderr, "Compile Time: %s \n", pkg.CompileTime); err != nil {
		panic(err)
	}

	if _, err := fmt.Fprintf(os.Stderr, "Commit: %s \n", pkg.Commit); err != nil {
		panic(err)
	}

	var opts pkg.Config
	if _, err := flags.Parse(&opts); err != nil {
		return
	}

	if err := opts.Exec(); err != nil {
		_, err1 := fmt.Fprintln(os.Stderr, err.Error())
		if err1 != nil {
			panic(err1)
		}
	}
}
