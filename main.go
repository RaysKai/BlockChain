package main

import (
	"fmt"
	"github.com/linkchain/config"
	)


func add(a int, b int) int{
	c := a+b;
	return c;
}

func main() {
	fmt.Printf("ret from add:%d\n", add(1,2));
	fmt.Printf("data:%d, %s\n", config.VarB, config.ConstB)
	config.Foo();
}
