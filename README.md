# Go Promise

An [A+ promises](https://promisesaplus.com/) implementation in Go (with minor adjustments).

The goal of this project is to implement promises similar to the ones in native JavaScript, in native Go.
You can use Go concurrency without promises but for JavaScript developers
it is often more convenient to think about asynchronous and parallel programming
in terms of promise resolution/rejection.

The key difference is with the error handling. In JavaScript the catch clause is called when a promise is rejected
or when an exception is thrown. In my implementation the catch/rejection clause is signified by returning an error.
If an error is returned then the resolution mechanism assumes the promise needs to be rejected.
Currently panic is not recovered by the catch clause.

## Installation

````go get  github.com/BorisKozo/gopromise````

## API

All the functions work with the ````Promise```` interface. You can implement
your own ````Promise```` and use it in conjunction with my implementation.

#### NewPromise(func) (equivalent to new Promise(func))
Signature: ````func NewPromise(callback func(resolve func(interface{}), reject func(error))) Promise ````

Create a new Promise by calling _func_. You may resolve the resulting promise by calling
 _resolve_ with the resolved value, or reject the returned promise by calling _reject_ with the error.
 Note that _func_ is called immediately (not async).
 
 ````go
promiseInstance := NewPromise(func(resolve func(interface{}), reject func(error)) {
        resolve("my value")
      })

//promiseInstance is resolved with the value "my value"

````
 
#### PromiseResolve(value) (equivalent to Promise.resolve(value))
Signature: ````func PromiseResolve(value interface{}) Promise ````

Create a new promise which is already resolved with the given _value_

```go
 promiseInstance := PromiseResolve("my value")
 //promiseInstance is resolved with the value "my value"
```

#### PromiseReject(error) (equivalent to Promise.reject(error))
Signature: ````func PromiseReject(err error) Promise````

Create a new promise which is already rejected with the given _error_

```go
promiseInstance := PromiseReject(fmt.Errorf("my error"))
//promiseInstance is rejected with the error "my error"
```

#### Promise.Then(func) (equivalent to promise.then(func) without the reject callback)
Signature: ````func (p Promise) Then(handler func(interface{}) interface{}) Promise ````

Registers a resolve handler on the promise. Returns a new promise. If the caller promise
is resolved (or if it is already resolved) call the _handler_. If the _handler_ returns a Promise, it
will be chained in front of the returned promise. If the _handler_ returns an error then the returned promise
is rejected with that error. For any other value the returned promise is resolved with that value.

```go
PromiseResolve("foo").Then(func(value interface{}) interface{} {
        //value == "foo"
        return "bar"
      }).Then(func(anotherValue interface{}) interface{} {
        //anotherValue == "bar"
        return "baz"
      })
```

#### Promise.Catch(func) (equivalent to promise.catch(func))
Signature: ```` func (p Promise) Catch(handler func(error) interface{}) Promise```` 

Registers a reject handler on the promise. Returns a new promise. If the caller promise is 
rejected (or if it is already rejected) call the _handler_. If the _handler_ returns a promise
it will be chained in front of the returned promise. If the _handler_ returns an error then the returned promise
is rejected with that error. For any other value the returned promise is resolved with that value.

```go
PromiseReject(fmt.Errorf("Bad Error")).Catch(func(err error) interface{} {
 //err.Error() == "Bad Error"
        return nil
      })
```

## Change Log

**1.0.0**
- Added initial implementation for Promise interface
- Promise has Then and Catch functions
- NewPromise function creates a promise
- PromiseResolve and PromiseReject


## License
MIT