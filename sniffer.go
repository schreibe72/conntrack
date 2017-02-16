package main

import (
	"net"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

type Sniffer struct {
	snapshotLen int32
	promiscuous bool
	timeout     time.Duration
	handle      *pcap.Handle
}

type Proto int

const (
	TCP Proto = iota
	UDP
)

var Protos = [...]string{
	"tcp",
	"udp",
}

func NewSniffer(device string, snapshot_len int32, promiscuous bool, timeout time.Duration) (Sniffer, error) {
	s := Sniffer{
		snapshotLen: snapshot_len,
		promiscuous: promiscuous,
		timeout:     timeout,
		handle:      nil,
	}
	handle, err := pcap.OpenLive(device, snapshot_len, promiscuous, timeout)
	if err != nil {
		return s, err
	}
	s.handle = handle
	return s, nil
}

func (s Sniffer) SetFilter(f string) error {
	return s.handle.SetBPFFilter(f)
}

func (s Sniffer) Sniff(callback func(proto Proto, srcip net.IP, srcport uint16, dstip net.IP, dstport uint16)) error {
	defer s.handle.Close()
	packetSource := gopacket.NewPacketSource(s.handle, s.handle.LinkType())
	for packet := range packetSource.Packets() {
		var ip4 *layers.IPv4
		if ip4Layer := packet.Layer(layers.LayerTypeIPv4); ip4Layer != nil {
			ip4 = ip4Layer.(*layers.IPv4)
		} else {
			continue
		}
		if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
			tcp, _ := tcpLayer.(*layers.TCP)
			callback(TCP, ip4.SrcIP, uint16(tcp.SrcPort), ip4.DstIP, uint16(tcp.DstPort))
		}
		if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
			udp, _ := udpLayer.(*layers.UDP)
			callback(UDP, ip4.SrcIP, uint16(udp.SrcPort), ip4.DstIP, uint16(udp.DstPort))
		}
	}
	return nil
}
