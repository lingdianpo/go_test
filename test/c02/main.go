package main

import "fmt"

func main() {
	// 第一种
	var i int
	i = 10
	fmt.Println(i)
	// 第二种
	var j int = 20
	fmt.Println(j)

	// 第三种
	k := 30
	fmt.Println(k)

	a, b, c := 1, 2, "q"
	fmt.Println(a, b, c)

	//集合类型
	var (
		age  int
		name string
	)
	fmt.Println(age)
	fmt.Println(name)

}
