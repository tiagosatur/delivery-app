package main

import (
	"fmt"

	"github.com/codeedu/imersaofsfc2-simulator/application/route"
)

func main() {
	route := route.Route{
		ID: "1",
		ClientID: "1",
	}

	route.LoadPositions()

	stringjson, _ := route.ExportJsonPositinos()
	fmt.Println(stringjson[0])
}