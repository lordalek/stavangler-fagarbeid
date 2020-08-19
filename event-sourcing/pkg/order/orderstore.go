package order

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/nats-io/stan.go"
)

type OrderStore interface {
	Save(*order) error
}

type StanAggregateStore struct {
	stan.Conn
}

type stanUpcastedEvent struct {
	Type string
	Data []byte
}

func parseEventInterface(e interface{}) ([]byte, error) {
	switch reflect.TypeOf(e).String() {
	case "*order.OrderCreated":
		event := e.(*OrderCreated)
		return json.Marshal(event)
	}
	return nil, errors.New("Unable to parse interface: " + reflect.TypeOf(e).String())
}

func (s *StanAggregateStore) Save(o *order) error {

	for _, e := range o.changes {
		eventBytes, err := parseEventInterface(e)

		if err != nil {
			return err
		}

		stanEventBytes, _ := json.Marshal(stanUpcastedEvent{
			reflect.TypeOf(e).String(),
			eventBytes,
		})
		fmt.Printf("Publishing %s, %v\n", reflect.TypeOf(e).String(), eventBytes)
		//implemnt changes to stan events mapping logic
		err = s.Conn.Publish("orders", stanEventBytes)

		if err != nil {
			return err
		}

	}

	return nil
}

func NewStore(conn stan.Conn) *StanAggregateStore {
	s := &StanAggregateStore{
		conn,
	}

	return s
}
