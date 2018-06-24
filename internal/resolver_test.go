package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tommy-muehle/garvis/internal/configuration"
	"github.com/tommy-muehle/garvis/internal/github"
	server2 "github.com/tommy-muehle/garvis/internal/server"
	"go.uber.org/zap"
)

func TestCanResolveWebhook(t *testing.T) {
	assert := assert.New(t)

	webhook := NewResolver(configuration.DefaultConfig()).ResolveWebhook()
	assert.IsType(new(github.Webhook), webhook)
}

func TestCanResolveAllocator(t *testing.T) {
	assert := assert.New(t)

	allocator := NewResolver(configuration.DefaultConfig()).ResolveAllocator()
	assert.IsType(new(github.Allocator), allocator)
}

func TestCanResolveLabeler(t *testing.T) {
	assert := assert.New(t)

	labeler := NewResolver(configuration.DefaultConfig()).ResolveLabeler()
	assert.IsType(new(github.Labeler), labeler)
}

func TestCanResolveLogger(t *testing.T) {
	assert := assert.New(t)

	logger := NewResolver(configuration.DefaultConfig()).ResolveLogger()
	assert.IsType(new(zap.Logger), logger)
}

func TestCanResolveGithubClient(t *testing.T) {
	assert := assert.New(t)

	client := NewResolver(configuration.DefaultConfig()).ResolveGithubClient()
	assert.IsType(new(github.GithubClient), client)
}

func TestCanResolveServer(t *testing.T) {
	assert := assert.New(t)

	server := NewResolver(configuration.DefaultConfig()).ResolveServer()
	assert.IsType(new(server2.Server), server)
}
