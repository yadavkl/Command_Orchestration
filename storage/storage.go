package storage

import(
	"sync"
)

type storage struct{
	stMap map[int]interface{}
	mux sync.Mutex
}

var DB *storage

func init()  {
	DB = NewStorage()
}

func NewStorage() *storage{
	return &storage{
		stMap : make(map[int]interface{}),
	}
}

func (s *storage)Add(key int, value interface{}){
	s.mux.Lock()
	defer s.mux.Unlock()
	s.stMap[key]=value
}

func (s *storage)Update(key int, value interface{})  {
	s.mux.Lock()
	defer s.mux.Unlock()
	_, ok := s.stMap[key]; if ok{
		s.stMap[key] = value
	}
}

func (s *storage)Delete(key int)  {
	s.mux.Lock()
	defer s.mux.Unlock()
	delete(s.stMap,key)	
}

func (s *storage)Get(key int) interface{} {
	s.mux.Lock()
	defer s.mux.Unlock()
	cmd, ok := s.stMap[key]; if !ok{
		return nil
	}
	return cmd
}