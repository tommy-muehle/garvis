package github

import (
	"context"
	"testing"

	"github.com/google/go-github/github"
)

func TestCanRequestReviewers(t *testing.T) {
	ctx := context.Background()
	event := new(eventData)
	event.repositoryLanguage = "PHP"

	users := []string{"foo", "bar"}

	client := new(fakeClient)
	client.
		On("RequestReviewers", ctx, event, github.ReviewersRequest{Reviewers: users}).
		Return(nil)

	allocator := NewAllocator(nil, client, []*Reviewers{&Reviewers{Language: "PHP", Users: users}})
	if err := allocator.RequestReviewers(ctx, event); err != nil {
		t.Fatal(err)
	}

	client.AssertExpectations(t)
}

func TestNoReviewersRequestedIfLanguageNotMatch(t *testing.T) {
	ctx := context.Background()
	event := new(eventData)
	event.repositoryLanguage = "Go"

	client := new(fakeClient)
	client.AssertNotCalled(t, "RequestReviewers")

	allocator := NewAllocator(nil, client, []*Reviewers{&Reviewers{Language: "PHP", Users: []string{}}})
	if err := allocator.RequestReviewers(ctx, event); err != nil {
		t.Fatal(err)
	}

	client.AssertExpectations(t)
}
