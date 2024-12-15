package usecase

import (
	"context"
	"simple-crud-clean-architecture/internal/entity"
	"simple-crud-clean-architecture/internal/helper"
	"simple-crud-clean-architecture/internal/model"
	"simple-crud-clean-architecture/internal/model/converter"
	"simple-crud-clean-architecture/internal/repository"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/redis/go-redis/v9"
)

type UserUseCase struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	Validate       *validator.Validate
	UserRepository *repository.UserRepository
	RedisClient    *redis.Client
	TokenHelper    *helper.TokenHelper
	EmailHelper    *helper.GomailSender
}

func NewUserUseCase(db *gorm.DB, logger *logrus.Logger, validate *validator.Validate,
	userRepository *repository.UserRepository, redisClient *redis.Client, tokenHelper *helper.TokenHelper, emailHelper *helper.GomailSender) *UserUseCase {
	return &UserUseCase{
		DB:             db,
		Log:            logger,
		Validate:       validate,
		RedisClient:    redisClient,
		UserRepository: userRepository,
		TokenHelper:    tokenHelper,
		EmailHelper:    emailHelper,
	}
}

func (c *UserUseCase) Create(ctx context.Context, request *model.RegisterUserRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// err := c.Validate.Struct(request)
	// if err != nil {
	// 	c.Log.Warnf("Invalid request body : %+v", err)
	// 	return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	// }

	helper.SanitiseStruct(request)

	total, err := c.UserRepository.CountByEmail(tx, request.Email)
	if err != nil {
		c.Log.Warnf("Failed count user from database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if total > 0 {
		c.Log.Warnf("User already exists : %+v", err)
		return nil, fiber.NewError(fiber.StatusConflict, "user already exists")
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Log.Warnf("Failed to generate bcrypt hash : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	user := &entity.User{
		UUID:       uuid.New(),
		Email:      request.Email,
		Password:   string(password),
		VerifiedAt: nil,
	}

	if err := c.UserRepository.Create(tx, user); err != nil {
		c.Log.Warnf("Failed create user to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	c.Log.Infof("User with email: %+v successfully created their account", request.Email)

	return converter.UserToResponse(user), nil
}

func (c *UserUseCase) RequestEmailVerification(ctx context.Context, email string, newUser bool) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	user := new(entity.User)
	if err := c.UserRepository.FindByEmail(tx, user, email); err != nil {
		c.Log.Warnf("Failed find user by email : %+v", err)
		return fiber.NewError(fiber.StatusNotFound, "failed find user by email")
	}

	if user.HasVerifiedEmail() {
		c.Log.Warnf("User already verified")
		return fiber.NewError(fiber.StatusBadRequest, "user already verified")
	}

	if !newUser {
		verificationToken := new(entity.VerificationToken)
		if err := c.UserRepository.GetVerificationTokenByEmail(tx, verificationToken, email); err != nil {
			c.Log.Warnf("Resend request is not valid: %+v", err)
			return fiber.NewError(fiber.StatusBadRequest, "resend request is not valid")
		}
	}

	token := uuid.New().String()
	if err := c.UserRepository.CreateOrUpdateVerificationToken(tx, email, token); err != nil {
		c.Log.Warnf("Failed to create verification token: %+v", err)
		return fiber.ErrInternalServerError
	}

	encryptedToken, err := c.TokenHelper.Encrypt(token)
	if err != nil {
		c.Log.Warnf("Failed to register user : %+v", err)
		return fiber.ErrInternalServerError
	}

	token = encryptedToken

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return fiber.ErrInternalServerError
	}

	if err := c.EmailHelper.SendEmail(email, token, "new email verification"); err != nil {
		c.Log.Warnf("Failed to send Email: %+v", err)
	}

	c.Log.Infof("User with email: %+v successfully send their verification email", email)

	return nil
}

func (c *UserUseCase) VerifyEmail(ctx context.Context, request *model.VerifyEmailUserRequest) error {
	tx := c.DB.WithContext(ctx).Begin()

	helper.SanitiseStruct(request)

	decryptedToken, err := c.TokenHelper.Decrypt(request.Token)
	if err != nil {
		c.Log.Warnf("failed to register user : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, "invalid token")
	}

	request.Token = decryptedToken
	verificationToken := new(entity.VerificationToken)
	if err := c.UserRepository.GetEmailVerificationToken(tx, verificationToken, request.Email, request.Token); err != nil {
		c.Log.Warnf("failed to find Verification Token: %+v", err)
		return fiber.ErrNotFound
	}

	if time.Since(verificationToken.UpdatedAt) > 15*time.Minute {
		c.Log.Warnf("Verification token expired")
		return fiber.NewError(fiber.StatusUnauthorized, "verification token expired")
	}

	if err := c.UserRepository.SetUserEmailValidated(tx, verificationToken); err != nil {
		c.Log.Warnf("Failed to set user email validated: %+v", err)
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return fiber.ErrInternalServerError
	}

	c.Log.Infof("User with email: %+v successfully verified their email", request.Email)

	return nil
}

func (c *UserUseCase) RequestResetPassword(ctx context.Context, email string, newUser bool) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	user := new(entity.User)
	if err := c.UserRepository.FindByEmail(tx, user, email); err != nil {
		c.Log.Warnf("Failed find user by email : %+v", err)
		return fiber.NewError(fiber.StatusNotFound, "failed find user by email")
	}

	token := uuid.New().String()
	if err := c.UserRepository.CreateOrUpdateResetPasswordToken(tx, email, token); err != nil {
		c.Log.Warnf("Failed to create reset password token: %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, "failed to create reset password token")
	}

	encryptedToken, err := c.TokenHelper.Encrypt(token)
	if err != nil {
		c.Log.Warnf("Failed to register user : %+v", err)
		return fiber.ErrInternalServerError
	}

	token = encryptedToken

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return fiber.ErrInternalServerError
	}

	if err := c.EmailHelper.SendEmail(email, token, "reset password"); err != nil {
		c.Log.Warnf("Failed to send Email: %+v", err)
	}

	c.Log.Infof("User with email: %+v successfully send a reset password request", email)

	return nil
}

func (c *UserUseCase) ValidateResetToken(ctx context.Context, request *model.ValidateResetTokenRequest) (bool, error) {
	tx := c.DB.WithContext(ctx).Begin()

	decryptedToken, err := c.TokenHelper.Decrypt(request.Token)
	if err != nil {
		c.Log.Warnf("Failed to register user : %+v", err)
		return false, fiber.ErrBadRequest
	}

	request.Token = decryptedToken

	resetPasswordToken := new(entity.ResetPasswordToken)
	if err := c.UserRepository.GetResetPasswordTokenByEmail(tx, resetPasswordToken, request.Email, request.Token); err != nil {
		c.Log.Warnf("Failed to find Verification Token: %+v", err)
		return false, nil
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return false, fiber.ErrInternalServerError
	}
	return true, nil
}

func (c *UserUseCase) ResetPassword(ctx context.Context, request *model.ResetPasswordUserRequest) error {
	tx := c.DB.WithContext(ctx).Begin()

	decryptedToken, err := c.TokenHelper.Decrypt(request.Token)
	if err != nil {
		c.Log.Warnf("Failed to register user : %+v", err)
		return err
	}

	request.Token = decryptedToken

	resetPasswordToken := new(entity.ResetPasswordToken)
	if err := c.UserRepository.GetResetPasswordTokenByEmail(tx, resetPasswordToken, request.Email, request.Token); err != nil {
		c.Log.Warnf("Failed to find Verification Token: %+v", err)
		return fiber.ErrNotFound
	}

	if time.Since(resetPasswordToken.UpdatedAt) > 15*time.Minute {
		c.Log.Warnf("Reset password token expired")
		return fiber.NewError(fiber.StatusUnauthorized, "reset password token expired")
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Log.Warnf("Failed to generate bcrype hash : %+v", err)
		return fiber.ErrInternalServerError
	}

	if err := c.UserRepository.ResetPassword(tx, resetPasswordToken, string(password)); err != nil {
		c.Log.Warnf("Failed to set user email validated: %+v", err)
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return fiber.ErrInternalServerError
	}

	c.Log.Infof("User with email: %+v successfully reset their password", request.Email)

	return nil
}

func (c *UserUseCase) Login(ctx context.Context, request *model.LoginUserRequest) (*model.UserResponse, error) {

	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	user := new(entity.User)
	if err := c.UserRepository.FindByEmail(tx, user, request.Email); err != nil {
		c.Log.Warnf("Failed find user by id : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "email or password is invalid")
	}

	if !user.HasVerifiedEmail() {
		c.Log.Warnf("Email is not verified")
		return nil, fiber.NewError(fiber.StatusBadRequest, "unverified email")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		c.Log.Warnf("Failed to compare user password with bcrype hash : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "email or password is invalid")
	}

	accessTokenDetails, err := c.TokenHelper.GenerateAccessToken(user.UUID)
	if err != nil {
		c.Log.Warnf("Failed to generate token : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	user.AccessToken = accessTokenDetails.Token

	refreshTokenDetails, err := c.TokenHelper.GenerateRefreshToken(user.UUID)
	if err != nil {
		c.Log.Warnf("Failed to generate token : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := c.RedisClient.Set(ctx, refreshTokenDetails.Token, user.ID, time.Until(refreshTokenDetails.ExpiresAt)).Err(); err != nil {
		c.Log.Warnf("Failed to save token to redis : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	user.RefreshToken = refreshTokenDetails.Token

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	c.Log.Infof("User with email: %+v successfully login", request.Email)

	return converter.UserToTokenResponse(user), nil
}

func (c *UserUseCase) Current(ctx context.Context, request *model.GetUserRequest) (*model.UserResponse, error) {

	user := new(entity.User)
	if err := c.UserRepository.FindByEmail(c.DB, user, request.Email); err != nil {
		c.Log.Warnf("Failed find user by email : %+v", err)
		return nil, fiber.ErrNotFound
	}

	return converter.UserToResponse(user), nil
}

func (c *UserUseCase) Verify(ctx context.Context, request *model.VerifyUserRequest) (*model.Auth, error) {

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	accessTokenDetails, err := c.TokenHelper.VerifyAccessToken(request.Token)
	if err != nil {
		c.Log.Warnf("Failed to verify access token : %+v", err)
		return nil, fiber.ErrUnauthorized
	}

	user := new(entity.User)
	if err := c.UserRepository.FindByUUID(c.DB, user, accessTokenDetails.UserUUID); err != nil {
		c.Log.Warnf("Failed find user by id : %+v", err)
		return nil, fiber.ErrNotFound
	}

	userUUID, err := c.RedisClient.Get(ctx, request.Token).Result()

	if userUUID != "" {
		c.Log.Warnf("Access token found in Redis, user has already signed out : %+v", err)
		return nil, fiber.ErrUnauthorized
	}

	c.Log.Infof("User with email: %+v successfully authenticated and authorized", user.Email)

	return &model.Auth{UUID: user.UUID, Id: user.ID, Email: user.Email, Token: request.Token}, nil
}

func (c *UserUseCase) Logout(ctx context.Context, request *model.LogoutUserRequest) (bool, error) {
	tx := c.DB.WithContext(ctx).Begin()

	defer tx.Rollback()

	user := new(entity.User)
	if err := c.UserRepository.FindByEmail(tx, user, request.Email); err != nil {
		c.Log.Warnf("Failed find user by email : %+v", err)
		return false, fiber.ErrNotFound
	}

	if err := c.RedisClient.Set(ctx, request.AccessToken, user.ID, time.Until(time.Now().Add(15*time.Minute))).Err(); err != nil {
		c.Log.Warnf("Failed to save token to redis : %+v", err)
		return false, fiber.ErrInternalServerError
	}

	if err := c.RedisClient.Del(ctx, request.RefreshToken).Err(); err != nil {
		c.Log.Warnf("Failed to delete refresh token in redis : %+v", err)
		return false, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return false, fiber.ErrInternalServerError
	}

	c.Log.Infof("User with email: %+v successfully logout", user.Email)

	return true, nil
}

func (c *UserUseCase) AccessTokenRequest(ctx context.Context, request *model.AccessTokenRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	refreshTokenDetails, err := c.TokenHelper.VerifyRefreshToken(request.Token)
	if err != nil {
		c.Log.Warnf("Invalid refresh token : %+v", err)
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid refresh token")
	}

	user := new(entity.User)
	if err := c.UserRepository.FindByUUID(c.DB.WithContext(ctx), user, refreshTokenDetails.UserUUID); err != nil {
		c.Log.Warnf("Failed to find user by ID : %+v", err)
		return nil, fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	userID, err := c.RedisClient.Get(ctx, request.Token).Result()
	if err != nil || userID == "" {
		c.Log.Warnf("Refresh token not found in Redis : %+v", err)
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid refresh token")
	}

	accessTokenDetails, err := c.TokenHelper.GenerateAccessToken(user.UUID)
	if err != nil {
		c.Log.Warnf("Failed to generate access token : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	user.AccessToken = accessTokenDetails.Token

	if time.Until(refreshTokenDetails.ExpiresAt) < time.Hour*24 { // 1 hari sebelum kadaluarsa
		newRefreshTokenDetails, err := c.TokenHelper.GenerateRefreshToken(user.UUID)
		if err != nil {
			c.Log.Warnf("Failed to generate new refresh token : %+v", err)
			return nil, fiber.ErrInternalServerError
		}

		if err := c.RedisClient.Set(ctx, newRefreshTokenDetails.Token, user.ID, time.Until(newRefreshTokenDetails.ExpiresAt)).Err(); err != nil {
			c.Log.Warnf("Failed to save new refresh token to Redis : %+v", err)
			return nil, fiber.ErrInternalServerError
		}

		user.RefreshToken = newRefreshTokenDetails.Token
	} else {
		user.RefreshToken = request.Token
	}

	c.Log.Infof("User with email: %+v successfully requested access token", user.Email)

	return converter.UserToTokenResponse(user), nil
}

func (c *UserUseCase) Update(ctx context.Context, request *model.UpdateUserRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	user := new(entity.User)
	if err := c.UserRepository.FindByEmail(tx, user, request.Email); err != nil {
		c.Log.Warnf("Failed find user by id : %+v", err)
		return nil, fiber.ErrNotFound
	}

	if request.Password != "" {
		password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			c.Log.Warnf("Failed to generate bcrype hash : %+v", err)
			return nil, fiber.ErrInternalServerError
		}
		user.Password = string(password)
	}

	if err := c.UserRepository.Update(tx, user); err != nil {
		c.Log.Warnf("Failed save user : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.UserToResponse(user), nil
}
