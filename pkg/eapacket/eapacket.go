package eapacket

import (
	"errors"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// EAPType represents a type of an eap packet
type EAPType uint8

// EAPTypes
const (
	Request EAPType = iota + 1
	Response
	Success
	Failure
)

// Packet represents an EAP packet
type Packet struct {
	eap *layers.EAP
	eth *layers.Ethernet
}

// Decode creates an eapacket.Packet from a byte array buffer
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

// Dest returns the ethernet layer's destination mac address from the packet
func (packet *Packet) Dest() net.HardwareAddr {
	return packet.eth.DstMAC
}

// Src returns the ethernet layer's source mac address from the packet
func (packet *Packet) Src() net.HardwareAddr {
	return packet.eth.SrcMAC
}

// Type returns the type in the eap layer
func (packet *Packet) Type() EAPType {
	return EAPType(packet.eap.Type)
}
