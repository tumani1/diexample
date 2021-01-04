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

type autoCompleteHandlerSuite struct {
	suite.Suite

	ctx              context.Context
	controller       *gomock.Controller
	httpServer       *echo.Echo
	logger           *zap.Logger
	validator        *validator.Validate
	autoCompleteRepo *mocks.MockIAutoCompleteRepository
}

func (s *autoCompleteHandlerSuite) SetupTest() {
	s.controller = gomock.NewController(s.T())

	s.ctx = context.Background()
	s.httpServer = echo.New()
	s.logger = zap.NewNop()
	s.validator = validator.New()
	s.autoCompleteRepo = mocks.NewMockIAutoCompleteRepository(s.controller)
}

func (s *autoCompleteHandlerSuite) TearDownTest() {
	s.controller.Finish()
	_ = s.httpServer.Shutdown(s.ctx)
}

func (s *autoCompleteHandlerSuite) TestValidateError() {
	c := NewAutoCompleteHandler(s.logger, s.validator, s.autoCompleteRepo)
	c.Serve(s.httpServer)

	req := httptest.NewRequest(http.MethodGet, autoCompletePathURL, nil)
	rec := httptest.NewRecorder()
	s.httpServer.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusBadRequest, rec.Code)
	require.Equal(s.T(), `{"message":"error validate request"}`, strings.Trim(rec.Body.String(), "\n"))
}

func (s *autoCompleteHandlerSuite) TestGetAutoCompleteDataError() {
	s.autoCompleteRepo.EXPECT().
		Find(s.ctx, &domain.AutoCompleteFilter{
			Query:  "test",
			Limit:  autoCompleteLimitEntries,
			Offset: 0,
		}).
		Return(nil, fmt.Errorf("some error"))

	c := NewAutoCompleteHandler(s.logger, s.validator, s.autoCompleteRepo)
	c.Serve(s.httpServer)

	req := httptest.NewRequest(http.MethodGet, autoCompletePathURL+"?query=test", nil)
	rec := httptest.NewRecorder()
	s.httpServer.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusBadRequest, rec.Code)
	require.Equal(s.T(), `{"message":"error find auto complete by query"}`, strings.Trim(rec.Body.String(), "\n"))
}

func (s *autoCompleteHandlerSuite) TestSuccess() {
	s.autoCompleteRepo.EXPECT().
		Find(s.ctx, &domain.AutoCompleteFilter{
			Query:  "test",
			Limit:  autoCompleteLimitEntries,
			Offset: 0,
		}).
		Return([]*domain.AutoCompleteEntry{
			{
				Name: "test1",
			},
			{
				Name: "test2",
			},
		}, nil)

	c := NewAutoCompleteHandler(s.logger, s.validator, s.autoCompleteRepo)
	c.Serve(s.httpServer)

	req := httptest.NewRequest(http.MethodGet, autoCompletePathURL+"?query=test", nil)
	rec := httptest.NewRecorder()
	s.httpServer.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusOK, rec.Code)
	require.Equal(s.T(), `["test1","test2"]`, strings.Trim(rec.Body.String(), "\n"))
}

func TestAutoCompleteHandlerSuite(t *testing.T) {
	suite.Run(t, new(autoCompleteHandlerSuite))
}
