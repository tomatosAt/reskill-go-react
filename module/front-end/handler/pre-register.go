package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/tomatosAt/reskill-go-react/module/front-end/dto"
	"github.com/tomatosAt/reskill-go-react/pkg/util"
)

func (h *Handler) PreRegisterHandler(ctx *fiber.Ctx) error {
	var payload dto.PreRegisterDataBasePayload
	if err := ctx.BodyParser(&payload); err != nil {
		return util.HttpError(ctx, http.StatusBadRequest, err.Error())
	}
	// TODO : check format
	// CheckFormatPreRegisterSVC
	if err := h.svc.CheckFormatPreRegisterSVC(&payload); err != nil {
		return util.HttpError(ctx, http.StatusBadRequest, err.Error())
	}
	// TODO : Process เก็บข้อมูล
	// Encrpy ชื่อ นามสกุล password
	return nil
}
