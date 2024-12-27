package usecase

import (
	"context"
	"simple-crud-clean-architecture/internal/entity"
	"simple-crud-clean-architecture/internal/helper"
	"simple-crud-clean-architecture/internal/model"
	"simple-crud-clean-architecture/internal/model/converter"
	"simple-crud-clean-architecture/internal/repository"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/redis/go-redis/v9"
)

type EmployeeUseCase struct {
	DB                 *gorm.DB
	Log                *logrus.Logger
	EmployeeRepository *repository.EmployeeRepository
	RedisClient        *redis.Client
	TokenHelper        *helper.TokenHelper
	EmailHelper        *helper.GomailSender
}

func NewEmployeeUseCase(db *gorm.DB, logger *logrus.Logger, employeeRepository *repository.EmployeeRepository,
	redisClient *redis.Client, tokenHelper *helper.TokenHelper, emailHelper *helper.GomailSender) *EmployeeUseCase {
	return &EmployeeUseCase{
		DB:                 db,
		Log:                logger,
		RedisClient:        redisClient,
		EmployeeRepository: employeeRepository,
		TokenHelper:        tokenHelper,
		EmailHelper:        emailHelper,
	}
}

func (c *EmployeeUseCase) Create(ctx context.Context, request *model.RegisterEmployeeRequest) (*model.EmployeeResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	helper.SanitiseStruct(request)

	total, err := c.EmployeeRepository.CountByEmail(tx, request.Email)
	if err != nil {
		c.Log.Warnf("Failed count employee from database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if total > 0 {
		c.Log.Warnf("Employee already exists : %+v", err)
		return nil, fiber.NewError(fiber.StatusConflict, "employee already exists")
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Log.Warnf("Failed to generate bcrypt hash : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	employee := &entity.Employee{
		UUID:     uuid.New(),
		Email:    request.Email,
		Password: string(password),
	}

	if err := c.EmployeeRepository.Create(tx, employee); err != nil {
		c.Log.Warnf("Failed create employee to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.EmployeeToResponse(employee), nil
}

func (c *EmployeeUseCase) Login(ctx context.Context, request *model.LoginEmployeeRequest) (*model.EmployeeResponse, error) {

	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	employee := new(entity.Employee)
	if err := c.EmployeeRepository.FindByEmail(tx, employee, request.Email); err != nil {
		c.Log.Warnf("Failed find employee by id : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "email or password is invalid")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(employee.Password), []byte(request.Password)); err != nil {
		c.Log.Warnf("Failed to compare employee password with bcrype hash : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "email or password is invalid")
	}

	accessTokenDetails, err := c.TokenHelper.GenerateEmployeeAccessToken(employee.UUID)
	if err != nil {
		c.Log.Warnf("Failed to generate token : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	employee.AccessToken = accessTokenDetails.Token

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.EmployeeToTokenResponse(employee), nil
}

func (c *EmployeeUseCase) Verify(ctx context.Context, request *model.VerifyEmployeeRequest) (*model.Auth, error) {

	accessTokenDetails, err := c.TokenHelper.VerifyEmployeeAccessToken(request.Token)
	if err != nil {
		c.Log.Warnf("Failed to verify access token : %+v", err)
		return nil, fiber.ErrUnauthorized
	}

	employee := new(entity.Employee)
	if err := c.EmployeeRepository.FindByUUID(c.DB, employee, accessTokenDetails.EmployeeUUID); err != nil {
		c.Log.Warnf("Failed find employee by uuid : %+v", err)
		return nil, fiber.ErrNotFound
	}

	employeeUUID, err := c.RedisClient.Get(ctx, request.Token).Result()

	if employeeUUID != "" {
		c.Log.Warnf("Access token found in Redis, employee has already signed out : %+v", err)
		return nil, fiber.ErrUnauthorized
	}

	return &model.Auth{UUID: employee.UUID, Id: employee.ID, Email: employee.Email, Token: request.Token}, nil
}

func (c *EmployeeUseCase) Current(ctx context.Context, request *model.GetEmployeeRequest) (*model.EmployeeResponse, error) {

	employee := new(entity.Employee)
	if err := c.EmployeeRepository.FindByEmail(c.DB, employee, request.Email); err != nil {
		c.Log.Warnf("Failed find employee by email : %+v", err)
		return nil, fiber.ErrNotFound
	}

	return converter.EmployeeToResponse(employee), nil
}

func (c *EmployeeUseCase) Logout(ctx context.Context, request *model.LogoutEmployeeRequest) (bool, error) {
	tx := c.DB.WithContext(ctx).Begin()

	defer tx.Rollback()

	employee := new(entity.Employee)
	if err := c.EmployeeRepository.FindByEmail(tx, employee, request.Email); err != nil {
		c.Log.Warnf("Failed find employee by email : %+v", err)
		return false, fiber.ErrNotFound
	}

	if err := c.RedisClient.Set(ctx, request.AccessToken, employee.ID, time.Until(time.Now().Add(1*24*time.Hour))).Err(); err != nil {
		c.Log.Warnf("Failed to save token to redis : %+v", err)
		return false, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return false, fiber.ErrInternalServerError
	}

	return true, nil
}
