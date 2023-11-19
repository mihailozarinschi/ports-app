package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mihailozarinschi/ports-app/internal/adapters/memory"
	"github.com/mihailozarinschi/ports-app/internal/port"
)

func TestHandler_ImportPorts(t *testing.T) {
	ctx := context.Background()

	portsRepo := memory.NewPortRepository()
	h := NewHandler(portsRepo)
	mux := NewServeMux(h)

	server := httptest.NewServer(mux)
	defer server.Close()

	// These ports' data matches what's in ./testdata/ports-data-valid.json
	mockPorts := map[string]port.Port{
		"ZAZBA": {
			ID:      "ZAZBA",
			Code:    "", // has no name
			Name:    "Coega",
			City:    "Coega",
			Country: "South Africa",
		},
		"BDCGP": {
			ID:      "BDCGP",
			Code:    "53827",
			Name:    "Chittagong",
			City:    "Chittagong",
			Country: "Bangladesh",
		},
		"AEAJM": {
			ID:      "AEAJM",
			Code:    "52000",
			Name:    "Ajman",
			City:    "Ajman",
			Country: "United Arab Emirates",
		},
	}

	t.Run("import ports successfully", func(t *testing.T) {

		f, err := os.Open("./testdata/ports-data-valid.json")
		require.NoError(t, err)
		defer f.Close()

		res, err := http.Post(server.URL+"/ports/import", "application/json", f)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)

		// Check data is in repository
		for portID, mockPort := range mockPorts {
			p, err := portsRepo.GetByID(ctx, portID)
			require.NoError(t, err)
			require.Equal(t, *p, mockPort)
		}
	})

	t.Run("import some ports again with some updates", func(t *testing.T) {

		f, err := os.Open("./testdata/ports-data-valid-updated.json")
		require.NoError(t, err)
		defer f.Close()
		// Port AEAJM has country updated, the import should update this information
		AEAJM := mockPorts["AEAJM"]
		AEAJM.Country = "UAE"
		mockPorts["AEAJM"] = AEAJM

		res, err := http.Post(server.URL+"/ports/import", "application/json", f)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)

		// Check new data is in repository
		for portID, mockPort := range mockPorts {
			p, err := portsRepo.GetByID(ctx, portID)
			require.NoError(t, err)
			require.Equal(t, *p, mockPort)
		}
	})

	t.Run("try import invalid data", func(t *testing.T) {
		f, err := os.Open("./testdata/ports-data-invalid.json")
		require.NoError(t, err)
		defer f.Close()

		res, err := http.Post(server.URL+"/ports/import", "application/json", f)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	//TODO: implement test for big import
}
