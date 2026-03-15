package osint

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/EmmmmDeee/One-Liner-OSINT/pkg/config"
	"github.com/EmmmmDeee/One-Liner-OSINT/pkg/logger"
	"github.com/EmmmmDeee/One-Liner-OSINT/pkg/output"
	"github.com/EmmmmDeee/One-Liner-OSINT/pkg/utils"
)

// IPModule provides IP address OSINT capabilities
type IPModule struct {
	config *config.Config
	logger *logger.Logger
	client *utils.HTTPClient
}

// NewIPModule creates a new IP module
func NewIPModule(cfg *config.Config) *IPModule {
	return &IPModule{
		config: cfg,
		logger: logger.NewLogger(),
		client: utils.NewHTTPClient(cfg.Timeout, cfg.UserAgent),
	}
}

// Name returns the module name
func (m *IPModule) Name() string {
	return "ip"
}

// Description returns the module description
func (m *IPModule) Description() string {
	return "IP address geolocation and analysis"
}

// Execute performs IP OSINT
func (m *IPModule) Execute(ctx context.Context, target string) ([]output.Result, error) {
	var results []output.Result

	// Validate IP address
	ip := net.ParseIP(target)
	if ip == nil {
		return nil, fmt.Errorf("invalid IP address: %s", target)
	}

	// Basic IP info
	basicInfo := m.getBasicInfo(target)
	results = append(results, basicInfo)

	// Reverse DNS lookup
	reverseResult, err := m.reverseLookup(target)
	if err != nil {
		m.logger.Debugf("Reverse lookup failed: %v", err)
	} else {
		results = append(results, reverseResult)
	}

	// Geolocation (if API key available)
	if m.config.APIKeys.IPGeolocation != "" {
		geoResult, err := m.getGeolocation(ctx, target)
		if err != nil {
			m.logger.Warnf("Geolocation lookup failed: %v", err)
		} else {
			results = append(results, geoResult)
		}
	}

	return results, nil
}

// getBasicInfo returns basic IP information
func (m *IPModule) getBasicInfo(ip string) output.Result {
	parsedIP := net.ParseIP(ip)
	data := map[string]interface{}{
		"ip":         ip,
		"is_private": isPrivateIP(parsedIP),
		"is_ipv4":    parsedIP.To4() != nil,
		"is_ipv6":    parsedIP.To4() == nil && parsedIP.To16() != nil,
	}

	if parsedIP.To4() != nil {
		data["version"] = 4
	} else if parsedIP.To16() != nil {
		data["version"] = 6
	}

	return output.Result{
		Timestamp: time.Now(),
		Source:    "ip",
		Type:      "basic",
		Target:    ip,
		Data:      data,
	}
}

// reverseLookup performs reverse DNS lookup
func (m *IPModule) reverseLookup(ip string) (output.Result, error) {
	names, err := net.LookupAddr(ip)
	if err != nil {
		return output.Result{}, err
	}

	return output.Result{
		Timestamp: time.Now(),
		Source:    "ip",
		Type:      "reverse_dns",
		Target:    ip,
		Data: map[string]interface{}{
			"hostnames": names,
			"count":     len(names),
		},
	}, nil
}

// getGeolocation retrieves geolocation information
func (m *IPModule) getGeolocation(ctx context.Context, ip string) (output.Result, error) {
	// Use ipapi.co free API if no API key is configured
	url := fmt.Sprintf("https://ipapi.co/%s/json/", ip)

	resp, err := m.client.Get(ctx, url)
	if err != nil {
		return output.Result{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return output.Result{}, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return output.Result{}, err
	}

	var geoData map[string]interface{}
	if err := json.Unmarshal(body, &geoData); err != nil {
		return output.Result{}, err
	}

	return output.Result{
		Timestamp: time.Now(),
		Source:    "ipapi",
		Type:      "geolocation",
		Target:    ip,
		Data:      geoData,
	}, nil
}

// isPrivateIP checks if an IP is private
func isPrivateIP(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return true
	}

	privateBlocks := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"fc00::/7",
	}

	for _, block := range privateBlocks {
		_, subnet, _ := net.ParseCIDR(block)
		if subnet.Contains(ip) {
			return true
		}
	}

	return false
}

// GetIPFromDomain resolves domain to IP
func (m *IPModule) GetIPFromDomain(domain string) ([]string, error) {
	ips, err := net.LookupIP(domain)
	if err != nil {
		return nil, err
	}

	ipStrings := make([]string, 0, len(ips))
	for _, ip := range ips {
		ipStrings = append(ipStrings, ip.String())
	}

	return ipStrings, nil
}
