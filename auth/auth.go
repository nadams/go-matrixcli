package auth

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

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

	return AccountAuth{}, nil
}
