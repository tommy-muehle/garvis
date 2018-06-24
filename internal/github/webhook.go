package github

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	origin "github.com/google/go-github/github"
	"go.uber.org/zap"
)

var (
	pullRequestEdited          string = "edited"
	pullRequestOpened          string = "opened"
	pullRequestReviewSubmitted string = "submitted"
	reviewStateApproved        string = "approved"
	reviewStateNeedsChanges    string = "changes_requested"
)

type Webhook struct {
	logger    *zap.Logger
	allocator *Allocator
	labeler   *Labeler
	secretKey []byte
}

func NewWebhook(secretKey string, labeler *Labeler, allocator *Allocator, logger *zap.Logger) *Webhook {
	return &Webhook{
		logger:    logger,
		allocator: allocator,
		labeler:   labeler,
		secretKey: []byte(secretKey),
	}
}

type eventData struct {
	action             string
	owner              string
	repository         string
	repositoryLanguage string
	number             int
}

func (h *Webhook) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.logger.Debug("payload received")

		payload, err := origin.ValidatePayload(r, h.secretKey)
		if err != nil {
			h.logger.Error("could not validate payload", zap.Error(err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		event, err := origin.ParseWebHook(origin.WebHookType(r), payload)
		if err != nil {
			h.logger.Error("could not parse event", zap.Error(err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelFunc()

		switch event := event.(type) {
		case *origin.PullRequestEvent:
			if strings.ToLower(*event.Action) == pullRequestOpened {
				h.processNewPullRequest(ctx, event)
			}
			if strings.ToLower(*event.Action) == pullRequestEdited {
				h.processUpdatedPullRequest(ctx, event)
			}
			break

		case *origin.PullRequestReviewEvent:
			if strings.ToLower(*event.Action) == pullRequestReviewSubmitted {
				h.processPullRequestReviewEvent(ctx, event)
			}
			break
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		if _, err := fmt.Fprintln(w, "Ok"); err != nil {
			h.logger.Error("could not respond", zap.Error(err))
		}
	})
}

func (h *Webhook) processUpdatedPullRequest(ctx context.Context, event *origin.PullRequestEvent) {
	data := &eventData{
		action:             event.GetAction(),
		owner:              event.GetRepo().GetOwner().GetLogin(),
		repository:         event.GetRepo().GetName(),
		number:             event.GetPullRequest().GetNumber(),
		repositoryLanguage: event.GetRepo().GetLanguage(),
	}

	if err := h.labeler.RemoveNeedsWork(ctx, data); err != nil {
		h.logger.Error("could not remove needs work label", zap.Error(err))
	}

	if err := h.labeler.AddNeedsReview(ctx, data); err != nil {
		h.logger.Error("could not add needs review label", zap.Error(err))
	}
}

func (h *Webhook) processNewPullRequest(ctx context.Context, event *origin.PullRequestEvent) {
	data := &eventData{
		action:             event.GetAction(),
		owner:              event.GetRepo().GetOwner().GetLogin(),
		repository:         event.GetRepo().GetName(),
		number:             event.GetPullRequest().GetNumber(),
		repositoryLanguage: event.GetRepo().GetLanguage(),
	}

	if err := h.labeler.AddNeedsReview(ctx, data); err != nil {
		h.logger.Error("could not add needs review label", zap.Error(err))
	}

	if err := h.allocator.RequestReviewers(ctx, data); err != nil {
		h.logger.Error("could not request reviewers", zap.Error(err))
	}
}

func (h *Webhook) processPullRequestReviewEvent(ctx context.Context, event *origin.PullRequestReviewEvent) {
	data := &eventData{
		action:     event.GetAction(),
		owner:      event.GetRepo().GetOwner().GetLogin(),
		repository: event.GetRepo().GetName(),
		number:     event.GetPullRequest().GetNumber(),
	}

	if err := h.labeler.RemoveNeedsReview(ctx, data); err != nil {
		h.logger.Error("could not remove needs review label", zap.Error(err))
	}

	reviewState := strings.ToLower(event.GetReview().GetState())

	switch reviewState {
	case reviewStateApproved:
		if err := h.labeler.RemoveNeedsWork(ctx, data); err != nil {
			h.logger.Error("could not remove needs work label", zap.Error(err))
		}
		if err := h.labeler.AddReviewed(ctx, data); err != nil {
			h.logger.Error("could not add reviewed label", zap.Error(err))
		}
		if err := h.allocator.DeleteReviewers(ctx, data); err != nil {
			h.logger.Error("could not delete reviewers", zap.Error(err))
		}

	case reviewStateNeedsChanges:
		if err := h.labeler.AddNeedsWork(ctx, data); err != nil {
			h.logger.Error("could not add needs work label", zap.Error(err))
		}
	}
}
