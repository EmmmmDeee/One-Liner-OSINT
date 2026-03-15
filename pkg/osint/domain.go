package osint

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/EmmmmDeee/One-Liner-OSINT/pkg/config"
	"github.com/EmmmmDeee/One-Liner-OSINT/pkg/logger"
	"github.com/EmmmmDeee/One-Liner-OSINT/pkg/output"
	"github.com/EmmmmDeee/One-Liner-OSINT/pkg/utils"
)

// DomainModule provides domain OSINT capabilities
type DomainModule struct {
	config *config.Config
	logger *logger.Logger
	client *utils.HTTPClient
}

// NewDomainModule creates a new domain module
func NewDomainModule(cfg *config.Config) *DomainModule {
	return &DomainModule{
		config: cfg,
		logger: logger.NewLogger(),
		client: utils.NewHTTPClient(cfg.Timeout, cfg.UserAgent),
	}
}

// Name returns the module name
func (m *DomainModule) Name() string {
	return "domain"
}

// Description returns the module description
func (m *DomainModule) Description() string {
	return "Domain reconnaissance and DNS enumeration"
}

// Execute performs domain OSINT
func (m *DomainModule) Execute(ctx context.Context, target string) ([]output.Result, error) {
	var results []output.Result

	// DNS A records
	aRecords, err := m.lookupA(target)
	if err != nil {
		m.logger.Debugf("A record lookup failed: %v", err)
	} else {
		results = append(results, aRecords)
	}

	// DNS MX records
	mxRecords, err := m.lookupMX(target)
	if err != nil {
		m.logger.Debugf("MX record lookup failed: %v", err)
	} else {
		results = append(results, mxRecords)
	}

	// DNS NS records
	nsRecords, err := m.lookupNS(target)
	if err != nil {
		m.logger.Debugf("NS record lookup failed: %v", err)
	} else {
		results = append(results, nsRecords)
	}

	// DNS TXT records
	txtRecords, err := m.lookupTXT(target)
	if err != nil {
		m.logger.Debugf("TXT record lookup failed: %v", err)
	} else {
		results = append(results, txtRecords)
	}

	// CNAME records
	cnameRecord, err := m.lookupCNAME(target)
	if err != nil {
		m.logger.Debugf("CNAME record lookup failed: %v", err)
	} else {
		results = append(results, cnameRecord)
	}

	return results, nil
}

// lookupA performs A record lookup
func (m *DomainModule) lookupA(domain string) (output.Result, error) {
	ips, err := net.LookupIP(domain)
	if err != nil {
		return output.Result{}, err
	}

	ipStrings := make([]string, 0, len(ips))
	for _, ip := range ips {
		if ip.To4() != nil {
			ipStrings = append(ipStrings, ip.String())
		}
	}

	return output.Result{
		Timestamp: time.Now(),
		Source:    "dns",
		Type:      "A",
		Target:    domain,
		Data: map[string]interface{}{
			"records": ipStrings,
			"count":   len(ipStrings),
		},
	}, nil
}

// lookupMX performs MX record lookup
func (m *DomainModule) lookupMX(domain string) (output.Result, error) {
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		return output.Result{}, err
	}

	records := make([]map[string]interface{}, 0, len(mxRecords))
	for _, mx := range mxRecords {
		records = append(records, map[string]interface{}{
			"host":     mx.Host,
			"priority": mx.Pref,
		})
	}

	return output.Result{
		Timestamp: time.Now(),
		Source:    "dns",
		Type:      "MX",
		Target:    domain,
		Data: map[string]interface{}{
			"records": records,
			"count":   len(records),
		},
	}, nil
}

// lookupNS performs NS record lookup
func (m *DomainModule) lookupNS(domain string) (output.Result, error) {
	nsRecords, err := net.LookupNS(domain)
	if err != nil {
		return output.Result{}, err
	}

	nameservers := make([]string, 0, len(nsRecords))
	for _, ns := range nsRecords {
		nameservers = append(nameservers, ns.Host)
	}

	return output.Result{
		Timestamp: time.Now(),
		Source:    "dns",
		Type:      "NS",
		Target:    domain,
		Data: map[string]interface{}{
			"nameservers": nameservers,
			"count":       len(nameservers),
		},
	}, nil
}

