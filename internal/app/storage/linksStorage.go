package storage

import "errors"

type LinksStorage interface {
	SetLink(id string, link string)
	GetLink(id string) (string, error)
}

type LinksStorageMap struct {
	links map[string]string
}

func NewLinksStorageMap() *LinksStorageMap {
	return &LinksStorageMap{
		links: map[string]string{},
	}
}

func (s *LinksStorageMap) SetLink(id string, link string) {
	s.links[id] = link
}

func (s *LinksStorageMap) GetLink(id string) (string, error) {
	if link, ok := s.links[id]; ok {
		return link, nil
	}

	return "", errors.New("link not found")
}
