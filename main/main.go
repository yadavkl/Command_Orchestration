package main

import(
	"fmt"
	"time"
	"encoding/json"	
	"../pubsub"
	"../order"
	_"../storage"
	"../worker"
	"../orchestrator"
	"../generator"	
)

func main()  {
	
	topic_create := pubsub.NewTopic("create")
	topic_payment := pubsub.NewTopic("payment")
	topic_schedule := pubsub.NewTopic("schedule")
	topic_done := pubsub.NewTopic("done")

	pubsub.PSQ.AddTopic(topic_create)
	pubsub.PSQ.AddTopic(topic_payment)
	pubsub.PSQ.AddTopic(topic_schedule)
	pubsub.PSQ.AddTopic(topic_done)

	//<Subscribe topics returns channel
	create_chan := pubsub.PSQ.Subscribe("create")
	payment_chan := pubsub.PSQ.Subscribe("payment")
	schedule_chan := pubsub.PSQ.Subscribe("schedule")
	//done_chan := pubsub.PSQ.Subscribe("create")

	//<create worker pool for each type to work on
	wpool_create := worker.NewWorkerpool(create_chan, "create", 1)
	wpool_payment := worker.NewWorkerpool(payment_chan, "payment", 2)
	wpool_schedule := worker.NewWorkerpool(schedule_chan, "schedule", 3)

	go wpool_create.CreateWorkerPool(10)
	go wpool_payment.CreateWorkerPool(10)
	go wpool_schedule.CreateWorkerPool(10)
		
	orch := orchestrator.NewOrchestrator()
	orch.Publish()
	orch.SubscribeAndUpdate("done")
	go func (){
		start := time.Now()
		for{		
			tcmd, rcmd, ccmd :=orch.GetCounts()
			now := time.Now()
			fmt.Printf("Time: %f Seconds Total: %v Pending: %d Processed: %d Processing: %d\n", now.Sub(start).Seconds(),tcmd,rcmd, ccmd, (tcmd-rcmd-ccmd))
			time.Sleep(time.Second)
		}
	}()

	in := generator.CommandGenerator()
	for cmd := range in{
		var command order.Order
		err := json.Unmarshal([]byte(cmd.(string)), &command)
		if err == nil{
			orch.AddCommand(command)
		}
	}	
}