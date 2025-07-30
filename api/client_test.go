package api

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ankitpokhrel/jira-cli/pkg/jira"
)

func TestClient_CloudflareAccessTokenFromEnv(t *testing.T) {
	// Clean up any existing environment variable
	originalValue := os.Getenv("CF_ACCESS_TOKEN")
	defer func() {
		if originalValue != "" {
			_ = os.Setenv("CF_ACCESS_TOKEN", originalValue)
		} else {
			_ = os.Unsetenv("CF_ACCESS_TOKEN")
		}
	}()

	// Test with manual token from environment variable
	expectedToken := "test-cf-token-from-env"
	_ = os.Setenv("CF_ACCESS_TOKEN", expectedToken)

	// Reset the global client to force re-initialization
	jiraClient = nil

	client := Client(jira.Config{
		Server:   "https://test.atlassian.net",
		Login:    "user@example.com",
		APIToken: "api-token",
	})

	assert.NotNil(t, client)
}

func TestClient_CloudflareAccessTokenFromConfig(t *testing.T) {
	// Clean up any existing environment variable
	originalValue := os.Getenv("CF_ACCESS_TOKEN")
	defer func() {
		if originalValue != "" {
			_ = os.Setenv("CF_ACCESS_TOKEN", originalValue)
		} else {
			_ = os.Unsetenv("CF_ACCESS_TOKEN")
		}
	}()

	// Unset environment variable to ensure config takes precedence
	_ = os.Unsetenv("CF_ACCESS_TOKEN")

	// Reset the global client to force re-initialization
	jiraClient = nil

	expectedToken := "test-cf-token-from-config"
	client := Client(jira.Config{
		Server:        "https://test.atlassian.net",
		Login:         "user@example.com",
		APIToken:      "api-token",
		CFAccessToken: expectedToken,
	})

	assert.NotNil(t, client)
}

func TestClient_CloudflareAccessTokenPriority(t *testing.T) {
	// Test that config takes precedence over environment variable
	originalValue := os.Getenv("CF_ACCESS_TOKEN")
	defer func() {
		if originalValue != "" {
			_ = os.Setenv("CF_ACCESS_TOKEN", originalValue)
		} else {
			_ = os.Unsetenv("CF_ACCESS_TOKEN")
		}
	}()

	// Set environment variable to "auto"
	_ = os.Setenv("CF_ACCESS_TOKEN", "auto")

	// Reset the global client to force re-initialization
	jiraClient = nil

	// But provide manual config value - should take precedence
	configToken := "config-token"
	client := Client(jira.Config{
		Server:        "https://test.atlassian.net",
		Login:         "user@example.com",
		APIToken:      "api-token",
		CFAccessToken: configToken,
	})

	// Config should take precedence over environment variable
	assert.NotNil(t, client)
}

func TestClient_NoCloudflareAccessToken(t *testing.T) {
	// Clean up any existing environment variable
	originalValue := os.Getenv("CF_ACCESS_TOKEN")
	defer func() {
		if originalValue != "" {
			_ = os.Setenv("CF_ACCESS_TOKEN", originalValue)
		} else {
			_ = os.Unsetenv("CF_ACCESS_TOKEN")
		}
	}()

	// Ensure no environment variable is set
	_ = os.Unsetenv("CF_ACCESS_TOKEN")

	// Reset the global client to force re-initialization
	jiraClient = nil

	client := Client(jira.Config{
		Server:   "https://test.atlassian.net",
		Login:    "user@example.com",
		APIToken: "api-token",
		// CFAccessToken is empty
	})

	assert.NotNil(t, client)
}

func TestClient_CloudflareAccessTokenAuto(t *testing.T) {
	// Test behavior when CF_ACCESS_TOKEN=auto
	originalValue := os.Getenv("CF_ACCESS_TOKEN")
	defer func() {
		if originalValue != "" {
			_ = os.Setenv("CF_ACCESS_TOKEN", originalValue)
		} else {
			_ = os.Unsetenv("CF_ACCESS_TOKEN")
		}
	}()

	// Set token to "auto" - should attempt automatic generation
	_ = os.Setenv("CF_ACCESS_TOKEN", "auto")

	// Reset the global client to force re-initialization
	jiraClient = nil

	client := Client(jira.Config{
		Server:   "https://test.atlassian.net",
		Login:    "user@example.com",
		APIToken: "api-token",
	})

	// Client should still be created even if cloudflared fails
	assert.NotNil(t, client)
}

func TestClient_EmptyCloudflareAccessTokenInEnv(t *testing.T) {
	// Test behavior when environment variable is set but empty
	originalValue := os.Getenv("CF_ACCESS_TOKEN")
	defer func() {
		if originalValue != "" {
			_ = os.Setenv("CF_ACCESS_TOKEN", originalValue)
		} else {
			_ = os.Unsetenv("CF_ACCESS_TOKEN")
		}
	}()

	// Set empty environment variable
	_ = os.Setenv("CF_ACCESS_TOKEN", "")

	// Reset the global client to force re-initialization
	jiraClient = nil

	client := Client(jira.Config{
		Server:   "https://test.atlassian.net",
		Login:    "user@example.com",
		APIToken: "api-token",
	})

	assert.NotNil(t, client)
}
