package main

import (
	"fmt"

	"github.com/figment-networks/indexing-engine/pipeline/example/custom/indexing"
)

func main() {
	if err := indexing.StartPipeline(); err != nil {
		fmt.Println("err", err)
	}
}
