package worker

import(
	"sync"
	"fmt"
	"time"
	"../order"
	"../pubsub"
)

type workerpool struct{
	In <-chan interface{}	
	Status int 
	Topic string
}

func NewWorkerpool(in <-chan interface{}, topic string, status int ) *workerpool{
	return &workerpool{
		In : in,		
		Topic : topic,
		Status : status,
	}
}

func (wp *workerpool)CreateWorkerPool(noOfWorkers int ) {
	var wg sync.WaitGroup
	for i:=0; i<noOfWorkers; i++{
		//fmt.Printf("WOREKER: %d Type: %s \n", i, wp.Topic)
		wg.Add(1)
		go wp.Worker(&wg)
	}
	wg.Wait()
}

func (wp *workerpool)Worker(wg *sync.WaitGroup){
	defer wg.Done()
	for job := range wp.In{
		fmt.Printf("[%s] Order Id: %d and Type: %s Done.\n",wp.Topic,job.(order.Order).Id, job.(order.Order).Type)
		result := order.Order{
			Id: job.(order.Order).Id,
			Type : job.(order.Order).Type,
			Status : wp.Status,
		}
		pubsub.PSQ.Publish("done", result)
		time.Sleep(time.Millisecond*100)
	}
}