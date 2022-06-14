package main

import "github.com/everactive/dmscore/iot-management/cmd/management"

func main() {
	err := management.Command.Execute()
	if err != nil {
		panic(err)
	}
}
