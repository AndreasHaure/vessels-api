package vesselsapi

import (
	"example.com/vesssels-api/pkg/vessels"
	"github.com/stretchr/testify/mock"
)

type MockStore struct {
	mock.Mock
}

// RunInTx runs the given function in a transaction. If the function returns an error, the transaction is rolled back.
// For mocking purposes, this is a no-op.
func (m *MockStore) RunInTx(f func(store Store) error) error {
	return f(m)
}

func (m *MockStore) UpdateVessel(imo int, vessel *vessels.UpdateVessel) error {
	args := m.Called(imo, vessel)
	return args.Error(0)
}

func (m *MockStore) GetVesselByIMO(imo int) (*vessels.Vessel, error) {
	args := m.Called(imo)
	return args.Get(0).(*vessels.Vessel), args.Error(1)
}

func (m *MockStore) GetVessels() ([]*vessels.Vessel, error) {
	args := m.Called()
	return args.Get(0).([]*vessels.Vessel), args.Error(1)
}
