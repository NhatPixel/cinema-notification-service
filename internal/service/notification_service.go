package service

import (
	"fmt"
	"sync"

	"github.com/google/uuid"

	"github.com/NhatPixel/cinema-notification-service/internal/model"
	"github.com/NhatPixel/cinema-notification-service/internal/repository"
	"github.com/NhatPixel/cinema-notification-service/internal/dto"
)

type NotificationService struct {
	repo *repository.NotificationRepository
	clients map[string][]chan model.Notification
	mu sync.RWMutex
}

func NewNotificationService(repo *repository.NotificationRepository) *NotificationService {
	return &NotificationService{
		repo:    repo,
		clients: make(map[string][]chan model.Notification),
	}
}

func (s *NotificationService) Subscribe(userID string) chan model.Notification {
	s.mu.Lock()
	defer s.mu.Unlock()

	ch := make(chan model.Notification, 10)
	s.clients[userID] = append(s.clients[userID], ch)

	return ch
}

func (s *NotificationService) Unsubscribe(userID string, ch chan model.Notification) {
	s.mu.Lock()
	defer s.mu.Unlock()

	chs := s.clients[userID]
	for i, c := range chs {
		if c == ch {
			s.clients[userID] = append(chs[:i], chs[i+1:]...)
			break
		}
	}

	if len(chs) == 0 {
		delete(s.clients, userID)
	} else {
		s.clients[userID] = chs
	}
}

func (s *NotificationService) Create(req dto.CreateRequest) error {
	n := req.ToModel()
	n.IsRead = false
	n.ID = uuid.Must(uuid.NewV7()).String()
	if err := s.repo.Create(&n); err != nil {
		return err
	}

	s.broadcast(n)

	return nil
}

func (s *NotificationService) CreateForUsers(reqs []dto.CreateRequest) error {
	var errs []error

	for _, req := range reqs {
		n := req.ToModel()
		n.IsRead = false
		n.ID = uuid.NewString()
		if err := s.repo.Create(&n); err != nil {
			errs = append(errs, err)
			continue
		}
		s.broadcast(n)
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to create %d notifications", len(errs))
	}
	return nil
}

func (s *NotificationService) broadcast(n model.Notification) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, ch := range s.clients[n.UserID] {
		select {
		case ch <- n:
		default:
		}
	}
}


func (s *NotificationService) FindByUserID(userID string) ([]model.Notification, error) {
	return s.repo.FindByUserID(userID)
}

func (s *NotificationService) UpdateReadStatus(id string) error {
	return s.repo.UpdateReadStatus(id)
}
