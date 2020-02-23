package rproxy

import (
	"fmt"
	"log"
	"net"
)

type RouteTable struct {
	table map[string]Backend
}

type Backend struct {
	addr  *net.UDPAddr
	Ready bool
}

var backends = map[string]string{
	"calculator": ":9101",
	// "127.0.0.1:9102",
}

var rt = InitRouteTable()

func InitRouteTable() RouteTable {
	rt := make(map[string]Backend)

	for k, v := range backends {
		addr, err := net.ResolveUDPAddr("udp", v)
		if err != nil {
			panic(err)
		}

		rt[k] = Backend{
			addr:  addr,
			Ready: true,
		}
	}

	return RouteTable{
		table: rt,
	}
}

func Start(port int) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	c, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}

	fmt.Printf("proxying CoAP requests on %s\n", c.LocalAddr())

	buffer := make([]byte, 64)
	for {
		fmt.Println("reading from udp conn")
		n, raddr, err := c.ReadFrom(buffer)
		if err != nil {
			panic(err)
		}

		fmt.Printf("read %d bytes from %v\n", n, raddr)

		fmt.Println("forwarding to handler")

		destConn, err := net.DialUDP("udp", nil, rt.table["calculator"].addr)
		if err != nil {
			log.Fatal(err)
		}

		_, err = destConn.Write(buffer[0:n])
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("reading response from calc handler...")
		buffer = make([]byte, 64)
		n, err = destConn.Read(buffer)
		if err != nil {
			log.Fatal(err)
		}
		err = destConn.Close()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("writing response from calc handler back to rmote device conn...")
		c.Write(buffer[0:n])

		buffer = make([]byte, 64)
	}
}
