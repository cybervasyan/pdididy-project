package part

import (
	"testing"

	"github.com/cybervasyan/pdididy-project/inventory/internal/repository/mocks"
	"github.com/stretchr/testify/suite"
)

type ServiceSuite struct {
	suite.Suite

	inventoryRepo *mocks.MockRepository

	service *service
}

func (s *ServiceSuite) SetupTest() {
	s.inventoryRepo = mocks.NewMockRepository(s.T())

	s.service = NewPartService(
		s.inventoryRepo,
	)
}

func TestAPI(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
