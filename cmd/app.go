package cmd

import "fmt"

// Version is the program version, set during compile time using:
//   -ldflags -X github.com/stevenxie/nightwatch/main.Version=$(VERSION)
var Version = "unset"

// Exec is the entrypoint to command nightwatch.
func Exec() {
	loadEnv() // load .env variables

	fmt.Printf("RGV: %s\n", Version)
}
