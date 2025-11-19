package network

import (
	"net"
	"reflect"
	"sort"
	"testing"
)

// MustParseCIDR is a helper to create a net from a CIDR string, panicking on error.
// This is useful for setting up test data.
func MustParseCIDR(s string) *net.IPNet {
	_, n, err := net.ParseCIDR(s)
	if err != nil {
		panic(`failed to parse CIDR "` + s + `": ` + err.Error())
	}
	return n
}

// sortNets sorts a slice of *net.IPNet by their string representation for consistent comparison.
func sortNets(nets []*net.IPNet) {
	sort.Slice(nets, func(i, j int) bool {
		return nets[i].String() < nets[j].String()
	})
}

func TestDeduplicate(t *testing.T) {
	testCases := []struct {
		name     string
		input    []*net.IPNet
		expected []*net.IPNet
	}{
		{
			name:     "No duplicates",
			input:    []*net.IPNet{MustParseCIDR("192.168.1.0/24"), MustParseCIDR("10.0.0.0/8")},
			expected: []*net.IPNet{MustParseCIDR("192.168.1.0/24"), MustParseCIDR("10.0.0.0/8")},
		},
		{
			name:     "With duplicates",
			input:    []*net.IPNet{MustParseCIDR("192.168.1.0/24"), MustParseCIDR("10.0.0.0/8"), MustParseCIDR("192.168.1.0/24")},
			expected: []*net.IPNet{MustParseCIDR("192.168.1.0/24"), MustParseCIDR("10.0.0.0/8")},
		},
		{
			name:     "Empty slice",
			input:    []*net.IPNet{},
			expected: []*net.IPNet{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Deduplicate(tc.input)
			sortNets(result)
			sortNets(tc.expected)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %v, but got %v", tc.expected, result)
			}
		})
	}
}

func TestAggregate(t *testing.T) {
	testCases := []struct {
		name     string
		input    []*net.IPNet
		expected []*net.IPNet
	}{
		{
			name:     "Simple aggregation",
			input:    []*net.IPNet{MustParseCIDR("10.0.0.0/8"), MustParseCIDR("10.1.0.0/16")},
			expected: []*net.IPNet{MustParseCIDR("10.0.0.0/8")},
		},
		{
			name:     "No aggregation needed",
			input:    []*net.IPNet{MustParseCIDR("192.168.0.0/24"), MustParseCIDR("192.168.1.0/24")},
			expected: []*net.IPNet{MustParseCIDR("192.168.0.0/24"), MustParseCIDR("192.168.1.0/24")},
		},
		{
			name:     "Multiple levels",
			input:    []*net.IPNet{MustParseCIDR("10.0.0.0/8"), MustParseCIDR("10.1.0.0/16"), MustParseCIDR("10.1.1.0/24")},
			expected: []*net.IPNet{MustParseCIDR("10.0.0.0/8")},
		},
		{
			name:     "Out of order input",
			input:    []*net.IPNet{MustParseCIDR("10.1.1.0/24"), MustParseCIDR("10.0.0.0/8")},
			expected: []*net.IPNet{MustParseCIDR("10.0.0.0/8")},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Aggregate(tc.input)
			sortNets(result)
			sortNets(tc.expected)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %v, but got %v", tc.expected, result)
			}
		})
	}
}

func TestMerge(t *testing.T) {
	testCases := []struct {
		name     string
		input    []*net.IPNet
		expected []*net.IPNet
	}{
		{
			name:     "Simple merge",
			input:    []*net.IPNet{MustParseCIDR("192.168.0.0/24"), MustParseCIDR("192.168.1.0/24")},
			expected: []*net.IPNet{MustParseCIDR("192.168.0.0/23")},
		},
		{
			name:     "Not adjacent",
			input:    []*net.IPNet{MustParseCIDR("192.168.0.0/24"), MustParseCIDR("192.168.2.0/24")},
			expected: []*net.IPNet{MustParseCIDR("192.168.0.0/24"), MustParseCIDR("192.168.2.0/24")},
		},
		{
			name:     "Different prefix lengths",
			input:    []*net.IPNet{MustParseCIDR("192.168.0.0/24"), MustParseCIDR("192.168.1.0/25")},
			expected: []*net.IPNet{MustParseCIDR("192.168.0.0/24"), MustParseCIDR("192.168.1.0/25")},
		},
		{
			name:     "Recursive merge",
			input:    []*net.IPNet{MustParseCIDR("192.168.0.0/24"), MustParseCIDR("192.168.1.0/24"), MustParseCIDR("192.168.2.0/24"), MustParseCIDR("192.168.3.0/24")},
			expected: []*net.IPNet{MustParseCIDR("192.168.0.0/22")},
		},
		{
			name:     "IPv6 merge",
			input:    []*net.IPNet{MustParseCIDR("2001:db8::/33"), MustParseCIDR("2001:db8:8000::/33")},
			expected: []*net.IPNet{MustParseCIDR("2001:db8::/32")},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Merge(tc.input)
			sortNets(result)
			sortNets(tc.expected)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %v, but got %v", tc.expected, result)
			}
		})
	}
}
