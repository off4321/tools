// Package protocol_analyzer provides HTTP protocol analyzer
package protocol_analyzer

import (
	"fmt"
)

// AnalyzeHTTP implements the ProtocolAnalyzer interface for HTTP
type AnalyzeHTTP struct {
	Pkt Packet
}

// NewAnalyzeHTTP creates a new HTTP analyzer
func NewAnalyzeHTTP(packet Packet) *AnalyzeHTTP {
	return &AnalyzeHTTP{
		Pkt: packet,
	}
}

// Analyze extracts HTTP information from a packet
func (a *AnalyzeHTTP) Analyze() map[string]interface{} {
	if !a.Pkt.LayerExists("http") {
		return nil
	}

	httpInfo := make(map[string]interface{})
	
	// Request line
	if requestLine, ok := a.Pkt.GetField("http", "request_line"); ok {
		httpInfo["request_line"] = requestLine
	}
	
	// Method (GET, POST, etc.)
	if method, ok := a.Pkt.GetField("http", "request_method"); ok {
		httpInfo["method"] = method
	}
	
	// Host
	if host, ok := a.Pkt.GetField("http", "host"); ok {
		httpInfo["host"] = host
	}
	
	// URI
	if uri, ok := a.Pkt.GetField("http", "request_uri"); ok {
		httpInfo["uri"] = uri
	}
	
	return httpInfo
}

// GetDisplayInfo returns formatted HTTP information
func (a *AnalyzeHTTP) GetDisplayInfo() string {
	info := a.Analyze()
	if info == nil {
		return "HTTP情報なし"
	}
	
	line, _ := info["request_line"].(string)
	method, _ := info["method"].(string)
	uri, _ := info["uri"].(string)
	
	if method != "" && uri != "" {
		return fmt.Sprintf("HTTP %s %s", method, uri)
	}
	if line != "" {
		return fmt.Sprintf("HTTP %s", line)
	}
	return "HTTP"
}

// GetSummary returns a summary of HTTP information
func (a *AnalyzeHTTP) GetSummary() string {
	info := a.Analyze()
	if info == nil {
		return "HTTP"
	}
	
	method, _ := info["method"].(string)
	if method == "" {
		method = "?"
	}
	
	uri, _ := info["uri"].(string)
	if uri == "" {
		uri = "?"
	}
	
	return fmt.Sprintf("HTTP %s %s", method, uri)
}