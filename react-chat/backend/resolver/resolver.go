package resolver

import (
	"context"
	"errors"
	"math"
	"strconv"

	"github.com/google/uuid"
	"github.com/theoremoon/kaidsuka/react-chat/backend/model"
	"github.com/theoremoon/kaidsuka/react-chat/backend/service"
)

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.
//go:generate gqlgen

const userKey = "github.com/theoremoon/kaidsuka/react-chat/server/service/resolver::user"

func AttachUser(ctx context.Context, user *model.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

type resolver struct {
	service service.Service
}

func New(s service.Service) *resolver {
	return &resolver{
		service: s,
	}
}

func (r *resolver) Message() MessageResolver {
	return &messageResolver{r}
}
func (r *resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *resolver) Query() QueryResolver {
	return &queryResolver{r}
}
func (r *resolver) Subscription() SubscriptionResolver {
	return &subscriptionResolver{r}
}

type messageResolver struct{ *resolver }

func (r *messageResolver) ID(ctx context.Context, obj *model.Message) (string, error) {
	return strconv.FormatUint(uint64(obj.ID), 10), nil
}
func (r *messageResolver) PostedAt(ctx context.Context, obj *model.Message) (int, error) {
	return int(obj.PostedAt), nil
}
func (r *messageResolver) User(ctx context.Context, obj *model.Message) (*model.User, error) {
	return r.service.GetUser(obj.UserID)
}

type mutationResolver struct{ *resolver }

func (r *mutationResolver) PostMessage(ctx context.Context, text string) (*model.Message, error) {
	user := ctx.Value(userKey).(*model.User)
	if user == nil {
		return nil, errors.New("Unauthorized")
	}
	return r.service.PostMessage(user.ID, text)
}
func (r *mutationResolver) EditMessage(ctx context.Context, messageID string, text string) (*model.Message, error) {
	user := ctx.Value(userKey).(*model.User)
	if user == nil {
		return nil, errors.New("Unauthorized")
	}
	id, err := strconv.Atoi(messageID)
	if err != nil {
		return nil, err
	}
	return r.service.UpdateMessage(user.ID, uint32(id), text)
}

type queryResolver struct{ *resolver }

func (r *queryResolver) User(ctx context.Context) (*model.User, error) {
	user := ctx.Value(userKey).(*model.User)
	if user == nil {
		return nil, errors.New("Unauthorized")
	}
	return user, nil
}
func (r *queryResolver) GetMessages(ctx context.Context, before *int, count *int) ([]*model.Message, error) {
	if before == nil {
		before = new(int)
		*before = math.MaxUint32
	}
	if count == nil {
		count = new(int)
		*count = 30
	}
	return r.service.ListMessages(uint32(*before), uint32(*count))
}

type subscriptionResolver struct{ *resolver }

func (r *subscriptionResolver) NewMessage(ctx context.Context) (<-chan *model.Message, error) {
	key := uuid.New().String()
	ch := r.service.AddPostSubscriber(key)
	go func() {
		<-ctx.Done()
		r.service.RemovePostSubscriber(key)
	}()
	return ch, nil
}
func (r *subscriptionResolver) UpdatedMessage(ctx context.Context) (<-chan *model.Message, error) {
	key := uuid.New().String()
	ch := r.service.AddUpdateSubscriber(key)
	go func() {
		<-ctx.Done()
		r.service.RemoveUpdateSubscriber(key)
	}()
	return ch, nil
}
