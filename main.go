package main

import (
	"fmt"
	"github.com/jianshao/parseJSON/parse"
)

func main()  {
	json := parse.NewParseJSON("./test/valid/test.json")
	err := json.Load()
	fmt.Println(err)
}
