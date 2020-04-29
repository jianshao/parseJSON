package main

import (
	"./parse"
	"fmt"
)

func main()  {
	json := parse.NewParseJSON("../test/valid/test.json")
	err := json.Load()
	fmt.Println(err)
}
