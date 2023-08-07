package postgres

import (
	"gitlab.com/Hofsiedge/l0/internal/domain"
)

// &MockOrders implements repo.Repo[domain.Order, string]
type MockOrders struct {
}

func (o *MockOrders) Get(id string) (domain.Order, error) {
	return domain.Order{}, nil
}

func (o *MockOrders) List() ([]string, error) {
	return make([]string, 0), nil
}

func (o *MockOrders) GetAll() ([]domain.Order, error) {
	return make([]domain.Order, 0), nil
}

func (o *MockOrders) Save(order domain.Order) error {
	return nil
}
