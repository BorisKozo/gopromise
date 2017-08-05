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
      ThenOrCatch(promise, func(value interface{}) interface{} {
        mutex.Lock()
        result[index] = value
        count++

        if count == total {
          mutex.Unlock()
          resolve(result)
        } else {
          mutex.Unlock()
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
