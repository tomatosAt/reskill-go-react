package services

import (
	"github.com/tomatosAt/reskill-go-react/module/open-api/ports"
	"github.com/tomatosAt/reskill-go-react/module/open-api/repositories"
)

/**
Business logic was here.
*/

type Service struct {
	repo ports.Repository
}

func New(repo *repositories.Repository) *Service {
	return &Service{repo: repo}
}
