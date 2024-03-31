package main
import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	// Uncomment this block to pass the first stage
	// "net"
	// "os"
)
const (
	OkResponse                        = "HTTP/1.1 200 OK\r\n"
	NotFoundResponse                  = "HTTP/1.1 404 Not Found\r\n"
	ContentTypeTextPlain              = "Content-Type: text/plain\r\n"
	ContentTypeApplicationOctetStream = "Content-Type: application/octet-stream\r\n"
)
func contentLength(str string) string {
	length := len(str)
	return "Content-Length: " + strconv.Itoa(length) + "\r\n\r\n"
}
func main() {
	defaultDirectory := "./"
	args := os.Args
	if len(args) > 2 {
		if args[1] == "--directory" {
			defaultDirectory = args[2]
		}
	}
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	fmt.Println("Server up listening on 0.0.0.0:4221")
	for {
		// Waiting for connection
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		fmt.Println("Accepted connection")
		go connHandler(conn, defaultDirectory)
	}
}
// Manage each connection
func connHandler(conn net.Conn, dir string) {
	receivedBytes := make([]byte, 1024)
	_, err := conn.Read(receivedBytes)
	if err != nil {
		fmt.Println("Error while reading data")
		os.Exit(1)
	}
	lines := strings.Split(string(receivedBytes), "\r\n")
	path := strings.Split(lines[0], " ")[1]
	switch {
	case path == "/":
		_, err = conn.Write([]byte(OkResponse + "\r\n"))
	case strings.Contains(path, "/echo/"):
		_, word, _ := strings.Cut(path, "/echo/")
		fmt.Println("Received ECHO for word: ", word, "length: ", strconv.Itoa(len(word)))
		_, err = conn.Write([]byte((OkResponse + ContentTypeTextPlain + contentLength(word) + word + "\r\n")))
	case strings.Contains(path, "/user-agent"):
		for _, line := range lines {
			if strings.Contains(line, "User-Agent") {
				userAgent := strings.Split(line, "User-Agent: ")[1]
				_, err = conn.Write([]byte((OkResponse + ContentTypeTextPlain + contentLength(userAgent) + userAgent + "\r\n")))
			}
		}
	case strings.Contains(path, "/files/"):
		_, file, _ := strings.Cut(path, "/files/")
		fmt.Println("Received GET FILES for file: ", file, "length: ", strconv.Itoa(len(file)), dir+file)
		fileContent, err := os.ReadFile(dir + file)
		if err != nil {
			log.Println(err)
			fmt.Println("File not found")
			_, err = conn.Write([]byte(NotFoundResponse + contentLength("") + "\r\n"))
		} else {
			fmt.Println("File found")
			_, err = conn.Write([]byte((OkResponse + ContentTypeApplicationOctetStream + contentLength(string(fileContent)) + string(fileContent) + "\r\n")))
		}
	default:		_, err = conn.Write([]byte(NotFoundResponse + contentLength("") + "\r\n"))
			fmt.Println("Received GET for ", path)
		}
	}