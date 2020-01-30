package main

import (
	"fmt"

	"github.com/coapcloud/demo/examplefuncs/calculator"
	"github.com/coapcloud/demo/rproxy"
	flag "github.com/spf13/pflag"
)

var port = flag.IntP("port", "p", 5683, "coap port to listen on")

func main() {
	flag.Parse()
	go startFuncs()
	rproxy.Start(*port)
}

func startFuncs() {
	fmt.Println("starting calculator func listener...")
	calculator.Run(9101)
}
