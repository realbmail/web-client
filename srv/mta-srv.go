package srv

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func srv() {
	listener, err := net.Listen("tcp", ":25")
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()

	fmt.Println("SMTP Server listening on port 25")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	writer.WriteString("220 Welcome to My SMTP Server\r\n")
	writer.Flush()

	var from, to string
	var data strings.Builder
	state := "init"

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading:", err)
			return
		}

		command := strings.ToUpper(strings.TrimSpace(line))
		fmt.Println("Command:", command)

		switch {
		case state == "init" && strings.HasPrefix(command, "HELO"):
			writer.WriteString("250 Hello\r\n")
			writer.Flush()
		case state == "init" && strings.HasPrefix(command, "MAIL FROM"):
			from = strings.TrimPrefix(command, "MAIL FROM:")
			writer.WriteString("250 OK\r\n")
			writer.Flush()
		case state == "init" && strings.HasPrefix(command, "RCPT TO"):
			to = strings.TrimPrefix(command, "RCPT TO:")
			writer.WriteString("250 OK\r\n")
			writer.Flush()
		case state == "init" && strings.HasPrefix(command, "DATA"):
			writer.WriteString("354 End data with <CRLF>.<CRLF>\r\n")
			writer.Flush()
			state = "data"
		case state == "data" && command == ".":
			writer.WriteString("250 Message accepted for delivery\r\n")
			writer.Flush()
			fmt.Println("From:", from)
			fmt.Println("To:", to)
			fmt.Println("Data:", data.String())
			return
		case strings.ToUpper(strings.TrimSpace(command)) == "QUIT":
			writer.WriteString("221 Bye\r\n")
			writer.Flush()
			return
		default:
			writer.WriteString("500 Syntax error, command unrecognized\r\n")
			writer.Flush()
		}
	}
}
