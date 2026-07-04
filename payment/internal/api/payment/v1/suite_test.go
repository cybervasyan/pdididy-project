package v1

import (
	"testing"

	"github.com/cybervasyan/pdididy-project/payment/internal/service/mocks"
	"github.com/stretchr/testify/suite"
)

type APISuite struct {
	suite.Suite

	paymentService *mocks.MockPayment

	api *api
}

func (s *APISuite) SetupTest() {
	s.paymentService = mocks.NewMockPayment(s.T())

	s.api = NewAPI(
		s.paymentService,
	)
}

func TestAPI(t *testing.T) {
	suite.Run(t, new(APISuite))
}
