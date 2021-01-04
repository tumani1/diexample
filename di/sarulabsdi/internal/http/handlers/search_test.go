package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"

	"github.com/tumani1/diexample/di/sarulabsdi/echo"
	"github.com/tumani1/diexample/di/sarulabsdi/internal/domain"
	"github.com/tumani1/diexample/di/sarulabsdi/internal/domain/mocks"
)

type searchHandlerSuite struct {
	suite.Suite

	ctx          context.Context
	controller   *gomock.Controller
	httpServer   *echo.Echo
	logger       *zap.Logger
	validator    *validator.Validate
	searchRepo   *mocks.MockISearchRepository
	calendarRepo *mocks.MockICalendarRepository
}

func (s *searchHandlerSuite) SetupTest() {
	s.controller = gomock.NewController(s.T())

	s.ctx = context.Background()
	s.httpServer = echo.New()
	s.logger = zap.NewNop()
	s.validator = validator.New()
	s.searchRepo = mocks.NewMockISearchRepository(s.controller)
	s.calendarRepo = mocks.NewMockICalendarRepository(s.controller)
}

func (s *searchHandlerSuite) TearDownTest() {
	s.controller.Finish()
	_ = s.httpServer.Shutdown(s.ctx)
}

func (s *searchHandlerSuite) TestValidateError() {
	c := NewSearchHandler(s.logger, s.validator, s.searchRepo, s.calendarRepo)
	c.Serve(s.httpServer)

	req := httptest.NewRequest(http.MethodGet, searchPathURL, nil)
	rec := httptest.NewRecorder()
	s.httpServer.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusBadRequest, rec.Code)
	require.Equal(s.T(), "{\"message\":\"error validate request\"}", strings.Trim(rec.Body.String(), "\n"))
}

func (s *searchHandlerSuite) TestGetSearchDataError() {
	s.searchRepo.EXPECT().
		Find(s.ctx, &domain.SearchFilter{
			Query:  "test",
			Limit:  searchLimitEntries,
			Offset: 0,
		}).
		Return(nil, fmt.Errorf("some error"))

	c := NewSearchHandler(s.logger, s.validator, s.searchRepo, s.calendarRepo)
	c.Serve(s.httpServer)

	req := httptest.NewRequest(http.MethodGet, searchPathURL+"?query=test", nil)
	rec := httptest.NewRecorder()
	s.httpServer.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusBadRequest, rec.Code)
	require.Equal(s.T(), `{"message":"error find by query"}`, strings.Trim(rec.Body.String(), "\n"))
}

func (s *searchHandlerSuite) TestGetAvailabilityDataError() {
	s.searchRepo.EXPECT().
		Find(s.ctx, &domain.SearchFilter{
			Query:  "test",
			Limit:  searchLimitEntries,
			Offset: 0,
		}).
		Return([]*domain.SearchEntry{
			{
				BoatID:      1,
				ModelName:   "model_name 1",
				BuilderName: "builder_name 1",
				FleetName:   "fleet_name 1",
			},
			{
				BoatID:      2,
				ModelName:   "model_name 2",
				BuilderName: "builder_name 2",
				FleetName:   "fleet_name 2",
			},
			{
				BoatID:      3,
				ModelName:   "model_name 3",
				BuilderName: "builder_name 3",
				FleetName:   "fleet_name 3",
			},
		}, nil)

	s.calendarRepo.EXPECT().
		GetAvailabilityByIDs(s.ctx, []int64{1, 2, 3}).
		Return(nil, fmt.Errorf("some error"))

	c := NewSearchHandler(s.logger, s.validator, s.searchRepo, s.calendarRepo)
	c.Serve(s.httpServer)

	req := httptest.NewRequest(http.MethodGet, searchPathURL+"?query=test", nil)
	rec := httptest.NewRecorder()
	s.httpServer.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusBadRequest, rec.Code)
	require.Equal(s.T(), `{"message":"error get additional data"}`, strings.Trim(rec.Body.String(), "\n"))
}

