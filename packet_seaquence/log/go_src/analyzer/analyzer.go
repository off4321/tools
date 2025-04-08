// Package analyzer provides the main packet sequence analysis functionality
package analyzer

import (
	"fmt"
	"os"
	
	"app/go_src/extractor"
	"app/go_src/models"
	"app/go_src/protocol_analyzer"
	"app/go_src/reader"
	"app/go_src/writer"
)

// PacketSequenceAnalyzer provides functionality to analyze packet sequences
type PacketSequenceAnalyzer struct {
	PcapFile    string
	OutputFile  string
	MaxEntries  int
	Options     map[string]interface{}
}

// PacketAdapter adapts reader.Packet to protocol_analyzer.Packet
type PacketAdapter struct {
	packet *reader.Packet
}

// LayerExists implements protocol_analyzer.Packet.LayerExists
func (pa *PacketAdapter) LayerExists(layerType string) bool {
	return pa.packet.LayerExists(layerType)
}

// GetField implements protocol_analyzer.Packet.GetField
func (pa *PacketAdapter) GetField(layer string, field string) (interface{}, bool) {
	return pa.packet.GetField(layer, field)
}

// Layer implements protocol_analyzer.Packet.Layer
func (pa *PacketAdapter) Layer(layerType string) interface{} {
	// For tshark-based implementation, we'll check if the layer exists and return the Info map
	// for that layer if it does
	if !pa.packet.LayerExists(layerType) {
		return nil
	}
	
	// Return the layer info if available
	if info, ok := pa.packet.Info[layerType]; ok {
		return info
	}
	
	// Otherwise just return an empty map to indicate the layer exists but has no details
	return map[string]interface{}{}
}

// GetProtocolName implements protocol_analyzer.Packet.GetProtocolName
func (pa *PacketAdapter) GetProtocolName() string {
	return pa.packet.GetProtocolName()
}

// GetTimestamp returns the packet timestamp
func (pa *PacketAdapter) GetTimestamp() string {
	return pa.packet.Timestamp.String()
}

// GetSrcIP returns the source IP
func (pa *PacketAdapter) GetSrcIP() string {
	return pa.packet.SrcIP
}

// GetDstIP returns the destination IP
func (pa *PacketAdapter) GetDstIP() string {
	return pa.packet.DstIP
}

// NewPacketSequenceAnalyzer creates a new PacketSequenceAnalyzer
func NewPacketSequenceAnalyzer(pcapFile, outputFile string, maxEntries int, options map[string]interface{}) (*PacketSequenceAnalyzer, error) {
	// Check if file exists
	if _, err := os.Stat(pcapFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("PCAPファイルが見つかりません: %s", pcapFile)
	}
	
	// If options is nil, initialize it to an empty map
	if options == nil {
		options = make(map[string]interface{})
	}
	
	return &PacketSequenceAnalyzer{
		PcapFile:    pcapFile,
		OutputFile:  outputFile,
		MaxEntries:  maxEntries,
		Options:     options,
	}, nil
}

// convertPackets converts []reader.Packet to []protocol_analyzer.Packet
func convertPackets(readerPackets []reader.Packet) []protocol_analyzer.Packet {
	protocolPackets := make([]protocol_analyzer.Packet, len(readerPackets))
	for i := range readerPackets {
		 // Create a new adapter for each packet
		protocolPackets[i] = &PacketAdapter{packet: &readerPackets[i]}
	}
	return protocolPackets
}

