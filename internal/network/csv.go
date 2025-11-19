package network

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

// ReadNetworksCSV reads a CSV file from the given path, filters out comments,
// and returns its content as a 2D string slice.
func ReadNetworksCSV(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open input file: %w", err)
	}
	defer file.Close()

	// Read line by line to filter out comments and empty lines
	var filteredContent strings.Builder
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "#") && len(line) > 0 {
			filteredContent.WriteString(line + "\n")
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan input file: %w", err)
	}

	// If after filtering, the content is empty, return an error.
	if filteredContent.Len() == 0 {
		return nil, fmt.Errorf("input file contains no valid data")
	}

	reader := csv.NewReader(strings.NewReader(filteredContent.String()))
	reader.FieldsPerRecord = -1 // Let the parser handle column validation

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read csv data: %w", err)
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("input file is empty or contains only comments")
	}

	return records, nil
}

// WriteNetworksCSV writes a slice of *net.IPNet to the specified writer.
// It writes a header row and then each network in "network,netmask" format.
func WriteNetworksCSV(writer io.Writer, networks []*net.IPNet) error {
	w := csv.NewWriter(writer)
	defer w.Flush()

	// Write header
	header := []string{"network", "netmask"}
	if err := w.Write(header); err != nil {
		return fmt.Errorf("failed to write csv header: %w", err)
	}

	// Write data rows
	for _, n := range networks {
		ip := n.IP.String()

		// Convert prefix length back to a netmask
		mask := net.IP(n.Mask).String()

		record := []string{ip, mask}
		if err := w.Write(record); err != nil {
			return fmt.Errorf("failed to write csv record for network %s: %w", n.String(), err)
		}
	}

	return w.Error()
}
