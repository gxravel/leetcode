/*
Упражнение 9.1. Добавьте функцию снятия со счета W ithdraw (am ount
i n t )
b ool в программу g o p l. io /c h 9 /b a n k l. Результат должен указывать, прошла ли
транзакция успешно или произошла ошибка из-за нехватки средств. Сообщение, от­
правляемое go-подпрограмме монитора, должно содержать как снимаемую сумму,
так и новый канал, по которому go-подпрограмма монитора сможет отправить булев
результат функции W ithdraw.
*/
package bank

import "fmt"

type withdraw struct {
	amount  int
	success chan bool
}

var deposits = make(chan int) // send amount to deposit
var balances = make(chan int) // receive balance
var withdraws = make(chan withdraw)

func Deposit(amount int) { deposits <- amount }
func Withdraw(amount int) bool {
	wd := withdraw{amount: amount, success: make(chan bool)}
	withdraws <- wd
	return <-wd.success
}
func Balance() int { return <-balances }

func teller() {
	var balance int // balance is confined to teller goroutine
	for {
		select {
		case amount := <-deposits:
			balance += amount
			fmt.Println("deposits, ", balance)
		case withdraw := <-withdraws:
			balance -= withdraw.amount
			fmt.Println("withdraws, ", balance)

			success := balance >= 0
			if !success {
				balance += withdraw.amount
				fmt.Println("withdraws (!success), ", balance)
			}
			withdraw.success <- success
		case balances <- balance:
		}
	}
}

func init() {
	go teller() // start the monitor goroutine
}

//!-
