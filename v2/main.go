package main

import (
	"fmt"

	"github.com/coapcloud/veetoo/api"
	"github.com/coapcloud/veetoo/rproxy"
	flag "github.com/spf13/pflag"
)

var port = flag.IntP("port", "p", 5683, "coap port to listen on")

func main() {
	flag.Parse()

	fmt.Println("starting api...")
	go api.Start()

	rproxy.Start(*port)
}
