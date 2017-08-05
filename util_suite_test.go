package Promise

import (
  . "github.com/onsi/ginkgo"
  "testing"
  "github.com/stretchr/testify/assert"
  "fmt"
)

func TestUtil(t *testing.T) {
  RunSpecs(t, "Promise Util Suite")
}

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
})
