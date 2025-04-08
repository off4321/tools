package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strings"
)

type Packet struct {
	Time        string
	Length      string
	Source      string
	Destination string
	TopProtocol string
	Stack       []string
	Info        string
}

func main() {
	// Parse command line arguments
	pcapFile := flag.String("file", "", "PCAP file to analyze")
	flag.Parse()

	if *pcapFile == "" {
		fmt.Println("Please specify a PCAP file using the -file flag")
		return
	}

	// Run tshark command
	cmd := exec.Command("tshark", 
		"-r", *pcapFile,
		"-T", "fields",
		"-e", "frame.time",
		"-e", "frame.len",
		"-e", "ip.src",
		"-e", "ip.dst",
		"-e", "frame.protocols",
		"-e", "_ws.col.Protocol",
		"-e", "_ws.col.Info",
		"-E", "header=n",
		"-E", "separator=|")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Error creating stdout pipe: %v\n", err)
		return
	}

	err = cmd.Start()
	if err != nil {
		fmt.Printf("Error starting tshark: %v\n", err)
		return
	}

	// Process the output
	reader := bufio.NewReader(stdout)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("Error reading stdout: %v\n", err)
			return
		}

		packet := parsePacket(line)
		printPacket(packet)
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Printf("Error waiting for tshark: %v\n", err)
		return
	}
}

func parsePacket(line string) Packet {
	fields := strings.Split(strings.TrimSpace(line), "|")
	
	// Default values in case some fields are missing
	packet := Packet{
		Time:        "Unknown",
		Length:      "0",
		Source:      "Unknown",
		Destination: "Unknown",
		TopProtocol: "Unknown",
		Stack:       []string{},
		Info:        "",
	}
	
	// Extract fields based on their position
	if len(fields) > 0 {
		packet.Time = fields[0]
	}
	if len(fields) > 1 {
		packet.Length = fields[1] + " bytes"
	}
	if len(fields) > 2 {
		packet.Source = fields[2]
	}
	if len(fields) > 3 {
		packet.Destination = fields[3]
	}
	if len(fields) > 4 {
		protocols := fields[4]
		protocolStack := strings.Split(protocols, ":")
		packet.Stack = protocolStack
		
		// Get top protocol (last in the stack)
		if len(protocolStack) > 0 {
			packet.TopProtocol = protocolStack[len(protocolStack)-1]
		}
	}
	if len(fields) > 5 {
		packet.TopProtocol = fields[5]
	}
	if len(fields) > 6 {
		packet.Info = fields[6]
	}
	
	return packet
}

func printPacket(packet Packet) {
	// Format the protocol stack for display
	var formattedStack string
	if len(packet.Stack) > 0 {
		formattedStack = strings.Join(packet.Stack, " > ")
	} else {
		formattedStack = "Unknown"
	}
	
	// Clean up the time string by removing nanoseconds precision indicators
	timeString := packet.Time
	r := regexp.MustCompile(`(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d+).*`)
	matches := r.FindStringSubmatch(timeString)
	if len(matches) > 1 {
		timeString = matches[1]
	}
	
	fmt.Printf("Time: %s UTC\n", timeString)
	fmt.Printf("Length: %s\n", packet.Length)
	fmt.Printf("Source: %s\n", packet.Source)
	fmt.Printf("Destination: %s\n", packet.Destination)
	fmt.Printf("Top Protocol: %s\n", packet.TopProtocol)
	fmt.Printf("Protocol Stack: %s\n", formattedStack)
	fmt.Printf("Info: %s\n\n", packet.Info)
}