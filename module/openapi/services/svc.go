package services

import (
	"github.com/tomatosAt/reskill-go-react/module/openapi/ports"
	"github.com/tomatosAt/reskill-go-react/module/openapi/repositories"
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
