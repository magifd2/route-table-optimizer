# Route Table Optimizer

A command-line tool written in Go to optimize IP network route tables by removing redundancies.

## Features

-   **Deduplication**: Removes duplicate network entries.
-   **Aggregation**: Removes networks that are already covered by a larger, less specific network entry.
-   **Merging**: Combines adjacent networks into a single, larger supernet (e.g., `192.168.0.0/24` and `192.168.1.0/24` become `192.168.0.0/23`).

## Installation

### From source

1.  Ensure you have Go installed (version 1.18 or later is recommended).
2.  Clone this repository:
    ```sh
    git clone https://github.com/user/route-table-optimizer.git
    cd route-table-optimizer
    ```
3.  Build the executable:
    ```sh
    go build
    ```
    This will create a `route-table-optimizer` executable in the current directory.

## Usage

```sh
./route-table-optimizer -i <input_file.csv> [-o <output_file.csv>]
```

### Flags

-   `-i, --input`: **(Required)** Path to the input CSV file.
-   `-o, --output`: (Optional) Path to the output CSV file. If not provided, the result will be printed to standard output.
-   `-h, --help`: Display the help message.

### Input File Format

The input CSV file can be in one of two formats. The tool will automatically detect the format. Comments (lines starting with `#`) and empty lines are ignored.

**1. CIDR Format**

A single column containing networks in CIDR notation. A header is optional.

```csv
cidr
192.168.0.0/24
192.168.1.0/24
10.0.0.0/8
```

**2. Network/Netmask Format**

Two columns, `network` and `netmask`. A header is optional.

```csv
network,netmask
192.168.0.0,255.255.255.0
192.168.1.0,255.255.255.0
10.0.0.0,255.0.0.0
```

### Output File Format

The output is always in the `network,netmask` format with a header.

```csv
network,netmask
10.0.0.0,255.0.0.0
192.168.0.0,255.255.254.0
```

## Development

### Building

```sh
go build
```

### Testing

Run all unit tests:
```sh
go test ./...
```
