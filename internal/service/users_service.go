package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ansonallard/users-service/internal/auth"
	"github.com/ansonallard/users-service/internal/constants"
	"github.com/ansonallard/users-service/internal/errors"
	"github.com/ansonallard/users-service/internal/keys"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const (
	AUTHORIZATION_CODE_WINDOW_SECONDS = 30
)

type UsersService struct {
	t  *TenantsService
	db *mongo.Client
}

func NewUsersService(tenantsService *TenantsService, mongoClient *mongo.Client) UsersService {
	return UsersService{t: tenantsService, db: mongoClient}
}

func (u *UsersService) Create(ctx context.Context, incomingUser IncomingUserRequest) error {
	_, err := u.t.Get(ctx, incomingUser.TenantId)
	if err != nil {
		return errors.TenantNotFoundError{Id: incomingUser.TenantId}
	}

	user, err := u.Get(ctx, incomingUser.Username, incomingUser.TenantId)
	if _, ok := err.(errors.UserNotFoundError); !ok && err != nil {
		return err
	}
	if user != nil {
		return errors.UserExistsError{}
	}

	hashedPassword, err := auth.HashPassword(incomingUser.Password)
	if err != nil {
		return err
	}

	encryptionKey, err := keys.GenerateKey()
	if err != nil {
		return err
	}

	user = &UserModel{
		EncodedHashParams: hashedPassword,
		Username:          incomingUser.Username,
		TenantId:          incomingUser.TenantId,
		EncryptionKey:     base64.StdEncoding.EncodeToString(encryptionKey),
	}
	_, err = u.db.Database("users-service").Collection("users").InsertOne(ctx, user)
	if err != nil {
		return err
	}
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
			return nil, errors.UserNotFoundError{
				Usr:      username,
				TenantId: tenantId,
			}
		}
		return nil, fmt.Errorf("error decoding tenant: %v", err)
	}
	return user, nil
}

func (u *UsersService) Login(ctx context.Context, input LoginInput) (*LoginResult, error) {
	user, err := u.Get(ctx, input.Username, input.TenantId)
	if err != nil {
		return nil, err
	}

	match, err := auth.VerifyPassword(input.Password, user.EncodedHashParams)
	if err != nil {
		return nil, err
	}
	if !match {
		return nil, errors.NotAuthorizedError{}
	}

	// encryptionKey, err := base64.StdEncoding.DecodeString(user.EncryptionKey)
	encryptionKey, err := keys.ReadKeyFromFile(constants.AUTHORIZATION_ENCRYPTION_FILENAME)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	authorizationCode := AuthorizationCode{
		Username: user.Username,
		TenantId: user.TenantId,
		Iat:      now,
		Exp:      now.Add(time.Duration(AUTHORIZATION_CODE_WINDOW_SECONDS) * time.Second),
	}

	authorizationCodeBytes, err := json.Marshal(authorizationCode)
	if err != nil {
		return nil, err
	}

	encryptedCode, err := keys.Encrypt(authorizationCodeBytes, encryptionKey)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		Code: encryptedCode,
	}, nil
}

type LoginResult struct {
	Code string
}

type LoginInput struct {
	Username      string
	Password      string
	RedirectUri   string
	TenantId      string
	ApplicationId string
}

type IncomingUserRequest struct {
	Username string
	Password string
	TenantId string
}

type UserModel struct {
	*auth.EncodedHashParams
	Username      string `bson:"username"`
	TenantId      string `bson:"tenant_id"`
	EncryptionKey string `bson:"encryption_key"`
}

type AuthorizationCode struct {
	Username string    `json:"username""`
	TenantId string    `json:"tenant_id"`
	Iat      time.Time `json:"iat"`
	Exp      time.Time `json:"exp"`
}
