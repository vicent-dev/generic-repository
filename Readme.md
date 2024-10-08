# Go repository

## About

Golang database access using generics and implementation with gorm.

Interface exposed:

```go
type Repository[T Entity] interface {
    Find(id string) (*T, error)
    FindWithRelations(id string) (*T, error)
    FindBy(fs ...Field) ([]*T, error)
    FindPaginated(pageSize, page int) (*Pagination, error)
    FindByWithRelations(fs ...Field) ([]*T, error)
    FindFirstBy(fs ...Field) (*T, error)
    CreateBulk(ts []T) error
    Create(t *T) error
    Update(t *T, fs ...Field) error
    Delete(t *T) error
}
```

## Install

```bash
go get github.com/vicent-dev/generic-repository
```

## Usage

Call of repository for an `MyEntity` entity using Gorm:

```go

import (
    repository "github.com/vicent-dev/generic-repository"
)

r := repository.GetGormRepository[MyEntity](db)

```

This function will return a repository of the struct type. 
Prepared for having only one instance in memory by struct type.
