package Promise

import (
  . "github.com/onsi/ginkgo"
  "testing"
  "github.com/stretchr/testify/assert"
  "fmt"
  "time"
)

func TestPromise(t *testing.T) {
  RunSpecs(t, "Promise Suite")
}

var _ = Describe("Promise", func() {
  var t = GinkgoT()
  BeforeEach(func() {
    t = GinkgoT()
  })
  Describe("Constructors", func() {
    It("should create a new promise", func() {
      done := false
      promiseInstance := NewPromise(func(resolve func(interface{}), reject func(error)) {
        done = true
      })
      promiseInternal := promiseInstance.(*promise)
      assert.NotNil(t, promiseInstance)
      assert.True(t, done)
      assert.Equal(t, promiseInternal.state, pendingState)
    })

    It("should create a resolved promise", func() {
      promiseInstance := PromiseResolve(nil)
      promiseInternal := promiseInstance.(*promise)
      assert.NotNil(t, promiseInternal)
      assert.Equal(t, promiseInternal.state, fulfilledState)
    })

    It("should create a rejected promise", func() {
      promiseInstance := PromiseReject(nil)
      promiseInternal := promiseInstance.(*promise)
      assert.NotNil(t, promiseInternal)
      assert.Equal(t, promiseInternal.state, rejectedState)
    })

    It("should reject if a rejected promise is resolved", func() {
      done := false
      promise1 := PromiseReject(fmt.Errorf("error"))
      NewPromise(func(resolve func(interface{}), reject func(error)) {
        resolve(promise1)
      }).Then(func(i interface{}) interface{} {
        assert.Fail(t, "should not be here")
        return nil
      }).Catch(func(i error) interface{} {
        done = true
        return nil
      })
      assert.True(t, done)
    })
  })

  Describe("Then", func() {
    It("should call Then callback on a resolved promise", func() {
      done := false
      promiseInstance := PromiseResolve("Foo")
      promiseInstance.Then(func(i interface{}) interface{} {
        foo := i.(string)
        assert.Equal(t, "Foo", foo)
        done = true
        return nil
      })
      assert.True(t, done)
    })

    It("should call Then callback with chaining", func() {
      promiseInstance := PromiseResolve("foo").Then(func(i interface{}) interface{} {
        assert.Equal(t, "foo", i)
        return "bar"
      }).Then(func(i interface{}) interface{} {
        assert.Equal(t, "bar", i)
        return "baz"
      })
      promiseInternal := promiseInstance.(*promise)
      assert.Equal(t, "baz", promiseInternal.resolveValue)
    })

    It("should call Then callback when a previous callback returned a promise", func() {
      promiseInstance := PromiseResolve("foo").Then(func(i interface{}) interface{} {
        assert.Equal(t, "foo", i)
        return PromiseResolve("bar")
      }).Then(func(i interface{}) interface{} {
        assert.Equal(t, "bar", i)
        return "baz"
      })
      promiseInternal := promiseInstance.(*promise)
      assert.Equal(t, "baz", promiseInternal.resolveValue)
    })

    It("should call Then callback when the original promise is resolved in the future", func() {
      doneChan := make(chan bool, 1)
      promiseInstance := NewPromise(func(resolve func(interface{}), reject func(error)) {
        go func() {
          time.Sleep(10 * time.Millisecond)
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
    })

    It("should call Then callback when a previous callback returned a promise in the future", func() {
      doneChan := make(chan bool, 1)
      promiseInstance := NewPromise(func(resolve func(interface{}), reject func(error)) {
        go func() {
          time.Sleep(10 * time.Millisecond)
          resolve("foo")
          doneChan <- true
        }()
      })
      promiseInstance.Then(func(i interface{}) interface{} {
        assert.Equal(t, "foo", i)
        return PromiseResolve("bar")
      }).Then(func(i interface{}) interface{} {
        assert.Equal(t, "bar", i)
        return "baz"
      })
      <-doneChan
    })

    It("should call Then callback when there is a chain of 10 callbacks in the future", func() {
      doneChan := make(chan bool, 1)
      promiseInstance := NewPromise(func(resolve func(interface{}), reject func(error)) {
        go func() {
          time.Sleep(10 * time.Millisecond)
          resolve(2)
          doneChan <- true
        }()
      })
      currentPromise := promiseInstance
      for i := 0; i < 10; i++ {
        currentPromise = currentPromise.Then(func(i interface{}) interface{} {
          return PromiseResolve(i).Then(func(i interface{}) interface{} {
            return i.(int) * 2
          })
        })
      }

      var result int
      currentPromise.Then(func(i interface{}) interface{} {
        result = i.(int)
        return nil
      })
      <-doneChan
      assert.Equal(t, 2048, result)
    })

    It("should call Then callback after Catch callback if an error value was returned", func() {
      promiseInstance := PromiseResolve("foo").Then(func(i interface{}) interface{} {
        assert.Equal(t, "foo", i)
        return fmt.Errorf("Error!")
      }).Catch(func(err error) interface{} {
        return "baz"
      }).Then(func(i interface{}) interface{} {
        return "bat"
      })
      promiseInternal := promiseInstance.(*promise)
      assert.Equal(t, "bat", promiseInternal.resolveValue)
    })

    It("should call Then callback after Catch callback if an error value was returned in the future", func() {
      doneChan := make(chan bool, 1)
      promiseInstance := NewPromise(func(resolve func(interface{}), reject func(error)) {
        go func() {
          time.Sleep(10 * time.Millisecond)
          resolve("foo")
          doneChan <- true
        }()
      }).Then(func(i interface{}) interface{} {
        assert.Equal(t, "foo", i)
        return fmt.Errorf("Error!")
      }).Catch(func(err error) interface{} {
        return "baz"
      }).Then(func(i interface{}) interface{} {
        return "bat"
      })
      <-doneChan
      promiseInternal := promiseInstance.(*promise)
      assert.Equal(t, "bat", promiseInternal.resolveValue)
    })

    It("should call Then callback after Catch callback if a value was returned", func() {
      promiseInstance := PromiseResolve("foo").Then(func(i interface{}) interface{} {
        assert.Equal(t, "foo", i)
        return "foo2"
      }).Catch(func(err error) interface{} {
        assert.Fail(t, "should not be here")
        return "baz"
      }).Then(func(i interface{}) interface{} {
        return "bat"
      })
      promiseInternal := promiseInstance.(*promise)
      assert.Equal(t, "bat", promiseInternal.resolveValue)
    })

    It("should call Then callback after Catch callback if a value was returned in the future", func() {
      doneChan := make(chan bool, 1)
      promiseInstance := NewPromise(func(resolve func(interface{}), reject func(error)) {
        go func() {
          time.Sleep(10 * time.Millisecond)
          resolve("foo")
          doneChan <- true
        }()
      }).Then(func(i interface{}) interface{} {
        assert.Equal(t, "foo", i)
        return "foo2"
      }).Catch(func(err error) interface{} {
        assert.Fail(t, "should not be here")
        return "baz"
      }).Then(func(i interface{}) interface{} {
        return "bat"
      })
      <-doneChan
      promiseInternal := promiseInstance.(*promise)
      assert.Equal(t, "bat", promiseInternal.resolveValue)
    })

    It("should not call Then callback of a rejected promise", func() {
      PromiseReject(fmt.Errorf("foo")).Then(func(i interface{}) interface{} {
        assert.Fail(t, "Should not call Then callback on rejected promise")
        return nil
      })
    })
  })

  Describe("Catch", func() {
    It("should call Catch on a rejected promise", func() {
      done := false
      promiseInstance := PromiseReject(fmt.Errorf("Error"))
      promiseInstance.Catch(func(i error) interface{} {
        assert.Equal(t, "Error", i.Error())
        done = true
        return nil
      })
      assert.True(t, done)
    })

    It("should call Catch on a rejected promise in the future", func() {
      doneChan := make(chan bool, 1)
      promiseInstance := NewPromise(func(resolve func(interface{}), reject func(error)) {
        go func() {
          time.Sleep(10 * time.Millisecond)
          reject(fmt.Errorf("foo"))
          doneChan <- true
        }()
      })

      <-doneChan
      promiseInternal := promiseInstance.(*promise)
      assert.Equal(t, "foo", promiseInternal.rejectValue.Error())
    })

    It("should call Catch on a rejected promise after Then", func() {
      done := false
      promiseInstance := PromiseReject(fmt.Errorf("Error"))
      promiseInstance.Catch(func(i error) interface{} {
        assert.Equal(t, "Error", i.Error())
        done = true
        return nil
      })
      assert.True(t, done)
    })

    It("should call Catch on a rejected promise after Then in the future", func() {
      doneChan := make(chan bool, 1)
      promiseInstance := NewPromise(func(resolve func(interface{}), reject func(error)) {
        go func() {
          time.Sleep(10 * time.Millisecond)
          reject(fmt.Errorf("foo"))
          doneChan <- true
        }()
      }).Then(func(i interface{}) interface{} {
        assert.Fail(t, "Should not be here")
        return nil
      }).Catch(func(i error) interface{} {
        return "foo"
      })

      <-doneChan
      promiseInternal := promiseInstance.(*promise)
      assert.Equal(t, "foo", promiseInternal.resolveValue)
    })

    It("should call Catch if an internal promise was rejected", func() {
      done := false
      promiseInstance := NewPromise(func(resolve func(interface{}), reject func(error)) {
        innerPromise := PromiseResolve("Hello").Then(func(i interface{}) interface{} {
          return fmt.Errorf("Inner error")
        })
        resolve(innerPromise)
      }).Catch(func(i error) interface{} {
        assert.Equal(t, "Inner error", i.Error())
        done = true
        return nil
      })
      assert.NotNil(t, promiseInstance)
      assert.True(t, done)
    })
  })
})
