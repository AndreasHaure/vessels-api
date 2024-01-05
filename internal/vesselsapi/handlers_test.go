package vesselsapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/vesssels-api/pkg/vessels"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type HandlerSuite struct {
	suite.Suite
	mockStore *MockStore
}

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}

func (s *HandlerSuite) SetupTest() {
	s.mockStore = &MockStore{}
}

func (s *HandlerSuite) TestGetVessels() {
	handler := Handler{
		Store: s.mockStore,
	}
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)

	r.GET("/v1/vessels", handler.GetVessels)
	req, err := http.NewRequestWithContext(ctx, "GET", "/v1/vessels", nil)
	require.NoError(s.T(), err)
	vessels := []*vessels.Vessel{{
		IMO:       1234567,
		Name:      "Test Vessel",
		Flag:      "US",
		YearBuilt: 2010,
		Owner:     "Test Owner",
	}}
	s.mockStore.On("GetVessels").Return(vessels, nil)
	r.ServeHTTP(w, req)
	require.Equal(s.T(), http.StatusOK, w.Code)
}
