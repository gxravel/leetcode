// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

package memo_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/gxravel/leetcode/exercises/gopl/ch9/ex-3/memotest"
	memo "github.com/gxravel/leetcode/exercises/gopl/ch9/ex-3/solution"
)

var httpGetBody = memotest.HTTPGetBody

func TestSequential(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	m := memo.New(httpGetBody)
	defer m.Close()
	memotest.Sequential(t, m)
}

func TestConcurrent(t *testing.T) {
	m := memo.New(httpGetBody)
	defer m.Close()
	memotest.Concurrent(t, m)
}
