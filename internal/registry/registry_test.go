package registry_test

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thanhpk/randstr"

	"github.com/orange-cloudavenue/kube-image-updater/internal/registry"
)

var (
	rdRepoName = strings.ToLower(randstr.String(16))
	versions   = []string{"v0.0.1", "v0.0.2", "v0.0.3", "v0.0.4", "v0.0.5"}
)

func buildImageRegistry() string {
	return fmt.Sprintf("registry.127.0.0.1.nip.io/%s", rdRepoName)
}

func buildImageRegistryWithTag(version string) string {
	return fmt.Sprintf("%s:%s", buildImageRegistry(), version)
}

func initImageRegistry() {
	for _, version := range versions {
		if err := exec.Command("docker", "build", "-f", "Dockerfile", "-t", buildImageRegistryWithTag(version), ".").Run(); err != nil { //nolint:gosec
			continue
		}

		if err := exec.Command("docker", "push", buildImageRegistryWithTag(version)).Run(); err != nil { //nolint:gosec
			continue
		}
	}
}

func cleanupImageRegistry() {
	for _, version := range versions {
		if err := exec.Command("docker", "rmi", buildImageRegistryWithTag(version)).Run(); err != nil { //nolint:gosec
			continue
		}
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		repo    string
		wantErr bool
	}{
		{
			name:    "valid repository",
			repo:    "alpine",
			wantErr: false,
		},
		{
			name:    "empty repository",
			repo:    "",
			wantErr: true,
		},
		{
			name:    "invalid registry",
			repo:    "gcr.oi:invalid/repo",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := registry.New(context.Background(), tt.repo, registry.Settings{})
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTags(t *testing.T) {
	tests := []struct {
		name        string
		repo        string
		credentials registry.Settings
		wantErr     bool
	}{
		{
			name:    "valid repository",
			repo:    "alpine",
			wantErr: false,
		},
		{
			name:    "empty repository",
			repo:    "",
			wantErr: true,
		},
		{
			name:    "valid non-dockerhub repository",
			repo:    "ghcr.io/traefik/whoami",
			wantErr: false,
		},
		{
			name: "valid repository with authentification",
			repo: buildImageRegistry(),
			credentials: registry.Settings{
				InsecureTLS: true,
				Username:    "myuser",
				Password:    "mypasswd",
			},
			wantErr: false,
		},
	}

	initImageRegistry()
	defer cleanupImageRegistry()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := registry.New(context.Background(), tt.repo, tt.credentials)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			t.Logf("Fetching tags for %s", tt.repo)

			tags, err := r.Tags()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, tags)
			}
		})
	}
}
