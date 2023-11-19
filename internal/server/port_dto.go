package server

import (
	"github.com/mihailozarinschi/ports-app/internal/port"
)

type PortDto struct {
	ID      string `json:"-"`
	Code    string `json:"code"`
	Name    string `json:"name"`
	City    string `json:"city"`
	Country string `json:"country"`
	// Add more fields when needed
}

func PortDtoToModel(dto PortDto) port.Port {
	return port.Port{
		ID:      dto.ID,
		Code:    dto.Code,
		Name:    dto.Name,
		City:    dto.City,
		Country: dto.Country,
	}
}
