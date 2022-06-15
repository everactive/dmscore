package main

import "github.com/everactive/dmscore/cmd"

func main() {
	err := cmd.Root.Execute()
	if err != nil {
		panic(err)
	}
}
