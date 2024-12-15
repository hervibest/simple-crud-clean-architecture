package repository

import (
	"simple-crud-clean-architecture/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository[T any] struct {
	DB *gorm.DB
}

func (r *Repository[T]) Create(db *gorm.DB, entity *T) error {
	return db.Create(entity).Error
}

func (r *Repository[T]) Update(db *gorm.DB, entity *T) error {
	return db.Save(entity).Error
}

func (r *Repository[T]) Delete(db *gorm.DB, entity *T) error {
	return db.Delete(entity).Error
}

func (r *Repository[T]) FindByUUID(db *gorm.DB, model *T, uuid uuid.UUID) error {
	return db.Where("uuid = ?", uuid).Take(model).Error
}

func (r *Repository[T]) FindById(db *gorm.DB, model *T, id int) error {
	return db.Where("id = ?", id).Take(model).Error
}

func (r *Repository[T]) FindByEmail(db *gorm.DB, model *T, email string) error {
	return db.Where("email = ?", email).Take(model).Error
}

func (r *Repository[T]) CountByEmail(db *gorm.DB, email any) (int64, error) {
	var total int64
	err := db.Model(new(T)).Where("email = ?", email).Count(&total).Error
	return total, err
}

func (r *Repository[T]) CountByName(db *gorm.DB, email string) (int64, error) {
	var total int64
	err := db.Model(new(T)).Where("name = ?", email).Count(&total).Error
	return total, err
}

func (r *Repository[T]) AttachUploadedFile(db *gorm.DB, model *T, file *entity.File) error {
	return db.Model(model).Association("File").Append(file)
}
