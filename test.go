package hentities

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/stretchr/testify/suite"
)

// RepositoryTestSuite is a testify suite.Suite that provides the following:
// - apply metadata before the test suite runs
// - apply migrations before a test in the suite runs
// - rollback migrations after a test in the suite is complete
// - rollback metadata after the test suite completes
//
// This idea for this is that for a test suite, the metadata will always be
// updated. Each test will also have a fresh set of database before it runs.
type RepositoryTestSuite struct {
	suite.Suite
	AppRoot     string
	Endpoint    string
	AdminSecret string
}

// SetupSuite applies the metadata and reloads the Hasura engine.
func (s *RepositoryTestSuite) SetupSuite() {
	s.T().Log("==========> Setting Up Test Suite")
	s.cmd("hasura migrate apply --up all --all-databases --endpoint %s --admin-secret %s", s.Endpoint, s.AdminSecret)
	s.cmd("hasura metadata apply --endpoint %s --admin-secret %s", s.Endpoint, s.AdminSecret)
	s.cmd("hasura metadata reload --endpoint %s --admin-secret %s", s.Endpoint, s.AdminSecret)
}

// SetupSuite rolls back the metadata and reloads the Hasura engine.
func (s *RepositoryTestSuite) TearDownSuite() {
	s.T().Log("==========> Tearing Down Test Suite")
	s.cmd("hasura metadata clear --endpoint %s --admin-secret %s", s.Endpoint, s.AdminSecret)
	s.cmd("hasura metadata reload --endpoint %s --admin-secret %s", s.Endpoint, s.AdminSecret)
}

// SetupTest applies the migrations and reloads the Hasura engine.
func (s *RepositoryTestSuite) SetupTest() {
	s.T().Log("==========> Setting Up Test")
	s.cmd("hasura migrate apply --up all --all-databases --endpoint %s --admin-secret %s", s.Endpoint, s.AdminSecret)
}

// TearDownTest rolls back the migrations and reloads the Hasura engine.
func (s *RepositoryTestSuite) TearDownTest() {
	s.T().Log("==========> Tearing Down Test")
	s.cmd("hasura migrate apply --down all --all-databases --endpoint %s --admin-secret %s", s.Endpoint, s.AdminSecret)
}

func (s *RepositoryTestSuite) cmd(command string, args ...interface{}) {
	fullCommand := fmt.Sprintf(command, args...)
	f := strings.Fields(fullCommand)
	s.Require().GreaterOrEqual(len(f), 2)
	e := exec.Command(f[0], f[1:]...)
	e.Dir = s.AppRoot
	e.Stdout = os.Stdout
	e.Stderr = os.Stderr
	err := e.Run()
	s.Require().NoError(err)
}
