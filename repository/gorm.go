package repository

import (
	"errors"
	"gorm.io/gorm"
)

type Gorm[T Entity] struct {
	db *gorm.DB
}

func GetRepository[T Entity](db *gorm.DB) Repository[T] {
	return Gorm[T]{db}
}

func (r Gorm[T]) Find(id int) (*T, error) {
	var t T

	result := r.db.Where(id).First(&t)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	return &t, nil
}

func (r Gorm[T]) FindBy(fs ...Field) ([]T, error) {
	var t []T

	whereClause := make(map[string]interface{}, len(fs))

	for _, f := range fs {
		whereClause[f.Column] = f.Value
	}

	result := r.db.Where(whereClause).Find(&t)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	return t, nil
}

func (r Gorm[T]) FindFirstBy(fs ...Field) (*T, error) {
	ts, err := r.FindBy(fs...)
	if err != nil {
		return nil, err
	}

	if len(ts) >= 1 {
		return &ts[0], nil
	}

	return nil, errors.New("Record not found")
}

func (r Gorm[T]) Create(t *T) error {
	return r.db.Create(t).Error
}

func (r Gorm[T]) CreateBulk(ts []T) error {
	return r.db.Create(&ts).Error
}

func (r Gorm[T]) Update(t *T, fs ...Field) error {
	updateFields := make(map[string]interface{}, len(fs))

	for _, f := range fs {
		updateFields[f.Column] = f.Value
	}

	return r.db.Model(t).Updates(updateFields).Error
}

func (r Gorm[T]) Delete(t *T) error {
	return r.db.Delete(t).Error
}
