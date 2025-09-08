package ftp

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestBogusDataIPWithUnspecified tests that 0.0.0.0 and :: are treated as bogus
func TestBogusDataIPWithUnspecified(t *testing.T) {
	tests := []struct {
		name     string
		cmdIP    string
		dataIP   string
		expected bool
	}{
		{
			name:     "0.0.0.0 should be bogus",
			cmdIP:    "127.0.0.1",
			dataIP:   "0.0.0.0",
			expected: true,
		},
		{
			name:     "IPv6 unspecified should be bogus",
			cmdIP:    "::1",
			dataIP:   "::",
			expected: true,
		},
		{
			name:     "0.0.0.0 from any IP should be bogus",
			cmdIP:    "192.168.1.1",
			dataIP:   "0.0.0.0",
			expected: true,
		},
		{
			name:     "Same private IP should not be bogus",
			cmdIP:    "192.168.1.1",
			dataIP:   "192.168.1.2",
			expected: false,
		},
		{
			name:     "Same loopback IP should not be bogus",
			cmdIP:    "127.0.0.1",
			dataIP:   "127.0.0.1",
			expected: false,
		},
		{
			name:     "Multicast should be bogus",
			cmdIP:    "127.0.0.1",
			dataIP:   "224.0.0.1",
			expected: true,
		},
		{
			name:     "Different network types should be bogus",
			cmdIP:    "127.0.0.1",   // loopback
			dataIP:   "192.168.1.1", // private
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmdIP := net.ParseIP(tt.cmdIP)
			dataIP := net.ParseIP(tt.dataIP)

			result := isBogusDataIP(cmdIP, dataIP)
			assert.Equal(t, tt.expected, result,
				"isBogusDataIP(%s, %s) should return %v", tt.cmdIP, tt.dataIP, tt.expected)
		})
	}
}
