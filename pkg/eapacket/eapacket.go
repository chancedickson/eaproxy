package eapacket

import (
	"errors"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type Packet struct {
	eap *layers.EAP
	eth *layers.Ethernet
}

func Decode(buf []byte) (*Packet, error) {
	packet := gopacket.NewPacket(buf, layers.LayerTypeEthernet, gopacket.Default)
	ethLayer := packet.Layer(layers.LayerTypeEthernet)
	if ethLayer == nil {
		return nil, errors.New("can't decode ethernet layer")
	}
	eth := ethLayer.(*layers.Ethernet)
	eapLayer := packet.Layer(layers.LayerTypeEAP)
	if eapLayer == nil {
		return nil, errors.New("can't decode eap layer")
	}
	eap := eapLayer.(*layers.EAP)
	return &Packet{eap, eth}, nil
}

func (packet *Packet) Dest() net.HardwareAddr {
	return packet.eth.DstMAC
}

func (packet *Packet) Src() net.HardwareAddr {
	return packet.eth.SrcMAC
}