// lookupTXT performs TXT record lookup
func (m *DomainModule) lookupTXT(domain string) (output.Result, error) {
	txtRecords, err := net.LookupTXT(domain)
	if err != nil {
		return output.Result{}, err
	}

	// Categorize TXT records
	categorized := map[string][]string{
		"spf":         []string{},
		"dmarc":       []string{},
		"dkim":        []string{},
		"verification": []string{},
		"other":       []string{},
	}

	for _, txt := range txtRecords {
		txt = strings.TrimSpace(txt)
		switch {
		case strings.HasPrefix(txt, "v=spf"):
			categorized["spf"] = append(categorized["spf"], txt)
		case strings.HasPrefix(txt, "v=DMARC"):
			categorized["dmarc"] = append(categorized["dmarc"], txt)
		case strings.Contains(txt, "dkim"):
			categorized["dkim"] = append(categorized["dkim"], txt)
		case strings.Contains(txt, "verification") || strings.Contains(txt, "verify"):
			categorized["verification"] = append(categorized["verification"], txt)
		default:
			categorized["other"] = append(categorized["other"], txt)
		}
	}

	return output.Result{
		Timestamp: time.Now(),
		Source:    "dns",
		Type:      "TXT",
		Target:    domain,
		Data: map[string]interface{}{
			"records":     txtRecords,
			"categorized": categorized,
			"count":       len(txtRecords),
		},
	}, nil
}

// lookupCNAME performs CNAME record lookup
func (m *DomainModule) lookupCNAME(domain string) (output.Result, error) {
	cname, err := net.LookupCNAME(domain)
	if err != nil {
		return output.Result{}, err
	}

	return output.Result{
		Timestamp: time.Now(),
		Source:    "dns",
		Type:      "CNAME",
		Target:    domain,
		Data: map[string]interface{}{
			"cname": cname,
		},
	}, nil
}

// LookupSubdomains attempts to find subdomains
func (m *DomainModule) LookupSubdomains(ctx context.Context, domain string) ([]string, error) {
	commonSubdomains := []string{
		"www", "mail", "ftp", "localhost", "webmail", "smtp", "pop", "ns1", "webdisk",
		"ns2", "cpanel", "whm", "autodiscover", "autoconfig", "m", "imap", "test",
		"ns", "blog", "pop3", "dev", "www2", "admin", "forum", "news", "vpn", "ns3",
		"mail2", "new", "mysql", "old", "lists", "support", "mobile", "mx", "static",
		"docs", "beta", "shop", "sql", "secure", "demo", "cp", "calendar", "wiki",
		"web", "media", "email", "images", "img", "www1", "intranet", "portal", "video",
		"sip", "dns2", "api", "cdn", "stats", "dns1", "ns4", "www3", "dns", "search",
		"staging", "server", "mx1", "chat", "wap", "my", "svn", "mail1", "sites",
		"proxy", "ads", "host", "crm", "cms", "backup", "mx2", "lyncdiscover", "info",
		"apps", "download", "remote", "db", "forums", "store", "relay", "files",
		"newsletter", "app", "live", "owa", "en", "start", "sms", "office", "exchange",
		"ipv4", "mail3", "help", "blogs", "helpdesk", "web1", "home", "library", "ftp2",
		"ntp", "monitor", "login", "service", "correo", "www4", "moodle", "it", "gateway",
		"gw", "i", "stat", "stage", "ldap", "tv", "ssl", "cloud", "development", "photos",
	}

	var subdomains []string
	for _, sub := range commonSubdomains {
		subdomain := fmt.Sprintf("%s.%s", sub, domain)
		_, err := net.LookupIP(subdomain)
		if err == nil {
			subdomains = append(subdomains, subdomain)
		}
	}

	return subdomains, nil
}
