package vesselsapi

import "example.com/vesssels-api/pkg/vessels"

type Store interface {
	GetVesselByIMO(imo int) (*vessels.Vessel, error)
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

func (s *inMemoryStore) GetVesselByIMO(imo int) (*vessels.Vessel, error) {
	if vessel, ok := s.data[imo]; ok {
		return vessel, nil
	}
	return nil, nil
}
