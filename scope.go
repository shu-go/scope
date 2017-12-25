// package scope provides
package scope

import (
	"reflect"
)

type scope struct {
	begin interface{}
	end   func()
}

func New(begin interface{}, end func()) scope {
	begint := reflect.TypeOf(begin)
	if begint.Kind() != reflect.Func {
		panic("begin must be a func")
	}

	return scope{
		begin: begin,
		end:   end,
	}
}

// usage:
//   var Transaction = scope.New(db.Begin, db.Commit)
//	 Transaction.Block(func(tx *sql.Tx /* the last value can be omitted*/) {
//	   tx.Exec(...)
//	   // auto commit!
//	 })
func (s scope) Block(body interface{}, argsForBegin ...interface{}) {
	bodyv := reflect.ValueOf(body)
	bodyt := bodyv.Type()
	if bodyt.Kind() != reflect.Func {
		panic("body must be a func")
	}

	beginv := reflect.ValueOf(s.begin)
	begint := beginv.Type()

	if len(argsForBegin) != begint.NumIn() {
		panic("mismatch of args and begin")
	}

	withoutErr := false
	if begint.NumOut() == bodyt.NumIn() {
		// nop
	} else if begint.NumOut() == bodyt.NumIn()+1 {
		withoutErr = true
	} else {
		panic("mismatch of begin and body")
	}

	var argsv []reflect.Value
	for _, i := range argsForBegin {
		argsv = append(argsv, reflect.ValueOf(i))
	}
	retv := beginv.Call(argsv)

	if withoutErr {
		retv = retv[:len(retv)-1]
	}

	bodyv.Call(retv)

	s.end()
}

// Done is for simple begin-body-end flow, whereas end is deferred.
//
// usage:
//   var Transaction = scope.New(db.Begin, db.Commit)
//   var DoneCommitted = Transaction.Done // do not call, this time
//   {
//	    defer DoneCommitted()() // DOUBLE (), first to call begin, second to defer end
//   }
func (s scope) Done(argsForBegin ...interface{}) func() {
	beginv := reflect.ValueOf(s.begin)
	begint := beginv.Type()

	if len(argsForBegin) != begint.NumIn() {
		panic("mismatch of args and begin")
	}

	var argsv []reflect.Value
	for _, i := range argsForBegin {
		argsv = append(argsv, reflect.ValueOf(i))
	}
	beginv.Call(argsv)

	return s.end
}
