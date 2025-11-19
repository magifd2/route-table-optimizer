# Route Table Optimizer

A command-line tool written in Go to optimize IP network route tables by removing redundancies.

## Features

-   **Deduplication**: Removes duplicate network entries.
-   **Aggregation**: Removes networks that are already covered by a larger, less specific network entry.
-   **Merging**: Combines adjacent networks into a single, larger supernet (e.g., `192.168.0.0/24` and `192.168.1.0/24` become `192.168.0.0/23`).

## Usage

Once you have built the executable using `make build`, you can run the tool from your terminal.

```sh
./route-table-optimizer -i <input_file.csv> [-o <output_file.csv>]
```

### Flags

-   `-i`: **(Required)** Path to the input CSV file containing the list of networks.
-   `-o`: (Optional) Path to the output CSV file. If this flag is not provided, the optimized list will be printed to the standard output (your terminal screen).
-   `-h`: Displays the help message.

### Input File Format

The tool accepts a CSV file containing a list of IP networks. Comments (lines starting with `#`) and empty lines are ignored. The tool automatically detects whether the input is in CIDR format or Network/Netmask format.

**1. CIDR Format**

A single column with a header (e.g., `cidr`).

```csv
# test_input.csv
cidr
192.168.0.0/24
192.168.1.0/24
10.0.0.0/8
10.0.1.0/24
```

**2. Network/Netmask Format**

Two columns with headers (e.g., `network,netmask`).

```csv
# test_input.csv
network,netmask
192.168.0.0,255.255.255.0
192.168.1.0,255.255.255.0
10.0.0.0,255.0.0.0
10.0.1.0,255.255.255.0
```

### Output File Format

The output is always a CSV file in `network,netmask` format with a header.

For the example input above, the optimized output would be:

```csv
# test_output.csv
network,netmask
10.0.0.0,255.0.0.0
192.168.0.0,255.255.254.0
```

## Building and Testing

This project uses a `Makefile` to simplify the build and test process.

### Prerequisites

-   Go (version 1.18 or later)
-   `make`

### Build the Executable

To build the executable for your current operating system, run:

```sh
make build
```

This will create a `route-table-optimizer` executable in the project root.

For cross-compilation (Linux, Windows, macOS), run:

```sh
make release
```

The executables will be placed in the `build/` directory.

### Run Tests

To run all unit tests, use:

```sh
make test
```

### Clean Build Artifacts

To remove all build artifacts and the compiled executables, run:

```sh
make clean
```
