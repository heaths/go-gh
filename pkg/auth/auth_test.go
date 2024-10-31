package auth

import (
	"strconv"
	"testing"

	"github.com/cli/go-gh/v2/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenForHost(t *testing.T) {
	tests := []struct {
		name                  string
		host                  string
		githubToken           string
		githubEnterpriseToken string
		ghToken               string
		ghEnterpriseToken     string
		codespaces            bool
		config                *config.Config
		wantToken             string
		wantSource            string
	}{
		{
			name:       "token for github.com with no env tokens and no config token",
			host:       "github.com",
			config:     testNoHostsConfig(),
			wantToken:  "",
			wantSource: defaultSource,
		},
		{
			name:       "token for enterprise.com with no env tokens and no config token",
			host:       "enterprise.com",
			config:     testNoHostsConfig(),
			wantToken:  "",
			wantSource: defaultSource,
		},
		{
			name:        "token for github.com with GH_TOKEN, GITHUB_TOKEN, and config token",
			host:        "github.com",
			ghToken:     "GH_TOKEN",
			githubToken: "GITHUB_TOKEN",
			config:      testHostsConfig(),
			wantToken:   "GH_TOKEN",
			wantSource:  ghToken,
		},
		{
			name:        "token for github.com with GITHUB_TOKEN, and config token",
			host:        "github.com",
			githubToken: "GITHUB_TOKEN",
			config:      testHostsConfig(),
			wantToken:   "GITHUB_TOKEN",
			wantSource:  githubToken,
		},
		{
			name:       "token for github.com with config token",
			host:       "github.com",
			config:     testHostsConfig(),
			wantToken:  "xxxxxxxxxxxxxxxxxxxx",
			wantSource: oauthToken,
		},
		{
			name:                  "token for enterprise.com with GH_ENTERPRISE_TOKEN, GITHUB_ENTERPRISE_TOKEN, and config token",
			host:                  "enterprise.com",
			ghEnterpriseToken:     "GH_ENTERPRISE_TOKEN",
			githubEnterpriseToken: "GITHUB_ENTERPRISE_TOKEN",
			config:                testHostsConfig(),
			wantToken:             "GH_ENTERPRISE_TOKEN",
			wantSource:            ghEnterpriseToken,
		},
		{
			name:                  "token for enterprise.com with GITHUB_ENTERPRISE_TOKEN, and config token",
			host:                  "enterprise.com",
			githubEnterpriseToken: "GITHUB_ENTERPRISE_TOKEN",
			config:                testHostsConfig(),
			wantToken:             "GITHUB_ENTERPRISE_TOKEN",
			wantSource:            githubEnterpriseToken,
		},
		{
			name:       "token for enterprise.com with config token",
			host:       "enterprise.com",
			config:     testHostsConfig(),
			wantToken:  "yyyyyyyyyyyyyyyyyyyy",
			wantSource: oauthToken,
		},
		{
			name:        "token for tenant with GH_TOKEN, GITHUB_TOKEN, and config token",
			host:        "tenant.ghe.com",
			ghToken:     "GH_TOKEN",
			githubToken: "GITHUB_TOKEN",
			config:      testHostsConfig(),
			wantToken:   "GH_TOKEN",
			wantSource:  ghToken,
		},
		{
			name:        "token for tenant with GITHUB_TOKEN, and config token",
			host:        "tenant.ghe.com",
			githubToken: "GITHUB_TOKEN",
			config:      testHostsConfig(),
			wantToken:   "GITHUB_TOKEN",
			wantSource:  githubToken,
		},
		{
			name:       "token for tenant with config token",
			host:       "tenant.ghe.com",
			config:     testHostsConfig(),
			wantToken:  "zzzzzzzzzzzzzzzzzzzz",
			wantSource: oauthToken,
		},
		{
			name:        "Token for non-github host in a codespace",
			host:        "doesnotmatter.com",
			config:      testNoHostsConfig(),
			githubToken: "GITHUB_TOKEN",
			codespaces:  true,
			wantToken:   "",
			wantSource:  defaultSource,
		},
		{
			name:        "Token for github.com in a codespace",
			host:        "github.com",
			config:      testNoHostsConfig(),
			githubToken: "GITHUB_TOKEN",
			codespaces:  true,
			wantToken:   "GITHUB_TOKEN",
			wantSource:  githubToken,
		},
		{
			// We are in a codespace (dotcom), and we have set our own GITHUB_TOKEN, not using the codespace one
			// and we are targeting tenant.ghe.com
			name:        "Token for tenant.ghe.com in a codespace",
			host:        "tenant.ghe.com",
			config:      testNoHostsConfig(),
			githubToken: "GITHUB_TOKEN",
			codespaces:  true,
			wantToken:   "GITHUB_TOKEN",
			wantSource:  githubToken,
		},
		{
			name:        "Token for github.localhost in a codespace",
			host:        "github.localhost",
			config:      testNoHostsConfig(),
			githubToken: "GITHUB_TOKEN",
			codespaces:  true,
			wantToken:   "GITHUB_TOKEN",
			wantSource:  githubToken,
		},
		{
			// We are in codespace (dotcom), and we have set a GITHUB_ENTERPRISE_TOKEN, and we are targeting GHES
			name:                  "Enterprise Token for GHES in a codespace",
			host:                  "enterprise.com",
			config:                testNoHostsConfig(),
			githubEnterpriseToken: "GITHUB_ENTERPRISE_TOKEN",
			codespaces:            true,
			wantToken:             "GITHUB_ENTERPRISE_TOKEN",
			wantSource:            githubEnterpriseToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("GITHUB_TOKEN", tt.githubToken)
			t.Setenv("GITHUB_ENTERPRISE_TOKEN", tt.githubEnterpriseToken)
			t.Setenv("GH_TOKEN", tt.ghToken)
			t.Setenv("GH_ENTERPRISE_TOKEN", tt.ghEnterpriseToken)
			t.Setenv("CODESPACES", strconv.FormatBool(tt.codespaces))
			token, source := tokenForHost(tt.config, tt.host)
			require.Equal(t, tt.wantToken, token, "Expected token for \"%s\" to be \"%s\", got \"%s\"", tt.host, tt.wantToken, token)
			require.Equal(t, tt.wantSource, source, "Expected source for \"%s\" to be \"%s\", got \"%s\"", tt.host, tt.wantSource, source)
		})
	}
}

