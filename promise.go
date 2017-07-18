package Promise

import "fmt"

const pendingState = "pending"
const fulfilledState = "fulfilled"
const rejectedState = "rejected"

type PromiseResolveCallback func(interface{}) interface{}
type PromiseRejectCallback func(error) interface{}

type Promise interface {
  Then(callback PromiseResolveCallback) Promise
  Catch(callback PromiseRejectCallback) Promise
}

type resolveCallbackData struct {
  callback PromiseResolveCallback
  //innerPromise Promise
  resolve func(interface{})
  reject  func(error)
}

type promise struct {
  state        string
  resolveValue interface{}
  rejectValue  error
  nextResolved []resolveCallbackData
  nextRejected []PromiseRejectCallback
}

func (p *promise) Then(callback PromiseResolveCallback) Promise {
  if p.state == fulfilledState {
    nextValue := callback(p.resolveValue)
    innerPromise, ok := nextValue.(Promise)
    if ok {
      return innerPromise
    }
    return PromiseResolve(nextValue)
  }

  if p.state == rejectedState {
    return PromiseReject(p.rejectValue)
  }

  callbackData := resolveCallbackData{callback: callback}
  innerPromise := NewPromise(func(resolve func(interface{}), reject func(error)) {
    callbackData.resolve = resolve
    callbackData.reject = reject
  })
  //callbackData.innerPromise = innerPromise
  p.nextResolved = append(p.nextResolved, callbackData)
  return innerPromise
}

func (p *promise) Catch(callback PromiseRejectCallback) Promise {
  if p.state == rejectedState {
    callback(p.rejectValue)
    return PromiseReject(p.rejectValue)
  }

  if p.state == fulfilledState {
    return PromiseResolve(p.resolveValue)
  }

  p.nextRejected = append(p.nextRejected, callback)
  return p
}

func resolveOrReject(value interface{}, callbackData resolveCallbackData) {
  err, isError := value.(error)
  if isError {
    callbackData.reject(err)
  } else {
    innerPromise, isPromise := value.(Promise)
    if isPromise {
      innerPromise.Then(func(innerValue interface{}) interface{} {
        resolveOrReject(innerValue, callbackData)
        return nil
      })
    } else {
      callbackData.resolve(value)
    }
  }
}

func (p *promise) handleResolve(value interface{}) {
  if p.state != pendingState {
    panic(fmt.Errorf("Trying to resolve a promise which is not pending but %v", p.state))
  }
  p.state = fulfilledState
  p.resolveValue = value
  for _, callbackData := range p.nextResolved {
    nextValue := callbackData.callback(value)
    resolveOrReject(nextValue, callbackData)
  }
}

func (p *promise) handleReject(err error) {
  if p.state != pendingState {
    panic(fmt.Errorf("Trying to reject a promise which is not pending but %v", p.state))
  }
  p.state = rejectedState
  p.rejectValue = err
  for _, callback := range p.nextRejected {
    callback(err)
  }
}

func defaultPromise() *promise {
  return &promise{
    state:        "pending",
    resolveValue: nil,
    rejectValue:  nil,
    nextResolved: []resolveCallbackData{},
    nextRejected: []PromiseRejectCallback{},
  }
}

func NewPromise(callback func(resolve func(interface{}), reject func(error))) Promise {
  result := defaultPromise()
  resolveFunc := func(value interface{}) {
    result.handleResolve(value)
  }

  rejectFunc := func(err error) {
    result.handleReject(err)
  }

  callback(resolveFunc, rejectFunc)
  return result
}

func PromiseResolve(value interface{}) Promise {
  result := defaultPromise()
  result.handleResolve(value)
  return result
}

func PromiseReject(err error) Promise {
  result := defaultPromise()
  result.handleReject(err)
  return result
}
