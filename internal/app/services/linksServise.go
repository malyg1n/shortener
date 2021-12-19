package services

import (
	"errors"
	"github.com/malyg1n/shortener/internal/app/storage"
	"math/rand"
	"net/url"
	"regexp"
)

type LinksService interface {
	SetLink(link string) (string, error)
	GetLink(id string) (string, error)
}

type DefaultLinksService struct {
	storage storage.LinksStorage
}

const linkPattern = `[a-zA-Z0-9]+`

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func NewDefaultLinksService(storage storage.LinksStorage) *DefaultLinksService {
	return &DefaultLinksService{
		storage: storage,
	}
}

func (s *DefaultLinksService) GetLink(id string) (string, error) {
	matched, _ := regexp.MatchString(linkPattern, id)
	if !matched {
		return "", errors.New("invalid link ID")
	}
	return s.storage.GetLink(id)
}

func (s *DefaultLinksService) SetLink(link string) (string, error) {
	_, err := url.ParseRequestURI(link)
	if err != nil {
		return "", errors.New("incorrect url param")
	}
	randString := randomString(6)
	s.storage.SetLink(randString, link)
	return randString, nil
}

func randomString(n uint) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
