package main

import (
	"strconv"
	"strings"
)

func minPartitions(n string) int {
	var result = 1
	for i := 9; i > 1; i-- {
		if strings.ContainsAny(n, strconv.Itoa(i)) {
			result = i
			break
		}
	}
	return result
}
