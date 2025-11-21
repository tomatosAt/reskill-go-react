package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tomatosAt/reskill-go-react/pkg/util"
)

func (h *Handler) GetProfilesHandler(ctx *fiber.Ctx) error {
	profiles, status, err := h.svc.GetProfilesService(ctx.UserContext())
	if err != nil {
		return util.HttpError(ctx, status, err.Error())
	}
	return util.HttpSuccess(ctx, 200, profiles)
}
