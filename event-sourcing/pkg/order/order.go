package order

import (
	"fmt"
	"reflect"
)

type order struct {
	changes []interface{}
	version int
}

func NewOrder(restaurant, id string) (*order, error) {
	o := &order{
		version: -1,
	}
	if err := o.apply(
		&OrderCreated{
			Restaurant: restaurant,
			Id:         id,
		},
	); err != nil {
		return nil, err
	}
	return o, nil
}

func (o *order) apply(e interface{}) error {
	o.when(e)
	if err := o.ensureValidState(); err != nil {
		return err
	}
	o.changes = append(o.changes, e)
	return nil
}

func (o *order) when(e interface{}) {
	switch reflect.TypeOf(e).String() {
	case "*order.OrderCreated":
		fmt.Printf("order created at %s\n", e)
	default:
		fmt.Printf("Default: %v, %s", e, reflect.TypeOf(e).String())
	}
}

func (o *order) ensureValidState() error {
	return nil
}
