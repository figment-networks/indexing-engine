package main

import (
	"fmt"

	"github.com/figment-networks/indexing-engine/example/default/indexing"
)

func main() {
	if err := indexing.StartPipeline(); err != nil {
		fmt.Println("err", err)
	}
}
