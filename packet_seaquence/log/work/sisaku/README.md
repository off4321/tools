# Packet Reader

A simple Go program that uses tshark to read packet captures and display them in a readable format.

## Requirements

- Go 1.16 or later
- tshark (Wireshark command-line interface)

## Usage

```bash
go run main.go -file /path/to/your/capture.pcap
```

## Output Format

The program outputs packet information in the following format:

```
Time: [Timestamp] UTC
Length: [Length] bytes
Source: [Source IP]
Destination: [Destination IP]
Top Protocol: [Protocol]
Protocol Stack: [eth > ethertype > ip > protocol]
Info: [Packet Info]
```

## Example

```
Time: May 13, 2004 10:17:25.216971000 UTC
Length: 54 bytes
Source: 65.208.228.223
Destination: 145.254.160.237
Top Protocol: tcp
Protocol Stack: eth > ethertype > ip > tcp
Info: 80 → 3372 [FIN, ACK] Seq=18365 Ack=480 Win=6432 Len=0
```