package orchestrator

import(
	"sync"
	"errors"
	//"time"
	"../pubsub"
	"../storage"
	"../order"
	//"log"
)


type orchestrator struct{
	commands []interface{}	
	cmdCount int
	complCount int
	mux sync.RWMutex
}

func NewOrchestrator() *orchestrator  {
	return &orchestrator{
		commands : make([]interface{},0),		
		//cmdCount :0,		
	}
}

func (os *orchestrator)AddCommand(command interface{})  {
	os.mux.Lock()
	defer os.mux.Unlock()
	os.commands = append(os.commands,command)
	os.cmdCount++;
}

func (os *orchestrator)AppendCommand(command interface{})  {
	os.mux.Lock()
	defer os.mux.Unlock()
	os.commands = append(os.commands,command)
}


//<Subscribe the completed task and update to database
//<Status can be read for inter relations

func (os *orchestrator) SubscribeAndUpdate(topic string){
	out := pubsub.PSQ.Subscribe(topic)
	go func (out <- chan interface{})  {
		for entity := range out{
			os.mux.Lock()
			storage.DB.Update(entity.(order.Order).Id, entity)
			os.complCount++
			os.mux.Unlock()
		}
	}(out)	
}

func (os *orchestrator) GetNextCommand() (interface{}, error){
	os.mux.Lock()
	defer os.mux.Unlock()

	if len(os.commands) > 0 {
		ret := os.commands[0]
		os.commands = os.commands[1:]		
		return ret, nil
	}
	return nil, errors.New("No more commands available...")
}


func (os *orchestrator)GetCounts() (int,int,int) {
	os.mux.RLock()
	tcount, rcount, compcount := os.cmdCount, len(os.commands), os.complCount
	//os.cmdCount=0
	//log.Print(os.cmdCount)
	os.mux.RUnlock()
	
	return tcount, rcount ,compcount
}


func (os *orchestrator)Publish()  {
	//<publish Commands available
	go func ()  {
		for{
			command, err := os.GetNextCommand()
			if err != nil{
				//log.Print(err.Error())
			}else{
				switch command.(order.Order).Type{
				case "create","Create","CREATE":
					pubsub.PSQ.Publish("create", command)
					storage.DB.Add(command.(order.Order).Id,command)
				case "payment","Payment","PAYMENT":
					//<check for creation theen publish
					//<otherwise put in pending list
					cmd := storage.DB.Get(command.(order.Order).Id)
					if cmd != nil{
						if cmd.(order.Order).Status == 1{
							pubsub.PSQ.Publish("payment", command)
						}else{
							os.AppendCommand(command)
						}
					}
					
				case "schedule","Schedule","SCHEDULE":
					cmd := storage.DB.Get(command.(order.Order).Id)
					if cmd != nil{
						if cmd.(order.Order).Status == 2{
							pubsub.PSQ.Publish("schedule", command)
						}else{
							os.AppendCommand(command)
						}
					}
				}				
			}
		}
		//time.Sleep(time.Millisecond*100)
	}()
}
