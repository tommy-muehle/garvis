package internal

import (
	"github.com/tommy-muehle/garvis/internal/configuration"
	"github.com/tommy-muehle/garvis/internal/github"
	"github.com/tommy-muehle/garvis/internal/server"

	origin "github.com/google/go-github/github"
	"go.uber.org/zap"
)

type Resolver struct {
	config       *configuration.Configuration
	logger       *zap.Logger
	githubClient github.Client
}

func NewResolver(config *configuration.Configuration) *Resolver {
	return &Resolver{
		config: config,
	}
}

func (r *Resolver) ResolveLabeler() *github.Labeler {
	return github.NewLabeler(
		r.ResolveGithubClient(),
		r.ResolveLogger(),
	)
}

func (r *Resolver) ResolveAllocator() *github.Allocator {
	reviewers := make([]*github.Reviewers, 0, len(r.config.Reviewers))

	for _, reviewer := range r.config.Reviewers {
		reviewers = append(reviewers, &github.Reviewers{
			Language: reviewer.Language,
			Users:    reviewer.Users,
		})
	}

	return github.NewAllocator(
		r.ResolveLogger(),
		r.ResolveGithubClient(),
		reviewers,
	)
}

func (r *Resolver) ResolveLogger() *zap.Logger {
	var err error

	if r.logger == nil {
		if r.config.Log.Debug {
			r.logger, err = zap.NewDevelopment()
		} else {
			r.logger, err = zap.NewProduction()
		}

		if err != nil {
			panic(err)
		}
	}

	return r.logger
}

func (r *Resolver) ResolveGithubClient() github.Client {
	if r.githubClient == nil {
		transport := origin.BasicAuthTransport{
			Username: r.config.Github.Username,
			Password: r.config.Github.Password,
		}
		r.githubClient = github.NewGithubClient(origin.NewClient(transport.Client()))
	}

	return r.githubClient
}

func (r *Resolver) ResolveWebhook() *github.Webhook {
	return github.NewWebhook(
		r.config.Github.SecretKey,
		r.ResolveLabeler(),
		r.ResolveAllocator(),
		r.ResolveLogger(),
	)
}

func (r *Resolver) ResolveServer() *server.Server {
	return server.New(
		r.config.Server.Addr,
		r.ResolveLogger(),
	)
}
