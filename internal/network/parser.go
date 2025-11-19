package network

import (
	"bytes"
	"fmt"
	"net"
	"strings"
)

// ParseCIDRs takes a 2D string slice (from a CSV) and parses it into a slice of *net.IPNet.
// It automatically detects the format:
// - ["ip/prefix"]
// - ["network", "netmask"]
func ParseCIDRs(records [][]string) ([]*net.IPNet, error) {
	var ipNets []*net.IPNet
	if len(records) == 0 {
		return nil, fmt.Errorf("no records to parse")
	}

	// Simple header detection: check if the first row contains "network", "ip", or "cidr"
	startRow := 0
	if len(records) > 0 {
		header := strings.ToLower(strings.Join(records[0], ","))
		if strings.Contains(header, "network") || strings.Contains(header, "ip") || strings.Contains(header, "cidr") {
			startRow = 1
		}
	}

	if len(records) <= startRow {
		return nil, fmt.Errorf("no data rows to parse")
	}

	for i, record := range records[startRow:] {
		rowNum := i + startRow + 1
		var cidrStr string

		// Trim whitespace from all fields in the record
		for j := range record {
			record[j] = strings.TrimSpace(record[j])
		}

		switch len(record) {
		case 1:
			// Format: "ip/prefix"
			cidrStr = record[0]
		case 2:
			// Format: "network", "netmask"
			ip := record[0]
			maskStr := record[1]

			maskIP := net.ParseIP(maskStr)
			if maskIP == nil {
				return nil, fmt.Errorf("row %d: invalid netmask format: %s", rowNum, maskStr)
			}
			
			var mask net.IPMask
			var bits int
			if maskIP.To4() != nil {
				mask = net.IPMask(maskIP.To4())
				_, bits = mask.Size()
			} else {
				mask = net.IPMask(maskIP.To16())
				_, bits = mask.Size()
			}

			ones, _ := mask.Size()
			
			// Validate that the mask is contiguous
			expectedMask := net.CIDRMask(ones, bits)
			if !bytes.Equal(expectedMask, mask) {
				return nil, fmt.Errorf("row %d: non-contiguous netmask: %s", rowNum, maskStr)
			}

			cidrStr = fmt.Sprintf("%s/%d", ip, ones)
		default:
			return nil, fmt.Errorf("row %d: unexpected number of columns: got %d, want 1 or 2", rowNum, len(record))
		}

		// Use ParseCIDR to get the canonical network address
		_, ipNet, err := net.ParseCIDR(cidrStr)
		if err != nil {
			// This handles cases like an IP not belonging to the subnet defined by the mask
			return nil, fmt.Errorf("row %d: failed to parse cidr '%s': %w", rowNum, cidrStr, err)
		}
		
		ipNets = append(ipNets, ipNet)
	}

	return ipNets, nil
}
