package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nats-io/stan.go"
	"gitlab.com/Hofsiedge/l0/internal/domain"
	"gitlab.com/Hofsiedge/l0/internal/repo"
	"gitlab.com/Hofsiedge/l0/internal/repo/cache"
)

// server uses cache internally
type Server struct {
	Orders repo.Repo[domain.Order, string]
	Stan   stan.Subscription
}

func NewServer(orders repo.Repo[domain.Order, string]) (*Server, error) {
	orderCache, err := cache.New(orders)
	if err != nil {
		return nil, err
	}
	srv := Server{
		Orders: orderCache,
	}
	return &srv, nil
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
}

func (s *Server) ListEndpoint(response http.ResponseWriter, request *http.Request) {
	orders, err := s.Orders.List()
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		log.Print(err)
	}
	data, err := json.Marshal(orders)
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
	order, err := s.Orders.Get(id)
	if err != nil {
		if errors.Is(err, cache.ErrorNotFound) {
			response.WriteHeader(http.StatusNotFound)
		} else {
			response.WriteHeader(http.StatusInternalServerError)
		}
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
