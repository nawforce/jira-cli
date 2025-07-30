package jira

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCloudflareAccessToken_ManualToken(t *testing.T) {
	// Clean up environment
	originalToken := os.Getenv("CF_ACCESS_TOKEN")
	defer func() {
		if originalToken != "" {
			_ = os.Setenv("CF_ACCESS_TOKEN", originalToken)
		} else {
			_ = os.Unsetenv("CF_ACCESS_TOKEN")
		}
	}()

	// Test manual token
	expectedToken := "manual-token-123"
	_ = os.Setenv("CF_ACCESS_TOKEN", expectedToken)

	token, err := GetCloudflareAccessToken("https://jira.example.com")
	assert.NoError(t, err)
	assert.Equal(t, expectedToken, token)
}

func TestGetCloudflareAccessToken_NoToken(t *testing.T) {
	// Clean up environment
	originalToken := os.Getenv("CF_ACCESS_TOKEN")
	defer func() {
		if originalToken != "" {
			_ = os.Setenv("CF_ACCESS_TOKEN", originalToken)
		} else {
			_ = os.Unsetenv("CF_ACCESS_TOKEN")
		}
	}()

	// No token set
	_ = os.Unsetenv("CF_ACCESS_TOKEN")

	token, err := GetCloudflareAccessToken("https://jira.example.com")
	assert.NoError(t, err)
	assert.Empty(t, token)
}

func TestGetCloudflareAccessToken_AutoWithoutCloudflared(t *testing.T) {
	// Clean up environment
	originalToken := os.Getenv("CF_ACCESS_TOKEN")
	defer func() {
		if originalToken != "" {
			_ = os.Setenv("CF_ACCESS_TOKEN", originalToken)
		} else {
			_ = os.Unsetenv("CF_ACCESS_TOKEN")
		}
	}()

	// Set token to "auto" to enable automatic generation
	_ = os.Setenv("CF_ACCESS_TOKEN", "auto")

	// This should fail because cloudflared is not available in test environment
	token, err := GetCloudflareAccessToken("https://jira.example.com")
	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Contains(t, err.Error(), "cloudflared command not found")
}

func TestParseCloudflareToken_Success(t *testing.T) {
	output := `Starting Cloudflare Access Login...
Opening browser...
Successfully fetched your token:

eyJhbGciOiJSUzI1NiIsImtpZCI6IjkxM2E0...token_continues_here

You can now use this token with your applications.`

	token, err := parseCloudflareToken(output)
	assert.NoError(t, err)
	assert.Equal(t, "eyJhbGciOiJSUzI1NiIsImtpZCI6IjkxM2E0...token_continues_here", token)
}

func TestParseCloudflareToken_WithWhitespace(t *testing.T) {
	output := `Successfully fetched your token:


   eyJhbGciOiJSUzI1NiIsImtpZCI6IjkxM2E0...token_with_whitespace   

`

	token, err := parseCloudflareToken(output)
	assert.NoError(t, err)
	assert.Equal(t, "eyJhbGciOiJSUzI1NiIsImtpZCI6IjkxM2E0...token_with_whitespace", token)
}

func TestParseCloudflareToken_NoToken(t *testing.T) {
	output := `Starting Cloudflare Access Login...
Error: Authentication failed`

	token, err := parseCloudflareToken(output)
	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Contains(t, err.Error(), "could not find token")
}

func TestParseCloudflareToken_EmptyToken(t *testing.T) {
	output := `Successfully fetched your token:


   

`

	token, err := parseCloudflareToken(output)
	assert.Error(t, err)
	assert.Empty(t, token)
	// The regex pattern won't match whitespace-only tokens, so it will say "could not find token"
	assert.Contains(t, err.Error(), "could not find token")
}

func TestParseCloudflareToken_MultipleLines(t *testing.T) {
	// Test with more realistic cloudflared output
	output := `⠋ Waiting for login...
⢿ Received access token from Cloudflare Access
Successfully fetched your token:

eyJhbGciOiJSUzI1NiIsImtpZCI6IjkxM2E0YTM4LWZhNjItNDU5Yy1hNzVlLWQ5ODg2NzEyNzg1MSIsInR5cCI6IkpXVCJ9.eyJhdWQiOlsiZXhhbXBsZS5jb20iXSwiZW1haWwiOiJ1c2VyQGV4YW1wbGUuY29tIiwiZXhwIjoxNjMwNTYwMDAwLCJpYXQiOjE2MzA1NTY0MDAsImlzcyI6Imh0dHBzOi8vYWNjZXNzLmNsb3VkZmxhcmUuY29tIiwibmJmIjoxNjMwNTU2NDAwLCJzdWIiOiJhYmMxMjMifQ.signature_here

You can now close this window and return to your terminal.`

	token, err := parseCloudflareToken(output)
	assert.NoError(t, err)
	assert.True(t, len(token) > 50)                  // JWT tokens are typically long
	assert.Contains(t, token, "eyJhbGciOiJSUzI1NiI") // Typical JWT header
}

func TestGenerateCloudflareToken_EmptyURL(t *testing.T) {
	token, err := generateCloudflareToken("")
	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Contains(t, err.Error(), "server URL is required")
}
