package handler

import (
	"github.com/tomatosAt/reskill-go-react/module/open-api/ports"
	"github.com/tomatosAt/reskill-go-react/module/open-api/services"
)

/**
API endpoint input/output controlling and validation
*/

type Handler struct {
	svc ports.Service
}

func NewHandler(svc *services.Service) *Handler {
	return &Handler{
		svc: svc,
	}
}
