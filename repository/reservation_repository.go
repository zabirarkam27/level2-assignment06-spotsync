package repository

import (
	"errors"

	"github.com/zabirarkam27/level2-assignment06-spotsync/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrZoneFull = errors.New("parking zone is full")

type ReservationRepository interface {
	CreateWithLock(userID, zoneID uint, licensePlate string) (*models.Reservation, error)
	FindByUserID(userID uint) ([]models.Reservation, error)
	FindAll() ([]models.Reservation, error)
	FindByID(id uint) (*models.Reservation, error)
	Cancel(id uint) error
}

type reservationRepository struct {
	db *gorm.DB
}

func NewReservationRepository(db *gorm.DB) ReservationRepository {
	return &reservationRepository{db: db}
}

func (r *reservationRepository) CreateWithLock(userID, zoneID uint, licensePlate string) (*models.Reservation, error) {
	var reservation models.Reservation

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var zone models.ParkingZone
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&zone, zoneID).Error; err != nil {
			return err
		}

		var activeCount int64
		if err := tx.Model(&models.Reservation{}).
			Where("zone_id = ? AND status = ?", zoneID, "active").
			Count(&activeCount).Error; err != nil {
			return err
		}

		if activeCount >= int64(zone.TotalCapacity) {
			return ErrZoneFull
		}

		reservation = models.Reservation{
			UserID:       userID,
			ZoneID:       zoneID,
			LicensePlate: licensePlate,
			Status:       "active",
		}

		return tx.Create(&reservation).Error
	})
	if err != nil {
		return nil, err
	}

	return &reservation, nil
}

func (r *reservationRepository) FindByUserID(userID uint) ([]models.Reservation, error) {
	var reservations []models.Reservation
	err := r.db.Preload("Zone").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&reservations).Error
	return reservations, err
}

func (r *reservationRepository) FindAll() ([]models.Reservation, error) {
	var reservations []models.Reservation
	err := r.db.Preload("User").
		Preload("Zone").
		Order("created_at DESC").
		Find(&reservations).Error
	return reservations, err
}

func (r *reservationRepository) FindByID(id uint) (*models.Reservation, error) {
	var reservation models.Reservation
	err := r.db.Preload("User").Preload("Zone").First(&reservation, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return &reservation, err
}

func (r *reservationRepository) Cancel(id uint) error {
	result := r.db.Model(&models.Reservation{}).
		Where("id = ?", id).
		Update("status", "cancelled")
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
