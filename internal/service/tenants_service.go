package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ansonallard/users-service/internal/api"
	"github.com/oklog/ulid/v2"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type TenantsService struct {
	mongoClient *mongo.Client
}

func NewTenantService(mongoClient *mongo.Client) TenantsService {
	return TenantsService{mongoClient: mongoClient}
}

func (t *TenantsService) Create(ctx context.Context) (*api.CreateTenantResponse, error) {
	id := ulid.Make().String()
	now := time.Now()
	model := TenantModel{
		Id:        id,
		CreatedAt: now,
		Version:   "0",
	}
	result, err := t.mongoClient.Database("users-service").Collection("tenants").InsertOne(ctx, model)
	if err != nil {
		return nil, fmt.Errorf("error: %+v", err)
	}
	fmt.Println(result)
	return &api.CreateTenantResponse{
		CreatedAt: model.CreatedAt,
		TenantId:  model.Id,
		Version:   model.Version,
	}, nil
}

func (t *TenantsService) Get(ctx context.Context, id string) (TenantModel, error) {
	var tenant TenantModel
	result := t.mongoClient.Database("users-service").Collection("tenants").FindOne(ctx, struct{ id string }{
		id: id,
	})
	fmt.Printf("%+v", result)
	err := result.Decode(&tenant)
	if err != nil {
		fmt.Printf("error: %+v", err)
	}
	return tenant, nil
}

type TenantModel struct {
	Id        string    `bson:"id"`
	CreatedAt time.Time `bson:"createdAt"`
	Version   string    `bson:"version"`
}
