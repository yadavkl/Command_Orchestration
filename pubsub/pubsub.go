package pubsub

import(
	//"fmt"
	//"time"
	"errors"
	"sync"
)

var PSQ *psq


func init(){
	PSQ = NewPSQ()
}

type topic struct{
	name string
	que []interface{}	
	mux sync.Mutex		
}

type psq struct{
	topicMap map[string]*topic	
	mux sync.Mutex
}


func NewTopic(topicname string) *topic  {
	return &topic{
		name : topicname,
		que: make([]interface{},0),		
	}
}

func (t *topic) Add(entity interface{}){
	t.mux.Lock()
	defer t.mux.Unlock()
	t.que = append(t.que, entity)
}

func (t *topic) GetNext() (interface{}, error){
	t.mux.Lock()
	defer t.mux.Unlock()
	if len(t.que) > 0{
		ret := t.que[0]
		t.que = t.que[1:]
		return ret, nil
	}	
	return nil,errors.New("No data in queue...")	
}

func NewPSQ() *psq{
	return &psq{
		topicMap : make(map[string]*topic),
		//subMap : make(map[string]*subscriber),
	}
}

func (p *psq)AddTopic(t *topic){
	p.mux.Lock()
	defer p.mux.Unlock()
	p.topicMap[t.name]=t
}

func (p *psq)Publish(topicname string,entity interface{})  {
	p.mux.Lock()
	defer p.mux.Unlock()
	topic := p.topicMap[topicname]
	topic.Add(entity)	
}


func (p *psq)Subscribe(topicname string) <-chan interface{}{
	out := make(chan interface{})	
	topic, ok := p.topicMap[topicname]; if ok{		
		go func(){
			defer close(out)
			for{
				ret, err := topic.GetNext()
				if err != nil{
					//fmt.Println("Entity does not exist...")
				}else{
					out <- ret
				}
				//time.Sleep(time.Millisecond*100)
			}			
		}()				
	}
	return out
}
