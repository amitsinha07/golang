package main

import "fmt"

func main()  {
	fmt.Println("Welcome to learn pointer")

	myNumber := 23;
	//var ptr  = &myNumber;

	Hii(&myNumber);
	fmt.Println("Outside Hii",myNumber)
}

func Hii(number *int)  {
	*number = 10;
	//fmt.Println("Inside Hii",number)
}