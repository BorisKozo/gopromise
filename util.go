package Promise

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