func TestDefaultHost(t *testing.T) {
	tests := []struct {
		name         string
		config       *config.Config
		ghHost       string
		wantHost     string
		wantSource   string
		wantNotFound bool
	}{
		{
			name:       "GH_HOST if set",
			config:     testHostsConfig(),
			ghHost:     "test.com",
			wantHost:   "test.com",
			wantSource: "GH_HOST",
		},
		{
			name:       "authenticated host if only one",
			config:     testSingleHostConfig(),
			wantHost:   "enterprise.com",
			wantSource: "hosts",
		},
		{
			name:         "default host if more than one authenticated host",
			config:       testHostsConfig(),
			wantHost:     "github.com",
			wantSource:   "default",
			wantNotFound: true,
		},
		{
			name:         "default host if no authenticated host",
			config:       testNoHostsConfig(),
			wantHost:     "github.com",
			wantSource:   "default",
			wantNotFound: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.ghHost != "" {
				t.Setenv("GH_HOST", tt.ghHost)
			}
			host, source := defaultHost(tt.config)
			assert.Equal(t, tt.wantHost, host)
			assert.Equal(t, tt.wantSource, source)
		})
	}
}

func TestKnownHosts(t *testing.T) {
	tests := []struct {
		name      string
		config    *config.Config
		ghHost    string
		ghToken   string
		wantHosts []string
	}{
		{
			name:      "no known hosts",
			config:    testNoHostsConfig(),
			wantHosts: []string{},
		},
		{
			name:      "includes GH_HOST",
			config:    testNoHostsConfig(),
			ghHost:    "test.com",
			wantHosts: []string{"test.com"},
		},
		{
			name:      "includes authenticated hosts",
			config:    testHostsConfig(),
			wantHosts: []string{"github.com", "enterprise.com", "tenant.ghe.com"},
		},
		{
			name:      "includes default host if environment auth token",
			config:    testNoHostsConfig(),
			ghToken:   "TOKEN",
			wantHosts: []string{"github.com"},
		},
		{
			name:      "deduplicates hosts",
			config:    testHostsConfig(),
			ghHost:    "test.com",
			ghToken:   "TOKEN",
			wantHosts: []string{"test.com", "github.com", "enterprise.com", "tenant.ghe.com"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.ghHost != "" {
				t.Setenv("GH_HOST", tt.ghHost)
			}
			if tt.ghToken != "" {
				t.Setenv("GH_TOKEN", tt.ghToken)
			}
			hosts := knownHosts(tt.config)
			assert.Equal(t, tt.wantHosts, hosts)
		})
	}
}

