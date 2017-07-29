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
})
