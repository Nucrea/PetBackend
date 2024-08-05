package services

import (
	"backend/src/repo"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type ShortlinkService interface {
	CreateLink(in string) (string, error)
	GetLink(id string) (string, error)
}

type NewShortlinkServiceParams struct {
	Endpoint string
	Cache    repo.Cache[string, string]
}

func NewShortlinkSevice(params NewShortlinkServiceParams) ShortlinkService {
	return &shortlinkService{
		cache: params.Cache,
	}
}

type shortlinkService struct {
	cache repo.Cache[string, string]
}

func (s *shortlinkService) randomStr() string {
	src := rand.NewSource(time.Now().UnixMicro())
	randGen := rand.New(src)

	builder := strings.Builder{}
	for i := 0; i < 9; i++ {
		offset := 0x41
		if randGen.Int()%2 == 1 {
			offset = 0x61
		}

		byte := offset + (randGen.Int() % 26)
		builder.WriteRune(rune(byte))
	}
	return builder.String()
}

func (s *shortlinkService) CreateLink(in string) (string, error) {
	str := s.randomStr()
	s.cache.Set(str, in, 7*24*60*60)
	return str, nil
}

func (s *shortlinkService) GetLink(id string) (string, error) {
	val, ok := s.cache.Get(id)
	if !ok {
		return "", fmt.Errorf("link does not exist or expired")
	}
	return val, nil
}
