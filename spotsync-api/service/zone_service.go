package service

import (
	"errors"

	"github.com/zabirarkam27/level2-assignment06-spotsync/dto"
	"github.com/zabirarkam27/level2-assignment06-spotsync/models"
	"github.com/zabirarkam27/level2-assignment06-spotsync/repository"
	"gorm.io/gorm"
)

type ZoneService interface {
	Create(req dto.CreateZoneRequest) (*dto.ZoneResponse, error)
	GetAll() ([]dto.ZoneResponse, error)
	GetByID(id uint) (*dto.ZoneResponse, error)
	Update(id uint, req dto.UpdateZoneRequest) (*dto.ZoneResponse, error)
	Delete(id uint) error
}

type zoneService struct {
	zoneRepo repository.ZoneRepository
}

func NewZoneService(zoneRepo repository.ZoneRepository) ZoneService {
	return &zoneService{zoneRepo: zoneRepo}
}

func (s *zoneService) Create(req dto.CreateZoneRequest) (*dto.ZoneResponse, error) {
	zone := models.ParkingZone{
		Name:          req.Name,
		Type:          req.Type,
		TotalCapacity: req.TotalCapacity,
		PricePerHour:  req.PricePerHour,
	}
	if err := s.zoneRepo.Create(&zone); err != nil {
		return nil, err
	}
	return s.zoneResponse(zone)
}

func (s *zoneService) GetAll() ([]dto.ZoneResponse, error) {
	zones, err := s.zoneRepo.FindAll()
	if err != nil {
		return nil, err
	}

	responses := make([]dto.ZoneResponse, 0, len(zones))
	for _, zone := range zones {
		response, err := s.zoneResponse(zone)
		if err != nil {
			return nil, err
		}
		responses = append(responses, *response)
	}
	return responses, nil
}

func (s *zoneService) GetByID(id uint) (*dto.ZoneResponse, error) {
	zone, err := s.zoneRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return s.zoneResponse(*zone)
}

func (s *zoneService) Update(id uint, req dto.UpdateZoneRequest) (*dto.ZoneResponse, error) {
	zone, err := s.zoneRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if req.Name != "" {
		zone.Name = req.Name
	}
	if req.Type != "" {
		zone.Type = req.Type
	}
	if req.TotalCapacity != 0 {
		zone.TotalCapacity = req.TotalCapacity
	}
	if req.PricePerHour != 0 {
		zone.PricePerHour = req.PricePerHour
	}

	if err := s.zoneRepo.Update(zone); err != nil {
		return nil, err
	}
	return s.zoneResponse(*zone)
}

func (s *zoneService) Delete(id uint) error {
	if err := s.zoneRepo.Delete(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (s *zoneService) zoneResponse(zone models.ParkingZone) (*dto.ZoneResponse, error) {
	activeCount, err := s.zoneRepo.CountActiveReservations(zone.ID)
	if err != nil {
		return nil, err
	}
	available := zone.TotalCapacity - int(activeCount)
	if available < 0 {
		available = 0
	}

	return &dto.ZoneResponse{
		ID:             zone.ID,
		Name:           zone.Name,
		Type:           zone.Type,
		TotalCapacity:  zone.TotalCapacity,
		AvailableSpots: available,
		PricePerHour:   zone.PricePerHour,
		CreatedAt:      zone.CreatedAt,
		UpdatedAt:      zone.UpdatedAt,
	}, nil
}
