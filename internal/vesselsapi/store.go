package vesselsapi

import (
	"database/sql"
	"fmt"

	"example.com/vesssels-api/pkg/vessels"
	"github.com/pkg/errors"
)

type Store interface {
	RunInTx(f func(store Store) error) error
	UpdateVessel(imo int, vessel *vessels.UpdateVessel) error
	GetVesselByIMO(imo int) (*vessels.Vessel, error)
	GetVessels() ([]*vessels.Vessel, error)
	DeleteVessel(imo int) error
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

func (s *inMemoryStore) DeleteVessel(imo int) error {
	delete(s.data, imo)
	return nil
}

type querier interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

type pgStore struct {
	querier    querier
	schemaName string
}

func NewPGStore(db *sql.DB, schemaName string) Store {
	return &pgStore{
		querier:    db,
		schemaName: schemaName,
	}
}

// RunInTx wraps the given function in a database transaction.
// If the function returns an error, the transaction is rolled back, otherwise it is committed.
func (s *pgStore) RunInTx(f func(store Store) error) error {
	// type assert querier to db, panic if fails
	db, ok := s.querier.(*sql.DB)
	if !ok {
		panic("Expected querier to be of type *sql.DB")
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	storeWithTx := &pgStore{querier: tx, schemaName: s.schemaName}
	if err = f(storeWithTx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (s *pgStore) UpdateVessel(imo int, vessel *vessels.UpdateVessel) error {
	upsertStmt := fmt.Sprintf(`
		INSERT INTO %[1]s.vessels (imo, name, flag, year_built, owner)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT ON CONSTRAINT vessels_pkey DO UPDATE SET
			name = $2,
			flag = $3,
			year_built = $4,
			owner = $5`,
		s.schemaName,
	)

	if _, err := s.querier.Exec(upsertStmt, imo, vessel.Name, vessel.Flag, vessel.YearBuilt, vessel.Owner); err != nil {
		return errors.Wrap(err, "unable to update vessel")
	}
	return nil
}

func (s *pgStore) GetVesselByIMO(imo int) (*vessels.Vessel, error) {
	query := fmt.Sprintf(`
		SELECT imo, name, flag, year_built, owner
		FROM %[1]s.vessels
		WHERE imo = $1`,
		s.schemaName,
	)

	row := s.querier.QueryRow(query, imo)

	var vessel vessels.Vessel
	err := row.Scan(&vessel.IMO, &vessel.Name, &vessel.Flag, &vessel.YearBuilt, &vessel.Owner)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "unable to get vessel by imo")
	}
	return &vessel, nil
}

func (s *pgStore) GetVessels() ([]*vessels.Vessel, error) {
	query := fmt.Sprintf(`
		SELECT imo, name, flag, year_built, owner
		FROM %[1]s.vessels`,
		s.schemaName,
	)

	result := []*vessels.Vessel{}
	rows, err := s.querier.Query(query)
	if err == sql.ErrNoRows {
		return result, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "unable to get vessels")
	}

	for rows.Next() {
		var vessel vessels.Vessel
		err := rows.Scan(&vessel.IMO, &vessel.Name, &vessel.Flag, &vessel.YearBuilt, &vessel.Owner)
		if err != nil {
			return nil, errors.Wrap(err, "unable to get vessels")
		}
		result = append(result, &vessel)
	}
	return result, nil
}
