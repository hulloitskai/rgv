package bot

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
	ess "github.com/unixpickle/essentials"
)

// creds holds credentials for authenticating with the Reddit API.
type creds struct {
	ClientID   string `split_words:"true"`
	Secret     string
	User, Pass string
}

// readCreds reads a creds configuration from environment variables.
func readCreds() (*creds, error) {
	a := new(creds)
	if err := envconfig.Process(Namespace, a); err != nil {
		return nil, ess.AddCtx("reading credentials from env", err)
	}
	return a, nil
}

func (a *creds) IsEmpty() bool {
	return (a == nil) || (a.ClientID == "") && (a.Secret == "") &&
		(a.User == "") && (a.Pass == "")
}

func (a *creds) Validate() error {
	if a.IsEmpty() {
		return nil // an empty set of credentials is valid
	}

	var field string
	switch "" {
	case a.ClientID:
		field = "ClientID"
	case a.Secret:
		field = "Secret"
	case a.User:
		field = "User"
	case a.Pass:
		field = "Pass"
	}
	if field != "" {
		return fmt.Errorf("not all fields are filled: field '%s' is empty", field)
	}
	return nil
}
