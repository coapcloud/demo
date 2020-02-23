package rproxy

import (
	"fmt"
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
	"calculator": "127.0.0.1:9101",
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
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("127.0.0.1:%d", port))
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
		n, raddr, err := c.ReadFrom(buffer)
		if err != nil {
			panic(err)
		}

		fmt.Printf("read %d bytes from %v\n", n, raddr)

		fmt.Println("forwarding to handler")

		pConn, err := net.DialUDP("udp", addr, rt.table["calculator"].addr)
		if err != nil {
			panic(err)
		}

		_, err = pConn.Write(buffer[0:n])
		if err != nil {
			panic(err)
		}

		buffer = make([]byte, 64)
		n, err = pConn.Read(buffer)
		if err != nil {
			panic(err)
		}

		respAddr, err := net.ResolveUDPAddr("udp", raddr.String())
		if err != nil {
			panic(err)
		}

		respConn, err := net.DialUDP("udp", addr, respAddr)
		if err != nil {
			panic(err)
		}

		respConn.Write(buffer[0:n])

		buffer = make([]byte, 64)
	}
}
