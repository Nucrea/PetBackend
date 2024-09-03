package services

import (
	"backend/src/cache"
	"backend/src/charsets"
	"backend/src/core/repos"
	"context"
	"fmt"
	"math/rand"
	"time"
)

var (
	ErrShortlinkNotexist = fmt.Errorf("shortlink does not exist or expired")
	ErrShortlinkExpired  = fmt.Errorf("shortlink expired")
)

type ShortlinkService interface {
	CreateShortlink(ctx context.Context, url string) (string, error)
	GetShortlink(ctx context.Context, id string) (string, error)
	ShortlinkRoutine(ctx context.Context) error
}

type NewShortlinkServiceParams struct {
	Cache cache.Cache[string, string]
	Repo  repos.ShortlinkRepo
}

func NewShortlinkSevice(params NewShortlinkServiceParams) ShortlinkService {
	return &shortlinkService{
		cache: params.Cache,
		repo:  params.Repo,
	}
}

type shortlinkService struct {
	cache cache.Cache[string, string]
	repo  repos.ShortlinkRepo
}

func (s *shortlinkService) CreateShortlink(ctx context.Context, url string) (string, error) {
	charset := charsets.GetCharset(charsets.CharsetTypeAll)

	src := rand.NewSource(time.Now().UnixMicro())
	randGen := rand.New(src)
	id := charset.RandomString(randGen, 10)

	expiration := time.Now().Add(7 * 24 * time.Hour)

	dto := repos.ShortlinkDTO{
		Id:         id,
		Url:        url,
		Expiration: expiration,
	}
	if err := s.repo.AddShortlink(ctx, dto); err != nil {
		return "", err
	}

	s.cache.Set(id, url, cache.Expiration{ExpiresAt: expiration})

	return id, nil
}

func (s *shortlinkService) GetShortlink(ctx context.Context, id string) (string, error) {
	if link, ok := s.cache.Get(id); ok {
		return link, nil
	}

	link, err := s.repo.GetShortlink(ctx, id)
	if err != nil {
		return "", err
	}
	if link == nil {
		return "", ErrShortlinkNotexist
	}
	if time.Now().After(link.Expiration) {
		return "", ErrShortlinkExpired
	}

	return link.Url, nil
}

func (s *shortlinkService) ShortlinkRoutine(ctx context.Context) error {
	return nil
}
