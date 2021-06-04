/*
#729. My Calendar I
Implement a.booking MyCalendar class to store your events. A new event can be added if adding the event will not cause a.booking double booking.

Your class will have the method, book(int start, int end). Formally, this represents a.booking booking on the half open interval [start, end),
the range of real numbers x such that start <= x < end.

A double booking happens when two events have some non-empty intersection (ie., there is some time that is common to both events.)

For each call to the method MyCalendar.book, return true if the event can be added to the calendar successfully without causing
a.booking double booking. Otherwise, return false and do not add the event to the calendar.
*/
package main

import "fmt"

type Booking struct {
	start, end int
}

func (a MyCalendar) Len() int           { return len(a.booking) }
func (a MyCalendar) Swap(i, j int)      { a.booking[i], a.booking[j] = a.booking[j], a.booking[i] }
func (a MyCalendar) Less(i, j int) bool { return a.booking[i].end < a.booking[j].start }

type MyCalendar struct {
	booking []Booking
}

func Constructor() MyCalendar {
	return MyCalendar{[]Booking{}}
}

func (this *MyCalendar) Book(start int, end int) bool {
	var iStart int = -1
	var iEnd int = -1
	for i, val := range this.booking {
		if start < val.end && end > val.start {
			return false
		} else if start == val.end {
			iEnd = i
		} else if end == val.start {
			iStart = i
		}
	}
	if iEnd >= 0 {
		this.booking[iEnd].end = end
	} else if iStart >= 0 {
		this.booking[iStart].start = start
	} else {
		this.booking = append(this.booking, Booking{start, end})
	}
	return true
}

// type Booking struct {
// 	start, end int
// }

// type MyCalendar struct {
// 	booking map[int]bool
// }

// func Constructor() MyCalendar {
// 	return MyCalendar{[]Booking{}}
// }

// func (this *MyCalendar) Book(start int, end int) bool {

// }

func main() {
	c := Constructor()
	fmt.Println(c.Book(10, 20))
	fmt.Println(c.Book(10, 20))
	fmt.Println(c.Book(15, 25))
	fmt.Println(c.Book(20, 30))
}

/**
 * Your MyCalendar object will be instantiated and called as such:
 * obj := Constructor();
 * param_1 := obj.Book(start,end);
 */
