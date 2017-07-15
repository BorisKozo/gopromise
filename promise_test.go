package Promise

import "testing"
import (
  "github.com/stretchr/testify/assert"
  "fmt"
  "time"
)

func TestNewPromise(t *testing.T) {
  done := false
  promiseInstance := NewPromise(func(resolve func(interface{}), reject func(error)) {
    done = true
  })
  promiseInternal := promiseInstance.(*promise)
  assert.NotNil(t, promiseInstance)
  assert.True(t, done)
  assert.Equal(t, promiseInternal.state, pendingState)
}

func TestPromiseResolve(t *testing.T) {
  promiseInstance := PromiseResolve(nil)
  promiseInternal := promiseInstance.(*promise)
  assert.NotNil(t, promiseInternal)
  assert.Equal(t, promiseInternal.state, resolvedState)
}

func TestPromiseReject(t *testing.T) {
  promiseInstance := PromiseReject(nil)
  promiseInternal := promiseInstance.(*promise)
  assert.NotNil(t, promiseInternal)
  assert.Equal(t, promiseInternal.state, rejectedState)
}

func TestPromise_Then(t *testing.T) {
  done := false
  promiseInstance := PromiseResolve("Foo")
  promiseInstance.Then(func(i interface{}) interface{} {
    foo := i.(string)
    assert.Equal(t, "Foo", foo)
    done = true
    return nil
  })
  assert.True(t, done)
}

func TestPromise_Catch(t *testing.T) {
  done := false
  promiseInstance := PromiseReject(fmt.Errorf("Error"))
  promiseInstance.Catch(func(i error) interface{} {
    assert.Equal(t, "Error", i.Error())
    done = true
    return nil
  })
  assert.True(t, done)
}

func TestPromiseChainThenThen(t *testing.T) {
  promiseInstance := PromiseResolve("foo").Then(func(i interface{}) interface{} {
    assert.Equal(t, "foo", i)
    return "bar"
  }).Then(func(i interface{}) interface{} {
    assert.Equal(t, "bar", i)
    return "baz"
  })
  promiseInternal := promiseInstance.(*promise)
  assert.Equal(t, "baz", promiseInternal.resolveValue)
}

func TestPromiseChainThenThenDelayed(t *testing.T) {
  doneChan := make(chan bool, 1)
  promiseInstance := NewPromise(func(resolve func(interface{}), reject func(error)) {
    go func() {
      time.Sleep(100 * time.Millisecond)
      resolve("foo")
      doneChan <- true
    }()
  })
  promiseInstance.Then(func(i interface{}) interface{} {
    assert.Equal(t, "foo", i)
    return "bar"
  }).Then(func(i interface{}) interface{} {
    assert.Equal(t, "bar", i)
    return "baz"
  })
  <-doneChan
}
