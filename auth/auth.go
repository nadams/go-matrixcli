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

const filename = "accounts.json"

var m sync.Mutex

type Auth struct {
	Accounts []AccountAuth `json:"accounts"`
}

type AccountAuth struct {
	Name   string `json:"name"`
	UserID string `json:"userId"`
	Token  string `json:"token"`
}

type accountAuths []AccountAuth

func (a accountAuths) Find(name string) (AccountAuth, bool) {
	for _, auth := range a {
		if auth.Name == name {
			return auth, true
		}
	}

	return AccountAuth{}, false
}

func (a accountAuths) Update(auth AccountAuth) accountAuths {
	var found bool

	for i, au := range a {
		if au.Name == auth.Name {
			found = true
			a[i] = auth

			break
		}
	}

	if !found {
		a = append(a, auth)
	}

	return a
}

type TokenStore struct {
	config   *config.Config
	client   *http.Client
	accounts accountAuths
}

func NewTokenStore(c *config.Config) (*TokenStore, error) {
	m.Lock()
	defer m.Unlock()

	accounts, err := loadFromFile(filepath.Join(c.CacheDir, filename))
	if err != nil {
		return nil, err
	}

	return &TokenStore{
		config:   c,
		client:   &http.Client{Timeout: time.Second * 30},
		accounts: accounts,
	}, nil
}

func (t *TokenStore) Login(name, homeserver, username, password string) (AccountAuth, error) {
	m.Lock()
	defer m.Unlock()

	resp, err := t.login(homeserver, username, password)
	if err != nil {
		return AccountAuth{}, err
	}

	aa := AccountAuth{
		Name:   name,
		UserID: resp.UserID,
		Token:  resp.AccessToken,
	}

	t.accounts = t.accounts.Update(aa)

	if err := t.persist(); err != nil {
		return AccountAuth{}, err
	}

	return aa, nil
}

func (t *TokenStore) Token(name string) (AccountAuth, error) {
	m.Lock()
	defer m.Unlock()

	if auth, ok := t.accounts.Find(name); ok {
		return auth, nil
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

	resp, err := t.login(account.Homeserver, account.Username, account.Password)
	if err != nil {
		return AccountAuth{}, err
	}

	aa := AccountAuth{
		Name:   name,
		UserID: resp.UserID,
		Token:  resp.AccessToken,
	}

	t.accounts = t.accounts.Update(aa)

	if err := t.persist(); err != nil {
		return AccountAuth{}, err
	}

	return aa, nil
}

func (t *TokenStore) login(homeserver, username, password string) (*gomatrix.RespLogin, error) {
	cl, err := gomatrix.NewClient(homeserver, "", "")
	if err != nil {
		return nil, err
	}

	return cl.Login(&gomatrix.ReqLogin{
		Type:     "m.login.password",
		User:     username,
		Password: password,
	})
}

func (t *TokenStore) persist() error {
	return saveFile(filepath.Join(t.config.CacheDir, filename), t.accounts)
}

func saveFile(path string, accounts []AccountAuth) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	defer f.Close()

	if err := os.Chmod(path, 0600); err != nil {
		return err
	}

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")

	return enc.Encode(accounts)
}

func loadFromFile(path string) ([]AccountAuth, error) {
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		return []AccountAuth{}, nil
	} else {
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}

		defer f.Close()

		var accounts []AccountAuth

		if err := json.NewDecoder(f).Decode(&accounts); err != nil {
			return nil, err
		}

		return accounts, nil
	}
}
