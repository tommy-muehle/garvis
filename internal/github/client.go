package github

import (
	"context"

	origin "github.com/google/go-github/github"
)

type Client interface {
	AddLabel(ctx context.Context, event *eventData, labels []string) error
	RemoveLabel(ctx context.Context, event *eventData, label string) error
	RequestReviewers(ctx context.Context, event *eventData, reviewers origin.ReviewersRequest) error
	DeleteReviewers(ctx context.Context, event *eventData) error
}

type GithubClient struct {
	origin *origin.Client
}

func NewGithubClient(origin *origin.Client) *GithubClient {
	return &GithubClient{origin: origin}
}

func (c *GithubClient) AddLabel(ctx context.Context, event *eventData, newLabels []string) error {
	labels, _, err := c.origin.Issues.ListLabelsByIssue(
		ctx, event.owner, event.repository, event.number, &origin.ListOptions{},
	)

	if err != nil {
		return err
	}

	for _, label := range labels {
		for _, newLabel := range newLabels {
			if label.GetName() == newLabel {
				return nil
			}
		}
	}

	_, _, err = c.origin.Issues.AddLabelsToIssue(
		ctx, event.owner, event.repository, event.number, newLabels,
	)

	return err
}

func (c *GithubClient) RemoveLabel(ctx context.Context, event *eventData, label string) error {
	labels, _, err := c.origin.Issues.ListLabelsByIssue(
		ctx, event.owner, event.repository, event.number, &origin.ListOptions{},
	)

	if err != nil {
		return err
	}

	labelExist := false

	for _, l := range labels {
		if l.GetName() == label {
			labelExist = true
		}
	}

	if !labelExist {
		return nil
	}

	_, err = c.origin.Issues.RemoveLabelForIssue(
		ctx, event.owner, event.repository, event.number, label,
	)

	return err
}

func (c *GithubClient) RequestReviewers(ctx context.Context, event *eventData, reviewers origin.ReviewersRequest) error {
	_, _, err := c.origin.PullRequests.RequestReviewers(
		ctx, event.owner, event.repository, event.number, reviewers,
	)

	return err
}

func (c *GithubClient) DeleteReviewers(ctx context.Context, event *eventData) error {
	reviewers, _, err := c.origin.PullRequests.ListReviewers(
		ctx, event.owner, event.repository, event.number, &origin.ListOptions{},
	)

	if err != nil {
		return err
	}

	reviews, _, err := c.origin.PullRequests.ListReviews(
		ctx, event.owner, event.repository, event.number, &origin.ListOptions{},
	)

	if err != nil {
		return err
	}

	removableReviewers := make([]string, 0)

	for _, reviewer := range reviewers.Users {
		for _, review := range reviews {
			if review.User.Login == reviewer.Login {
				continue
			}

			removableReviewers = append(removableReviewers, *reviewer.Login)
		}
	}

	_, err = c.origin.PullRequests.RemoveReviewers(
		ctx, event.owner, event.repository, event.number, origin.ReviewersRequest{
			Reviewers: removableReviewers,
		},
	)

	return err
}
