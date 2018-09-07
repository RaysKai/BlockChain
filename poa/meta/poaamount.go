package meta

import "strconv"

type POAAmount struct {
	Value int
}

func (a *POAAmount)GetInt() int  {
	return a.Value
}

func (a *POAAmount)GetFloat() float32  {
	return float32(a.Value)
}

func (a *POAAmount)GetString() string  {
	return strconv.Itoa(a.Value)
}