# Application Specification: Route Table Optimizer

This document outlines the detailed specifications for the Go-based CLI tool, "Route Table Optimizer".

## 1. Overview

A Go-based CLI tool that optimizes a list of IP networks by performing deduplication, aggregation of contained networks, and merging of adjacent networks.

## 2. Input/Output

-   **Input (`--input` or `-i`):**
    -   Format: CSV
    -   Headers: The tool supports two formats, which it will automatically detect:
        1.  `network,netmask` (e.g., `192.168.0.0,255.255.255.0`)
        2.  A single column with CIDR notation `ip/prefix` (e.g., `192.168.0.0/24`)
    -   Character Encoding: UTF-8

-   **Output (`--output` or `-o`):**
    -   Format: CSV
    -   Header: `network,netmask`
    -   Character Encoding: UTF-8
    -   Standard Output: If `-o` is not specified, the results will be printed to standard output.

## 3. CLI Specification

`route-optimizer -i <input_file_path> -o <output_file_path>`

-   **Flags:**
    -   `--input, -i`: (Required) Path to the input CSV file.
    -   `--output, -o`: (Optional) Path to the output CSV file. Defaults to standard output if omitted.
    -   `--help, -h`: Displays the help message.

## 4. Optimization Logic

Internally, all network entries are handled uniformly in CIDR format using `net.IPNet`.

-   **Processing Steps:**
    1.  **Read and Parse:** Read the input CSV and convert each row into a slice of `net.IPNet` objects. Report an error for any malformed rows.
    2.  **Deduplication:** Remove identical network entries (e.g., `192.168.0.0/24` and `192.168.0.0/24`).
    3.  **Aggregate Contained Networks:**
        -   Sort the list of networks in ascending order of prefix length (from larger networks to smaller ones).
        -   Iterate through the list. If a network is found to be completely contained within a preceding, larger network, remove the smaller (more specific) network from the list.
        -   Example: If `10.0.0.0/8` exists, `10.1.2.0/24` will be removed.
    4.  **Merge Adjacent Networks:**
        -   Repeatedly perform the following process until no more merges occur.
        -   Sort the list to efficiently find potentially adjacent pairs.
        -   Merge two networks if they meet both of the following conditions:
            -   They have the same prefix length.
            -   They are bitwise-adjacent and can be represented by a single, larger network with a prefix length one bit shorter.
        -   Example: `192.168.0.0/24` and `192.168.1.0/24` merge into `192.168.0.0/23`.
    5.  **Output:** Write the final list of networks to a CSV file in the specified format.

## 5. Error Handling

-   File read/write errors.
-   CSV parsing errors.
-   Invalid IP address or netmask format errors.
-   In case of any of the above errors, print a detailed error message to standard error and exit the program with a non-zero exit code.
