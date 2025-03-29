# Overview
This file provides an `overview of the packet_sequence tool`, including `its purpose` and `usage`.
It is intended for users who want to analyze packet sequences from tcpdump files and visualize them in a sequence chart format.

## Purpose
The packet_sequence tool is designed to convert tcpdump files into sequence charts.
This tool makes it easier to analyze the packet sequence of a tcpdump file, providing a visual representation of the data flow.
This tool converts a `tcpdump file(.pcap)` into a `mermaid.js` sequence chart.

## Technical Details
 - **Implementation Language**: Python x.x
 - **Dependencies**: 
   - `scapy` for packet manipulation and analysis.
   - `mermaid` for generating sequence charts.

## Usage
```bash
packet_sequence <tcpdump_file.pcap> [options]
```