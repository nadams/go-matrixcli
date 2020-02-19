package auth

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

var m sync.Mutex

type Auth struct {
	Accounts []AccountAuth `json:"accounts"`
}

type AccountAuth struct {
	Name  string `json:"name"`
	Token string `json:"token"`
}

func Token(name, cacheDir string, client http.Client) (string, error) {
	m.Lock()
	defer m.Unlock()

	path := filepath.Join(cacheDir, "auth.json")

	var auth *Auth

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		auth, _ = loadAuth(path)
	}

	for _, a := range auth.Accounts {
		if a.Name == name {
			return a.Token, nil
		}
	}

	// TODO: fetch a new token if account doesn't exist

	return "", nil
}

func loadAuth(path string) (*Auth, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	var auth Auth

	if err := json.NewDecoder(f).Decode(&auth); err != nil {
		return nil, err
	}

	return &auth, nil
}
