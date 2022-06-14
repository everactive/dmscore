package main

import (
	"github.com/everactive/dmscore/iot-devicetwin/cmd/devicetwin"
)

func main() {
	err := devicetwin.Root.Execute()
	if err != nil {
		panic(err)
	}
}
