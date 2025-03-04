package repo

import "gorm.io/gorm"

type Repository[T any] struct {
	db *gorm.DB
}

func (r *Repository[T]) Create(db *gorm.DB, entity *T) error {
	return r.db.Create(entity).Error
}

func (r *Repository[T]) Update(db *gorm.DB, entity *T) error {
	return r.db.Save(entity).Error
}

func (r *Repository[T]) Delete(db *gorm.DB, entity *T) error {
	return r.db.Delete(entity).Error
}

func (r *Repository[T]) FindById(db *gorm.DB, entity *T, id uint64) error {
	return r.db.Where("id = ?", id).Take(entity).Error
}
