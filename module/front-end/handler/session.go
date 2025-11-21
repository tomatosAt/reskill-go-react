package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/tomatosAt/reskill-go-react/module/front-end/dto"
	"github.com/tomatosAt/reskill-go-react/pkg/util"
)

func (h *Handler) CreateSessionHandler(ctx *fiber.Ctx) error {
	var payload dto.ResponsePreRegister
	if err := ctx.BodyParser(&payload); err != nil {
		return util.HttpError(ctx, http.StatusBadRequest, err.Error())
	}
	if _, err := util.ValidatorStruct(&payload); err != nil {
		return err
	}
	claimPreRegData, err := h.svc.CreateSessionService(ctx.UserContext(), payload.Code)
	if err != nil {
		return util.HttpError(ctx, http.StatusBadRequest, err.Error())
	}
	return util.HttpSuccess(ctx, 200, claimPreRegData)
}

func (h *Handler) GetSessionsHandler(ctx *fiber.Ctx) error {
	return util.HttpSuccess(ctx, fiber.StatusOK, "success")
}
