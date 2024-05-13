package server

import (
	"log"
	"net"
	"syscall"

	"github.com/rajchandvaniya/orangedb/config"
	"github.com/rajchandvaniya/orangedb/core"
)

func StartAsynchronousTCPServer() {
	log.Println("Starting asynchronous TCP server on ", config.Host, config.Port)
	conClients := 0

	var events []syscall.EpollEvent = make([]syscall.EpollEvent, config.MaxConnections)

	// create a socket
	serverFD, err := syscall.Socket(syscall.AF_INET, syscall.O_NONBLOCK|syscall.SOCK_STREAM, 0)
	if err != nil {
		log.Fatal(err)
	}
	defer syscall.Close(serverFD)

	// set socket to operate in non-blocking mode
	if err := syscall.SetNonblock(serverFD, true); err != nil {
		log.Fatal(err)
	}

	// bind IP and port to it
	ip4 := net.ParseIP(config.Host)
	if err := syscall.Bind(serverFD, &syscall.SockaddrInet4{
		Port: config.Port,
		Addr: [4]byte{ip4[0], ip4[1], ip4[2], ip4[3]},
	}); err != nil {
		log.Fatal(err)
	}

	// start listening on server socket
	if err := syscall.Listen(serverFD, config.MaxConnections); err != nil {
		log.Fatal(err)
	}

	// starting async io by setting up EPOLL
	epollFD, err := syscall.EpollCreate1(0)
	if err != nil {
		log.Fatal(err)
	}
	defer syscall.Close(epollFD)

	// specify events we want epoll to monitor and registering it
	var socketServerEvent syscall.EpollEvent = syscall.EpollEvent{
		Events: syscall.EPOLLIN,
		Fd:     int32(serverFD),
	}
	if err := syscall.EpollCtl(epollFD, syscall.EPOLL_CTL_ADD, serverFD, &socketServerEvent); err != nil {
		log.Fatal(err)
	}

	for {
		// see if any FD is available for new I/O
		nevents, err := syscall.EpollWait(epollFD, events[:], -1)
		if err != nil {
			continue
		}

		for i := 0; i < nevents; i++ {
			if int(events[i].Fd) == serverFD {
				// server is ready to accept a new connection

				// accept the connection and register with epoll
				fd, _, err := syscall.Accept(serverFD)
				if err != nil {
					log.Println(err)
					continue
				}

				conClients++
				syscall.SetNonblock(fd, true)

				var socketClientEvent syscall.EpollEvent = syscall.EpollEvent{
					Events: syscall.EPOLLIN,
					Fd:     int32(fd),
				}
				if err := syscall.EpollCtl(epollFD, syscall.EPOLL_CTL_ADD, fd, &socketClientEvent); err != nil {
					log.Fatal(err)
				}
			} else {
				// client is ready with I/O
				comm := core.FDComm{Fd: int(events[i].Fd)}
				cmd, err := readCommand(comm)
				if err != nil {
					syscall.Close(comm.Fd)
					conClients--
					continue
				}
				respond(cmd, comm)
			}
		}
	}

}
