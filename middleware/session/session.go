package session

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/tinh-tinh/tinhtinh/internal/memory"
)

type Config struct {
	store       *memory.Store
	cookie      http.Cookie
	GeneratorID func() string
	Secret      string
	ExpiresIn   time.Duration
}

type Options struct {
	StoreOptions memory.Options
	GeneratorID  func() string
	Secret       string
	// Default is 1 hour.
	ExpiresIn time.Duration
}

func New(opt Options) *Config {
	session := &Config{
		Secret: opt.Secret,
		store:  memory.New(opt.StoreOptions),
	}
	if session.ExpiresIn == 0 {
		session.ExpiresIn = time.Hour
	}
	if opt.GeneratorID != nil {
		session.GeneratorID = opt.GeneratorID
	}

	return session
}

func (s *Config) Get(key string) interface{} {
	return s.store.Get(s.Hash(key))
}

func (s *Config) Set(key string, val interface{}) *http.Cookie {
	var ID string
	if s.GeneratorID != nil {
		ID = s.GeneratorID()
	} else {
		ID = s.DefaultGenerateID()
	}

	s.cookie = http.Cookie{
		Name:     key,
		Value:    ID,
		HttpOnly: true,
		MaxAge:   int(s.ExpiresIn),
		Secure:   true,
	}
	s.store.Set(s.Hash(ID), val, s.ExpiresIn)
	return &s.cookie
}

func (s *Config) Hash(data string) string {
	hmac := hmac.New(sha256.New, []byte(s.Secret))

	hmac.Write([]byte(data))
	dataHMac := hmac.Sum(nil)

	return hex.EncodeToString(dataHMac)
}

func (s *Config) DefaultGenerateID() string {
	newUUID, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(newUUID)
}
