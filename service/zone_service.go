package service

import (
	"spotsync/dto"
	"spotsync/models"
	"spotsync/repository"
)

type ZoneService interface {
	CreateZone(req dto.CreateZoneRequest) (*dto.ZoneResponse, error)
	GetAllZones() ([]dto.ZoneResponse, error)
	GetZoneByID(id uint) (*dto.ZoneResponse, error)
	UpdateZone(id uint, req dto.UpdateZoneRequest) (*dto.ZoneResponse, error)
	DeleteZone(id uint) error
}

type zoneService struct {
	zoneRepo repository.ZoneRepository
}

func NewZoneService(zoneRepo repository.ZoneRepository) ZoneService {
	return &zoneService{zoneRepo: zoneRepo}
}

func (s *zoneService) CreateZone(req dto.CreateZoneRequest) (*dto.ZoneResponse, error) {
	zone := models.ParkingZone{
		Name:          req.Name,
		Type:          req.Type,
		TotalCapacity: req.TotalCapacity,
		PricePerHour:  req.PricePerHour,
	}
	if err := s.zoneRepo.Create(&zone); err != nil {
		return nil, err
	}
	return s.toZoneResponse(zone, 0), nil
}

func (s *zoneService) GetAllZones() ([]dto.ZoneResponse, error) {
	zones, err := s.zoneRepo.FindAll()
	if err != nil {
		return nil, err
	}

	responses := make([]dto.ZoneResponse, 0, len(zones))
	for _, z := range zones {
		active, err := s.zoneRepo.CountActiveReservations(z.ID)
		if err != nil {
			return nil, err
		}
		available := z.TotalCapacity - int(active)
		if available < 0 {
			available = 0
		}
		responses = append(responses, *s.toZoneResponse(z, available))
	}
	return responses, nil
}

func (s *zoneService) GetZoneByID(id uint) (*dto.ZoneResponse, error) {
	zone, err := s.zoneRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	active, err := s.zoneRepo.CountActiveReservations(zone.ID)
	if err != nil {
		return nil, err
	}
	available := zone.TotalCapacity - int(active)
	if available < 0 {
		available = 0
	}
	return s.toZoneResponse(*zone, available), nil
}

func (s *zoneService) UpdateZone(id uint, req dto.UpdateZoneRequest) (*dto.ZoneResponse, error) {
	zone, err := s.zoneRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		zone.Name = req.Name
	}
	if req.Type != "" {
		zone.Type = req.Type
	}
	if req.TotalCapacity > 0 {
		zone.TotalCapacity = req.TotalCapacity
	}
	if req.PricePerHour > 0 {
		zone.PricePerHour = req.PricePerHour
	}
	if err := s.zoneRepo.Update(zone); err != nil {
		return nil, err
	}
	active, _ := s.zoneRepo.CountActiveReservations(zone.ID)
	available := zone.TotalCapacity - int(active)
	if available < 0 {
		available = 0
	}
	return s.toZoneResponse(*zone, available), nil
}

func (s *zoneService) DeleteZone(id uint) error {
	_, err := s.zoneRepo.FindByID(id)
	if err != nil {
		return err
	}
	return s.zoneRepo.Delete(id)
}

func (s *zoneService) toZoneResponse(z models.ParkingZone, availableSpots int) *dto.ZoneResponse {
	return &dto.ZoneResponse{
		ID:             z.ID,
		Name:           z.Name,
		Type:           z.Type,
		TotalCapacity:  z.TotalCapacity,
		AvailableSpots: availableSpots,
		PricePerHour:   z.PricePerHour,
		CreatedAt:      z.CreatedAt,
		UpdatedAt:      z.UpdatedAt,
	}
}
