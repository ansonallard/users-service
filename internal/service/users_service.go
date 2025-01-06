package service

import (
	"context"
	"fmt"

	"github.com/ansonallard/users-service/internal/auth"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UsersService struct {
	t  *TenantsService
	db *mongo.Client
}

func NewUsersService(tenantsService *TenantsService, mongoClient *mongo.Client) UsersService {
	return UsersService{t: tenantsService, db: mongoClient}
}

func (u *UsersService) Create(ctx context.Context, incomingUser IncomingUserRequest) error {

	tenant, err := u.t.Get(ctx, incomingUser.TenantId)
	if err != nil {
		return fmt.Errorf("error: Tenant '%s' not found", incomingUser.TenantId)
	}
	fmt.Println(tenant)

	user, err := u.Get(ctx, incomingUser.Username, incomingUser.TenantId)
	if user != nil {
		return fmt.Errorf("error: user already exists")
	}

	hashedPassword, err := auth.HashPassword(incomingUser.Password)
	return nil

}

func (u *UsersService) Get(ctx context.Context, username, tenantId string) (user *UserModel, err error) {
	response := u.db.Database("users-service").Collection("users").FindOne(ctx, bson.M{
		"username":  username,
		"tenant_id": tenantId,
	})
	err = response.Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user '%s' not found in tenant '%s'", username)
		}
		return nil, fmt.Errorf("error decoding tenant: %v", err)
	}
}

type IncomingUserRequest struct {
	Username string
	Password string
	TenantId string
}

type UserModel struct {
	Username       string `bson:"username"`
	HashedPassword string `bson:"hashed_password"`
	TenantId       string `bson:"tenant_id"`
}
