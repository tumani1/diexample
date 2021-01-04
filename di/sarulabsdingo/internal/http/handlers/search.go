package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"

	"github.com/tumani1/diexample/di/sarulabsdingo/echo"
	echoErrors "github.com/tumani1/diexample/di/sarulabsdingo/echo/errors"
	"github.com/tumani1/diexample/di/sarulabsdingo/internal/domain"
)

const (
	searchPathURL      = "/api/search"
	searchLimitEntries = 20
)

type (
	searchRequest struct {
		Query string `query:"query" validate:"required,min=1,max=255"`
	}

	searchResponse struct {
		ModelName     string    `json:"model_name"`
		BuilderName   string    `json:"builder_name"`
		FleetName     string    `json:"fleet_name"`
		Available     bool      `json:"available"`
		AvailableFrom time.Time `json:"available_from,omitempty"`
		AvailableTo   time.Time `json:"available_to,omitempty"`
	}

	searchAdditionalData struct {
		AvailableMap map[int64]bool
		UpcomingMap  map[int64]*domain.CalendarEntry
	}

	searchHandler struct {
		logger       *zap.Logger
		validator    *validator.Validate
		searchRepo   domain.ISearchRepository
		calendarRepo domain.ICalendarRepository
	}
)

func NewSearchHandler(
	logger *zap.Logger,
	validator *validator.Validate,
	searchRepo domain.ISearchRepository,
	calendarRepo domain.ICalendarRepository,
) echo.Handler {
	return &searchHandler{
		logger:       logger,
		validator:    validator,
		searchRepo:   searchRepo,
		calendarRepo: calendarRepo,
	}
}

func (h *searchHandler) Serve(e *echo.Echo) {
	e.GET(searchPathURL, h.searchHandler)
}

func (h *searchHandler) searchHandler(c echo.Context) (err error) {
	var request = new(searchRequest)
	if err = c.Bind(request); err != nil {
		return echoErrors.Wrap(http.StatusBadRequest, err, "error bind data")
	}

	if err = h.validator.Struct(request); err != nil {
		return echoErrors.Wrap(http.StatusBadRequest, err, "error validate request")
	}

	ctx := context.Background()
	entities, err := h.searchRepo.Find(ctx, &domain.SearchFilter{
		Query:  request.Query,
		Limit:  searchLimitEntries,
		Offset: 0,
	})
	if err != nil {
		h.logger.Error("error find by query",
			zap.String("query", request.Query), zap.Error(err),
		)

		return echoErrors.Wrap(http.StatusBadRequest, err, "error find by query")
	}

	var ids = make([]int64, 0)
	for _, entry := range entities {
		ids = append(ids, entry.BoatID)
	}

	additionalData, err := h.getAdditionalData(ctx, ids)
	if err != nil {
		h.logger.Error("error get additional data",
			zap.String("query", request.Query),
			zap.Int64s("ids", ids),
			zap.Error(err),
		)

		return echoErrors.Wrap(http.StatusBadRequest, err, "error get additional data")
	}

	var response = make([]*searchResponse, len(entities))
	for index, entity := range entities {
		response[index] = h.adaptForResponse(entity, additionalData)
	}

	return c.JSON(http.StatusOK, response)
}

func (h *searchHandler) getAdditionalData(ctx context.Context, ids []int64) (_ *searchAdditionalData, err error) {
	var data = &searchAdditionalData{}
	if data.AvailableMap, err = h.calendarRepo.GetAvailabilityByIDs(ctx, ids); err != nil {
		return nil, errors.Wrap(err, "error get availability by ids")
	}

	if data.UpcomingMap, err = h.calendarRepo.GetUpcomingAvailabilityDatesByIDs(ctx, ids); err != nil {
		return nil, errors.Wrap(err, "error get availability by ids")
	}

	return data, nil
}

func (h *searchHandler) adaptForResponse(
	entry *domain.SearchEntry, additionalData *searchAdditionalData,
) *searchResponse {
	respEntry := &searchResponse{
		ModelName:   entry.ModelName,
		BuilderName: entry.BuilderName,
		FleetName:   entry.FleetName,
		Available:   false,
	}

	if available, ok := additionalData.AvailableMap[entry.BoatID]; ok {
		respEntry.Available = available
	}

	if upcomingEntry, ok := additionalData.UpcomingMap[entry.BoatID]; ok {
		respEntry.AvailableFrom = upcomingEntry.DateFrom
		respEntry.AvailableTo = upcomingEntry.DateTo
	}

	return respEntry
}
