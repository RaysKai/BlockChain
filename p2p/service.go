package p2p

import (
	"fmt"
	//"github.com/linkchain/common"
)

type Service struct{
}

func (s *Service) Init(i interface{}) bool{
	fmt.Println("p2p service init...");
	return true
}

func (s *Service) Start() bool{
	fmt.Println("p2p service start...");
	return true
}

func (s *Service) Stop(){
	fmt.Println("p2p service stop...");
}