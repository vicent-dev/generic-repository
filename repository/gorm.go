package repository

import (
	"errors"
	"math"
	"reflect"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Gorm[T Entity] struct {
	db *gorm.DB
}

func GetGormRepository[T Entity](db *gorm.DB) Repository[T] {
	name := reflect.TypeOf((*T)(nil)).Elem().Name()

	if rs == nil {
		rs = make(map[string]interface{})
	}

	if r, ok := rs[name]; ok {
		return r.(Repository[T])
	}

	rs[name] = Gorm[T]{db}

	return rs[name].(Repository[T])
}

func (r Gorm[T]) Find(id string) (*T, error) {
	var t T

	result := r.db.First(&t, "id = ?", id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	return &t, nil
}

func (r Gorm[T]) FindPaginated(pageSize, page int) (*Pagination, error) {
	p := &Pagination{
		Limit: pageSize,
		Page:  page,
		Sort:  "created_at desc",
	}

	var totalRows int64

	var t T

	r.db.Model(t).Count(&totalRows)

	p.TotalRows = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(p.Limit)))
	p.TotalPages = totalPages

	var ts []*T
	r.db.Scopes(func(db *gorm.DB) *gorm.DB {
		return db.Offset(p.GetOffset()).Limit(p.GetLimit()).Order(p.GetSort())
	}).Find(&ts)

	p.Rows = ts

	return p, nil
}

func (r Gorm[T]) FindWithRelations(id string) (*T, error) {
	var t T

	result := r.db.Preload(clause.Associations).First(&t, "id = ?", id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	return &t, nil
}

func (r Gorm[T]) FindByWithRelations(fs ...Field) ([]*T, error) {
	var t []*T

	whereClause := make(map[string]interface{}, len(fs))

	for _, f := range fs {
		whereClause[f.Column] = f.Value
	}

	result := r.db.Preload(clause.Associations).Where(whereClause).Find(&t)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	return t, nil
}

func (r Gorm[T]) FindBy(fs ...Field) ([]*T, error) {
	var t []*T

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
		return ts[0], nil
	}

	return nil, errors.New("record not found")
}

func (r Gorm[T]) Create(t *T) error {
	return r.db.Create(t).Error
}

func (r Gorm[T]) CreateBulk(ts []T) error {
	return r.db.Create(&ts).Error
}

func (r Gorm[T]) Update(t *T, fs ...Field) error {
	updateClause := make(map[string]interface{}, len(fs))

	for _, f := range fs {
		updateClause[f.Column] = f.Value
	}

	return r.db.Model(t).Updates(updateClause).Error
}

func (r Gorm[T]) Delete(t *T) error {
	return r.db.Delete(t).Error
}
