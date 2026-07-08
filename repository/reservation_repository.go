package repository

import (
	"errors"
	"spotsync/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ReservationRepository interface {
	CreateWithTransaction(userID, zoneID uint, licensePlate string) (*models.Reservation, error)
	FindByUserID(userID uint) ([]models.Reservation, error)
	FindByID(id uint) (*models.Reservation, error)
	Cancel(id uint) error
	FindAll() ([]models.Reservation, error)
}

type reservationRepository struct {
	db *gorm.DB
}

func NewReservationRepository(db *gorm.DB) ReservationRepository {
	return &reservationRepository{db: db}
}

// ErrZoneFull is returned when a zone has reached its capacity.
var ErrZoneFull = errors.New("parking zone is full")

// CreateWithTransaction performs a safe, atomic reservation creation using
// row-level locking (SELECT ... FOR UPDATE) to prevent overbooking.
func (r *reservationRepository) CreateWithTransaction(userID, zoneID uint, licensePlate string) (*models.Reservation, error) {
	var created models.Reservation

	err := r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Lock the parking zone row to prevent concurrent overbooking
		var zone models.ParkingZone
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&zone, zoneID).Error; err != nil {
			return err
		}

		// 2. Count current active reservations
		var activeCount int64
		if err := tx.Model(&models.Reservation{}).
			Where("zone_id = ? AND status = ?", zoneID, "active").
			Count(&activeCount).Error; err != nil {
			return err
		}

		// 3. Enforce capacity
		if int(activeCount) >= zone.TotalCapacity {
			return ErrZoneFull
		}

		// 4. Create the reservation
		reservation := models.Reservation{
			UserID:       userID,
			ZoneID:       zoneID,
			LicensePlate: licensePlate,
			Status:       "active",
		}
		if err := tx.Create(&reservation).Error; err != nil {
			return err
		}

		created = reservation
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &created, nil
}

func (r *reservationRepository) FindByUserID(userID uint) ([]models.Reservation, error) {
	var reservations []models.Reservation
	if err := r.db.Preload("Zone").
		Where("user_id = ?", userID).
		Find(&reservations).Error; err != nil {
		return nil, err
	}
	return reservations, nil
}

func (r *reservationRepository) FindByID(id uint) (*models.Reservation, error) {
	var reservation models.Reservation
	if err := r.db.First(&reservation, id).Error; err != nil {
		return nil, err
	}
	return &reservation, nil
}

func (r *reservationRepository) Cancel(id uint) error {
	return r.db.Model(&models.Reservation{}).
		Where("id = ?", id).
		Update("status", "cancelled").Error
}

func (r *reservationRepository) FindAll() ([]models.Reservation, error) {
	var reservations []models.Reservation
	if err := r.db.Preload("User").Preload("Zone").
		Find(&reservations).Error; err != nil {
		return nil, err
	}
	return reservations, nil
}
