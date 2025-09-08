package ftp

import (
	"net"
	"net/textproto"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPASVZeroIPIntegration tests the complete PASV flow when server returns 0,0,0,0
func TestPASVZeroIPIntegration(t *testing.T) {
	// Create a mock server that returns 0,0,0,0 in PASV
	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer l.Close()

	serverAddr := l.Addr().String()

	// Start a simple FTP server that returns 0,0,0,0 in PASV
	go func() {
		conn, err := l.Accept()
		if err != nil {
			return
		}
		defer conn.Close()

		proto := textproto.NewConn(conn)
		defer proto.Close()

		// Send welcome message
		proto.PrintfLine("220 Test FTP Server")

		// Handle commands
		for {
			line, err := proto.ReadLine()
			if err != nil {
				break
			}

			parts := strings.Fields(line)
			if len(parts) == 0 {
				continue
			}

			cmd := strings.ToUpper(parts[0])
			switch cmd {
			case "USER":
				proto.PrintfLine("331 User name okay, need password")
			case "PASS":
				proto.PrintfLine("230 User logged in, proceed")
			case "FEAT":
				proto.PrintfLine("211-Features:\r\n PASV\r\n EPSV\r\n211 End")
			case "TYPE":
				proto.PrintfLine("200 Type set")
			case "OPTS":
				proto.PrintfLine("200 Command okay")
			case "PASV":
				// Return 0,0,0,0 as IP with a dummy port
				proto.PrintfLine("227 Entering Passive Mode (0,0,0,0,20,21).")
			case "QUIT":
				proto.PrintfLine("221 Goodbye")
				return
			default:
				proto.PrintfLine("500 Unknown command")
			}
		}
	}()

	// Connect to our test server
	c, err := Dial(serverAddr)
	require.NoError(t, err)
	defer c.Quit()

	err = c.Login("test", "test")
	require.NoError(t, err)

	// Disable EPSV to force PASV
	c.options.disableEPSV = true

	// Test getDataConnPort which internally calls pasv()
	host, port, err := c.getDataConnPort()
	require.NoError(t, err)

	// Should return 127.0.0.1 (the connection IP) instead of 0.0.0.0
	assert.Equal(t, "127.0.0.1", host, "Should use connection IP when PASV returns 0,0,0,0")
	assert.Equal(t, 5141, port, "Should return correct port (20*256 + 21)")

	t.Logf("PASV with 0,0,0,0 correctly returned host=%s, port=%d", host, port)
}
