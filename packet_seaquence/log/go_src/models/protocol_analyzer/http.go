// Package protocol_analyzer provides HTTP protocol analyzer
package protocol_analyzer

import (
	"fmt"
)

// AnalyzeHTTP implements the ProtocolAnalyzer interface for HTTP
type AnalyzeHTTP struct {
	BaseProtocolAnalyzer
}

// NewAnalyzeHTTP creates a new HTTP analyzer
func NewAnalyzeHTTP(packet Packet) *AnalyzeHTTP {
	return &AnalyzeHTTP{
		BaseProtocolAnalyzer: NewBaseProtocolAnalyzer(packet),
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
	
	// Request method (GET, POST, etc.)
	if method, ok := a.Pkt.GetField("http", "request_method"); ok {
		httpInfo["method"] = method
	}
	
	// Host
	if host, ok := a.Pkt.GetField("http", "host"); ok {
		httpInfo["host"] = host
	}
	
	// Request URI
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
	
	line := getOrDefault(info["request_line"], "")
	method := getOrDefault(info["method"], "")
	uri := getOrDefault(info["uri"], "")
	
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
	
	method := getOrDefault(info["method"], "?")
	uri := getOrDefault(info["uri"], "?")
	
	return fmt.Sprintf("HTTP %s %s", method, uri)
}