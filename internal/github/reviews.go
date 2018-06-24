package github

import (
	"context"
	"strings"

	"github.com/google/go-github/github"
	"go.uber.org/zap"
)

type Allocator struct {
	logger       *zap.Logger
	githubClient Client
	reviewers    []*Reviewers
}

type Reviewers struct {
	Language string
	Users    []string
}

func NewAllocator(logger *zap.Logger, githubClient Client, reviewers []*Reviewers) *Allocator {
	return &Allocator{
		logger:       logger,
		githubClient: githubClient,
		reviewers:    reviewers,
	}
}

func (a *Allocator) RequestReviewers(ctx context.Context, event *eventData) error {
	for _, reviewers := range a.reviewers {
		if strings.ToUpper(reviewers.Language) != strings.ToUpper(event.repositoryLanguage) {
			continue
		}

		return a.githubClient.RequestReviewers(ctx, event, github.ReviewersRequest{
			Reviewers: reviewers.Users,
		})
	}

	return nil
}

func (a *Allocator) DeleteReviewers(ctx context.Context, event *eventData) error {
	return a.githubClient.DeleteReviewers(ctx, event)
}
