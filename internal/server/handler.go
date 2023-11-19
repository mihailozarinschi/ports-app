package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/mihailozarinschi/ports-app/internal/port"
)

type Handler struct {
	portsRepo port.Repository
}

func NewHandler(portsRepo port.Repository) *Handler {
	return &Handler{portsRepo: portsRepo}
}

func (h *Handler) ImportPorts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse ports
	var portsDto map[string]PortDto // map[ID]PortDto
	// TODO: improve decoding logic to not do it all at once, it is expected to received huge amount of data,
	//  try to ingest, parse, and store the data gradually as it comes in.
	err := json.NewDecoder(r.Body).Decode(&portsDto)
	if err != nil {
		http.Error(w, fmt.Sprintf("error decoding request body: %s", err), http.StatusBadRequest)
		return
	}

	// Create or update ports
	// TODO: This would have to change when decoding changes behavior,
	//  we should call CreateOrUpdate as soon as we parse a single/few ports from the json payload.
	for portID, portDto := range portsDto {
		portDto.ID = portID
		portModel := PortDtoToModel(portDto)

		_, err = h.portsRepo.CreateOrUpdate(ctx, portModel)
		if err != nil {
			log.Println(fmt.Errorf("erorr diring CreateOrUpdate of port with ID %q: %w", portID, err))
			switch {
			case errors.Is(err, ctx.Err()):
				http.Error(w, err.Error(), 499) // Client canceled request
			case errors.Is(err, port.ErrConnectionFailed):
				http.Error(w, err.Error(), http.StatusBadGateway)
			default:
				http.Error(w, err.Error(), http.StatusServiceUnavailable)
			}
			return
		}
	}

	log.Printf("imported %d ports", len(portsDto))
}
