package jira

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// GetCloudflareAccessToken retrieves a Cloudflare access token.
// If CF_ACCESS_TOKEN="auto", it generates a token via cloudflared.
// Otherwise, it returns the CF_ACCESS_TOKEN value as-is.
func GetCloudflareAccessToken(serverURL string) (string, error) {
	token := os.Getenv("CF_ACCESS_TOKEN")

	// If token is "auto", generate it automatically
	if token == "auto" {
		return generateCloudflareToken(serverURL)
	}

	// Return manual token (or empty string if not set)
	return token, nil
}

// generateCloudflareToken runs cloudflared access login and extracts the token from stdout.
func generateCloudflareToken(serverURL string) (string, error) {
	if serverURL == "" {
		return "", fmt.Errorf("server URL is required for automatic Cloudflare token generation")
	}

	// Check if cloudflared is available
	if _, err := exec.LookPath("cloudflared"); err != nil {
		return "", fmt.Errorf("cloudflared command not found: %w", err)
	}

	// Set timeout to prevent hanging
	const cloudflaredTimeout = 30 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), cloudflaredTimeout)
	defer cancel()

	// Run cloudflared access login
	cmd := exec.CommandContext(ctx, "cloudflared", "access", "login", serverURL)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("cloudflared access login failed: %w\nOutput: %s", err, string(output))
	}

	// Parse token from output
	token, err := parseCloudflareToken(string(output))
	if err != nil {
		return "", fmt.Errorf("failed to parse token from cloudflared output: %w", err)
	}

	return token, nil
}

// parseCloudflareToken extracts the token from cloudflared output.
// Expected format: "Successfully fetched your token:\n\n<token>".
func parseCloudflareToken(output string) (string, error) {
	// Look for the success message pattern
	re := regexp.MustCompile(`Successfully fetched your token:\s*\n\s*\n\s*([^\s\n]+)`)
	matches := re.FindStringSubmatch(output)

	if len(matches) < 2 {
		return "", fmt.Errorf("could not find token in cloudflared output")
	}

	token := strings.TrimSpace(matches[1])
	if token == "" {
		return "", fmt.Errorf("extracted token is empty")
	}

	return token, nil
}
