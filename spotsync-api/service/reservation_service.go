package service

import (
	"errors"

	"github.com/zabirarkam27/level2-assignment06-spotsync/dto"
	"github.com/zabirarkam27/level2-assignment06-spotsync/models"
	"github.com/zabirarkam27/level2-assignment06-spotsync/repository"
	"gorm.io/gorm"
)

type ReservationService interface {
	Create(userID uint, req dto.CreateReservationRequest) (*dto.ReservationResponse, error)
	GetMine(userID uint) ([]dto.ReservationResponse, error)
	GetAll() ([]dto.ReservationResponse, error)
	Cancel(userID, reservationID uint, role string) error
}

type reservationService struct {
	reservationRepo repository.ReservationRepository
	zoneRepo        repository.ZoneRepository
}

func NewReservationService(reservationRepo repository.ReservationRepository, zoneRepo repository.ZoneRepository) ReservationService {
	return &reservationService{reservationRepo: reservationRepo, zoneRepo: zoneRepo}
}

func (s *reservationService) Create(userID uint, req dto.CreateReservationRequest) (*dto.ReservationResponse, error) {
	reservation, err := s.reservationRepo.CreateWithLock(userID, req.ZoneID, req.LicensePlate)
	if err != nil {
		if errors.Is(err, repository.ErrZoneFull) {
			return nil, ErrZoneFull
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return mapReservation(*reservation, false, false), nil
}

func (s *reservationService) GetMine(userID uint) ([]dto.ReservationResponse, error) {
	reservations, err := s.reservationRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	return mapReservations(reservations, false, true), nil
}

func (s *reservationService) GetAll() ([]dto.ReservationResponse, error) {
	reservations, err := s.reservationRepo.FindAll()
	if err != nil {
		return nil, err
	}
	return mapReservations(reservations, true, true), nil
}

func (s *reservationService) Cancel(userID, reservationID uint, role string) error {
	reservation, err := s.reservationRepo.FindByID(reservationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}

	if role != "admin" && reservation.UserID != userID {
		return ErrForbidden
	}

	return s.reservationRepo.Cancel(reservationID)
}

func mapReservations(reservations []models.Reservation, includeUser, includeZone bool) []dto.ReservationResponse {
	responses := make([]dto.ReservationResponse, 0, len(reservations))
	for _, reservation := range reservations {
		responses = append(responses, *mapReservation(reservation, includeUser, includeZone))
	}
	return responses
}

func mapReservation(reservation models.Reservation, includeUser, includeZone bool) *dto.ReservationResponse {
	response := &dto.ReservationResponse{
		ID:           reservation.ID,
		UserID:       reservation.UserID,
		ZoneID:       reservation.ZoneID,
		LicensePlate: reservation.LicensePlate,
		Status:       reservation.Status,
		CreatedAt:    reservation.CreatedAt,
		UpdatedAt:    reservation.UpdatedAt,
	}

	if includeUser && reservation.User.ID != 0 {
		response.User = &dto.ReservationUserInfo{
			ID:    reservation.User.ID,
			Name:  reservation.User.Name,
			Email: reservation.User.Email,
			Role:  reservation.User.Role,
		}
	}

	if includeZone && reservation.Zone.ID != 0 {
		response.Zone = &dto.ZoneInfo{
			ID:   reservation.Zone.ID,
			Name: reservation.Zone.Name,
			Type: reservation.Zone.Type,
		}
	}

	return response
}
