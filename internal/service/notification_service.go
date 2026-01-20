package service

import (
	"sync"

	"github.com/NhatPixel/cinema-notification-service/internal/model"
	"github.com/NhatPixel/cinema-notification-service/internal/repository"
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

func (s *NotificationService) Notify(n model.Notification) error {
	if err := s.repo.Create(&n); err != nil {
		return err
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, ch := range s.clients[n.UserID] {
		select {
		case ch <- n:
		default:
		}
	}

	return nil
}

