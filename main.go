package main

import (
	"fmt"
	"github.com/tangchen2018/go-utils/http"
)

func main() {

	c := http.NewClient().Get("http://www.baidu.com")
	_, err := c.DoMustSuccess()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(c.Result))
}
