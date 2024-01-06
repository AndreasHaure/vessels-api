//go:build integration
// +build integration

package vesselsapi

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	_ "github.com/lib/pq" // Import Postgres driver

	"example.com/vesssels-api/pkg/vessels"
	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type StoreSuite struct {
	suite.Suite
	db         *sql.DB
	schemaName string
}

func TestStoreSuite(t *testing.T) {
	suite.Run(t, new(StoreSuite))
}

func (s *StoreSuite) SetupSuite() {
	type testConfig struct {
		Postgres struct {
			User       string `default:"postgres"`
			Password   string `default:"postgres"`
			Host       string `default:"localhost"` //`default:"db"`
			Port       int    `default:"5432"`
			DBName     string `envconfig:"DB_NAME" default:"postgres"`
			SchemaName string `envconfig:"SCHEMA_NAME" default:"vessels"`
		}
	}

	var tc testConfig
	if err := envconfig.Process("", &tc); err != nil {
		s.T().Fatalf("Unable to process config: %s", err)
	}

	db, err := sql.Open(
		"postgres",
		fmt.Sprintf(
			`host=%s port=%d user=%s password=%s dbname=%s options=--search_path=%s sslmode=disable`,
			tc.Postgres.Host, tc.Postgres.Port, tc.Postgres.User, tc.Postgres.Password, tc.Postgres.DBName, tc.Postgres.SchemaName,
		),
	)
	if err != nil {
		s.T().Fatalf("Unable to setup postgres connection (host=%s, port=%d, user=%s, password=%s, dbname=%s): %s", tc.Postgres.Host, tc.Postgres.Port, tc.Postgres.User, tc.Postgres.Password, tc.Postgres.DBName, err)
	}
	s.db = db

	s.T().Log(fmt.Sprintf("Connected (host=%s, port=%d, user=%s, password=%s, dbname=%s, schema=%s)", tc.Postgres.Host, tc.Postgres.Port, tc.Postgres.User, tc.Postgres.Password, tc.Postgres.DBName, tc.Postgres.SchemaName))

	// wait for the test database to be ready
	var success bool
	for attempt := 0; attempt < 5; attempt++ {
		if err := db.Ping(); err != nil {
			s.T().Log("Database not alive yet, waiting to try again...")
			time.Sleep(1 * time.Second)
			continue
		}
		success = true
		break
	}

	if !success {
		s.T().Fatal("Giving up waiting for database to come alive.")
	}
}

func (s *StoreSuite) SetupTest() {
	s.clearDB()
}

func (s *StoreSuite) TestGetVessels() {
	store := NewPGStore(s.db, s.schemaName)

	vessels := []vessels.Vessel{
		{
			IMO:       1,
			Name:      "Test Vessel 1",
			Flag:      "US",
			YearBuilt: 2000,
			Owner:     "Test Owner 1",
		},
		{
			IMO:       2,
			Name:      "Test Vessel 2",
			Flag:      "US",
			YearBuilt: 2000,
			Owner:     "Test Owner 2",
		},
	}
	s.addVesselsToDB(vessels, store)

	insertedVessels, err := store.GetVessels()
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), vessels, insertedVessels)
}

func (s *StoreSuite) TearDownTest() {
	s.db.Close()
}

func (s *StoreSuite) clearDB() {
	// Clear databases
	_, err := s.db.Exec(`DELETE FROM vessels.vessels`)
	require.NoError(s.T(), err)
}

func (s *StoreSuite) addVesselsToDB(insertVessels []vessels.Vessel, store Store) {
	for _, vessel := range insertVessels {
		err := store.UpdateVessel(int(vessel.IMO), &vessels.UpdateVessel{
			Name:      vessel.Name,
			Flag:      vessel.Flag,
			YearBuilt: vessel.YearBuilt,
			Owner:     vessel.Owner,
		})
		require.NoError(s.T(), err)
	}
}
