// Упражнение 9.3. Расширьте тип Func и метод (*Memo) .G et так, чтобы вызыва­
// ющая функция могла предоставить необязательный канал done, с помощью которо­
// го можно было бы отменить операцию (раздел 8.9). Результаты отмененного вызова
// Func кешироваться не должны.
package memo

import "fmt"

//!+Func

// Func is the type of the function to memoize.
type Func func(key string, done chan struct{}) (interface{}, error)

// A result is the result of calling a Func.
type result struct {
	value interface{}
	err   error
}

type entry struct {
	res   result
	ready chan struct{} // closed when res is ready
}

//!-Func

//!+get

// A request is a message requesting that the Func be applied to key.
type request struct {
	key      string
	response chan<- result // the client wants a single result
	done     chan struct{}
}

type Memo struct{ requests chan request }

// New returns a memoization of f.  Clients must subsequently call Close.
func New(f Func) *Memo {
	memo := &Memo{requests: make(chan request)}
	go memo.server(f)
	return memo
}

func (memo *Memo) Get(key string, done chan struct{}) (interface{}, error) {
	response := make(chan result)
	memo.requests <- request{key, response, done}

	res := <-response
	return res.value, res.err
}

func (memo *Memo) Close() { close(memo.requests) }

//!-get

//!+monitor

func (memo *Memo) server(f Func) {
	cache := make(map[string]*entry)
	for req := range memo.requests {
		e := cache[req.key]
		if e == nil {
			// This is the first request for this key.
			e = &entry{ready: make(chan struct{})}
			cache[req.key] = e
			go e.call(f, req.key, req.done) // call f(key)
		}
		go e.deliver(req.response, req.done)
	}
}

func (e *entry) call(f Func, key string, done chan struct{}) {
	// Evaluate the function.
	e.res.value, e.res.err = f(key, done)
	// Broadcast the ready condition.
	select {
	case <-done:
		e.res.value, e.res.err = nil, fmt.Errorf("%s was canceled", key)
	default:
	}
	close(e.ready)
}

func (e *entry) deliver(response chan<- result, done chan struct{}) {
	// Wait for the ready condition.
	<-e.ready
	// Send the result to the client.
	response <- e.res
	select {
	case <-done:
		e = nil
	default:
	}
}

//!-monitor
