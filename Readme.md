# Go repository

Golang database access using generics and implementation with gorm.

Copy and paste in your project (WIP library?)

Interface exposed:

```go
type Repository[T Entity] interface {
    Find(id int) (*T, error)
    FindBy(fs ...Field) ([]T, error)
    FindFirstBy(fs ...Field) (*T, error)
    CreateBulk(ts []T) error
    Create(t *T) error
    Update(t *T, data map[string]interface{}) error
    Delete(t *T) error
}
```

Call of repository for an `Example` entity using Gorm:

```go
r := repository.GetRepository[MyEntity](db)
```