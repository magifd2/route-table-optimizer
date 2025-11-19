package network

import (
	"bytes"
	"net"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestReadNetworksCSV(t *testing.T) {
	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "test_*.csv")
	if err != nil {
		t.Fatal("Failed to create temp file:", err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	content := "cidr\n192.168.1.0/24\n10.0.0.0/8"
	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal("Failed to write to temp file:", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal("Failed to close temp file:", err)
	}

	records, err := ReadNetworksCSV(tmpfile.Name())
	if err != nil {
		t.Fatalf("ReadNetworksCSV failed: %v", err)
	}

	expected := [][]string{
		{"cidr"},
		{"192.168.1.0/24"},
		{"10.0.0.0/8"},
	}
	if !reflect.DeepEqual(records, expected) {
		t.Errorf("Expected %v, got %v", expected, records)
	}
}

func TestReadNetworksCSV_FileNotExist(t *testing.T) {
	_, err := ReadNetworksCSV("nonexistent.file")
	if err == nil {
		t.Error("Expected an error for non-existent file, but got nil")
	}
}

func TestWriteNetworksCSV(t *testing.T) {
	_, net1, _ := net.ParseCIDR("192.168.0.0/23")
	_, net2, _ := net.ParseCIDR("10.0.0.0/8")
	networks := []*net.IPNet{net1, net2}

	var buf bytes.Buffer
	err := WriteNetworksCSV(&buf, networks)
	if err != nil {
		t.Fatalf("WriteNetworksCSV failed: %v", err)
	}

	// The CSV writer might add a trailing newline, so we check with TrimSpace
	expected := "network,netmask\n" +
		"192.168.0.0,255.255.254.0\n" +
		"10.0.0.0,255.0.0.0"

	if strings.TrimSpace(buf.String()) != strings.TrimSpace(expected) {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, buf.String())
	}
}
