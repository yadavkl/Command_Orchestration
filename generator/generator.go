package generator

import(
	"encoding/json"
	"time"
	//"math/rand"
)

type command struct{
	Id int
	Type string
}

func CommandGenerator() <-chan interface{}  {
	out := make(chan interface{})
	cmdmap := []string{"create","payment","schedule"}
	counter := 1
	go func (counter int)  {
		defer close(out)
		for{
			key := counter
			counter++
			for i:= 0 ; i<3; i++{
				cmd :=command{
					Id : key,
					Type : cmdmap[i],
				}
				jsonval, err := json.Marshal(cmd)
				if err == nil{
					out <- string(jsonval)
				}
			}
			time.Sleep(time.Millisecond*100)
		}
		
	}(counter)
	return out
}