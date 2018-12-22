package bot

import (
	"github.com/kelseyhightower/envconfig"
	ess "github.com/unixpickle/essentials"
)

// creds holds credentials for authenticating with the Reddit API.
type creds struct {
	ClientID string `split_words:"true"`
	Secret   string
}

// readCreds reads a creds configuration from environment variables.
func readCreds() (*creds, error) {
	a := new(creds)
	if err := envconfig.Process(Namespace, a); err != nil {
		return nil, ess.AddCtx("bot: reading credentials from env", err)
	}
	return a, nil
}

func (a *creds) IsValid() bool {
	return (a != nil) && (a.ClientID != "") && (a.Secret != "")
}
