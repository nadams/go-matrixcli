package auth

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/matrix-org/gomatrix"
)

const filename = "accounts.json"

var m sync.Mutex

type Auth struct {
	Accounts []AccountAuth `json:"accounts"`
}

type AccountAuth struct {
	Name       string `json:"name"`
	Homeserver string `json:"homeserver"`
	UserID     string `json:"userId"`
	Token      string `json:"token"`
}

func (a AccountAuth) Domain() (string, error) {
	u, err := url.Parse(a.Homeserver)
	if err != nil {
		return "", err
	}

	return u.Hostname(), nil
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
	dir      string
	client   *http.Client
	accounts accountAuths
}

func NewTokenStore(dir string) (*TokenStore, error) {
	m.Lock()
	defer m.Unlock()

	accounts, err := loadFromFile(filepath.Join(dir, filename))
	if err != nil {
		return nil, err
	}

	return &TokenStore{
		dir:      dir,
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
		Name:       name,
		Homeserver: homeserver,
		UserID:     resp.UserID,
		Token:      resp.AccessToken,
	}

	t.accounts = t.accounts.Update(aa)

	if err := t.persist(); err != nil {
		return AccountAuth{}, err
	}

	return aa, nil
}

func (t *TokenStore) Find(name string) (AccountAuth, error) {
	m.Lock()
	defer m.Unlock()

	if auth, ok := t.accounts.Find(name); ok {
		return auth, nil
	}

	return AccountAuth{}, ErrAccountNotFound
}

func (t *TokenStore) First() (AccountAuth, error) {
	m.Lock()
	defer m.Unlock()

	if len(t.accounts) == 0 {
		return AccountAuth{}, ErrNoAccounts
	}

	return t.accounts[0], nil
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
	return saveFile(filepath.Join(t.dir, filename), t.accounts)
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
