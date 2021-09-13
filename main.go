package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/lwch/runtime"
)

func main() {
	listen := flag.Uint("listen", 0, "监听端口")
	remote := flag.String("addr", "", "转发IP或域名")
	port := flag.Uint("port", 0, "转发端口号")
	flag.Parse()

	if *listen == 0 {
		fmt.Println("缺少listen参数")
		os.Exit(1)
	}

	if len(*remote) == 0 {
		fmt.Println("缺少remote参数")
		os.Exit(1)
	}

	if *port == 0 {
		fmt.Println("缺少port参数")
		os.Exit(1)
	}

	l, err := net.ListenTCP("tcp", &net.TCPAddr{
		Port: int(*listen),
	})
	runtime.Assert(err)
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}
		go forward(conn, *remote, uint16(*port))
	}
}

func forward(local net.Conn, addr string, port uint16) {
	defer local.Close()
	ip, err := net.ResolveIPAddr("ip", addr)
	runtime.Assert(err)
	remote, err := net.DialTCP("tcp", nil, &net.TCPAddr{
		IP:   ip.IP,
		Port: int(port),
		Zone: ip.Zone,
	})
	runtime.Assert(err)
	go io.Copy(remote, local)
	io.Copy(local, remote)
}
