package mocks

import "github.com/stretchr/testify/mock"

type MockRateService struct {
	mock.Mock
}

func (m *MockRateService) Get() (*float64, error) {
	args := m.Called()
	if args.Get(0) != nil {
		rate := args.Get(0).(float64)
		return &rate, args.Error(1)
	}
	return nil, args.Error(1)
}
