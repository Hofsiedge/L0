package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nats-io/stan.go"
	"gitlab.com/Hofsiedge/l0/internal/domain"
	"gitlab.com/Hofsiedge/l0/internal/repo"
)

type Server struct {
	Orders repo.Repo[domain.Order, string]
	Stan   stan.Subscription
	Cache  map[string]domain.Order
}

func NewServer(orders repo.Repo[domain.Order, string]) (*Server, error) {
	srv := Server{
		Orders: orders,
		Cache:  make(map[string]domain.Order),
	}
	if err := srv.FillCache(); err != nil {
		return nil, err
	}
	return &srv, nil
}

func (s *Server) FillCache() error {
	orders, err := s.Orders.GetAll()
	if err != nil {
		return err
	}
	for _, order := range orders {
		s.Cache[order.OrderUid] = order
	}
	return nil
}

func (s *Server) HandleMessage(msg *stan.Msg) {
	data := msg.Data
	var order domain.Order
	if err := json.Unmarshal(data, &order); err != nil {
		log.Printf("could not unmarshal a message: %v", data)
		return
	}
	log.Printf("received a message with an Order: %v\n", order)
	if err := s.Orders.Save(order); err != nil {
		err = fmt.Errorf("error saving an order: %w", err)
		log.Println(err)
	}
	s.Cache[order.OrderUid] = order
}

func (s *Server) ListEndpoint(response http.ResponseWriter, request *http.Request) {
	data, err := json.Marshal(s.Cache)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		err = fmt.Errorf("error marshalling cache: %w", err)
		log.Print(err)
		return
	}
	response.Write(data)
}

func (s *Server) GetByIdEndpoint(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, ok := vars["id"]
	if !ok {
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	order, ok := s.Cache[id]
	if !ok {
		response.WriteHeader(http.StatusNotFound)
		return
	}
	data, err := json.Marshal(order)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		err = fmt.Errorf("error marshalling Order: %w", err)
		log.Print(err)
		return
	}
	response.Write(data)
}
