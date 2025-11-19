package network

import (
	"net"
	"reflect"
	"strings"
	"testing"
)

func TestParseCIDRs(t *testing.T) {
	_, expectedNet1, _ := net.ParseCIDR("192.168.1.0/24")
	_, expectedNet2, _ := net.ParseCIDR("10.0.0.0/8")
	_, expectedNet3, _ := net.ParseCIDR("2001:db8::/32")

	testCases := []struct {
		name          string
		records       [][]string
		expected      []*net.IPNet
		expectError   bool
		errorContains string
	}{
		{
			name: "CIDR format with header",
			records: [][]string{
				{"cidr"},
				{"192.168.1.0/24"},
				{"10.0.0.0/8"},
			},
			expected:    []*net.IPNet{expectedNet1, expectedNet2},
			expectError: false,
		},
		{
			name: "CIDR format without header",
			records: [][]string{
				{"192.168.1.0/24"},
				{"10.0.0.0/8"},
			},
			expected:    []*net.IPNet{expectedNet1, expectedNet2},
			expectError: false,
		},
		{
			name: "Network,Netmask format with header",
			records: [][]string{
				{"network", "netmask"},
				{"192.168.1.0", "255.255.255.0"},
				{"10.0.0.0", "255.0.0.0"},
			},
			expected:    []*net.IPNet{expectedNet1, expectedNet2},
			expectError: false,
		},
		{
			name: "IPv6 CIDR format",
			records: [][]string{
				{"2001:db8::/32"},
			},
			expected:    []*net.IPNet{expectedNet3},
			expectError: false,
		},
		{
			name: "Host address normalization",
			records: [][]string{
				{"192.168.1.123/24"},
			},
			expected:    []*net.IPNet{expectedNet1},
			expectError: false,
		},
		{
			name: "Invalid CIDR",
			records: [][]string{
				{"192.168.1.0/33"},
			},
			expectError:   true,
			errorContains: "invalid CIDR address",
		},
		{
			name: "Invalid Netmask format",
			records: [][]string{
				{"192.168.1.0", "255.255.255"},
			},
			expectError:   true,
			errorContains: "invalid netmask format",
		},
		{
			name: "Non-contiguous Netmask",
			records: [][]string{
				{"192.168.1.0", "255.255.0.255"},
			},
			expectError:   true,
			errorContains: "non-contiguous netmask",
		},
		{
			name: "Wrong number of columns",
			records: [][]string{
				{"192.168.1.0", "255.255.255.0", "extra"},
			},
			expectError:   true,
			errorContains: "unexpected number of columns",
		},
		{
			name:          "Empty records",
			records:       [][]string{},
			expectError:   true,
			errorContains: "no records to parse",
		},
		{
			name: "Header only",
			records: [][]string{
				{"network", "netmask"},
			},
			expectError:   true,
			errorContains: "no data rows to parse",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ParseCIDRs(tc.records)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected an error, but got nil")
				} else if !strings.Contains(err.Error(), tc.errorContains) {
					t.Errorf("Expected error to contain '%s', but got '%s'", tc.errorContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Did not expect an error, but got: %v", err)
				}
				if !reflect.DeepEqual(result, tc.expected) {
					t.Errorf("Expected %v, but got %v", tc.expected, result)
				}
			}
		})
	}
}
