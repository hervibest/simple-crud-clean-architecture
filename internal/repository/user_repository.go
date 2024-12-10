package repository

import (
	"simple-crud-clean-architecture/internal/entity"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepository struct {
	Repository[entity.User]
	Log *logrus.Logger
}

func NewUserRepository(log *logrus.Logger) *UserRepository {
	return &UserRepository{
		Log: log,
	}
}

func (r *UserRepository) CountByEmail(db *gorm.DB, email any) (int64, error) {
	var total int64
	err := db.Model(new(entity.User)).Where("email = ?", email).Count(&total).Error
	return total, err
}

func (r *UserRepository) FindByEmail(db *gorm.DB, user *entity.User, email string) error {
	return db.Where("email = ?", email).Take(user).Error
}

func (r *UserRepository) FindByUUID(db *gorm.DB, user *entity.User, userUUID uuid.UUID) error {
	return db.Where("uuid = ?", userUUID).Take(user).Error
}

func (r *UserRepository) GetVerificationTokenByEmail(
	db *gorm.DB,
	verifToken *entity.VerificationToken,
	email string,
) error {
	return db.Model(&verifToken).
		Where("email = ?", email).
		Take(&verifToken).
		Error
}

func (r *UserRepository) GetEmailVerificationToken(
	db *gorm.DB,
	verifToken *entity.VerificationToken,
	email string,
	token string,
) error {
	return db.Model(&verifToken).
		Where("email = ? AND token = ?", email, token).
		Take(&verifToken).
		Error
}

func (r *UserRepository) CreateOrUpdateVerificationToken(db *gorm.DB, email string, token string) error {
	verificationToken := entity.VerificationToken{
		UserEmail: email,
		Token:     token,
	}

	if err := db.Save(&verificationToken).Error; err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) SetUserEmailValidated(db *gorm.DB, verificationToken *entity.VerificationToken) error {
	user := new(entity.User)
	if err := db.Where("email = ?", verificationToken.UserEmail).Take(user).Error; err != nil {
		return err
	}

	if err := db.Model(user).Update("verified_at", time.Now()).Error; err != nil {
		return err
	}

	if err := db.Model(verificationToken).Where("email = ? AND token = ?", verificationToken.UserEmail, verificationToken.Token).Delete(nil).Error; err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) CreateOrUpdateResetPasswordToken(db *gorm.DB, email string, token string) error {
	resetPasswordToken := entity.ResetPasswordToken{
		UserEmail: email,
		Token:     token,
	}

	if err := db.Save(&resetPasswordToken).Error; err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetResetPasswordTokenByEmail(db *gorm.DB, resetPasswordToken *entity.ResetPasswordToken, email string, token string) error {
	return db.Model(resetPasswordToken).Where("email = ? AND token = ?", email, token).Take(resetPasswordToken).Error
}

func (r *UserRepository) ResetPassword(db *gorm.DB, resetPasswordTokenDetails *entity.ResetPasswordToken, newPassword string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		user := new(entity.User)
		if err := tx.Where("email = ?", resetPasswordTokenDetails.UserEmail).Take(user).Error; err != nil {
			return err
		}

		if err := tx.Model(user).Update("password", newPassword).Error; err != nil {
			return err
		}

		if err := tx.Model(resetPasswordTokenDetails).Where("email = ? AND token = ?", resetPasswordTokenDetails.UserEmail, resetPasswordTokenDetails.Token).Delete(nil).Error; err != nil {
			return err
		}

		return nil
	})
}
