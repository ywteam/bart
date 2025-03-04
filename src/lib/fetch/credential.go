package fetch

import (
	"strings"
	"encoding/base64"
)

type Credential struct {
	Secret []byte
}

func NewCrential(value string) *Credential {
	return &Credential{
		Secret: []byte(value),
	}
}
func (s *Credential) Get() string {
	return string(s.Secret)
}
func (s *Credential) Header() []string {
	return strings.Split(s.Get(), ":")
}
func (s *Credential) Rotate(value string) *Credential{
	s.Secret = []byte(value)
	return s
}

type ICredentailStore interface {
	AddBasicAuth(username, password string) *CredentialStore
	AddBasic64Auth(username, password string) *CredentialStore
	AddBearerToken(token string) *CredentialStore
	AddAuthzHeader(value string) *CredentialStore
	AddCustomHeader(header, value string) *CredentialStore
	IsEmpty() bool
	Next() *Credential
	Rotate(value string)	*Credential
}

// var store []Credential = make([]Credential, 0)

type CredentialStore struct {
	current int
	store   []Credential
}

func NewCredentialStore() ICredentailStore {
	return &CredentialStore{
		current: 0,
		store:   make([]Credential, 0),
	}
}
func (c *CredentialStore) Rotate(value string) *Credential {
	return c.store[c.current].Rotate(value)
}
func (c *CredentialStore) AddBasic64Auth(username, password string) *CredentialStore {
	encoded := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	c.AddAuthzHeader("Basic " + encoded)
	return c
}
func (c *CredentialStore) AddBasicAuth(username, password string) *CredentialStore {
	c.AddAuthzHeader("Basic " + username + ":" + password)
	return c
}
func (c *CredentialStore) AddBearerToken(token string) *CredentialStore {
	c.AddAuthzHeader("Bearer " + token)
	return c
}
func (c *CredentialStore) AddAuthzHeader(value string) *CredentialStore {
	c.AddCustomHeader("Authorization", value)
	return c
}
func (c *CredentialStore) AddCustomHeader(header, value string) *CredentialStore {
	c.store = append(c.store, Credential{Secret: []byte(header + ":" + value)})
	return c
}
func (c *CredentialStore) IsEmpty() bool {
	return len(c.store) == 0
}
func (c *CredentialStore) Next() *Credential {
	if c.current >= len(c.store) {
		c.current = 0
	}
	secret := c.store[c.current]
	c.current++
	return &secret
}
