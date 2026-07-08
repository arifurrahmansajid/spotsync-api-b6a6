package service

import (
	"errors"
	"spotsync/dto"
	"spotsync/repository"

	"gorm.io/gorm"
)

// ErrZoneFull indicates a parking zone has reached capacity.
var ErrZoneFull = errors.New("parking zone is full")

type ReservationService interface {
	CreateReservation(userID uint, req dto.CreateReservationRequest) (*dto.ReservationResponse, error)
	GetMyReservations(userID uint) ([]dto.MyReservationResponse, error)
	CancelReservation(reservationID, userID uint, role string) error
	GetAllReservations() ([]dto.AdminReservationResponse, error)
}

type reservationService struct {
	reservationRepo repository.ReservationRepository
	zoneRepo        repository.ZoneRepository
}

func NewReservationService(reservationRepo repository.ReservationRepository, zoneRepo repository.ZoneRepository) ReservationService {
	return &reservationService{
		reservationRepo: reservationRepo,
		zoneRepo:        zoneRepo,
	}
}

func (s *reservationService) CreateReservation(userID uint, req dto.CreateReservationRequest) (*dto.ReservationResponse, error) {
	// Verify zone exists before entering transaction
	_, err := s.zoneRepo.FindByID(req.ZoneID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	reservation, err := s.reservationRepo.CreateWithTransaction(userID, req.ZoneID, req.LicensePlate)
	if err != nil {
		if errors.Is(err, repository.ErrZoneFull) {
			return nil, ErrZoneFull
		}
		return nil, err
	}

	return &dto.ReservationResponse{
		ID:           reservation.ID,
		UserID:       reservation.UserID,
		ZoneID:       reservation.ZoneID,
		LicensePlate: reservation.LicensePlate,
		Status:       reservation.Status,
		CreatedAt:    reservation.CreatedAt,
		UpdatedAt:    reservation.UpdatedAt,
	}, nil
}

func (s *reservationService) GetMyReservations(userID uint) ([]dto.MyReservationResponse, error) {
	reservations, err := s.reservationRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.MyReservationResponse, 0, len(reservations))
	for _, r := range reservations {
		responses = append(responses, dto.MyReservationResponse{
			ID:           r.ID,
			LicensePlate: r.LicensePlate,
			Status:       r.Status,
			Zone: dto.ZoneInfo{
				ID:   r.Zone.ID,
				Name: r.Zone.Name,
				Type: r.Zone.Type,
			},
			CreatedAt: r.CreatedAt,
		})
	}
	return responses, nil
}

func (s *reservationService) CancelReservation(reservationID, userID uint, role string) error {
	reservation, err := s.reservationRepo.FindByID(reservationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gorm.ErrRecordNotFound
		}
		return err
	}

	// Drivers can only cancel their own reservations
	if role != "admin" && reservation.UserID != userID {
		return errors.New("forbidden")
	}

	if reservation.Status == "cancelled" {
		return errors.New("reservation already cancelled")
	}

	return s.reservationRepo.Cancel(reservationID)
}

func (s *reservationService) GetAllReservations() ([]dto.AdminReservationResponse, error) {
	reservations, err := s.reservationRepo.FindAll()
	if err != nil {
		return nil, err
	}

	responses := make([]dto.AdminReservationResponse, 0, len(reservations))
	for _, r := range reservations {
		responses = append(responses, dto.AdminReservationResponse{
			ID:           r.ID,
			LicensePlate: r.LicensePlate,
			Status:       r.Status,
			User: dto.LoginUserResponse{
				ID:    r.User.ID,
				Name:  r.User.Name,
				Email: r.User.Email,
				Role:  r.User.Role,
			},
			Zone: dto.ZoneInfo{
				ID:   r.Zone.ID,
				Name: r.Zone.Name,
				Type: r.Zone.Type,
			},
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
		})
	}
	return responses, nil
}