func (s *searchHandlerSuite) TestGetUpcomingAvailabilityDataError() {
	s.searchRepo.EXPECT().
		Find(s.ctx, &domain.SearchFilter{
			Query:  "test",
			Limit:  searchLimitEntries,
			Offset: 0,
		}).
		Return([]*domain.SearchEntry{
			{
				BoatID:      1,
				ModelName:   "model_name 1",
				BuilderName: "builder_name 1",
				FleetName:   "fleet_name 1",
			},
			{
				BoatID:      2,
				ModelName:   "model_name 2",
				BuilderName: "builder_name 2",
				FleetName:   "fleet_name 2",
			},
			{
				BoatID:      3,
				ModelName:   "model_name 3",
				BuilderName: "builder_name 3",
				FleetName:   "fleet_name 3",
			},
		}, nil)

	s.calendarRepo.EXPECT().
		GetAvailabilityByIDs(s.ctx, []int64{1, 2, 3}).
		Return(map[int64]bool{}, nil)

	s.calendarRepo.EXPECT().
		GetUpcomingAvailabilityDatesByIDs(s.ctx, []int64{1, 2, 3}).
		Return(nil, fmt.Errorf("some error"))

	c := NewSearchHandler(s.logger, s.validator, s.searchRepo, s.calendarRepo)
	c.Serve(s.httpServer)

	req := httptest.NewRequest(http.MethodGet, searchPathURL+"?query=test", nil)
	rec := httptest.NewRecorder()
	s.httpServer.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusBadRequest, rec.Code)
	require.Equal(s.T(), `{"message":"error get additional data"}`, strings.Trim(rec.Body.String(), "\n"))
}

func (s *searchHandlerSuite) TestSuccess() {
	s.searchRepo.EXPECT().
		Find(s.ctx, &domain.SearchFilter{
			Query:  "test",
			Limit:  searchLimitEntries,
			Offset: 0,
		}).
		Return([]*domain.SearchEntry{
			{
				BoatID:      1,
				ModelName:   "model_name 1",
				BuilderName: "builder_name 1",
				FleetName:   "fleet_name 1",
			},
			{
				BoatID:      2,
				ModelName:   "model_name 2",
				BuilderName: "builder_name 2",
				FleetName:   "fleet_name 2",
			},
			{
				BoatID:      3,
				ModelName:   "model_name 3",
				BuilderName: "builder_name 3",
				FleetName:   "fleet_name 3",
			},
		}, nil)

	s.calendarRepo.EXPECT().
		GetAvailabilityByIDs(s.ctx, []int64{1, 2, 3}).
		Return(map[int64]bool{
			1: false,
			2: true,
		}, nil)

	s.calendarRepo.EXPECT().
		GetUpcomingAvailabilityDatesByIDs(s.ctx, []int64{1, 2, 3}).
		Return(map[int64]*domain.CalendarEntry{
			1: {
				BoatID:    1,
				Available: true,
			},
			2: {
				BoatID:    2,
				Available: true,
			},
			3: {
				BoatID:    3,
				Available: true,
			},
		}, nil)

	c := NewSearchHandler(s.logger, s.validator, s.searchRepo, s.calendarRepo)
	c.Serve(s.httpServer)

	req := httptest.NewRequest(http.MethodGet, searchPathURL+"?query=test", nil)
	rec := httptest.NewRecorder()
	s.httpServer.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusOK, rec.Code)
	require.Equal(s.T(),
		`[{"model_name":"model_name 1","builder_name":"builder_name 1","fleet_name":"fleet_name 1","available":false,"available_from":"0001-01-01T00:00:00Z","available_to":"0001-01-01T00:00:00Z"},{"model_name":"model_name 2","builder_name":"builder_name 2","fleet_name":"fleet_name 2","available":true,"available_from":"0001-01-01T00:00:00Z","available_to":"0001-01-01T00:00:00Z"},{"model_name":"model_name 3","builder_name":"builder_name 3","fleet_name":"fleet_name 3","available":false,"available_from":"0001-01-01T00:00:00Z","available_to":"0001-01-01T00:00:00Z"}]`,
		strings.Trim(rec.Body.String(), "\n"),
	)
}

func TestSearchHandlerSuite(t *testing.T) {
	suite.Run(t, new(searchHandlerSuite))
}
