package main
import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)
const (
	HOST = "localhost"
	PORT = "4221"
	TYPE = "tcp"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
	listener, err := net.Listen(TYPE, HOST+":"+PORT)
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	path := getPath(conn)
	var response []byte
	if path != "/" {
		response = []byte("HTTP/1.1 404 Not Found\r\n\r\n")
	} else {
		response = []byte("HTTP/1.1 200 OK\r\n\r\n")
	}
	_, errWrite := conn.Write(response)
	if errWrite != nil {
		log.Fatal(errWrite)
	}
}
func getPath(conn net.Conn) string {
	buffer := make([]byte, 1024)
	buffN, errRd := conn.Read(buffer)
	if errRd != nil {
		log.Fatal(errRd)
	}

	request := string(buffer[:buffN])
	fmt.Println("REQUEST: ", request)
	startLine := strings.Split(request, "\n")[0]
	fmt.Println("START LINE: ", startLine)
	path := strings.Split(startLine, " ")[1]
	fmt.Println("PATH: ", path)
	return path
}
