package main

import (
	"fmt"
	"net"
	"log"
	"bufio"
	"os"
)

type MessageProtocol struct {
	Username [16]byte
	Type   byte // 4 bits of type (stored in the first byte)
	Length byte // 12 bits for length (stored in the first 2 bytes)
	Data   [62]byte  // Variable length data
}

var DataType = map[byte]string{
	1: "ClientMessage",
	2: "ServerMessage",
}

func Encode(username string, messagetype int, data string) MessageProtocol {
	message := MessageProtocol{}
	copy(message.Username[:], []byte(username));
	message.Type = byte(messagetype)
	message.Length = byte(len(data))
	copy(message.Data[:], []byte(data));
	
	return message
}

func Decode(data []byte) MessageProtocol {
		protocol := MessageProtocol{}
		copy(protocol.Username[:16], data[0:16])
		protocol.Type = data[17]
		protocol.Length = data[18]
		dataLength := 62
		copy(protocol.Data[:dataLength], data[18:18+dataLength])
		
		return protocol
	};

func SerializeMessage(message MessageProtocol) []byte {
	// Ensure username is exactly 16 bytes
    usernameBytes := make([]byte, 16)
    copy(usernameBytes, message.Username[:]) // Copy what fits
	
	serializedMessage := append([]byte{}, usernameBytes...);
	serializedMessage = append(serializedMessage, message.Type)
	serializedMessage = append(serializedMessage, message.Length)
	serializedMessage = append(serializedMessage, message.Data[:]...)
	
	return serializedMessage
}

func main() {
	reader := bufio.NewReader(os.Stdin)
    conn, err := net.Dial("tcp", "localhost:9000")
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()
    
    
    Username, _ := reader.ReadString('\n')
    
    // Goroutine to receive messages
    go func() {
        buf := make([]byte, 64)
        for {
            n, err := conn.Read(buf)
            if err != nil {
                return
            }
            if n > 0 {
                data := Decode(buf[:n])
                fmt.Print(string(data.Username[:len(data.Username)]), string(data.Data[:data.Length]))
            }
        }
    }()

    // Main loop to send messages
    
    for {
        input, _ := reader.ReadString('\n')
        conn.Write(SerializeMessage(Encode(Username, 1, input)))
    }
}
