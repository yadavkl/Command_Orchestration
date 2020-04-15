package subscriber

import(
	"../order"
)

type subscriber struct{
	In chan order.Order
	Out chan order.Order
}

func NewSubscriber() *subscriber{
	return &subscriber{
		In : make(chan order.Order),
		Out : make(chan order.Order, 10),
	}
}

func (s* subscriber) Notify(order order.Order){
	s.In<- order
}