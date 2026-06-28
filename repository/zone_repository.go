package repository

import (
	"errors"

	"github.com/zabirarkam27/level2-assignment06-spotsync/models"
	"gorm.io/gorm"
)

type ZoneRepository interface {
	Create(zone *models.ParkingZone) error
	FindAll() ([]models.ParkingZone, error)
	FindByID(id uint) (*models.ParkingZone, error)
	Update(zone *models.ParkingZone) error
	Delete(id uint) error
	CountActiveReservations(zoneID uint) (int64, error)
}

type zoneRepository struct {
	db *gorm.DB
}

func NewZoneRepository(db *gorm.DB) ZoneRepository {
	return &zoneRepository{db: db}
}

func (r *zoneRepository) Create(zone *models.ParkingZone) error {
	return r.db.Create(zone).Error
}

func (r *zoneRepository) FindAll() ([]models.ParkingZone, error) {
	var zones []models.ParkingZone
	err := r.db.Order("id ASC").Find(&zones).Error
	return zones, err
}

func (r *zoneRepository) FindByID(id uint) (*models.ParkingZone, error) {
	var zone models.ParkingZone
	err := r.db.First(&zone, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return &zone, err
}

func (r *zoneRepository) Update(zone *models.ParkingZone) error {
	return r.db.Save(zone).Error
}

func (r *zoneRepository) Delete(id uint) error {
	result := r.db.Delete(&models.ParkingZone{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *zoneRepository) CountActiveReservations(zoneID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Reservation{}).
		Where("zone_id = ? AND status = ?", zoneID, "active").
		Count(&count).Error
	return count, err
}
