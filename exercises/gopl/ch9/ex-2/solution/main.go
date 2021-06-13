// Упражнение 9.2. Перепишите пример PopCount из раздела 2.6.2 так, чтобы он
// инициализировал таблицу поиска с использованием sync .Once при первом к ней об­
// ращении. (В реальности стоимость синхронизации для таких малых и высокооптими-
// зированных функций, как PopCount, является чрезмерно высокой.)
package popcount

import "sync"

var popOnce sync.Once

// pc[i] is the population count of i.
var pc [256]byte

// PopCount returns the population count (number of set bits) of x.
func PopCount(x uint64) int {
	popOnce.Do(func() {
		for i := range pc {
			pc[i] = pc[i/2] + byte(i&1)
		}
	})
	return int(pc[byte(x>>(0*8))] +
		pc[byte(x>>(1*8))] +
		pc[byte(x>>(2*8))] +
		pc[byte(x>>(3*8))] +
		pc[byte(x>>(4*8))] +
		pc[byte(x>>(5*8))] +
		pc[byte(x>>(6*8))] +
		pc[byte(x>>(7*8))])
}

//!-
