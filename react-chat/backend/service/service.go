package service

import (
	"sync"

	"github.com/theoremoon/kaidsuka/react-chat/backend/model"
	"github.com/theoremoon/kaidsuka/react-chat/backend/repository"
)

type Service interface {
	GetUser(id uint32) (*model.User, error)
	LoginUser(username string) (*model.User, error)
	RegisterUser(username string) (*model.User, error)

	ListMessages(olderThan uint32, limit uint32) ([]*model.Message, error)
	PostMessage(userID uint32, text string) (*model.Message, error)
	UpdateMessage(userID, messageID uint32, text string) (*model.Message, error)

	AddPostSubscriber(key string) <-chan *model.Message
	RemovePostSubscriber(key string)
	AddUpdateSubscriber(key string) <-chan *model.Message
	RemoveUpdateSubscriber(key string)

	Close() error
}

type service struct {
	repo              repository.Repository
	postSubscribers   map[string]chan<- *model.Message
	updateSubscribers map[string]chan<- *model.Message
	sync.Mutex
}

func New(repo repository.Repository) (Service, error) {
	return &service{
		repo:              repo,
		postSubscribers:   make(map[string]chan<- *model.Message),
		updateSubscribers: make(map[string]chan<- *model.Message),
	}, nil
}

func (s *service) RegisterUser(username string) (*model.User, error) {
	return s.repo.RegisterUser(username)
}

func (s *service) LoginUser(username string) (*model.User, error) {
	return s.repo.GetUser(username)
}

func (s *service) GetUser(id uint32) (*model.User, error) {
	return s.repo.GetUserByID(id)
}

func (s *service) ListMessages(olderThan uint32, limit uint32) ([]*model.Message, error) {
	return s.repo.ListMessages(olderThan, limit)
}

func (s *service) PostMessage(userID uint32, text string) (*model.Message, error) {
	msg, err := s.repo.PostMessage(userID, text)
	if err == nil {
		s.Mutex.Lock()
		for _, ch := range s.postSubscribers {
			ch <- msg
		}
		s.Mutex.Unlock()
	}
	return msg, err
}

func (s *service) UpdateMessage(userID uint32, messageID uint32, text string) (*model.Message, error) {
	msg, err := s.repo.UpdateMessage(userID, messageID, text)
	if err == nil {
		s.Mutex.Lock()
		for _, ch := range s.postSubscribers {
			ch <- msg
		}
		s.Mutex.Unlock()
	}
	return msg, err
}

func (s *service) AddPostSubscriber(key string) <-chan *model.Message {
	ch := make(chan *model.Message)
	s.Mutex.Lock()
	s.postSubscribers[key] = ch
	s.Mutex.Unlock()
	return ch
}
func (s *service) RemovePostSubscriber(key string) {
	s.Mutex.Lock()
	close(s.postSubscribers[key])
	delete(s.postSubscribers, key)
	s.Mutex.Unlock()
}

func (s *service) AddUpdateSubscriber(key string) <-chan *model.Message {
	ch := make(chan *model.Message)
	s.Mutex.Lock()
	s.updateSubscribers[key] = ch
	s.Mutex.Unlock()
	return ch
}
func (s *service) RemoveUpdateSubscriber(key string) {
	s.Mutex.Lock()
	close(s.updateSubscribers[key])
	delete(s.updateSubscribers, key)
	s.Mutex.Unlock()
}

func (s *service) Close() error {
	return s.repo.Close()
}
