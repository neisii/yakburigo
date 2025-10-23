package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/neisii/yakburigo/internal/dto"
	"github.com/neisii/yakburigo/internal/service"
	"github.com/neisii/yakburigo/pkg/response"
)

type LocationHandler struct {
	service service.LocationService
}

func NewLocationHandler(service service.LocationService) *LocationHandler {
	return &LocationHandler{service: service}
}

func (h *LocationHandler) GetNearbyLocations(c *gin.Context) {
	var req dto.NearbyLocationsRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request parameters")
		return
	}

	if req.Radius == 0 {
		req.Radius = 3000
	}
	if req.Type == "" {
		req.Type = "all"
	}

	locations, err := h.service.GetNearbyLocations(req.Latitude, req.Longitude, req.Radius, req.Type)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	results := make([]dto.LocationResponse, len(locations))
	for i, loc := range locations {
		results[i] = dto.LocationResponse{
			ID:             loc.ID,
			Name:           loc.Name,
			Type:           loc.Type,
			Address:        loc.Address,
			Latitude:       loc.Latitude,
			Longitude:      loc.Longitude,
			Distance:       loc.Distance,
			OperatingHours: loc.OperatingHours,
			Contact:        loc.Contact,
		}
	}

	response.SuccessWithCount(c, http.StatusOK, results, len(results))
}

func (h *LocationHandler) GetLocationByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid location ID")
		return
	}

	location, err := h.service.GetLocationByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "Location not found")
		return
	}

	result := dto.LocationResponse{
		ID:             location.ID,
		Name:           location.Name,
		Type:           location.Type,
		Address:        location.Address,
		Latitude:       location.Latitude,
		Longitude:      location.Longitude,
		Distance:       0,
		OperatingHours: location.OperatingHours,
		Contact:        location.Contact,
	}

	response.Success(c, http.StatusOK, result)
}
