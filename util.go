package Promise

import "sync"

func ThenOrCatch(promise Promise, resolveHandler PromiseResolveCallback, rejectHandler PromiseRejectCallback) Promise {
  return NewPromise(func(resolve func(interface{}), reject func(error)) {
    promise.Then(func(value interface{}) interface{} {
      resolve(resolveHandler(value))
      return nil
    })
    promise.Catch(func(value error) interface{} {
      resolve(rejectHandler(value))
      return nil
    })
  })
}

func All(promises []Promise) Promise {
  return NewPromise(func(resolve func(interface{}), reject func(error)) {
    total := len(promises)
    result := make([]interface{}, total)
    if total == 0 {
      resolve(result)
      return
    }
    count := 0
    hadError := false
    mutex := sync.Mutex{}
    for index, promise := range promises {
      innerIndex := index
      ThenOrCatch(promise, func(value interface{}) interface{} {
        result[innerIndex] = value
        mutex.Lock()
        count++
        equalLen := count == total
        mutex.Unlock()
        if equalLen {
          resolve(result)
        }
        return nil
      }, func(err error) interface{} {
        mutex.Lock()
        if !hadError {
          hadError = true
          mutex.Unlock()
          reject(err)
        } else {
          mutex.Unlock()
        }

        return nil
      })
    }
  })
}

func Race(promises []Promise) Promise {
  return NewPromise(func(resolve func(interface{}), reject func(error)) {
    anyReturned := false
    mutex := sync.Mutex{}
    for _, promise := range promises {
      ThenOrCatch(promise, func(value interface{}) interface{} {
        mutex.Lock()
        if anyReturned {
          mutex.Unlock()
          return nil
        }
        anyReturned = true
        mutex.Unlock()
        resolve(value)
        return nil
      }, func(err error) interface{} {
        mutex.Lock()
        if anyReturned {
          mutex.Unlock()
          return nil
        }
        anyReturned = true
        mutex.Unlock()
        reject(err)
        return nil
      })
    }
  })
}

func Every(promises []Promise) Promise {
  return NewPromise(func(resolve func(interface{}), reject func(error)) {
    total := len(promises)
    var results = make([]interface{}, total)
    var count = 0
    mutex := sync.Mutex{}
    for index, promise := range promises {
      innerIndex := index
      ThenOrCatch(promise, func(value interface{}) interface{} {
        results[innerIndex] = value
        mutex.Lock()
        count++
        equalLen := count == total
        mutex.Unlock()
        if equalLen {
          resolve(results)
        }
        return nil
      }, func(err error) interface{} {
        results[innerIndex] = err
        mutex.Lock()
        count++
        equalLen := count == total
        mutex.Unlock()
        if equalLen {
          resolve(results)
        }
        return nil
      })
    }
  })
}

func Run(fn func() interface{}) Promise {
  return NewPromise(func(resolve func(interface{}), reject func(error)) {
    go func() {
      result := fn()
      err, ok := result.(error)
      if ok {
        reject(err)
      } else {
        resolve(result)
      }
    }()
  })
}
