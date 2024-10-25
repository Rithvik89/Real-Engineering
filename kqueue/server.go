package main
​
import (
	"log"
	"net"
​
	"golang.org/x/sys/unix"
)
​
const (
	port       = ":6379"
	bufferSize = 4096
	maxEvents  = 10
)
​
// getFD extracts the file descriptor from a net.Listener or net.Conn
func getFd(conn net.Conn, listener net.Listener) int {
	if conn != nil {
		file, err := conn.(*net.TCPConn).File()
		if err != nil {
			log.Fatalf("Failed to get file descriptor for connection: %v", err)
		}
		return int(file.Fd())
	} else {
		file, err := listener.(*net.TCPListener).File()
		if err != nil {
			log.Fatalf("Failed to get file descriptor for listener: %v", err)
		}
		return int(file.Fd())
	}
}
​
func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Unable to listen on %s: %v", port, err)
	}
	defer listener.Close()
​
	listenerFd := getFd(nil, listener)
	if err := unix.SetNonblock(listenerFd, true); err != nil {
		log.Fatalf("Failed to set listener to non-blocking mode: %v", err)
	}
	// Create a kqueue instance
	kqueueFd, err := unix.Kqueue()
	if err != nil {
		log.Fatalf("Failed to create kqueue instance: %v", err)
	}
	defer unix.Close(kqueueFd)
	// Add the listener file descriptor to the kqueue instance
	event := unix.Kevent_t{
		Ident:  uint64(listenerFd),
		Filter: unix.EVFILT_READ,
		Flags:  unix.EV_ADD | unix.EV_ENABLE,
	}
	_, err = unix.Kevent(kqueueFd, []unix.Kevent_t{event}, nil, nil)
	if err != nil {
		log.Fatalf("Failed to add listener fd to kqueue: %v", err)
	}
	// Create a buffer for events
	events := make([]unix.Kevent_t, maxEvents)
	for {
		// Wait for events
		n, err := unix.Kevent(kqueueFd, nil, events, nil)
		if err != nil {
			log.Fatalf("Kevent error: %v", err)
		}
		// Handle events
		for i := 0; i < n; i++ {
			if int(events[i].Ident) == listenerFd {
				// Accept new connection
				conn, err := listener.Accept()
				if err != nil {
					log.Printf("Failed to accept connection: %v", err)
					continue
				}
				// Get the file descriptor for the connection
				connFd := getFd(conn, nil)
				// Set the connection to non-blocking mode
				if err := unix.SetNonblock(connFd, true); err != nil {
					log.Printf("Failed to set connection to non-blocking mode: %v", err)
					conn.Close()
					continue
				}
				// Add the connection file descriptor to the kqueue instance
				event := unix.Kevent_t{
					Ident:  uint64(connFd),
					Filter: unix.EVFILT_READ,
					Flags:  unix.EV_ADD | unix.EV_ENABLE,
				}
				_, err = unix.Kevent(kqueueFd, []unix.Kevent_t{event}, nil, nil)
				if err != nil {
					log.Printf("Failed to add connection fd to kqueue: %v", err)
					conn.Close()
					continue
				}
			} else {
				// Handle data from an existing connection
				connFd := int(events[i].Ident)
				buffer := make([]byte, bufferSize)
				n, err := unix.Read(connFd, buffer)
				if err != nil || n == 0 {
​
					unix.Kevent(kqueueFd, []unix.Kevent_t{{
						Ident:  uint64(connFd),
						Filter: unix.EVFILT_READ,
						Flags:  unix.EV_DELETE,
					}}, nil, nil)
					unix.Close(connFd)
					continue
				}
				log.Printf("Received %d bytes: %s", n, string(buffer[:n]))
				_, err = unix.Write(connFd, []byte("Message received."))
				if err != nil {
​
					unix.Kevent(kqueueFd, []unix.Kevent_t{{
						Ident:  uint64(connFd),
						Filter: unix.EVFILT_READ,
						Flags:  unix.EV_DELETE,
					}}, nil, nil)
					unix.Close(connFd)
				}
			}
		}
	}
}