package handler

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/tomatosAt/reskill-go-react/pkg/util"
)

func (h *Handler) ConvertMP4ToHLS(ctx *fiber.Ctx) error {
	fileHeader, err := ctx.FormFile("data")
	if err != nil {
		return util.HttpError(ctx, http.StatusBadRequest, "error uploading file")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return util.HttpError(ctx, http.StatusInternalServerError, "cannot open file")
	}
	defer file.Close()

	// ใช้ service ย่อย
	uploadDir := "react/module/openapi/uploads"
	uploadPath, err := h.svc.SaveFile(ctx.Context(), file, fileHeader.Filename, uploadDir)
	if err != nil {
		return util.HttpError(ctx, http.StatusInternalServerError, err.Error())
	}

	outputDir := "hls"
	outputPath, err := h.svc.ConvertToHLS(ctx.Context(), uploadPath, outputDir)
	if err != nil {
		return util.HttpError(ctx, http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(fiber.Map{
		"message": "success",
		"hls_url": fmt.Sprintf("/%s/%s", outputDir, filepath.Base(outputPath)),
	})
}
