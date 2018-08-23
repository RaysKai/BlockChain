package consensus

import (
	"fmt"
)

type Service struct{
}

func (s *Service) Init(i interface{}) bool{
	fmt.Println("consensus service init...");
	return true
}

func (s *Service) Start() bool{
	fmt.Println("consensus service start...");
	return true
}

func (s *Service) Stop(){
	fmt.Println("consensus service stop...");
}