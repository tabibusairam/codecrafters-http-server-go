

package main
import (
	"fmt"
	// Uncomment this block to pass the first stage
	"net"
	"os"
	"strings"
)
func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")
	// Uncomment this block to pass the first stage
	//
	listener, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConnection(conn)
	}
}
func handleConnection(connection net.Conn) {

	defer connection.Close()
	buffer := make([]byte, 1024)
	connection.Read(buffer)
	request := string(buffer)
	firstLine := request[:strings.Index(request, "\n")]
	path := firstLine[strings.Index(firstLine, "/") : strings.LastIndex(firstLine, "HTTP")-1]
	fmt.Println("Path: ", path)
	if path == "/" {
		response := "HTTP/1.1 200 OK\r\n\r\n"
		_, err := connection.Write([]byte(response))
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
			fmt.Println("Failed to response", err.Error())
		}
		connection.Close()
	} else if strings.HasPrefix(path, "/echo") {
		echoParam := path[6:]
		body := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + fmt.Sprint(len(echoParam)) + "\r\n\r\n" + echoParam
		_, err := connection.Write([]byte(body))
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
			fmt.Println("Failed to response", err.Error())
		}
		connection.Close()
	} else if strings.HasPrefix(path, "/user-agent") {
		requestLines := strings.Split(request, "\n")
		userAgent := ""
		for _, line := range requestLines {
			if strings.HasPrefix(line, "User-Agent") {
				userAgent = line[strings.Index(line, ":")+2 : len(line)-1]
				break
			}
		}
		body := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + fmt.Sprint(len(userAgent)) + "\r\n\r\n" + userAgent
		_, err := connection.Write([]byte(body))
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
			fmt.Println("Failed to response", err.Error())
		}
		connection.Close()
	} else {
		response := "HTTP/1.1 404 Not Found\r\n\r\n"
		_, err := connection.Write([]byte(response))
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
			fmt.Println("Failed to response", err.Error())
		}
		connection.Close()
	}
}

