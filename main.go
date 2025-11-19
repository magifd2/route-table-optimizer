package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"route-table-optimizer/internal/network"
)

func main() {
	// 1. Define and parse CLI flags
	inputFile := flag.String("i", "", "Path to the input CSV file. (required)")
	outputFile := flag.String("o", "", "Path to the output CSV file. (optional, defaults to stdout)")
	help := flag.Bool("h", false, "Display this help message")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s -i <input_file> [-o <output_file>]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "This tool optimizes a list of IP network routes by merging and removing redundant entries.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *help || len(os.Args) == 1 {
		flag.Usage()
		os.Exit(0)
	}
	if *inputFile == "" {
		fmt.Fprintf(os.Stderr, "Error: input file path is required.\n\n")
		flag.Usage()
		os.Exit(1)
	}
	
	log.SetFlags(0) // Remove timestamps from log output

	// 2. Read input file
	log.Printf("--- Reading input file: %s ---", *inputFile)
	records, err := network.ReadNetworksCSV(*inputFile)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// 3. Parse records into network objects
	log.Println("--- Parsing networks ---")
	nets, err := network.ParseCIDRs(records)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// 4. Run optimization steps
	log.Printf("Original number of networks: %d", len(nets))
	
	nets = network.Deduplicate(nets)
	log.Printf("After deduplication: %d", len(nets))

	// The order is important: aggregate first, then merge.
	nets = network.Aggregate(nets)
	log.Printf("After aggregation: %d", len(nets))

	nets = network.Merge(nets)
	log.Printf("After merging: %d", len(nets))

	// 5. Determine output writer
	writer := os.Stdout
	if *outputFile != "" {
		file, err := os.Create(*outputFile)
		if err != nil {
			log.Fatalf("Error creating output file: %v", err)
		}
		defer file.Close()
		writer = file
		log.Printf("--- Writing output to: %s ---", *outputFile)
	} else {
		log.Println("--- Writing output to stdout ---")
	}

	// 6. Write the optimized list to CSV
	if err := network.WriteNetworksCSV(writer, nets); err != nil {
		log.Fatalf("Error writing CSV: %v", err)
	}

	log.Println("--- Optimization complete ---")
}
