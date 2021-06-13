// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

package bank_test

import (
	"fmt"
	"testing"

	bank "github.com/gxravel/leetcode/exercises/gopl/ch9/ex-1/solution"
)

func TestBank(t *testing.T) {
	done := make(chan struct{})

	// Alice
	go func() {
		bank.Deposit(200)
		fmt.Println("=", bank.Balance())
		done <- struct{}{}
	}()

	// Bob
	go func() {
		bank.Deposit(100)
		done <- struct{}{}
	}()

	// Wait for both transactions.
	<-done
	<-done

	if got, want := bank.Balance(), 300; got != want {
		t.Errorf("Balance = %d, want %d", got, want)
	}
}

func TestWithdrawOverflow(t *testing.T) {
	done := make(chan struct{})

	bank.Withdraw(bank.Balance())

	// Alice
	go func() {
		bank.Deposit(200)
		fmt.Println("=", bank.Balance())
		done <- struct{}{}
	}()

	// Bob
	go func() {
		bank.Withdraw(300)
		done <- struct{}{}
	}()

	// Wait for both transactions.
	<-done
	<-done

	if got, want := bank.Balance(), 200; got != want {
		t.Errorf("Balance = %d, want %d", got, want)
	}
}

func TestWithdrawOK(t *testing.T) {
	done := make(chan struct{})
	var success bool

	bank.Withdraw(bank.Balance())

	// Alice
	go func() {
		bank.Deposit(200)
		fmt.Println("=", bank.Balance())
		done <- struct{}{}
	}()

	// Bob
	go func() {
		success = bank.Withdraw(100)
		done <- struct{}{}
	}()

	// Wait for both transactions.
	<-done
	<-done

	if got, want := bank.Balance(), 200; !success && got != want {
		t.Errorf("Not success, Balance = %d, want %d", got, want)
	}

	if got, want := bank.Balance(), 100; success && got != want {
		t.Errorf("Success, Balance = %d, want %d", got, want)
	}
}
