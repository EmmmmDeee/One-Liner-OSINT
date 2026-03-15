package osint

import (
	"testing"

	"github.com/EmmmmDeee/One-Liner-OSINT/pkg/config"
)

func TestEmailValidation(t *testing.T) {
	cfg := &config.Config{}
	module := NewEmailModule(cfg)

	tests := []struct {
		email string
		valid bool
	}{
		{"user@example.com", true},
		{"user.name@example.com", true},
		{"user+tag@example.co.uk", true},
		{"invalid-email", false},
		{"@example.com", false},
		{"user@", false},
		{"", false},
	}

	for _, tt := range tests {
		result := module.isValidEmail(tt.email)
		if result != tt.valid {
			t.Errorf("Email %s: expected %v, got %v", tt.email, tt.valid, result)
		}
	}
}

func TestExtractEmailsFromText(t *testing.T) {
	cfg := &config.Config{}
	module := NewEmailModule(cfg)

	text := `
		Contact us at support@example.com or sales@example.com
		You can also reach admin@test.org
		Duplicate: support@example.com
	`

	emails := module.ExtractEmailsFromText(text)

	expectedCount := 3
	if len(emails) != expectedCount {
		t.Errorf("Expected %d unique emails, got %d", expectedCount, len(emails))
	}

	// Check for expected emails
	expectedEmails := map[string]bool{
		"support@example.com": true,
		"sales@example.com":   true,
		"admin@test.org":      true,
	}

	for _, email := range emails {
		if !expectedEmails[email] {
			t.Errorf("Unexpected email found: %s", email)
		}
	}
}
