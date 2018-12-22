package bot_test

import (
	"os"
	"testing"

	"github.com/stevenxie/rgv/pkg/reddit"

	"github.com/stevenxie/rgv/pkg/bot"
)

func TestBot_creds(t *testing.T) {
	// Test without env config variables.
	err := os.Unsetenv("REDDIT_CLIENT_ID")
	err = os.Unsetenv("REDDIT_SECRET")
	if err != nil {
		t.Fatalf("Failed to unset env variable: %v", err)
	}

	b, err := bot.New(receiver{}, nil)
	berr, ok := err.(*bot.Error)

	if !ok {
		t.Fatalf("Expected bot to return an error of type *bot.Error. Instead, "+
			"got: %v", err)
	}
	if berr.Code != bot.InvalidConfig {
		t.Error("Expected error to have the code bot.InvalidConfig.")
	}
	if b == nil {
		t.Error("Expected bot (made without env vars) to be non-nil.")
	}

	// Test with env config variables.
	err = os.Setenv("REDDIT_CLIENT_ID", "someid")
	err = os.Setenv("REDDIT_SECRET", "somesecret")
	if err != nil {
		t.Fatalf("Failed to set env variable: %v", err)
	}

	if b, err = bot.New(receiver{}, nil); err != nil {
		t.Errorf("Did not expect bot creation to return an error, but got: %v", err)
	}
	if b == nil {
		t.Error("Expected bot (made with env vars) to be non-nil.")
	}
}

type receiver struct{}

func (r receiver) ReceivePost(p *reddit.Post) error {
	return nil
}

func (r receiver) ReceiveComment(c *reddit.Comment) error {
	return nil
}
