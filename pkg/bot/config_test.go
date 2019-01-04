package bot_test

import (
	"os"
	"strings"
	"testing"

	"github.com/stevenxie/rgv/pkg/bot"
	"github.com/stevenxie/rgv/pkg/reddit"
)

func TestBot_creds(t *testing.T) {
	// Test without env config.
	envVars := []string{
		"REDDIT_CLIENT_ID", "REDDIT_SECRET", "REDDIT_USER", "REDDIT_PASS",
	}
	for _, key := range envVars {
		if err := os.Unsetenv(key); err != nil {
			t.Fatalf("Failed to unset env variable '%s': %v", key, err)
		}
	}

	b, err := bot.New(receiver{}, nil)
	if err != nil {
		t.Error("Expected empty-env bot creation to return without errors.")
	}
	if b == nil {
		t.Error("Expected bot (made without env vars) to be non-nil.")
	}

	// Test with incomplete env config.
	if err = os.Setenv("REDDIT_USER", "testuser"); err != nil {
		t.Fatalf("Failed to set env variable: %v", err)
	}
	if b, err = bot.New(receiver{}, nil); err == nil {
		t.Errorf("Expected invalid-env bot creation to return an error.")
	}
	if !strings.Contains(err.Error(), "not all fields are filled") {
		t.Errorf("Unexpcted bot creation error: %v", err)
	} else {
		t.Logf("Got bot creation error (expected): %v", err)
	}
	if b != nil {
		t.Error("Expected bot (made with invalid env) to be nil.")
	}
}

type receiver struct{}

func (r receiver) ReceivePost(p *reddit.Post) error {
	return nil
}

func (r receiver) ReceiveComment(c *reddit.Comment) error {
	return nil
}
