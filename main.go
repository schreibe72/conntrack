package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"path/filepath"
	"strconv"
	"time"
)

var (
	snapshot_len int32         = 1024
	promiscuous  bool          = false
	timeout      time.Duration = 30 * time.Second
	filter       string        = "tcp[tcpflags] & (tcp-syn) != 0 and tcp[tcpflags] & (tcp-ack) == 0 or udp"
)

type tracker struct {
	device       string
	localdevice  Interface
	path         string
	skipelse     bool
	udpports     []uint16
	seenudpports []uint16
}

type udpports []uint16

func (t *tracker) connection(proto Proto, srcip net.IP, srcport uint16, dstip net.IP, dstport uint16) {
	var direction string
	var fileParts struct {
		srcip   string
		srcport string
		dstip   string
		dstport string
	}

	switch {
	case t.localdevice.isLocalIP(srcip):
		direction = "out"
		if isInternet(dstip) {
			fileParts.dstip = "INTERNET"
		} else {
			fileParts.dstip = dstip.String()
		}
		fileParts.dstport = fmt.Sprintf("%d", dstport)
		fileParts.srcport = "PORT"
		fileParts.srcip = srcip.String()
		if proto == UDP && uint16InSlice(t.udpports, srcport) {
			return
		}
		if proto == UDP && dstport <= 1024 {
			t.seenudpports = append(t.seenudpports, dstport)
		}
	case t.localdevice.isLocalIP(dstip):
		direction = "in"
		if isInternet(srcip) {
			fileParts.srcip = "INTERNET"
		} else {
			fileParts.srcip = srcip.String()
		}
		fileParts.dstport = fmt.Sprintf("%d", dstport)
		fileParts.srcport = "PORT"
		fileParts.dstip = dstip.String()
		if proto == UDP && uint16InSlice(t.seenudpports, srcport) {
			return
		}
	default:
		if t.skipelse {
			return
		}
		direction = "else"
		fileParts.srcip = srcip.String()
		fileParts.srcport = fmt.Sprintf("%d", srcport)
		fileParts.dstip = dstip.String()
		fileParts.dstport = fmt.Sprintf("%d", dstport)
	}
	var path string
	path = filepath.Join(t.path, direction, Protos[proto])
	b := []byte{}
	file := fmt.Sprintf("%s-%s-%s-%s", fileParts.srcip, fileParts.srcport, fileParts.dstip, fileParts.dstport)
	err := ioutil.WriteFile(filepath.Join(path, file), b, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func (u *udpports) Set(value string) error {
	tmp, err := strconv.Atoi(value)
	if err != nil {
		*u = append(*u, 0)
	} else {
		*u = append(*u, uint16(tmp))
	}
	return nil
}

func (u *udpports) String() string {
	return fmt.Sprintf("%d", *u)
}

func main() {
	var t tracker
	var u udpports
	flag.StringVar(&t.device, "d", "eth0", "Network Device")
	flag.StringVar(&t.path, "p", "", "Path to Store")
	flag.BoolVar(&t.skipelse, "e", false, "Skip else connection")
	flag.Var(&u, "u", "List of udp Ports")
	flag.Parse()
	t.udpports = []uint16(u)
	ld, err := getLocalInterface(t.device)
	if err != nil {
		log.Fatal(err)
	}
	createDataPath(t.path)
	t.localdevice = ld
	sniffer, err := NewSniffer(t.device, snapshot_len, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	sniffer.SetFilter(filter)
	sniffer.Sniff(t.connection)
}
