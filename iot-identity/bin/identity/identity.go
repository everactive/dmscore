package main

import "github.com/everactive/dmscore/iot-identity/cmd/identity"

func main() {
	err := identity.Root.Execute()
	if err != nil {
		panic(err)
	}
}
