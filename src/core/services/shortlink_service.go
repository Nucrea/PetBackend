package services

import (
	"backend/src/core/repos"
	"backend/src/core/utils"
	"fmt"
)

type ShortlinkService interface {
	CreateLink(in string) (string, error)
	GetLink(id string) (string, error)
}

type NewShortlinkServiceParams struct {
	Endpoint string
	Cache    repos.Cache[string, string]
}

func NewShortlinkSevice(params NewShortlinkServiceParams) ShortlinkService {
	return &shortlinkService{
		randomUtil: *utils.NewRand(),
		cache:      params.Cache,
	}
}

type shortlinkService struct {
	randomUtil utils.RandomUtil
	cache      repos.Cache[string, string]
}

func (s *shortlinkService) CreateLink(in string) (string, error) {
	str := s.randomUtil.RandomID(10, utils.CharsetAll)
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
