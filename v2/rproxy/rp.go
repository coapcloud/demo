package rproxy

import (
	"encoding/gob"
	"io"
	"net"
	"sync"
)

type RouteTable struct {
	table map[string]Backend
	lock  *sync.RWMutex
}

type Backend struct {
	addr  *net.UDPAddr
	Ready bool
}

func NewRouteTable() RouteTable {
	return RouteTable{
		table: make(map[string]Backend),
		lock:  &sync.RWMutex{},
	}
}

func (rt RouteTable) Serialize(w io.Writer) error {
	err := gob.NewEncoder(w).Encode(rt)
	if err != nil {
		return err
	}

	return nil
}

func Unserialize(r io.Reader) (RouteTable, error) {
	var table RouteTable
	err := gob.NewDecoder(r).Decode(&table)
	if err != nil {
		return RouteTable{}, err
	}

	return table, nil
}

func Start(port int) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("127.0.0.1:%d", port)
	if err != nil {
		panic(err)
	}
	
	c, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}
	
	fmt.Printf("serving CoAP requests on %d\n", c.LocalAddr())

	buffer := make([]byte, 65535)

	for {
		n, addr, err := c.ReadFrom(buffer)
		if err != nil {
			panic(err)
		}

		fmt.Printf("read %d bytes from %v\n", n, addr)
		fmt.Println(string(buffer[0:n]))
	}
}