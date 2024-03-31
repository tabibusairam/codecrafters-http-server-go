package main
import (
	"bytes"
	"fmt"
	"net"
	"os"
)
func main() {
	fmt.Println("Logs from your program will appear here!")
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	c, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	requestBuffer := make([]byte, 1024)
	c.Read(requestBuffer)
	contents := bytes.Split(requestBuffer, []byte("\r\n"))
	path := bytes.SplitN(contents[0], []byte{' '}, 3)[1]
	if bytes.Equal(path, []byte{'/'}) {
		c.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	} else if bytes.HasPrefix(path, []byte("/echo/")) {
		c.Write([]byte("HTTP/1.1 200 OK\r\n"))
		c.Write([]byte("Content-Type: text/plain\r\n"))
		res := bytes.SplitN(path, []byte{'/'}, 3)[2]
		a := "Content-Length: " + fmt.Sprint(len(res)) + "\r\n\r\n"
		c.Write([]byte(a))
		c.Write(bytes.SplitN(path, []byte{'/'}, 3)[2])
	} else {
		c.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
	c.Close()
}