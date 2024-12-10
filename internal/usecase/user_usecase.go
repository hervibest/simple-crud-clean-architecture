package usecase

import (
	"context"
	"fmt"
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

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

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

	fmt.Println(user)

	if err := c.UserRepository.Create(tx, user); err != nil {
		c.Log.Warnf("Failed create user to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

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
	fmt.Println(token)
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

	return nil
}

func (c *UserUseCase) VerifyEmail(ctx context.Context, request *model.VerifyEmailUserRequest) error {
	tx := c.DB.WithContext(ctx).Begin()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("failed to validate request body")
		return fiber.NewError(fiber.StatusBadRequest, "failed to validate request body")
	}

	helper.SanitiseStruct(request)

	decryptedToken, err := c.TokenHelper.Decrypt(request.Token)
	if err != nil {
		c.Log.Warnf("failed to register user : %+v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "failed to validate request body")
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

	return nil
}

func (c *UserUseCase) ValidateResetToken(ctx context.Context, request *model.ValidateResetTokenRequest) (bool, error) {
	tx := c.DB.WithContext(ctx).Begin()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("failed to validate request body")
		return false, fiber.ErrBadRequest
	}

	decryptedToken, err := c.TokenHelper.Decrypt(request.Token)
	if err != nil {
		c.Log.Warnf("Failed to register user : %+v", err)
		return false, fiber.ErrBadRequest
	}

	request.Token = decryptedToken
	fmt.Println(request)

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

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("failed to validate request body")
		return fiber.ErrBadRequest
	}

	decryptedToken, err := c.TokenHelper.Decrypt(request.Token)
	if err != nil {
		c.Log.Warnf("Failed to register user : %+v", err)
		return err
	}

	request.Token = decryptedToken
	fmt.Println(request)

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
	return nil
}
