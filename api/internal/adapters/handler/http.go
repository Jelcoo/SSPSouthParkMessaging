package handler

import (
	"net/http"
	"sp-messaging/api/internal/core/domain"
	"sp-messaging/api/internal/core/ports"

	"github.com/gin-gonic/gin"
)

type HTTPHandler struct {
	svc ports.MessengerService
}

func NewHTTPHandler(svc ports.MessengerService) *HTTPHandler {
	return &HTTPHandler{
		svc: svc,
	}
}

type postMessageRequest struct {
	Author string `json:"author" binding:"required"`
	Body   string `json:"body" binding:"required"`
}

func (h *HTTPHandler) PostMessage(ctx *gin.Context) {
	var req postMessageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	posted, err := h.svc.PostMessage(domain.Message{
		Author: req.Author,
		Body:   req.Body,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusAccepted, posted)
}
