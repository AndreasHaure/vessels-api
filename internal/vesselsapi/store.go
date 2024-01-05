package vesselsapi

import "example.com/vesssels-api/pkg/vessels"

type Store interface {
	RunInTx(f func(store Store) error) error
	UpdateVessel(imo int, vessel *vessels.UpdateVessel) error
	GetVesselByIMO(imo int) (*vessels.Vessel, error)
	GetVessels() ([]*vessels.Vessel, error)
}

type inMemoryStore struct {
	data map[int]*vessels.Vessel
}

func NewInMemoryStore() Store {
	return &inMemoryStore{
		data: map[int]*vessels.Vessel{
			1234567: {
				IMO:       1234567,
				Name:      "Test Vessel",
				Flag:      "US",
				YearBuilt: 2010,
				Owner:     "Test Owner",
			},
		},
	}
}

// RunInTx runs the given function in a transaction. If the function returns an error, the transaction is rolled back.
// For in-memory store, this is a no-op.
func (s *inMemoryStore) RunInTx(f func(store Store) error) error {
	return f(s)
}

func (s *inMemoryStore) UpdateVessel(imo int, vessel *vessels.UpdateVessel) error {
	s.data[imo] = &vessels.Vessel{
		IMO:       int64(imo),
		Name:      vessel.Name,
		Flag:      vessel.Flag,
		YearBuilt: vessel.YearBuilt,
		Owner:     vessel.Owner,
	}
	return nil
}

func (s *inMemoryStore) GetVesselByIMO(imo int) (*vessels.Vessel, error) {
	if vessel, ok := s.data[imo]; ok {
		return vessel, nil
	}
	return nil, nil
}

func (s *inMemoryStore) GetVessels() ([]*vessels.Vessel, error) {
	var vessels []*vessels.Vessel
	for _, vessel := range s.data {
		vessels = append(vessels, vessel)
	}
	return vessels, nil
}
