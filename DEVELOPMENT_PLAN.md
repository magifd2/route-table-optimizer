# Development Plan: Route Table Optimizer

This document outlines the development plan for the "Route Table Optimizer" tool.

## 1. Milestones

-   **M1: Project Foundation Setup**
    -   [ ] Initialize Go module (`go mod init`).
    -   [ ] Implement the basic CLI skeleton with flags (`-i`, `-o`, `-h`) using the `flag` package in `main.go`.
    -   [ ] Decide on the project structure (e.g., using an `internal` package).

-   **M2: CSV I/O and Data Structures**
    -   [ ] Create an `internal/network` package.
    -   [ ] Implement functionality to read a CSV file and obtain a string slice in either `network,netmask` or `ip/prefix` format.
    -   [ ] Implement a parser to convert the string slices into a slice of `net.IPNet` objects.
    -   [ ] Implement functionality to write a slice of `net.IPNet` objects to a CSV file.
    -   [ ] Create unit tests for the above functions.

-   **M3: Implementation of Optimization Logic**
    -   [ ] **Deduplication:** Implement a function to remove duplicates from a slice of `net.IPNet`.
    -   [ ] **Aggregation:** Implement a function to remove contained (sub-) networks.
    -   [ ] **Merging:** Implement a function to find and merge adjacent networks.
    -   [ ] Implement unit tests for each optimization function, paying special attention to edge cases.

-   **M4: Integration and Finalization**
    -   [ ] Integrate the modules from M2 and M3 into `main.go` to complete the overall processing flow.
    -   [ ] Create end-to-end tests that verify the entire workflow from input to output using actual CSV files.
    -   [ ] Write `README.md` with instructions on how to use the tool, its specifications, and how to build it.
    -   [ ] Prepare a build script (e.g., Makefile) for cross-compilation.

## 2. Testing Strategy

-   Utilize Go's standard `testing` package.
-   Write unit tests for each logical component (parser, deduplication, aggregation, merging) covering normal cases, error cases, and boundary conditions.
-   Ensure CLI functionality with integration tests using real CSV files.

## 3. Deliverables

-   `route-optimizer` executable binary.
-   Complete source code.
-   `SPECIFICATION.md` (This document).
-   `DEVELOPMENT_PLAN.md` (This document).
-   `README.md` (Usage documentation).
