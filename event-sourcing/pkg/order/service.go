package order

import (
	context "context"
	"fmt"

	"github.com/nats-io/stan.go"
)

type orderServer struct {
	UnimplementedOrderServer
	store OrderStore
	db    *inMemProjection
}

func NewServer() (*orderServer, error) {
	s := &orderServer{}

	sc, err := stan.Connect("test-cluster", "order-server-1")

	if err != nil {
		return nil, err
	}

	s.store = NewStore(sc)
	s.db = NewProjectionDb(sc)

	return s, nil
}

func (s *orderServer) GetOrder(c context.Context, r *FindOrderRequest) (*OrderResponse, error) {
	fmt.Printf("Getting order: %s\n", r.Id)

	res := s.db.Get(r.Id)

	return res, nil
}

func (s *orderServer) CreateOrder(c context.Context, r *CreateOrderRequest) (*ServerResponse, error) {
	fmt.Printf("%v\n", r)

	o, err := NewOrder(r.Restaurant, r.Id)

	if err != nil {
		return &ServerResponse{
			Result: "Oppsie",
		}, err
	}

	err = s.store.Save(o)

	return &ServerResponse{
		Result: "Ok, I guess",
	}, err
}