// Analyze performs the packet sequence analysis
func (a *PacketSequenceAnalyzer) Analyze() error {
	// 1. Read the PCAP file
	pcapReader := reader.NewPCAPReader(a.PcapFile)
	err := pcapReader.Open()
	if err != nil {
		return fmt.Errorf("PCAPファイルオープン中にエラーが発生しました: %w", err)
	}
	defer pcapReader.Close()
	
	// verbose flag from options
	verbose := false
	if v, ok := a.Options["verbose"].(bool); ok {
		verbose = v
	}
	
	// Read all packets (maxPackets=0 means no limit)
	packets, err := pcapReader.Read(0, verbose)
	if err != nil {
		return fmt.Errorf("PCAPファイル読み込み中にエラーが発生しました: %w", err)
	}
	
	// 2. Extract data from the packets
	// Convert packets to protocol_analyzer.Packet type
	protocolPackets := convertPackets(packets)
	dataExtractor := extractor.NewExtractData(protocolPackets, len(protocolPackets)) // Extract all packets
	packetData := dataExtractor.Extract()
	
	// Debug output if verbose is enabled
	if verbose, ok := a.Options["verbose"].(bool); ok && verbose {
		a.debugPacketSequenceData(packetData)
	}
	
	// 3. Get all packets and divide them into chunks
	allPackets := packetData.GetPackets()
	totalPackets := len(allPackets)
	
	// Calculate chunk size
	chunkSize := a.MaxEntries
	if chunkSize <= 0 {
		chunkSize = totalPackets
	}
	
	// Divide packets into chunks
	chunks := make([]*models.PacketSequenceData, 0)
	for i := 0; i < totalPackets; i += chunkSize {
		// Calculate end index for this chunk
		end := i + chunkSize
		if end > totalPackets {
			end = totalPackets
		}
		
		// Create a new PacketSequenceData for this chunk
		chunkData := models.NewPacketSequenceData()
		for j := i; j < end; j++ {
			chunkData.AddPacket(allPackets[j])
		}
		
		chunks = append(chunks, chunkData)
	}
	
	// 4. Generate Mermaid diagrams for each chunk
	mermaidContents := make([]string, 0, len(chunks))
	for _, chunkData := range chunks {
		mermaidWriter := writer.NewWriteMermaid(chunkData)
		mermaidContent := mermaidWriter.Generate()
		mermaidContents = append(mermaidContents, mermaidContent)
	}
	
	// 5. Write the Mermaid diagrams to files
	fileWriter := writer.NewWriteFile(a.OutputFile)
	
	if len(mermaidContents) == 1 {
		// If there's only one chunk, write it directly
		if err := fileWriter.Write(mermaidContents[0]); err != nil {
			return fmt.Errorf("ファイル出力中にエラーが発生しました: %w", err)
		}
		fmt.Printf("解析が完了しました。出力ファイル: %s\n", a.OutputFile)
	} else {
		// If there are multiple chunks, write them to separate files
		outputFiles, err := fileWriter.WriteSplit(mermaidContents)
		if err != nil {
			return fmt.Errorf("ファイル出力中にエラーが発生しました: %w", err)
		}
		fmt.Printf("解析が完了しました。%d個のファイルに分割出力しました。\n", len(outputFiles))
		for i, file := range outputFiles {
			fmt.Printf("  %d. %s\n", i+1, file)
		}
	}
	
	return nil
}

// debugPacketSequenceData outputs debug information about the PacketSequenceData
func (a *PacketSequenceAnalyzer) debugPacketSequenceData(data *models.PacketSequenceData) {
	fmt.Println("\n===== PacketSequenceData DEBUG =====")
	
	// Get packets and display basic info
	packets := data.GetPackets()
	fmt.Printf("Total packets: %d\n", len(packets))
	
	// Display details of the first few packets
	maxPacketsToShow := 3
	if len(packets) < maxPacketsToShow {
		maxPacketsToShow = len(packets)
	}
	
	for i := 0; i < maxPacketsToShow; i++ {
		packet := packets[i]
		fmt.Printf("\nPacket %d:\n", i+1)
		fmt.Printf("  Protocol: %s\n", packet.Protocol)
		fmt.Printf("  ProtocolNameFromPacket: %s\n", packet.ProtocolNameFromPacket)
		fmt.Printf("  Src: %s\n", packet.Source)
		fmt.Printf("  Dst: %s\n", packet.Dest)
		fmt.Printf("  Time: %v\n", packet.Time)
		
		// Display info keys
		infoKeys := make([]string, 0, len(packet.Info))
		for key := range packet.Info {
			infoKeys = append(infoKeys, key)
		}
		fmt.Printf("  Info keys: %v\n", infoKeys)
		
		// Display protocol-specific info
		for key, value := range packet.Info {
			fmt.Printf("  %s: %v\n", key, value)
		}
	}
	
	fmt.Println("=================================\n")
	
	// Debug the WriteMermaid._build_message method
	fmt.Println("\n===== WriteMermaid DEBUG =====")
	mermaidWriter := writer.NewWriteMermaid(models.NewPacketSequenceData())
	
	for i := 0; i < maxPacketsToShow; i++ {
		if i >= len(packets) {
			break
		}
		
		packet := packets[i]
		message := mermaidWriter.BuildMessage(packet)
		fmt.Printf("Packet %d message: '%s'\n", i+1, message)
	}
	
	fmt.Println("=================================\n")
}