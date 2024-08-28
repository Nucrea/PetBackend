package services

import (
	"backend/src/cache"
	"backend/src/charsets"
	"fmt"
	"math/rand"
	"time"
)

type ShortlinkService interface {
	CreateLink(in string) (string, error)
	GetLink(id string) (string, error)
}

type NewShortlinkServiceParams struct {
	Endpoint string
	Cache    cache.Cache[string, string]
}

func NewShortlinkSevice(params NewShortlinkServiceParams) ShortlinkService {
	return &shortlinkService{
		cache: params.Cache,
	}
}

type shortlinkService struct {
	cache cache.Cache[string, string]
}

func (s *shortlinkService) CreateLink(in string) (string, error) {
	charset := charsets.GetCharset(charsets.CharsetTypeAll)

	src := rand.NewSource(time.Now().UnixMicro())
	randGen := rand.New(src)
	str := charset.RandomString(randGen, 10)

	s.cache.Set(str, in, cache.Expiration{Ttl: 7 * 24 * time.Hour})
	return str, nil
}

func (s *shortlinkService) GetLink(id string) (string, error) {
	val, ok := s.cache.Get(id)
	if !ok {
		return "", fmt.Errorf("link does not exist or expired")
	}
	return val, nil
}
