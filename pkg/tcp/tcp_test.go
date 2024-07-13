package tcp

import (
	"bytes"
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTCPClient(port uint16) (*net.TCPConn, error) {
  serverAddr := fmt.Sprintf("127.0.0.1:%d", port)
  tcpAddr, err := net.ResolveTCPAddr("tcp", serverAddr)
  if err != nil {
    return nil, err
  }

  conn, err := net.DialTCP("tcp", nil, tcpAddr)
  if err != nil {
    return nil, err
  }

  return conn, nil
}

func TestTCPServer(t *testing.T) {
  port := uint16(5040)
  server, err := NewTcpServer(port)
  if err != nil {
    t.Fatalf("Error creating TCP server test: %s", err)
  }

  client, err := createTCPClient(uint16(5040))
  if err != nil {
    t.Fatalf("Error creating TCP client: %s", err)
  }

  client2, err := createTCPClient(uint16(5040))
  if err != nil {
    t.Fatalf("Error creating TCP client: %s", err)
  }

  cmd := TCPCommand{
    Command: "t",
    Data: "Hello World",
  }

  _, err = client.Write(cmd.Bytes())
  if err != nil {
    t.Fatalf("Error writing cmd to the client: %s", err)
  }

  c := <-server.FromSockets
  assert.Equal(t, c, cmd)

  cmd2 := TCPCommand{
    Command: "t",
    Data: "69:420",
  }

  server.Send(cmd2)

  clientCmd := commandParser(client)
  clientCmd2 := commandParser(client2)

  out := <-clientCmd
  out2 := <-clientCmd2

  assert.Equal(t, out, cmd2)
  assert.Equal(t, out2, cmd2)

  client.Close()

  server.Send(cmd)
  out2 = <-clientCmd2
  assert.Equal(t, out2, cmd)
}

func TestCommandParser(t *testing.T) {
  cmd := TCPCommand{
    Command: "t",
    Data: "Hello World",
  }

  cmd2 := TCPCommand{
    Command: "t",
    Data: "Goodbye",
  }

  b:= cmd.Bytes()
  b2 := cmd2.Bytes()

  reader := bytes.NewReader(append(b, b2...))

  parsedCmd := commandParser(reader)
  
  c := <-parsedCmd
  c2 := <-parsedCmd
  assert.Equal(t, c, cmd)
  assert.Equal(t, c2, cmd2)
}
