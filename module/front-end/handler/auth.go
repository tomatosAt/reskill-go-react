package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/tomatosAt/reskill-go-react/module/front-end/dto"
	"github.com/tomatosAt/reskill-go-react/pkg/util"
)

func (h *Handler) LoginUserPassHandler(ctx *fiber.Ctx) error {
	var payload dto.UserPasswordPayload
	if err := ctx.BodyParser(&payload); err != nil {
		return util.HttpError(ctx, http.StatusBadRequest, err.Error())
	}
	//  : Process Login
	res, status, err := h.svc.LoginUserPassService(ctx.UserContext(), payload)
	if err != nil {
		return util.HttpError(ctx, status, err.Error())
	}
	return util.HttpSuccess(ctx, status, res)
}

func (h *Handler) LogoutHandler(ctx *fiber.Ctx) error {
	status, err := h.svc.LogoutService(ctx.UserContext())
	if err != nil {
		return util.HttpError(ctx, status, err.Error())
	}
	return util.HttpSuccess(ctx, status, "logout successful")
}
