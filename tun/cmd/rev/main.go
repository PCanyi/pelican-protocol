package main

import (
	"fmt"
	tun "github.com/mailgun/pelican-protocol/tun"
	"time"
)

func main() {

	//rdest := web.Cfg.Listen
	rdest := tun.NewAddr1("127.0.0.1:8080")
	rlsn := tun.NewAddr1("127.0.0.1:9999")

	fmt.Printf("rev starting: '%#v' -> '%#v'\n", rlsn, rdest)

	rev := tun.NewReverseProxy(tun.ReverseProxyConfig{Dest: rdest, Listen: rlsn})
	rev.Start()

	time.Sleep(30 * time.Second)
	rev.Stop()
	fmt.Printf("rev stopped.\n")
	time.Sleep(600 * time.Minute)
}