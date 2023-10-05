package handler

import (
	"errors"
	"fmt"
	"multi-output-data-processor/internal/config"
	"multi-output-data-processor/internal/entity"
	"multi-output-data-processor/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	srv    service.Pipeliner
	config config.Config
}

func NewHandler(srv service.Pipeliner, config config.Config) *Handler {
	h := &Handler{
		srv:    srv,
		config: config,
	}

	return h
}

// InitRoutes initializes routes for the service.
func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	{
		router.POST("/process", h.processInputData)
	}
	return router
}

// ProcessInputData processes input data.
func (h *Handler) processInputData(c *gin.Context) {
	var input entity.InputData
	if err := c.ShouldBindJSON(&input); err != nil {
		entity.NewErrorResponse(c, http.StatusBadRequest, "invalid input")
		return
	}

	err := h.srv.ValidateInput(c.Request.Context(), input)
	if err != nil {
		h.checkErrHandler(c, err)
		return
	}

	outputCh := h.srv.SelectOutputCh(c, h.config, input)

	go h.srv.Process(c, input, outputCh)

	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "OK",
	})
}

// CheckErrHandler returns an error response if the error matches any of the predefined errors.
func (h *Handler) checkErrHandler(c *gin.Context, err error) error {
	switch {
	case errors.Is(err, entity.ErrEmptyTagParameter):
		entity.NewErrorResponse(c, http.StatusInternalServerError, entity.EmptyTagParameterReason)
	case errors.Is(err, entity.ErrInvalidTagParameter):
		entity.NewErrorResponse(c, http.StatusInternalServerError, entity.InvalidTagParameterReason)
	case errors.Is(err, entity.ErrEmptyDataParameter):
		entity.NewErrorResponse(c, http.StatusInternalServerError, entity.EmptyDataParameterReason)
	}

	fmt.Println("err:", err)

	return err
}
