package main

import (
	"fmt"

	"github.com/AbdulHaseebAhmad/New folder/Internal/Configurator"
)

func main() {
	cfg := Configurator.LoadConfiguration()
	fmt.Println("All good")
}
