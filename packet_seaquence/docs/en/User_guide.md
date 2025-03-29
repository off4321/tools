# User guide
This document provides a `user guide for the packet_sequence tool`, including `installation instructions`, `usage examples`, and `troubleshooting tips`.
It is intended for users who want to analyze packet sequences from tcpdump files and visualize them in a sequence chart format.

## Installation
To install the packet_sequence tool, follow these steps:
1. Clone the repository:
   ```bash
   git clone <URL>
   ```
2. Navigate to the project directory:
   ```bash
   cd packet_sequence
   cp -p conv_packet_sequence /bin
   ```

## Usage
To use the packet_sequence tool, run the following command:
```bash
packet_sequence <tcpdump_file.pcap> [options]
```
### Options
- `-h`, `--help`: Show help message and exit.
- `-o`, `--output`: Specify the output file name. Default is `output.mmd`.
- `-f`, `--format`: Specify the output format. Options are `mermaid` or `text`. Default is `mermaid`.
- `-v`, `--verbose`: Enable verbose output.
- `-version`: Show the version of the tool.
- `-i`, `--input`: Specify the input file name. Default is `input.pcap`.

## Troubleshooting
