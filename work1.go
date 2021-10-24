package main

import (
	"fmt"
)

func main() {
	c1, c2, c3, c4 := make(chan bool), make(chan bool), make(chan bool), make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			fmt.Println("张三")
			c1 <- true
		}()
		<-c1

		go func() {
			fmt.Println("李四")
			c2 <- true
		}()
		<-c2

		go func() {
			fmt.Println("王五")
			c3 <- true
		}()
		<-c3

		go func() {
			fmt.Println("赵六")
			c4 <- true
		}()
		<-c4
	}
}
