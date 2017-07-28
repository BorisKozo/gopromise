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

type resolveRejector interface {
  resolve(interface{})
  reject(error)
}
type resolveCallbackData struct {
  callback PromiseResolveCallback
  //innerPromise Promise
  resolveFunc func(interface{})
  rejectFunc  func(error)
}

func (r resolveCallbackData) reject(err error) {
  r.rejectFunc(err)
}

func (r resolveCallbackData) resolve(value interface{}) {
  r.resolveFunc(value)
}

type rejectCallbackData struct {
  callback PromiseRejectCallback
  //innerPromise Promise
  resolveFunc func(interface{})
  rejectFunc  func(error)
}

func (r rejectCallbackData) resolve(value interface{}) {
  r.resolveFunc(value)
}

func (r rejectCallbackData) reject(err error) {
  r.rejectFunc(err)
}

type promise struct {
  state        string
  resolveValue interface{}
  rejectValue  error
  nextResolved []resolveCallbackData
  nextRejected []rejectCallbackData
}

func (p *promise) Then(callback PromiseResolveCallback) Promise {
  if p.state == fulfilledState {
    nextValue := callback(p.resolveValue)
    if innerPromise, ok := nextValue.(Promise); ok {
      return innerPromise
    }
    if innerError, ok := nextValue.(error); ok {
      return PromiseReject(innerError)
    }
    return PromiseResolve(nextValue)
  }

  if p.state == rejectedState {
    return PromiseReject(p.rejectValue)
  }

  callbackData := resolveCallbackData{callback: callback}
  innerPromise := NewPromise(func(resolve func(interface{}), reject func(error)) {
    callbackData.resolveFunc = resolve
    callbackData.rejectFunc = reject
  })
  //callbackData.innerPromise = innerPromise
  p.nextResolved = append(p.nextResolved, callbackData)
  return innerPromise
}

func (p *promise) Catch(callback PromiseRejectCallback) Promise {
  if p.state == rejectedState {
    nextValue := callback(p.rejectValue)
    if innerPromise, ok := nextValue.(Promise); ok {
      return innerPromise
    }
    if innerError, ok := nextValue.(error); ok {
      return PromiseReject(innerError)
    }
    return PromiseResolve(nextValue)
  }

  if p.state == fulfilledState {
    return PromiseResolve(p.resolveValue)
  }

  callbackData := rejectCallbackData{callback: callback}
  innerPromise := NewPromise(func(resolve func(interface{}), reject func(error)) {
    callbackData.resolveFunc = resolve
    callbackData.rejectFunc = reject
  })
  //callbackData.innerPromise = innerPromise
  p.nextRejected = append(p.nextRejected, callbackData)
  return innerPromise

}

func resolveOrReject(value interface{}, resolveRejector resolveRejector) {
  err, isError := value.(error)
  if isError {
    resolveRejector.reject(err)
  } else {
    innerPromise, isPromise := value.(Promise)
    if isPromise {
      innerPromise.Then(func(innerValue interface{}) interface{} {
        resolveOrReject(innerValue, resolveRejector)
        return nil
      })
    } else {
      resolveRejector.resolve(value)
    }
  }
}

func (p *promise) handleResolve(value interface{}) {
  if p.state != pendingState {
    panic(fmt.Errorf("Trying to resolve a promise which is not pending but %v", p.state))
  }
  innerPromise, isPromise := value.(Promise)
  if isPromise {
    innerPromise.Then(func(innerValue interface{}) interface{} {
      p.handleResolve(innerValue)
      return nil
    })
    innerPromise.Catch(func(innerError error) interface{} {
      p.handleReject(innerError)
      return nil
    })
    return
  }
  p.state = fulfilledState
  p.resolveValue = value
  for _, callbackData := range p.nextResolved {
    nextValue := callbackData.callback(value)
    resolveOrReject(nextValue, callbackData)
  }

  for _, callbackData := range p.nextRejected {
    callbackData.resolve(value)
  }

}

func (p *promise) handleReject(err error) {
  if p.state != pendingState {
    panic(fmt.Errorf("Trying to reject a promise which is not pending but %v", p.state))
  }
  p.state = rejectedState
  p.rejectValue = err
  for _, callbackData := range p.nextRejected {
    nextValue := callbackData.callback(err)
    resolveOrReject(nextValue, callbackData)
  }

  for _, callbackData := range p.nextResolved {
    callbackData.reject(err)
  }
}

func defaultPromise() *promise {
  return &promise{
    state:        "pending",
    resolveValue: nil,
    rejectValue:  nil,
    nextResolved: []resolveCallbackData{},
    nextRejected: []rejectCallbackData{},
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
