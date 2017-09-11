package Promise

import (
  . "github.com/onsi/ginkgo"
  "github.com/stretchr/testify/assert"
  "fmt"
  "time"
)

var _ = Describe("Util", func() {
  var t = GinkgoT()
  BeforeEach(func() {
    t = GinkgoT()
  })
  Describe("ThenOrCatch", func() {
    It("should call Then of ThenOrCatch", func() {
      done := false
      promise := Resolve("Foo")
      ThenOrCatch(promise, func(i interface{}) interface{} {
        assert.Equal(t, "Foo", i)
        done = true
        return nil
      }, func(i error) interface{} {
        assert.Fail(t, "Should not be here")
        return nil
      })

      assert.True(t, done)
    })

    It("should call Catch of ThenOrCatch", func() {
      done := false
      promise := Reject(fmt.Errorf("Foo"))
      ThenOrCatch(promise, func(i interface{}) interface{} {
        assert.Fail(t, "Should not be here")
        return nil
      }, func(i error) interface{} {
        assert.Equal(t, "Foo", i.Error())
        done = true
        return nil
      })

      assert.True(t, done)

    })

    It("should call Then after ThenOrCatch", func() {
      done := false
      promise := Resolve("Foo")
      ThenOrCatch(promise, func(i interface{}) interface{} {
        assert.Equal(t, "Foo", i)
        return "Bar"
      }, func(i error) interface{} {
        assert.Fail(t, "Should not be here")
        return nil
      }).Then(func(value interface{}) interface{} {
        done = true
        assert.Equal(t, "Bar", value)
        return nil
      })

      assert.True(t, done)
    })

    It("should call Catch after ThenOrCatch", func() {
      done := false
      promise := Resolve("Foo")
      ThenOrCatch(promise, func(i interface{}) interface{} {
        assert.Equal(t, "Foo", i)
        return fmt.Errorf("Bar")
      }, func(i error) interface{} {
        assert.Fail(t, "Should not be here")
        return nil
      }).Catch(func(err error) interface{} {
        done = true
        assert.Equal(t, "Bar", err.Error())
        return nil
      })

      assert.True(t, done)
    })
  })

  Describe("All", func() {
    It("should resolve if all promises were resolved", func() {
      promise1 := Resolve(1)
      promise2 := Resolve(2)
      done := false
      All([]Promise{promise1, promise2}).Then(func(values interface{}) interface{} {
        results := values.([]interface{})
        assert.Len(t, results, 2)
        assert.Equal(t, 1, results[0])
        assert.Equal(t, 2, results[1])
        done = true
        return nil
      }).Catch(func(err error) interface{} {
        assert.Fail(t, "should not be here")
        return nil
      })
      assert.True(t, done)
    })

    It("should reject if any promise rejects", func() {
      promise1 := Resolve(1)
      promise2 := Resolve(2)
      promise3 := Reject(fmt.Errorf("Error!"))
      done := false
      All([]Promise{promise1, promise2, promise3}).Then(func(values interface{}) interface{} {
        assert.Fail(t, "should not be here")
        return nil
      }).Catch(func(err error) interface{} {
        assert.Equal(t, "Error!", err.Error())
        done = true
        return nil
      })
      assert.True(t, done)
    })

    It("should reject with only one promise if any promise rejects", func() {
      promise1 := Resolve(1)
      promise2 := Reject(fmt.Errorf("Error 1"))
      promise3 := Reject(fmt.Errorf("Error 2"))
      done := false
      All([]Promise{promise1, promise2, promise3}).Then(func(values interface{}) interface{} {
        assert.Fail(t, "should not be here")
        return nil
      }).Catch(func(err error) interface{} {
        assert.Equal(t, "Error 1", err.Error())
        done = true
        return nil
      })
      assert.True(t, done)
    })

    It("should resolve if no promises are passed", func() {
      done := false
      All([]Promise{}).Then(func(values interface{}) interface{} {
        results := values.([]interface{})
        assert.Len(t, results, 0)
        done = true
        return nil
      }).Catch(func(err error) interface{} {
        assert.Fail(t, "should not be here")
        return nil
      })
      assert.True(t, done)
    })
  })

  Describe("Race", func() {
    It("should resolve with the first resolved promise", func() {
      promise1 := Resolve(1)
      promise2 := Resolve(2)
      done := false
      Race([]Promise{promise1, promise2}).Then(func(value interface{}) interface{} {
        assert.Equal(t, 1, value)
        done = true
        return nil
      }).Catch(func(value error) interface{} {
        assert.Fail(t, "should not be here")
        return nil
      })
      assert.True(t, done)
    })

    It("should resolve with the first resolved promise even if the others reject", func() {
      promise1 := Resolve(1)
      promise2 := Reject(fmt.Errorf("err"))
      done := false
      Race([]Promise{promise1, promise2}).Then(func(value interface{}) interface{} {
        assert.Equal(t, 1, value)
        done = true
        return nil
      }).Catch(func(value error) interface{} {
        assert.Fail(t, "should not be here")
        return nil
      })
      assert.True(t, done)
    })

    It("should reject with the first rejected promise", func() {
      promise1 := Resolve(1)
      promise2 := Reject(fmt.Errorf("err"))
      done := false
      Race([]Promise{promise2, promise1}).Then(func(value interface{}) interface{} {
        assert.Fail(t, "should not be here")
        return nil
      }).Catch(func(value error) interface{} {
        assert.Equal(t, "err", value.Error())
        done = true
        return nil
      })
      assert.True(t, done)
    })
  })

  Describe("Every", func() {
    It("should resolve if all promises were resolved", func() {
      promise1 := Resolve(1)
      promise2 := Resolve(2)
      done := false
      Every([]Promise{promise1, promise2}).Then(func(values interface{}) interface{} {
        results := values.([]interface{})
        assert.Len(t, results, 2)
        assert.Equal(t, 1, results[0])
        assert.Equal(t, 2, results[1])
        done = true
        return nil
      }).Catch(func(err error) interface{} {
        assert.Fail(t, "should not be here")
        return nil
      })
      assert.True(t, done)
    })

    It("should resolve even if a promise rejects", func() {
      promise1 := Resolve(1)
      promise2 := Resolve(2)
      promise3 := Reject(fmt.Errorf("Error!"))
      done := false
      Every([]Promise{promise1, promise2, promise3}).Then(func(values interface{}) interface{} {
        results := values.([]interface{})
        assert.Len(t, results, 3)
        assert.Equal(t, 1, results[0])
        assert.Equal(t, 2, results[1])
        assert.Equal(t, "Error!", (results[2].(error)).Error())
        done = true
        return nil
      }).Catch(func(err error) interface{} {
        assert.Fail(t, "should not be here")
        return nil
      })
      assert.True(t, done)
    })
  })

  Describe("Run", func() {
    It("should run an async function and report the result", func() {
      endChan := make(chan bool)
      done := 1
      Run(func() interface{} {
        time.Sleep(2 * time.Millisecond)
        done = 2
        endChan <- true
        return "AAA"
      }).Then(func(i interface{}) interface{} {
        assert.Equal(t, i, "AAA")
        done = 3
        return nil
      })
      assert.Equal(t, 1, done)
      <-endChan
      time.Sleep(2 * time.Millisecond)
      assert.Equal(t, 3, done)
    })

    It("should run an async function and reject if there was an error", func() {
      endChan := make(chan bool)
      done := 1
      Run(func() interface{} {
        time.Sleep(2 * time.Millisecond)
        done = 2
        endChan <- true
        return fmt.Errorf("Oh no")
      }).Catch(func(i error) interface{} {
        assert.Equal(t, i.Error(), "Oh no")
        done = 3
        return nil
      })
      assert.Equal(t, 1, done)
      <-endChan
      time.Sleep(2 * time.Millisecond)
      assert.Equal(t, 3, done)
    })
  })
})
