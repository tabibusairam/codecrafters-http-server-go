package main
import (
	"fmt"
	"strings"
	// Uncomment this block to pass the first stage
	"net"
	"os"
)
func getHeaders(req_parts []string) map[string]string {
	headers := make(map[string]string)
	for _, line := range req_parts {
		parts := strings.Split(line, ": ")
		if len(parts) == 2 {
			headers[parts[0]] = parts[1]
		}
	}
	return headers
}
func handleConnection(conn net.Conn) {
	var response []byte
	buffer := make([]byte, 1024)
	buffN, _ := conn.Read(buffer)
	// fmt.Println(buffN)
	request := string(buffer[:buffN])
	// fmt.Println("REQUEST: ", request)

	req_parts := strings.Split(request, "\r\n")
	// fmt.Println(req_parts)
	req_path_method := strings.Split(req_parts[0], " ")
	headers := getHeaders(req_parts)
	// fmt.Println(headers, req_path_method)
	if req_path_method[1] == "/" {
		response = []byte("HTTP/1.1 200 OK\r\nContent-Length: 13\r\nContent-Type: text/plain\r\n\r\nHello, world!")
		} else if strings.HasPrefix(req_path_method[1], "/echo/") {
			path := strings.Split(req_path_method[1], "/echo/")
			resource := path[len(path)-1]
			// fmt.Println(resource)
			res_len := len(resource)
			message := "HTTP/1.1 200 OK\r\nContent-Length:" + fmt.Sprint(res_len) + "\r\nContent-Type: text/plain\r\n\r\n" + resource
			response = []byte(message)
		} else if strings.HasPrefix(req_path_method[1], "/user-agent") {
			user_agent, exists := headers["User-Agent"]
			if !exists {
				response = []byte("HTTP/1.1 200 OK\r\nContent-Length: 25\r\nContent-Type: text/plain\r\n\r\nNo user-agent in headers")
			} else {
				message := "HTTP/1.1 200 OK\r\nContent-Length:" + fmt.Sprint(len(user_agent)) + "\r\nContent-Type: text/plain\r\n\r\n" + user_agent
				response = []byte(message)
			}
		} else {
			response = []byte("HTTP/1.1 404 Not Found\r\nContent-Length: 11\r\nContent-Type: text/plain\r\n\r\nNot Found")
		}
		_, res_err := conn.Write(response)
		if res_err != nil {
			fmt.Println("Error sending data", res_err.Error())
		}
		// return response
	}
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
		//
		defer listener.Close()
		for {
			conn, con_err := listener.Accept()
			if con_err != nil {
				fmt.Println("Error accepting connection: ", err.Error())
				os.Exit(1)
			}
			// response := []byte("HTTP/1.1 200 OK\r\n\r\n test output")
			handleConnection(conn)
			conn.Close()
		}
	}