package dbutils

import (
	"gorm.io/gorm"
)

func FindAll[T any](db *gorm.DB, out *[]T) error {
	return db.Find(out).Error
}

func Create[T any](db *gorm.DB, obj *T) error {
	return db.Create(obj).Error
}

func UpdateByID[T any](db *gorm.DB, id int, column string, value interface{}) error {
	return db.Model(new(T)).Where("id = ?", id).Update(column, value).Error
}

func DeleteByID[T any](db *gorm.DB, id int) error {
	return db.Delete(new(T), id).Error
}
