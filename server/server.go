	package server

	import (
		"net"
		"log"
		"fmt"
		"sync"
	)
	
	type ConnList struct {
		mu sync.Mutex
		Clients map[net.Conn]bool
	}
	
	type MessageProtocol struct {
	Username [16]byte
	Type   byte
	Length byte
	Data   [62]byte 
	}
	
	var DataType = map[byte]string{
		1: "ClientMessage",
		2: "ServerMessage",
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
	
	func Broadcast(clients *ConnList, message []byte, bl net.Conn) {
		clients.mu.Lock();
		defer clients.mu.Unlock();
		fmt.Print(string(message[2:]));
		fmt.Println("Broadcasting to :");
		
		for conn := range clients.Clients {
			if conn != bl {
				fmt.Println(conn.RemoteAddr().String());
				_, err := conn.Write(message)
				if err != nil {
					fmt.Print("err:", err);
					delete(clients.Clients, conn);
					conn.Close();
				}
			} else if conn == bl {
				fmt.Println("skipping", conn.RemoteAddr().String());
			}
		}
		
	};
	
	func Server(host string, port int) {
		listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
		fmt.Printf("Started listening on port %d", port)
		if err != nil {
			log.Fatal(err)
		}
		defer listener.Close()
		
		clients := ConnList{
			Clients: make(map[net.Conn]bool),
		};
		
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Fatal(err)
			}
			
			clients.mu.Lock();
			clients.Clients[conn] = true
			clients.mu.Unlock();
			
			go handleconn(conn, &clients);
		}
	}

	func handleconn(conn net.Conn, cl *ConnList) {
		message := make([]byte, 80);
		for {
			_, err := conn.Read(message);
			
			if err != nil {
				closerr := conn.Close();
				if closerr != nil {
					return
				}
				
				return
			}
			
			if len(message) > 0 {
				data := Decode(message);
				fmt.Println("");
				fmt.Print(conn.RemoteAddr().String());
				fmt.Print(" ", string(data.Username[:]));
				fmt.Print(" ", DataType[data.Type]);
				fmt.Print(" ",data.Length);
				fmt.Print(" ",string(data.Data[:]));
				
				
				message[0] = 2
				
				Broadcast(cl, message[:len(message)], conn);
			}
		}
	}
