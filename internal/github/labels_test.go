package github

import (
	"context"
	"testing"
)

func Test_CanRemoveLabels(t *testing.T) {
	ctx := context.Background()
	event := new(eventData)

	t.Run("Can remove 'Needs review' label", func(t *testing.T) {
		client := new(fakeClient)
		client.On("RemoveLabel", ctx, event, labelNeedsReview).Return(nil)

		labeler := NewLabeler(client, nil)
		if err := labeler.RemoveNeedsReview(ctx, event); err != nil {
			t.Fatal(err)
		}

		client.AssertExpectations(t)
	})

	t.Run("Can remove 'Needs work' label", func(t *testing.T) {
		client := new(fakeClient)
		client.On("RemoveLabel", ctx, event, labelNeedsWork).Return(nil)

		labeler := NewLabeler(client, nil)
		if err := labeler.RemoveNeedsWork(ctx, event); err != nil {
			t.Fatal(err)
		}

		client.AssertExpectations(t)
	})

	t.Run("Can remove 'Reviewed' label", func(t *testing.T) {
		client := new(fakeClient)
		client.On("RemoveLabel", ctx, event, labelReviewed).Return(nil)

		labeler := NewLabeler(client, nil)
		if err := labeler.RemoveReviewed(ctx, event); err != nil {
			t.Fatal(err)
		}

		client.AssertExpectations(t)
	})
}

func Test_CanAddLabels(t *testing.T) {
	ctx := context.Background()
	event := new(eventData)

	t.Run("Can add 'Needs review' label", func(t *testing.T) {
		client := new(fakeClient)
		client.On("AddLabel", ctx, event, []string{labelNeedsReview}).Return(nil)

		labeler := NewLabeler(client, nil)
		if err := labeler.AddNeedsReview(ctx, event); err != nil {
			t.Fatal(err)
		}

		client.AssertExpectations(t)
	})

	t.Run("Can add 'Needs work' label", func(t *testing.T) {
		client := new(fakeClient)
		client.On("AddLabel", ctx, event, []string{labelNeedsWork}).Return(nil)

		labeler := NewLabeler(client, nil)
		if err := labeler.AddNeedsWork(ctx, event); err != nil {
			t.Fatal(err)
		}

		client.AssertExpectations(t)
	})

	t.Run("Can add 'Reviewed' label", func(t *testing.T) {
		client := new(fakeClient)
		client.On("AddLabel", ctx, event, []string{labelReviewed}).Return(nil)

		labeler := NewLabeler(client, nil)
		if err := labeler.AddReviewed(ctx, event); err != nil {
			t.Fatal(err)
		}

		client.AssertExpectations(t)
	})
}
