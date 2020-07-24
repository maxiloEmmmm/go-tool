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

	fmt.Println(lib.ArrayFirst(&test))
	fmt.Println(lib.InArray([]string{"!", "#"}, "#"))

	m := lib.ArrayKeyBy(&test, "Age")
	for key, mm := range m {
		fmt.Printf("group key: %d\n", key)
		fmt.Printf("items: \n")
		for _, item := range mm {
			fmt.Printf("- %s\n", item.(t).Name)
		}
	}

	// 计算年龄为x的人总年龄(没啥用的例子)
	// 1. array分组后 将分组结果单另计算
	fmt.Println(lib.MapMap(&m, func(d interface{}) interface{} {
		// 此处返回的为slice 所以不能再&d 不然就是取指针的地址了 造成非slice panic
		return lib.ArrayReduce(d, func(count float64, d interface{}) float64 {
			return count + float64(d.(t).Age)
		}, 0)
	}))

	// 2. array 分组且直接计算结果
	fmt.Println(lib.ArrayKeyByFunc(&test, "Age", func(old interface{}, d interface{}) interface{} {
		tmp := d.(t)
		// 代表第一项
		if old == nil {
			return tmp.Age
		}

		return old.(int) + tmp.Age
	}))

	fmt.Println(lib.ArrayMakeKey(&test, "Age"))

	fmt.Println(lib.ArrayMakeKeyFunc(test, "Age", func(d interface{}) interface{} {
		return d.(t).Name
	}))

	x := []int{1, 2, 3, 4, 5}
	fmt.Println(lib.ArrayReduce(lib.ArrayFilter(lib.ArrayMap(&x, func(d interface{}) interface{} {
		return d.(int) + 1
	}), func(d interface{}) bool {
		return d.(int) > 4
	}), func(count float64, d interface{}) float64 {
		return count + float64(d.(int))*2
	}, 0))

	fmt.Println(lib.ArrayPick(&test, []string{"Age"}))
	fmt.Println(lib.ArrayOmit(&test, []string{"Age"}))

	fmt.Println(lib.Has(&test, "4.Age1"))
	fmt.Println(lib.Has(&test, "4.Age"))

	value, exist := lib.Get(&test, "4.Age")

	if exist {
		fmt.Println(value)
	} else {
		fmt.Println("4.Age不存在")
	}

	fmt.Println(lib.AssetsReturn(false, "1", nil))

	testMap := map[string]int{
		"a": 1, "b": 2, "c": 3,
	}

	fmt.Println(lib.MapKeyFilter(&testMap, func(d interface{}) bool {
		return lib.InArray([]string{"a", "b"}, d.(string))
	}))

	fmt.Println(lib.MapKeyFilterToArray(&testMap, func(d interface{}) bool {
		return lib.InArray([]string{"a", "b"}, d.(string))
	}))

	fmt.Println(lib.MapValueFilter(&testMap, func(d interface{}) bool {
		return d.(int) >= 2
	}))

	fmt.Println(lib.MapValueFilterToArray(&testMap, func(d interface{}) bool {
		return d.(int) >= 2
	}))

	fmt.Println(lib.MapKeys(&testMap, nil))

	fmt.Println(lib.MapKeys(&testMap, func(d interface{}) interface{} {
		return lib.StringJoin(d.(string), "-haha")
	}))

	fmt.Println(lib.MapValues(&testMap, func(d interface{}) interface{} {
		return d.(int) + 1
	}))

	testMapSet := map[string]string{}

	lib.AssetMapSet(true, &testMapSet, "a", "哈")
	lib.AssetMapSet(false, &testMapSet, "b", "嘻")
	fmt.Println(testMapSet)
}
