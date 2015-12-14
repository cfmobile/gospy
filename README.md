#GoSpy [![Build Status](https://travis-ci.org/cfmobile/gospy.svg?branch=master)](https://travis-ci.org/cfmobile/gospy)

Go testing utility that lets you monitor calls to a function, verify which arguments were passed in and mock its behaviour.

This was created with unit testing in mind, to make it easier to verify interactions with dependencies and isolate components. Inspired by [Counterfeiter](https://github.com/maxbrunsfeld/counterfeiter) and [Cedar's Doubles](https://github.com/pivotal/cedar/wiki/Writing-specs#doubles)

Dependencies: [gmock](http://github.com/cfmobile/gmock)

Compatibility verified with Go 1.4 and up.

1. [Installation](#installation)
2. [Usage](#usage)
3. [Examples](#examples)
  1. [Using Spy()](#using-basic-constructor-spy)
  2. [Using SpyAndFake()](#using-spyandfake)
  3. [Using SpyAndFakeWithReturn()](#using-spyandfakewithreturn)
  4. [Using SpyAndFakeWithFunc()](#using-spyandfakewithfunc)
4. [API](#api)
  1. [Types](#types)
    1. [GoSpy](#gospy)
    2. [ArgList](#arglist)
    3. [CallList](#calllist)
  2. [Constructors](#constructors)
    1. [Spy()](#spy)
    2. [SpyAndFake()](#spyandfake)
    3. [SpyAndFakeWithReturn()](#spyandfakewithreturn)
    4. [SpyAndFakeWithFunc()](#spyandfakewithfunc)
  3. [GoSpy Methods](#gospy-methods)
    1. [GoSpy.Called()](#gospycalled)
    2. [GoSpy.CallCount()](#gospycallcount)
    3. [GoSpy.Calls()](#gospycalls)
    4. [GoSpy.ArgsForCall()](#gospyargsforcall)
    5. [GoSpy.Reset()](#gospyreset)
    6. [GoSpy.Restore()](#gospyrestore)

##Installation

Just use go get

```
  go get github.com/cfmobile/gospy
```

##Usage

- Import the header
```go
import (
  //...
  "github.com/cfmobile/gospy"
)
```

1. Call one of the GoSpy constructors by passing **a pointer to a Func** that you'd like to monitor. If you also wish to modify the behaviour of the target while the spy monitors it, pick one of the constructors that allow you to pass in a fake return value or fake implementation. **Store the returned `*GoSpy` object**
2. Run your code that makes calls to the target function.
3. Use the `*GoSpy` object that you stored to inspect the calls that were made to the target. You can check whether it was called, the number of calls that it received, and the arguments (if any) that were used in each of the calls.
4. If at any point you would like to start monitoring from scratch again, call `Reset()` in the `*GoSpy`. It just empties the object from the calls that have been recorded already, but doesn't change the fake behaviour (if any) that has been set.
5. When you're done monitoring your target, call `Restore()` on the `*GoSpy` object. **This is very important, otherwise the behaviour of your target will remain modified for the lifetime of your application.**. After restoring, it's safe to let the `*GoSpy` object get garbage-collected.

##Examples
### Using basic constructor `Spy()`
Showcases the main functionality of GoSpy - `Spy()` constructor, and GoSpy methods `Called()`, `CallCount()`, `Calls()`, `ArgsForCall()`, `Reset()` and `Restore()`.
```go
package main
import (
  "fmt"
  "github.com/cfmobile/gospy"
)

func main() {
  MyFunction := func(s string, i int) {
    // Do something
  }

  spy := gospy.Spy(&MyFunction)

  fmt.Println("** Spy Created **")
  fmt.Println("MyFunction was called? ", spy.Called())
  fmt.Println("Number of calls to MyFunction: ", spy.CallCount())
  fmt.Println("Calls to MyFunction: ", spy.Calls())

  MyFunction("abc", 0)
  MyFunction("bcd", 1)

  fmt.Println("** After Calls **")
  fmt.Println("MyFunction was called? ", spy.Called())
  fmt.Println("Number of calls to MyFunction: ", spy.CallCount())
  fmt.Println("Calls to MyFunction: ", spy.Calls())
  fmt.Println("Args for first call to MyFunction: ", spy.ArgsForCall(0))
  fmt.Println("Args for second call to MyFunction: ", spy.ArgsForCall(1))

  spy.Reset()
  MyFunction("cde", 2)

  fmt.Println("** After Reset **")
  fmt.Println("MyFunction was called? ", spy.Called())
  fmt.Println("Number of calls to MyFunction: ", spy.CallCount())
  fmt.Println("Calls to MyFunction: ", spy.Calls())
  fmt.Println("Args for first call to MyFunction: ", spy.ArgsForCall(0))

  spy.Restore()
  MyFunction("def", 3)

  fmt.Println("** After Restore **")
  fmt.Println("MyFunction was called? ", spy.Called())
  fmt.Println("Number of calls to MyFunction: ", spy.CallCount())
  fmt.Println("Args for calls to MyFunction: ", spy.Calls())
  fmt.Println("Args for first call to MyFunction: ", spy.ArgsForCall(0))
}
```
The execution of the example above should print:
```
** Spy Created **
MyFunction was called?  false
Number of calls to MyFunction:  0
Calls to MyFunction:  []
** After Calls **
MyFunction was called?  true
Number of calls to MyFunction:  2
Calls to MyFunction:  [[abc 0] [bcd 1]]
Args for first call to MyFunction:  [abc 0]
Args for second call to MyFunction:  [bcd 1]
** After Reset **
MyFunction was called?  true
Number of calls to MyFunction:  1
Calls to MyFunction:  [[cde 2]]
Args for first call to MyFunction:  [cde 2]
** After Restore **
MyFunction was called?  true
Number of calls to MyFunction:  1
Args for calls to MyFunction:  [[cde 2]]
Args for first call to MyFunction:  [cde 2]
```
**Note**
If you have a package function that you would like to spy on, like this:
```go
package main
import (
  "fmt"
  "github.com/cfmobile/gospy"
)

func MyFunction(s string, i int) {
  // Do something
}

func main() {
  //...
}
```
The Go compiler doesn't let you get a pointer to that `MyFunction` directly (`./main.go:13: cannot take the address of MyFunction`). In that case, we suggest using a func var that points to the desired function, and using that var in your code to invoke it:
```go
package main
import (
  "fmt"
  "github.com/cfmobile/gospy"
)

func MyFunction(s string, i int) {
  // Do something
}

func main() {
  MyPackageFunction := MyFunction

  spy := Spy(&MyPackageFunction)

  MyPackageFunction("abc", 0)
  //...
}
```

### Using `SpyAndFake()`
`SpyAndFake()` creates a default mock function that returns the Zero-value for each of the return values.
```go
package main
import (
  "fmt"
  "github.com/cfmobile/gospy"
)

func main() {
  MyFunction := func() (string, int) {
    return "Hello world", 1
  }

  text, num := MyFunction()

  fmt.Println("** Before Spy Creation **")
  fmt.Println("MyFunction() returns: ", text, num)

  spy := gospy.SpyAndFake(&MyFunction)
  text, num = MyFunction()

  fmt.Println("** After Spy Creation **")
  fmt.Println("MyFunction() returns: ", text, num)

  spy.Restore()
  text, num = MyFunction()

  fmt.Println("** After Spy Restore **")
  fmt.Println("MyFunction() returns: ", text, num)
}
```
Running the above example produces:
```
** Before Spy Creation **
MyFunction() returns:  Hello world 1
** After Spy Creation **
MyFunction() returns:   0
** After Spy Restore **
MyFunction() returns:  Hello world 1
```
Note that the second run (After Spy Creation) returns an empty string and the number 0, both are the zero-values for each type (string and int).

### Using `SpyAndFakeWithReturn()`
`SpyAndFakeWithReturn()` creates a mock function that returns the specified mock values. The constructor has to be called with the exact same number of arguments, types and order as the return values of the target, otherwise it panics. The only exception to that is to pass no arguments, which would cause this function to generate zero-values for each return type (in essence running like `SpyAndFake()`)
```go
package main
import (
  "fmt"
  "github.com/cfmobile/gospy"
)

func main() {
  MyFunction := func() (string, int) {
    return "Hello world", 1
  }

  text, num := MyFunction()

  fmt.Println("** Before Spy Creation **")
  fmt.Println("MyFunction() returns: ", text, num)

  spy := gospy.SpyAndFakeWithReturn(&MyFunction, "Mock return value", 101)
  text, num = MyFunction()

  fmt.Println("** After Spy Creation **")
  fmt.Println("MyFunction() returns: ", text, num)

  spy.Restore()
  text, num = MyFunction()

  fmt.Println("** After Spy Restore **")
  fmt.Println("MyFunction() returns: ", text, num)
}
```
Running the above example produces:
```
** Before Spy Creation **
MyFunction() returns:  Hello world 1
** After Spy Creation **
MyFunction() returns:  Mock return value 101
** After Spy Restore **
MyFunction() returns:  Hello world 1
```

### Using `SpyAndFakeWithFunc()`
`SpyAndFakeWithFunc()` allows you to replace the behaviour of the target with a custom mock function (with the same signature).
```go
package main
import (
  "fmt"
  "github.com/cfmobile/gospy"
)

func main() {
  MyFunction := func() (string, int) {
    return "Hello world", 1
  }

  text, num := MyFunction()

  fmt.Println("** Before Spy Creation **")
  fmt.Println("MyFunction() returns: ", text, num)

  spy := gospy.SpyAndFakeWithFunc(&MyFunction, func()(string, int) {
    return text + " fake", num + 1
  })
  text, num = MyFunction()

  fmt.Println("** After Spy Creation **")
  fmt.Println("MyFunction() returns: ", text, num)

  spy.Restore()
  text, num = MyFunction()

  fmt.Println("** After Spy Restore **")
  fmt.Println("MyFunction() returns: ", text, num)
}
```
Running the above example should print:
```
** Before Spy Creation **
MyFunction() returns:  Hello world 1
** After Spy Creation **
MyFunction() returns:  Hello world fake 2
** After Spy Restore **
MyFunction() returns:  Hello world 1
```

##API

###Types

#####GoSpy
```go
type GoSpy struct
```

The object that monitors a target function, records calls to the target and their arguments, and is responsible for mocking and restoring the target function's behaviour.

#####ArgList
```go
type ArgList []interface{}
```

Represents a list of arguments in a function call.

#####CallList
```go
type CallList []ArgList
```

Represents a list of function calls, complete with all arguments used.

###Constructors

#####Spy()
```go
func Spy(targetFuncPtr interface{}) *GoSpy
```

Basic constructor of a `GoSpy` object.
This constructor doesn't modify the behaviour of the target function.

**`targetFuncPtr`** has to be a pointer to a function. Any other type will cause the constructor to panic.

**Returns:** a pointer to a new `GoSpy` object. After the constructor returns, the target function will be monitored for calls until `spy.Restore()` is called.


#####SpyAndFake()
```go
func SpyAndFake(targetFuncPtr interface{}) *GoSpy
```

Constructor of a `GoSpy` object that modifies the target's behaviour to just return the Zero value for each of it's return values.

**`targetFuncPtr`** has to be a pointer to a function. Any other type will cause the constructor to panic.

**Returns:** a pointer to a new `GoSpy` object. After the constructor returns, the target function have it's behaviour modified and will be monitored for calls until `spy.Restore()` is called.

#####SpyAndFakeWithReturn()
```go
func SpyAndFakeWithReturn(targetFuncPtr interface{}, fakeReturnValues ...interface{}) *GoSpy
```

Constructor of a `GoSpy` object that modifies the target's behaviour to return a matching set of mock return values specified by the user.

**`targetFuncPtr`** has to be a pointer to a function. Any other type will cause the constructor to panic.

**`fakeReturnValues`** (variadic) has to be a set of values of the same type and appearing in the same order as the return values of the target.

**Note 1:** Passing a mock value of a non-matching type, passing an incomplete list of mock values (one for each return value) or passing the list in the wrong order will all cause this constructor to panic.

**Note 2:** Passing no mock values at all will cause this constructor to behave like `SpyAndFake()`

**Returns:** a pointer to a new `GoSpy` object. After the constructor returns, the target function have it's behaviour modified and will be monitored for calls until `spy.Restore()` is called.

#####SpyAndFakeWithFunc()
```go
func SpyAndFakeWithFunc(targetFuncPtr interface{}, mockFunc interface{}) *GoSpy
```

Constructor of a `GoSpy` object that modifies the target's behaviour to be replaced by another specified mock function (with identical signature).

**`targetFuncPtr`** has to be a pointer to a function. Any other type will cause the constructor to panic.

**`mockFunc`** has to be a function with the exact same signature as the target's, with alternate behaviour.

**Note 1:** Passing nil or a mock function with a different signature will cause this constructor to panic.

**Returns:** a pointer to a new `GoSpy` object. After the constructor returns, the target function have it's behaviour modified and will be monitored for calls until `spy.Restore()` is called.

###GoSpy Methods

#####GoSpy.Called()
```go
func (self *GoSpy) Called() bool
```

**Returns:** `bool` indicating whether the target has been called since the spy was constructed, or since the last call to `Reset()`.

#####GoSpy.CallCount()
```go
func (self *GoSpy) CallCount() int
```

**Returns:** `int` containing number of calls to the target recorded since the spy was constructed, or since the last call to `Reset()`

#####GoSpy.Calls()
```go
func (self *GoSpy) Calls() CallList
```

**Returns:** `CallList` containing all the calls that were made to the target recorded since the spy was constructed, or since the last call to `Reset()`. Returns `nil` if no calls have been recorded.

The `CallList` returned will contain all the arguments to all of the calls, preserving the order that they were made.

**Note:** Variadic functions (i.e. `func(args ...interface{})`) will store the variadic arguments as a single argument of type **`[]interface{}`**. When no arguments are passed in it records the equivalent to `[]interface{}{nil}`.

#####GoSpy.ArgsForCall()
```go
func (self *GoSpy) ArgsForCall(callIndex uint) ArgList
```

Retrieves the `ArgList` for a single call specified by `callIndex`.

**`callIndex`** is a zero-based index of the call you would like to inspect. If the `callIndex` is invalid (out of bounds), the function will panic.

**Returns:** `ArgList` containing the arguments for the selected call.

**Note:** This function effectively returns one specific entry that would be available in `Calls()`.

#####GoSpy.Reset()
```go
func (self *GoSpy) Reset()
```

Deletes the internal records of the `GoSpy` object. Immediate subsequent calls to: `Called()` would return `false`; `CallCount()` would return `0`;  `Calls()` would return `nil`; and `ArgsForCall()` with any index would panic.

Any subsequent calls to the target would be recorded normally.

**Note:** `Reset()` doesn't affect the monitoring of subsequent calls nor the modified behaviour of the target function. Use `Restore()` for that purpose.

#####GoSpy.Restore()
```go
func (self *GoSpy) Restore()
```

Restores the target function back to normal. The `GoSpy` will no longer record any subsequent calls and any modified behaviour in the function would be reverted back.

**Note:** Doesn't clear the `GoSpy` object from the calls that have been recorded up to that point.

**IMPORTANT: Always, always, ALWAYS Restore your function once you're done monitoring it, otherwise the changes to your target are permanent for the lifetime of your application.**
