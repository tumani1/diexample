package handlers

import (
	"context"
	"net/http"

	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"

	"github.com/tumani1/diexample/di/sarulabsdingo/echo"
	"github.com/tumani1/diexample/di/sarulabsdingo/echo/errors"
	"github.com/tumani1/diexample/di/sarulabsdingo/internal/domain"
)

const (
	autoCompletePathURL      = "/api/autocomplete"
	autoCompleteLimitEntries = 10
)

type (
	autoCompleteRequest struct {
		Query string `query:"query" validate:"required,min=1,max=255"`
	}

	autoCompleteHandler struct {
		logger           *zap.Logger
		validator        *validator.Validate
		autoCompleteRepo domain.IAutoCompleteRepository
	}
)

func NewAutoCompleteHandler(
	logger *zap.Logger,
	validator *validator.Validate,
	autoCompleteRepo domain.IAutoCompleteRepository,
) echo.Handler {
	return &autoCompleteHandler{
		logger:           logger,
		validator:        validator,
		autoCompleteRepo: autoCompleteRepo,
	}
}

func (h *autoCompleteHandler) Serve(e *echo.Echo) {
	e.GET(autoCompletePathURL, h.autoCompleteHandler)
}

func (h *autoCompleteHandler) autoCompleteHandler(c echo.Context) (err error) {
	var request = new(autoCompleteRequest)
	if err = c.Bind(request); err != nil {
		return errors.Wrap(http.StatusBadRequest, err, "error bind data")
	}

	if err = h.validator.Struct(request); err != nil {
		return errors.Wrap(http.StatusBadRequest, err, "error validate request")
	}

	ctx := context.Background()
	entries, err := h.autoCompleteRepo.Find(ctx, &domain.AutoCompleteFilter{
		Query:  request.Query,
		Limit:  autoCompleteLimitEntries,
		Offset: 0,
	})
	if err != nil {
		h.logger.Error("error find auto complete by query",
			zap.String("query", request.Query), zap.Error(err),
		)
		return errors.Wrap(http.StatusBadRequest, err, "error find auto complete by query")
	}

	var response = make([]string, len(entries))
	for index, entry := range entries {
		response[index] = entry.Name
	}

	return c.JSON(http.StatusOK, response)
}
