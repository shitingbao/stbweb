package stbserver

import (
	"context"
	"stbweb/lib/external_service/stbserver"
)

const (
	posrt = ":5000"
)

type stbserve struct{}

func (s *stbserve) GetSummonerInfo(context.Context, *stbserver.Identity) (*stbserver.Character, error) {
	return nil, nil
}
func (s *stbserve) PutSummonerInfo(stbserver.StbServer_PutSummonerInfoServer) error {
	return nil
}
func (s *stbserve) GetAllSummonerInfo(*stbserver.Identity, stbserver.StbServer_GetAllSummonerInfoServer) error {
	return nil
}
func (s *stbserve) ShareSummonerInfo(stbserver.StbServer_ShareSummonerInfoServer) error {
	return nil
}
