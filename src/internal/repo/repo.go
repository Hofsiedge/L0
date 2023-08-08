package repo

// generic repo. T - main type, I - id type
type Repo[T any, I any] interface {
	Get(id I) (T, error)
	List() ([]I, error)
	GetAll() ([]T, error)
	Save(T) error
}

// &Mock[T, I] implements repo.Repo[T, I]
type Mock[T any, I any] struct {
}

func (o *Mock[T, I]) Close() error {
	return nil
}

func (o *Mock[T, I]) Get(_ I) (T, error) {
	var object T
	return object, nil
}

func (o *Mock[T, I]) List() ([]I, error) {
	return make([]I, 0), nil
}

func (o *Mock[T, I]) GetAll() ([]T, error) {
	return make([]T, 0), nil
}

func (o *Mock[T, I]) Save(_ T) error {
	return nil
}
