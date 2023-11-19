package port

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/stretchr/testify/suite"
)

// TestRepositorySuite contains a set of tests to cover a Repository's expected behaviors
// ignoring the implementation details.
// To run use:
//
//	suite.Run(t, &port.TestRepositorySuite{
//	    Repo: yourRepoImplementation,
//	})
type TestRepositorySuite struct {
	suite.Suite

	ctx    context.Context
	cancel context.CancelFunc

	Repo Repository
}

func (s *TestRepositorySuite) SetupTest() {
	s.ctx, s.cancel = context.WithCancel(context.Background())
}

func (s *TestRepositorySuite) TearDownTest() {
	s.cancel()
}

func (s *TestRepositorySuite) TestGetCreateUpdate() {
	// Make some ports with some mock data
	portsA := mockPorts(3, "A")

	s.Run("ports not found before creation", func() {
		for _, port := range portsA {
			p, err := s.Repo.GetByID(s.ctx, port.ID)
			s.Require().ErrorIs(err, ErrPortNotFound)
			s.Require().Nil(p)
		}
	})

	s.Run("create ports first time", func() {
		for _, port := range portsA {
			created, err := s.Repo.CreateOrUpdate(s.ctx, port)
			s.Require().NoError(err)
			s.Require().True(created) // created=true
		}
	})

	s.Run("update same ports again", func() {
		for _, port := range portsA {
			created, err := s.Repo.CreateOrUpdate(s.ctx, port)
			s.Require().NoError(err)
			s.Require().False(created) // created=false
		}
	})

	s.Run("get created/updated ports", func() {
		for _, port := range portsA {
			p, err := s.Repo.GetByID(s.ctx, port.ID)
			s.Require().NoError(err)
			s.Require().Equal(port, *p)
		}
	})

	// Make the same ports again, but with some data changed
	portsB := mockPorts(3, "B")
	s.Run("update ports", func() {
		for _, port := range portsB {
			created, err := s.Repo.CreateOrUpdate(s.ctx, port)
			s.Require().NoError(err)
			s.Require().False(created) // created=false
		}
	})

	s.Run("get updated ports", func() {
		for _, port := range portsB {
			p, err := s.Repo.GetByID(s.ctx, port.ID)
			s.Require().NoError(err)
			s.Require().Equal(port, *p)
		}
	})
}

// mockPorts created len Ports with mock data, plus some randomized data given as input
func mockPorts(len int, randData string) []Port {
	var ports = make([]Port, 0, len)
	for i := 0; i < len; i++ {
		ports = append(ports, Port{
			ID:      fmt.Sprintf("id-%d", i),
			Code:    fmt.Sprintf("code-%d-%s", i, randData),
			Name:    fmt.Sprintf("name-%d-%s", i, randData),
			City:    fmt.Sprintf("city-%d-%s", i, randData),
			Country: fmt.Sprintf("country-%d-%s", i, randData),
		})
	}
	return ports
}
