


### Example

```go
package main

import (
	"fmt"
	"github.com/tangchen2018/go-utils/http"
)

func main() {
	req := http.HttpRequest{
		Url: "https://www.baidu.com",
	}
	if err := req.Do(); err != nil {
		panic(err)
	}

	fmt.Println(string(req.Result))
}
```
