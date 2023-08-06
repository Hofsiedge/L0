package repo

// generic repo. T - main type, I - id type
type Repo[T any, I any] interface {
	Get(id I) (T, error)
	List() ([]I, error)
	GetAll() ([]T, error)
	Save(T) error
}
