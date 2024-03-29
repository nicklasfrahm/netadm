package nsdp

import (
	"context"
	"fmt"
	"net"
)

// Send is a low-level API that allows it to send and receive messages directly.
// It is recommended to set an explicit destination IP as otherwise the message
// will be sent to the global broadcast address, which is often filtered out by
// routers.
func Send(ctx context.Context, iface *net.Interface, dst *net.IP, request *Message) ([]Message, error) {
	// Create a UDP socket to listen for incoming packets.
	socketAddr := net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
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
					if response.Header.Result == 0 {
						errs <- ErrInvalidResponse
						return
					}
				}

				// Check operation result status code.
				// I assume all non-zero values are bad.
				if response.Header.Result != 0 {
					switch response.Header.Result {
					case uint16(ResponseCodeInvalidRecordLength):
						errs <- ErrInvalidRecordLength
					case uint16(ResponseCodeInvalidPassword):
						errs <- ErrInvalidPassword
					case uint16(ResponseCodeInvalidPasswordLockdown):
						errs <- ErrInvalidPasswordLockdown
					default:
						errs <- fmt.Errorf("operation failed with status code 0x%04X", response.Header.Result)
					}
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
	if dst != nil {
		deviceAddr.IP = *dst
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
