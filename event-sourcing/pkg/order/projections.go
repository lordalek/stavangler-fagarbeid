package order

import (
	"encoding/json"
	"fmt"
	"log"
	sync "sync"

	"github.com/nats-io/stan.go"
)

type inMemProjection struct {
	orders map[string]*OrderResponse
	sync.RWMutex
	conn stan.Conn
}

func NewProjectionDb(conn stan.Conn) *inMemProjection {
	db := &inMemProjection{
		make(map[string]*OrderResponse, 0),
		sync.RWMutex{},
		conn,
	}
	_, err := db.conn.Subscribe(
		"orders",
		db.handleMsg,
		stan.DeliverAllAvailable(),
	)

	if err != nil {
		log.Fatal(err)
	}

	return db
}

func (s *inMemProjection) Get(id string) *OrderResponse {
	s.RLock()
	defer s.RUnlock()
	or := s.orders[id]
	fmt.Printf("Returning %s at %s", or.Id, or.Restaurant)
	return or
}

func (s *inMemProjection) handleMsg(m *stan.Msg) {
	fmt.Printf("Processing msg seq %d\n", m.Sequence)
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	var stanUpcastedEvent stanUpcastedEvent

	err := json.Unmarshal(m.Data, &stanUpcastedEvent)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Processing %s, %v\n", stanUpcastedEvent.Type, stanUpcastedEvent.Data)

	switch stanUpcastedEvent.Type {
	case "*order.OrderCreated":
		var oc *OrderCreated
		err = json.Unmarshal(stanUpcastedEvent.Data, &oc)
		if err != nil {
			fmt.Println(err)
			return
		}

		s.orders[oc.Id] = &OrderResponse{
			Id:         oc.Id,
			Restaurant: oc.Restaurant,
			Orderlines: "none",
		}
		fmt.Printf("OrderId: %s at %s written to db\n", oc.Id, oc.Restaurant)
	default:
		fmt.Printf("did not parse %s\n", stanUpcastedEvent.Type)
	}

}
