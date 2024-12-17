package repository

import (
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

func (r *Repository[T]) CountByName(db *gorm.DB, name string) (int64, error) {
	var total int64
	err := db.Model(new(T)).Where("LOWER(name) = LOWER(?)", name).Count(&total).Error
	return total, err
}

func (r *Repository[T]) FindManyByUUIDs(db *gorm.DB, uuids []uuid.UUID) ([]*T, error) {
	var model []*T

	if err := db.Where("uuid IN ?", uuids).Find(&model).Error; err != nil {
		return nil, err
	}

	return model, nil
}
