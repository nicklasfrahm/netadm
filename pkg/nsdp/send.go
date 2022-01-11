package nsdp

import (
	"context"
	"errors"
	"fmt"
	"net"
)

// Send is a low-level API that allows it to send and receive messages directly.
func Send(ctx context.Context, iface *net.Interface, request *Message) ([]Message, error) {
	// Check if the provided interface has a valid configuration.
	ip, err := GetInterfaceIP(iface)
	if err != nil {
		return nil, err
	}

	// Create a UDP socket to listen for incoming packets.
	socketAddr := net.UDPAddr{
		IP:   *ip,
		Port: ClientPort,
	}
	socket, err := net.ListenUDP("udp", &socketAddr)
	if err != nil {
		return nil, err
	}
	defer socket.Close()

	responses := make([]Message, 0)
	errs := make(chan error, 1)

	// Create a goroutine to listen for incoming packets.
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				buf := make([]byte, 1500)

				n, err := socket.Read(buf)
				if err != nil {
					errs <- err
					return
				}

				response := new(Message)
				if err := response.UnmarshalBinary(buf[:n]); err != nil {
					errs <- errors.New("malformed response message")
					return
				}

				// Check operation result status code.
				// I assume all non-zero values are bad.
				if response.Header.Result != 0 {
					errs <- fmt.Errorf("operation failed with status code 0x%04X", response.Header.Result)
					return
				}

				responses = append(responses, *response)
			}
		}
	}()

	// Create discovery message and encode it into its binary form.
	payload, err := request.MarshalBinary()
	if err != nil {
		return nil, err
	}

	// Send the message to the broadcast address.
	deviceAddr := net.UDPAddr{
		IP:   net.IPv4bcast,
		Port: ServerPort,
	}
	if _, err := socket.WriteToUDP(payload, &deviceAddr); err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return responses, nil
	case err := <-errs:
		return nil, err
	}
}
