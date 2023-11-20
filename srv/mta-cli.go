package srv

import (
	"bufio"
	"fmt"
	"net"
)

func cli() {
	conn, err := net.Dial("tcp", "localhost:25")
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	readResponse(reader)

	cmds := []string{
		"HELO localhost\r\n",
		"MAIL FROM: <sender@example.com>\r\n",
		"RCPT TO: <recipient@example.com>\r\n",
		"DATA\r\n",
		"Subject: Test Email\r\n\r\n",
		"This is a test email message.\r\n.\r\n",
		"QUIT\r\n",
	}

	for _, cmd := range cmds {
		sendCommand(writer, cmd)
		readResponse(reader)
	}
}

func sendCommand(writer *bufio.Writer, command string) {
	_, err := writer.WriteString(command)
	if err != nil {
		fmt.Println("Error writing command:", err)
		return
	}
	writer.Flush()
}

func readResponse(reader *bufio.Reader) {
	response, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}
	fmt.Print("Server:", response)
}
