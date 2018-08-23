package common

type IService interface{
	Init(i interface{}) bool
	Start() bool
	Stop()
}