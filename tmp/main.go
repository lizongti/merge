package main

import (
	"fmt"

	"github.com/cloudlibraries/merge"
)

func chanToSlice[t any](c chan t) []t {
	s := make([]t, len(c))
	for index, length := 0, len(c); index < length; index++ {
		s[index] = <-c
	}
	for index, length := 0, len(c); index < length; index++ {
		c <- s[index]
	}
	return s
}

func main() {
	c1 := make(chan int, 2)
	c2 := make(chan int, 3)
	c1 <- 1
	c1 <- 2
	c2 <- 3
	c2 <- 4
	c2 <- 5
	fmt.Println(chanToSlice(merge.MustMerge(c1, c2,
		merge.WithChanStrategy(merge.ChanStrategyIgnore),
	).(chan int)))
}
