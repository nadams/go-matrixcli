package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/matrix-org/gomatrix"
	"github.com/nadams/go-matrixcli/config"
)

var m sync.Mutex

type Auth struct {
	Accounts []AccountAuth `json:"accounts"`
}

type AccountAuth struct {
	Name   string `json:"name"`
	UserID string `json:"userId"`
	Token  string `json:"token"`
}

type TokenStore struct {
	config   *config.Config
	client   *http.Client
	accounts []AccountAuth
}

func NewTokenStore(c *config.Config) (*TokenStore, error) {
	t := &TokenStore{
		config: c,
		client: &http.Client{Timeout: time.Second * 30},
	}

	fi, err := os.Stat(filepath.Join(c.CacheDir, "accounts.json"))
	if os.IsNotExist(err) {
		t.accounts = []AccountAuth{}
	} else {
		f, err := os.Open(filepath.Join(c.CacheDir, fi.Name()))
		if err != nil {
			return nil, err
		}

		defer f.Close()

		if err := json.NewDecoder(f).Decode(&t.accounts); err != nil {
			return nil, err
		}
	}

	return t, nil
}

func (t *TokenStore) Token(name string) (AccountAuth, error) {
	m.Lock()
	defer m.Unlock()

	for _, account := range t.accounts {
		if account.Name == name {
			return account, nil
		}
	}

	var found bool
	var account config.Account

	for _, a := range t.config.Accounts {
		if a.Name == name {
			found = true
			account = a

			break
		}
	}

	if !found {
		return AccountAuth{}, fmt.Errorf("could not found account in config: %s", name)
	}

	cl, err := gomatrix.NewClient(account.Homeserver, "", "")
	if err != nil {
		return AccountAuth{}, err
	}

	resp, err := cl.Login(&gomatrix.ReqLogin{
		Type:     "m.login.password",
		User:     account.Username,
		Password: account.Password,
	})

	if err != nil {
		return AccountAuth{}, err
	}

	aa := AccountAuth{
		Name:   name,
		UserID: resp.UserID,
		Token:  resp.AccessToken,
	}

	t.accounts = append(t.accounts, aa)

	f, err := os.Create(filepath.Join(t.config.CacheDir, "accounts.json"))
	if err != nil {
		return AccountAuth{}, nil
	}

	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")

	if err := enc.Encode(t.accounts); err != nil {
		return AccountAuth{}, err
	}

	return aa, nil
}
