package Promise

import "fmt"

const pendingState = "pending"
const resolvedState = "resolved"
const rejectedState = "rejected"

type PromiseResolveCallback func(interface{})
type PromiseRejectCallback func(error)
type Promise interface {
  Then(callback PromiseResolveCallback) Promise
  Catch(callback PromiseRejectCallback) Promise
}

type promise struct {
  state        string
  resolveValue interface{}
  rejectValue  error
  nextResolved []PromiseResolveCallback
  nextRejected []PromiseRejectCallback
}

func (p *promise) Then(callback PromiseResolveCallback) Promise {
  if p.state == resolvedState {
    callback(p.resolveValue)
    return PromiseResolve(p.resolveValue)
  }

  if p.state == rejectedState {
    return PromiseReject(p.rejectValue)
  }

  p.nextResolved = append(p.nextResolved, callback)
  return p
}

func (p *promise) Catch(callback PromiseRejectCallback) Promise {
  if p.state == rejectedState {
    callback(p.rejectValue)
    return PromiseReject(p.rejectValue)
  }

  if p.state == resolvedState {
    return PromiseResolve(p.resolveValue)
  }

  p.nextRejected = append(p.nextRejected, callback)
  return p
}

func (p *promise) handleResolve(value interface{}) {
  if p.state != pendingState {
    panic(fmt.Errorf("Trying to resolve a promise which is not pending but %v", p.state))
  }
  p.state = resolvedState
  p.resolveValue = value
  for _, callback := range p.nextResolved {
    callback(value)
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
    nextResolved: make([]PromiseResolveCallback, 1),
    nextRejected: make([]PromiseRejectCallback, 1),
  }
}

func NewPromise(callback func(resolve PromiseResolveCallback, reject PromiseRejectCallback)) Promise {
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
