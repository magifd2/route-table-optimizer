# Route Table Optimizer

A command-line tool written in Go to optimize IP network route tables by removing redundancies.

## Features

-   **Deduplication**: Removes duplicate network entries.
-   **Aggregation**: Removes networks that are already covered by a larger, less specific network entry.
-   **Merging**: Combines adjacent networks into a single, larger supernet (e.g., `192.168.0.0/24` and `192.168.1.0/24` become `192.168.0.0/23`).

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
