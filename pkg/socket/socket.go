package socket

import (
	"errors"
	"net"
	"syscall"

	"golang.org/x/sys/unix"
)

// Socket represents an eap socket
type Socket struct {
	ifname string
	fd     int
	closed bool
}

// New builds a socket ready for eap
func New(ifname string, addrSlice []byte) (*Socket, error) {
	if len(addrSlice) > 8 {
		return nil, errors.New("addrSlice too long")
	}
	iface, err := net.InterfaceByName(ifname)
	if err != nil {
		return nil, err
	}
	fd, err := unix.Socket(unix.AF_PACKET, unix.SOCK_RAW, syscall.ETH_P_PAE)
	if err != nil {
		return nil, err
	}
	err = unix.SetNonblock(fd, true)
	if err != nil {
		return nil, err
	}
	unix.BindToDevice(fd, ifname)
	var addr [8]uint8
	copy(addr[:], addrSlice)
	packet := unix.PacketMreq{
		Ifindex: int32(iface.Index),
		Type:    unix.PACKET_MR_MULTICAST,
		Alen:    uint16(len(addrSlice)),
		Address: addr,
	}
	unix.SetsockoptPacketMreq(fd, unix.SOL_PACKET, unix.PACKET_ADD_MEMBERSHIP, &packet)
	socket := &Socket{
		fd:     fd,
		ifname: ifname,
		closed: false,
	}
	return socket, nil
}

// GetFileDescriptor returns the current file descriptor
func (sock *Socket) GetFileDescriptor() int {
	return sock.fd
}

// GetIfname returns the interface name of this socket
func (sock *Socket) GetIfname() string {
	return sock.ifname
}

// Write writes to the file descriptor
func (sock *Socket) Write(buf []byte) error {
	if sock.closed == true {
		return nil
	}
	_, err := unix.Write(sock.fd, buf)
	return err
}

// WriteWithErrorSignal writes to the file descriptor but sends any errors
// over a channel instead of returning them
func (sock *Socket) WriteWithErrorSignal(buf []byte, errSignal chan error) {
	err := sock.Write(buf)
	if err != nil {
		errSignal <- err
	}
}

// Close closes the underlying file descriptor
func (sock *Socket) Close() error {
	sock.closed = true
	return unix.Close(sock.fd)
}
