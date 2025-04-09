package main

import (
	"cjson"
	"fmt"

	"github.com/goplus/lib/c"
)

func main() {
	jsonStr := c.Str(`{"name":"ZhangSan","age":20,"city":"Beijing","male":true}`)
	root := cjson.Parse(jsonStr)
	if root == nil {
		errPtr := cjson.GetErrorPtr()
		if errPtr != nil {
			fmt.Printf("parse error: %s\n", c.GoString(errPtr))
		}
		return
	}
	defer root.Delete()

	for child := root.Child; child != nil; child = child.Next {
		key := c.GoString(child.String)
		switch {
		case child.IsString() != 0:
			val := c.GoString(child.Valuestring)
			fmt.Printf("Key = %s, Value(String) = %s\n", key, val)
		case child.IsNumber() != 0:
			val := child.GetNumberValue()
			fmt.Printf("Key = %s, Value(Number) = %f\n", key, val)
		case child.IsBool() != 0:
			isTrue := child.IsTrue()
			fmt.Printf("Key = %s, Value(Bool) = %v\n", key, (isTrue != 0))
		default:
			fmt.Printf("Key = %s, Value(OtherType)\n", key)
		}
	}
}
