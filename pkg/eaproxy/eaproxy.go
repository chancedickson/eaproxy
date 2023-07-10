package eaproxy

import (
	"eaproxy/pkg/dhcp"
	"eaproxy/pkg/eapacket"
	"eaproxy/pkg/socket"
	"errors"
	"fmt"

	"golang.org/x/sys/unix"
)

// Eaproxy represents an eaproxy instance
type Eaproxy struct {
	routerSocket *socket.Socket
	wanSocket    *socket.Socket
	ifVlan       string
}

// New makes an eaproxy
func New(routerSocket *socket.Socket, wanSocket *socket.Socket, vlanID int) *Eaproxy {
	return &Eaproxy{
		routerSocket,
		wanSocket,
		fmt.Sprintf("%s.%d", wanSocket.GetIfname(), vlanID),
	}
}

// Start starts the eaproxy
func (eaproxy *Eaproxy) Start() error {
	epoll, err := unix.EpollCreate1(0)
	if err != nil {
		return err
	}
	defer unix.Close(epoll)
	routerFd := eaproxy.routerSocket.GetFileDescriptor()
	routerEvent := &unix.EpollEvent{
		Events: unix.EPOLLIN | unix.EPOLLOUT,
		Fd:     int32(routerFd),
	}
	err = unix.EpollCtl(epoll, unix.EPOLL_CTL_ADD, routerFd, routerEvent)
	if err != nil {
		return err
	}
	wanFd := eaproxy.wanSocket.GetFileDescriptor()
	wanEvent := &unix.EpollEvent{
		Events: unix.EPOLLIN | unix.EPOLLOUT,
		Fd:     int32(wanFd),
	}
	err = unix.EpollCtl(epoll, unix.EPOLL_CTL_ADD, wanFd, wanEvent)
	if err != nil {
		return err
	}
	stopSignal := make(chan bool, 1)
	errSignal := make(chan error)
	epollWait := eaproxy.newEpollWaitChannel(epoll, stopSignal, errSignal)
	for {
		select {
		case events := <-epollWait:
			for _, event := range events {
				if event.Events&unix.EPOLLIN != 0 {
					return errors.New("event not EPOLLIN")
				}
				go eaproxy.handleEvent(event, errSignal)
			}
		case err := <-errSignal:
			return err
		}
	}
}

func (eaproxy *Eaproxy) newEpollWaitChannel(epoll int, stopSignal chan bool, errSignal chan error) chan []unix.EpollEvent {
	eventChannel := make(chan []unix.EpollEvent)
	var eventBuf [32]unix.EpollEvent
	go func() {
		for {
			n, err := unix.EpollWait(epoll, eventBuf[:], 250)
			if err != nil {
				errSignal <- err
				return
			}
			if n > 0 {
				events := make([]unix.EpollEvent, n)
				copy(events[:], eventBuf[:n])
				eventChannel <- events
			}
			if len(stopSignal) == 1 {
				stop := <-stopSignal
				if stop {
					return
				}
			}
		}
	}()
	return eventChannel
}

func (eaproxy *Eaproxy) handleEvent(event unix.EpollEvent, errSignal chan error) {
	fd := int(event.Fd)
	var buf []byte
	_, err := unix.Read(fd, buf)
	if err != nil {
		errSignal <- err
		return
	}
	packet, err := eapacket.Decode(buf)
	if err != nil {
		errSignal <- err
		return
	}
	var handle func(packet *eapacket.Packet) error
	var write func(buf []byte, errSignal chan error)
	if fd == eaproxy.routerSocket.GetFileDescriptor() {
		handle = eaproxy.handleRouterPacket
		write = eaproxy.wanSocket.WriteWithErrorSignal
	} else if fd == eaproxy.wanSocket.GetFileDescriptor() {
		handle = eaproxy.handleWanPacket
		write = eaproxy.routerSocket.WriteWithErrorSignal
	} else {
		// received event for unknown file descriptor.. something went really wrong
		return
	}
	// proxy the packet in a separate goroutine
	go write(buf, errSignal)
	err = handle(packet)
	if err != nil {
		errSignal <- err
		return
	}
}

// no-op currently. May be used later.
func (eaproxy *Eaproxy) handleRouterPacket(packet *eapacket.Packet) error {
	return nil
}

func (eaproxy *Eaproxy) handleWanPacket(packet *eapacket.Packet) error {
	if packet.Type() == eapacket.Success {
		return dhcp.Restart(eaproxy.ifVlan)
	}
	return nil
}
