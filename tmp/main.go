package main

import (
	"fmt"

	"github.com/cloudlibraries/merge"
)

func main() {
	merge.NewGroup(1, 2).MustMerge()
}
