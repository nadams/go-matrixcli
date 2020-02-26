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
	Current  string       `json:"current_account"`
	Accounts AccountAuths `json:"accounts"`
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

type AccountAuths []AccountAuth

func (a AccountAuths) Find(name string) (AccountAuth, bool) {
	for _, auth := range a {
		if auth.Name == name {
			return auth, true
		}
	}

	return AccountAuth{}, false
}

func (a AccountAuths) Update(auth AccountAuth) AccountAuths {
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
	dir    string
	client *http.Client
	auth   Auth
}

func NewTokenStore(dir string, client *http.Client) (*TokenStore, error) {
	m.Lock()
	defer m.Unlock()

	auth, err := loadFromFile(filepath.Join(dir, filename))
	if err != nil {
		return nil, err
	}

	if client == nil {
		client = &http.Client{Timeout: time.Second * 30}
	}

	return &TokenStore{
		dir:    dir,
		client: client,
		auth:   auth,
	}, nil
}

func (t *TokenStore) List() AccountAuths {
	m.Lock()
	defer m.Unlock()

	return t.auth.Accounts
}

func (t *TokenStore) CurrentName() string {
	m.Lock()
	defer m.Unlock()

	return t.auth.Current
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

	t.auth.Accounts = t.auth.Accounts.Update(aa)

	if len(t.auth.Accounts) == 1 {
		t.auth.Current = aa.Name
	}

	if err := t.persist(); err != nil {
		return AccountAuth{}, err
	}

	return aa, nil
}

func (t *TokenStore) Find(name string) (AccountAuth, error) {
	m.Lock()
	defer m.Unlock()

	return t.find(name)
}

func (t *TokenStore) Current() (AccountAuth, error) {
	m.Lock()
	defer m.Unlock()

	return t.current()
}

func (t *TokenStore) SetCurrent(name string) error {
	m.Lock()
	defer m.Unlock()

	if _, ok := t.auth.Accounts.Find(name); !ok {
		return ErrAccountNotFound
	}

	t.auth.Current = name

	return t.persist()
}

func (t *TokenStore) Remove(name string) error {
	m.Lock()
	defer m.Unlock()

	t.remove(name)

	return t.persist()
}

func (t *TokenStore) remove(name string) {
	idx := -1

	for i, aa := range t.auth.Accounts {
		if aa.Name == name {
			idx = i
			break
		}
	}

	if idx > -1 {
		copy(t.auth.Accounts[idx:], t.auth.Accounts[idx+1:])
		t.auth.Accounts = t.auth.Accounts[:len(t.auth.Accounts)-1]
	}

	if len(t.auth.Accounts) > 0 {
		t.auth.Current = t.auth.Accounts[0].Name
	}
}

func (t *TokenStore) find(name string) (AccountAuth, error) {
	if auth, ok := t.auth.Accounts.Find(name); ok {
		return auth, nil
	}

	return AccountAuth{}, ErrAccountNotFound
}

func (t *TokenStore) current() (AccountAuth, error) {
	if len(t.auth.Accounts) == 0 {
		return AccountAuth{}, ErrNoAccounts
	}

	if t.auth.Current == "" {
		return t.auth.Accounts[0], nil
	}

	auth, ok := t.auth.Accounts.Find(t.auth.Current)
	if !ok {
		return AccountAuth{}, ErrAccountNotFound
	}

	return auth, nil
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
	return saveFile(filepath.Join(t.dir, filename), t.auth)
}

func saveFile(path string, auth Auth) error {
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

	return enc.Encode(auth)
}

func loadFromFile(path string) (Auth, error) {
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		return Auth{}, nil
	} else {
		f, err := os.Open(path)
		if err != nil {
			return Auth{}, err
		}

		defer f.Close()

		var auth Auth

		if err := json.NewDecoder(f).Decode(&auth); err != nil {
			return Auth{}, err
		}

		return auth, nil
	}
}
