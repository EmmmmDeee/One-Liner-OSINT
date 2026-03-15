package osint

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/mail"
	"regexp"
	"strings"
	"time"

	"github.com/EmmmmDeee/One-Liner-OSINT/pkg/config"
	"github.com/EmmmmDeee/One-Liner-OSINT/pkg/logger"
	"github.com/EmmmmDeee/One-Liner-OSINT/pkg/output"
	"github.com/EmmmmDeee/One-Liner-OSINT/pkg/utils"
)

// EmailModule provides email OSINT capabilities
type EmailModule struct {
	config *config.Config
	logger *logger.Logger
	client *utils.HTTPClient
}

// NewEmailModule creates a new email module
func NewEmailModule(cfg *config.Config) *EmailModule {
	return &EmailModule{
		config: cfg,
		logger: logger.NewLogger(),
		client: utils.NewHTTPClient(cfg.Timeout, cfg.UserAgent),
	}
}

// Name returns the module name
func (m *EmailModule) Name() string {
	return "email"
}

// Description returns the module description
func (m *EmailModule) Description() string {
	return "Email validation and breach checking"
}

// Execute performs email OSINT
func (m *EmailModule) Execute(ctx context.Context, target string) ([]output.Result, error) {
	var results []output.Result

	// Validate email format
	if !m.isValidEmail(target) {
		return nil, fmt.Errorf("invalid email format: %s", target)
	}

	// Check email format
	formatResult := m.checkEmailFormat(target)
	results = append(results, formatResult)

	// Check for breaches if API key is available
	if m.config.APIKeys.HaveIBeenPwned != "" {
		breachResult, err := m.checkBreaches(ctx, target)
		if err != nil {
			m.logger.Warnf("Breach check failed: %v", err)
		} else {
			results = append(results, breachResult)
		}
	}

	// Extract domain info
	domainResult := m.extractDomainInfo(target)
	results = append(results, domainResult)

	return results, nil
}

// isValidEmail validates email format
func (m *EmailModule) isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// checkEmailFormat checks email format and structure
func (m *EmailModule) checkEmailFormat(email string) output.Result {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return output.Result{
			Timestamp: time.Now(),
			Source:    "email",
			Type:      "format",
			Target:    email,
			Error:     "invalid email format",
		}
	}

	return output.Result{
		Timestamp: time.Now(),
		Source:    "email",
		Type:      "format",
		Target:    email,
		Data: map[string]interface{}{
			"username": parts[0],
			"domain":   parts[1],
			"valid":    true,
		},
	}
}

// extractDomainInfo extracts domain information from email
func (m *EmailModule) extractDomainInfo(email string) output.Result {
	parts := strings.Split(email, "@")
	domain := parts[1]

	// Common email providers
	commonProviders := map[string]string{
		"gmail.com":     "Google Gmail",
		"yahoo.com":     "Yahoo Mail",
		"outlook.com":   "Microsoft Outlook",
		"hotmail.com":   "Microsoft Hotmail",
		"icloud.com":    "Apple iCloud",
		"protonmail.com": "ProtonMail",
		"aol.com":       "AOL Mail",
		"mail.com":      "Mail.com",
	}

	data := map[string]interface{}{
		"domain":       domain,
		"is_freemail":  false,
		"provider":     "Unknown",
	}

	if provider, ok := commonProviders[domain]; ok {
		data["is_freemail"] = true
		data["provider"] = provider
	}

	return output.Result{
		Timestamp: time.Now(),
		Source:    "email",
		Type:      "domain",
		Target:    email,
		Data:      data,
	}
}

// checkBreaches checks if email is in any breaches
func (m *EmailModule) checkBreaches(ctx context.Context, email string) (output.Result, error) {
	url := fmt.Sprintf("https://haveibeenpwned.com/api/v3/breachedaccount/%s", email)

	m.client.SetHeader("hibp-api-key", m.config.APIKeys.HaveIBeenPwned)
	m.client.SetHeader("User-Agent", m.config.UserAgent)

	resp, err := m.client.Get(ctx, url)
	if err != nil {
		return output.Result{}, err
	}
	defer resp.Body.Close()

	result := output.Result{
		Timestamp: time.Now(),
		Source:    "haveibeenpwned",
		Type:      "breach",
		Target:    email,
	}

	if resp.StatusCode == 404 {
		// No breaches found
		result.Data = map[string]interface{}{
			"breached":      false,
			"breach_count":  0,
			"breaches":      []string{},
		}
		return result, nil
	}

	if resp.StatusCode != 200 {
		return result, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	var breaches []map[string]interface{}
	if err := json.Unmarshal(body, &breaches); err != nil {
		return result, err
	}

	breachNames := make([]string, 0, len(breaches))
	for _, breach := range breaches {
		if name, ok := breach["Name"].(string); ok {
			breachNames = append(breachNames, name)
		}
	}

	result.Data = map[string]interface{}{
		"breached":     true,
		"breach_count": len(breaches),
		"breaches":     breachNames,
	}

	return result, nil
}

// ExtractEmailsFromText extracts email addresses from text
func (m *EmailModule) ExtractEmailsFromText(text string) []string {
	emailRegex := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	emails := emailRegex.FindAllString(text, -1)

	// Deduplicate
	seen := make(map[string]bool)
	unique := make([]string, 0)
	for _, email := range emails {
		email = strings.ToLower(email)
		if !seen[email] {
			seen[email] = true
			unique = append(unique, email)
		}
	}

	return unique
}
