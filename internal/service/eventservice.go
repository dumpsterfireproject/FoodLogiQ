package service

import (
	"context"
	"net/http"

	"github.com/dumpsterfireproject/FoodLogiQ/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReturnCode struct {
	Status int
	Err    error
}

func (r ReturnCode) IsSuccess() bool {
	return r.Status == 200 || r.Status == 201
}

type EventService interface {
	CreateEvent(context.Context, *User, *model.Event) ReturnCode
	DeleteEvent(context.Context, *User, string) ReturnCode
	GetEvent(context.Context, *User, string) (*model.Event, ReturnCode)
	ListEvents(context.Context, *User) ([]*model.Event, ReturnCode)
}

type EventHandlerServiceImpl struct {
	mongo          *mongo.Client
	dbName         string
	collectionName string
}

func WithClient(client *mongo.Client) func(*EventHandlerServiceImpl) {
	return func(s *EventHandlerServiceImpl) {
		s.mongo = client
	}
}

func WithDbName(name string) func(*EventHandlerServiceImpl) {
	return func(s *EventHandlerServiceImpl) {
		s.dbName = name
	}
}

func WithCollectionName(name string) func(*EventHandlerServiceImpl) {
	return func(s *EventHandlerServiceImpl) {
		s.collectionName = name
	}
}

func NewEventHandlerService(options ...func(*EventHandlerServiceImpl)) *EventHandlerServiceImpl {
	service := &EventHandlerServiceImpl{}
	for _, o := range options {
		o(service)
	}
	return service
}

func (s *EventHandlerServiceImpl) collection() *mongo.Collection {
	return s.mongo.Database(s.dbName).Collection(s.collectionName)
}

func (s *EventHandlerServiceImpl) CreateEvent(ctx context.Context, user *User, event *model.Event) ReturnCode {
	event.CreatedBy = &user.UserID
	// TODO: Handle defaults for rest of the event
	_, err := s.collection().InsertOne(ctx, event)
	if err != nil {
		return ReturnCode{http.StatusInternalServerError, err}
	}
	return ReturnCode{http.StatusCreated, nil}
}

func (s *EventHandlerServiceImpl) DeleteEvent(ctx context.Context, user *User, id string) ReturnCode {
	_, err := s.collection().UpdateOne(ctx, bson.M{"_id": id, "createdBy": user.UserID}, bson.M{"isDeleted": true})
	if err != nil {
		return ReturnCode{http.StatusInternalServerError, err}
	}
	return ReturnCode{200, nil}
}

func (s *EventHandlerServiceImpl) GetEvent(ctx context.Context, user *User, id string) (*model.Event, ReturnCode) {
	var event *model.Event
	err := s.collection().FindOne(ctx, bson.M{"_id": id, "createdBy": user.UserID}).Decode(event)
	if err != nil {
		return nil, ReturnCode{http.StatusInternalServerError, err}
	}
	return event, ReturnCode{200, nil}
}

func (s *EventHandlerServiceImpl) ListEvents(ctx context.Context, user *User) ([]*model.Event, ReturnCode) {
	var events []*model.Event
	cursor, err := s.collection().Find(ctx, bson.M{"createdBy": user.UserID})
	if err != nil {
		return []*model.Event{}, ReturnCode{http.StatusInternalServerError, err}
	}

	if err = cursor.All(ctx, &events); err != nil {
		return []*model.Event{}, ReturnCode{http.StatusInternalServerError, err}
	}

	cursor.Close(ctx)

	return events, ReturnCode{200, nil}
}
