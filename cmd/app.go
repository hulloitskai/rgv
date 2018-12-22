package cmd

import (
	"fmt"

	"github.com/stevenxie/rgv/cmd/info"
)

// Exec is the entrypoint to command rgv.
func Exec() {
	loadEnv() // load .env variables

	fmt.Printf("RGV: %s\n", info.Version)
}