func TestIsEnterprise(t *testing.T) {
	tests := []struct {
		name    string
		host    string
		wantOut bool
	}{
		{
			name:    "github",
			host:    "github.com",
			wantOut: false,
		},
		{
			name:    "github API",
			host:    "api.github.com",
			wantOut: false,
		},
		{
			name:    "localhost",
			host:    "github.localhost",
			wantOut: false,
		},
		{
			name:    "localhost API",
			host:    "api.github.localhost",
			wantOut: false,
		},
		{
			name:    "enterprise",
			host:    "mygithub.com",
			wantOut: true,
		},
		{
			name:    "tenant",
			host:    "tenant.ghe.com",
			wantOut: false,
		},
		{
			name:    "tenant API",
			host:    "api.tenant.ghe.com",
			wantOut: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := IsEnterprise(tt.host)
			assert.Equal(t, tt.wantOut, out)
		})
	}
}

func TestIsTenancy(t *testing.T) {
	tests := []struct {
		name    string
		host    string
		wantOut bool
	}{
		{
			name:    "github",
			host:    "github.com",
			wantOut: false,
		},
		{
			name:    "github API",
			host:    "api.github.com",
			wantOut: false,
		},
		{
			name:    "localhost",
			host:    "github.localhost",
			wantOut: false,
		},
		{
			name:    "localhost API",
			host:    "api.github.localhost",
			wantOut: false,
		},
		{
			name:    "enterprise",
			host:    "mygithub.com",
			wantOut: false,
		},
		{
			name:    "tenant",
			host:    "tenant.ghe.com",
			wantOut: true,
		},
		{
			name:    "tenant API",
			host:    "api.tenant.ghe.com",
			wantOut: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := IsTenancy(tt.host)
			assert.Equal(t, tt.wantOut, out)
		})
	}
}

func TestNormalizeHostname(t *testing.T) {
	tests := []struct {
		name     string
		host     string
		wantHost string
	}{
		{
			name:     "github domain",
			host:     "test.github.com",
			wantHost: "github.com",
		},
		{
			name:     "capitalized",
			host:     "GitHub.com",
			wantHost: "github.com",
		},
		{
			name:     "localhost domain",
			host:     "test.github.localhost",
			wantHost: "github.localhost",
		},
		{
			name:     "enterprise domain",
			host:     "mygithub.com",
			wantHost: "mygithub.com",
		},
		{
			name:     "bare tenant",
			host:     "tenant.ghe.com",
			wantHost: "tenant.ghe.com",
		},
		{
			name:     "subdomained tenant",
			host:     "api.tenant.ghe.com",
			wantHost: "tenant.ghe.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			normalized := NormalizeHostname(tt.host)
			assert.Equal(t, tt.wantHost, normalized)
		})
	}
}

func testNoHostsConfig() *config.Config {
	var data = ``
	return config.ReadFromString(data)
}

func testSingleHostConfig() *config.Config {
	var data = `
hosts:
  enterprise.com:
    user: user2
    oauth_token: yyyyyyyyyyyyyyyyyyyy
    git_protocol: https
`
	return config.ReadFromString(data)
}

func testHostsConfig() *config.Config {
	var data = `
hosts:
  github.com:
    user: user1
    oauth_token: xxxxxxxxxxxxxxxxxxxx
    git_protocol: ssh
  enterprise.com:
    user: user2
    oauth_token: yyyyyyyyyyyyyyyyyyyy
    git_protocol: https
  tenant.ghe.com:
    user: user3
    oauth_token: zzzzzzzzzzzzzzzzzzzz
    git_protocol: https
`
	return config.ReadFromString(data)
}
