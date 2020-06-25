package main

import (
	"fmt"
	lib "github.com/maxiloEmmmm/go-tool"
)

func main() {
	type t struct {
		Name string
		Age  int
	}

	test := []t{
		{Name: "hxm", Age: 12},
		{Name: "hxm1", Age: 13},
		{Name: "hxm2", Age: 12},
		{Name: "hxm3", Age: 13},
		{Name: "hxm4", Age: 15},
	}
	m := lib.ArrayKeyBy(test, "Age")

	for key, mm := range m {
		fmt.Printf("group key: %d\n", key)
		fmt.Printf("items: \n")
		for _, item := range mm {
			fmt.Printf("- %s\n", item.(t).Name)
		}
	}

	fmt.Println(lib.ArrayReduce(lib.ArrayFilter(lib.ArrayMap([]int{1, 2, 3, 4, 5}, func(d interface{}) interface{} {
		return d.(int) + 1
	}), func(d interface{}) bool {
		return d.(int) > 4
	}), func(count float64, d interface{}) float64 {
		return count + float64(d.(int))*2
	}, 0))

	fmt.Println(lib.ArrayPick(test, []string{"Age"}))
	fmt.Println(lib.ArrayOmit(test, []string{"Age"}))

	fmt.Println(lib.Has(test, "4.Age1"))
	fmt.Println(lib.Has(test, "4.Age"))
}
