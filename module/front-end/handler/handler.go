package handler

import (
	"github.com/tomatosAt/reskill-go-react/module/front-end/ports"
	"github.com/tomatosAt/reskill-go-react/module/front-end/services"
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
