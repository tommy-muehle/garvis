package github

import (
	"context"

	"go.uber.org/zap"
)

var (
	labelNeedsReview = "Status: Needs review"
	labelReviewed    = "Status: Reviewed"
	labelNeedsWork   = "Status: Needs Work"
)

type Labeler struct {
	logger       *zap.Logger
	githubClient Client
}

func NewLabeler(githubClient Client, logger *zap.Logger) *Labeler {
	return &Labeler{
		logger:       logger,
		githubClient: githubClient,
	}
}

func (l *Labeler) AddNeedsReview(ctx context.Context, event *eventData) error {
	return l.githubClient.AddLabel(ctx, event, []string{labelNeedsReview})
}

func (l *Labeler) RemoveNeedsReview(ctx context.Context, event *eventData) error {
	return l.githubClient.RemoveLabel(ctx, event, labelNeedsReview)
}

func (l *Labeler) AddReviewed(ctx context.Context, event *eventData) error {
	return l.githubClient.AddLabel(ctx, event, []string{labelReviewed})
}

func (l *Labeler) RemoveReviewed(ctx context.Context, event *eventData) error {
	return l.githubClient.RemoveLabel(ctx, event, labelReviewed)
}

func (l *Labeler) AddNeedsWork(ctx context.Context, event *eventData) error {
	return l.githubClient.AddLabel(ctx, event, []string{labelNeedsWork})
}

func (l *Labeler) RemoveNeedsWork(ctx context.Context, event *eventData) error {
	return l.githubClient.RemoveLabel(ctx, event, labelNeedsWork)
}
