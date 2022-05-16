# result

`result` is an easier, more concise way to work with values that may exist at runtime, or may not exist due to an error.
It's best suited for holding the result of a function, thus the name. If you apprectiate go's explicit error handling,
but you're still a little tired to typing...
```go
x, err := foo()
if err != nil {
    return 0, fmt.Errorf("Couldn't foo: %v", err)
}
```
...then you're in the right place.

With `result`, you can replace that boilerplate error handling with 
```go
x := foo().
    OrError("Couldn't foo")
```

`result` is
* **faster to write.** You specify higher-level logic that abstracts out repetitive coding
* **easier to read.** Error handling is still handled explicitly, but better sepearated from the non-error flow of your
  code.
* **simpler to maintain.** The signature of a function is better seperated from its implementation. Error cases don't
  have to worry about what's returned on success cases and vice-versa.
* **harder to add bugs.** It's impossible to get the return value of a function without specifying exactly what to do in
  the event of an error.

## How it works

Imagine we want to generate out a very small *madlib*-style story. We need to get the specifics of the story from
functions that may return errors. If any of those functions return an error, we'll substitute in a default value. We
then want to send the story over a network, which also might result in a error. In regular go, it looks something like
this:

```go
func main() {
    n, err := name()
    if err != nil {
        n = "Aaron"
    }
    a, err := animal()
    if err != nil {
        a = "sheep"
    }
    v, err := verb()
    if err != nil {
        v = "painted"
    }
    o, err := object()
    if err != nil {
        o = "rainbow"
    }

    s := fmt.Sprintf("%v the %v %v a %v.", n, a, v, o)
    err = send(s)
    if err != nil {
        fmt.Printf("Couldn't send: %v\n", err)
    }
}

func name() (string, error) { /* ... */ }

func animal() (string, error) { /* ... */ }

func verb() (string, error) { /* ... */ }

func object() (string, error) { /* ... */ }

func send(s string) error { /* ... */ }
```

Okay, but there's a lot of boilerplate code in there.

With `result`, we change the return values of our worker functions from `(T, error)` to `result.Val[T]`, or from `error`
to `result.Status`. That allows us to significantly simplify the code that calls them. `Val` and `Status` support common
functions that return the underlying value when there's no error, and execute whatever backup logic you need if there
is.

In the following code for example, `name().OrUse("Aaron")` will return whatever the result of `name()` was, as long as
it wasnt' an error. If it was an error, we get the default "Aaron". Because it only returns a single value, we can
simplify `main()` down to 2 statements:

```go
func main() {
    s := fmt.Sprintf(
        "%v the %v %v a %v.",
        name().OrUse("Aaron"),
        animal().OrUse("sheep"),
        verb().OrUse("painted"),
        object().OrUse("rainbow"),
    )
    send(s).OrDo(func(e error) {
        fmt.Printf("Couldn't send: %v\n", e)
    })
}

func name() result.Val[string] { /* ... */ }

func animal() result.Val[string] { /* ... */ }

func verb() result.Val[string] { /* ... */ }

func object() result.Val[string] { /* ... */ }

func send(s string) result.Status { /* ... */ }
```

## Passing up errors

Go has great errors that tell you exactly what went wrong, especially in well-written code. Sadly, that well-written
code can be a bit repetitive and tedious to write. Let's go back to our previous example, but factor the main logic into
its own function that returns an error. Instead of using default values, if we get an error while filling in the story,
we'll pass the error back up to `main()`.

Without `result`:

```go
func main() {
    err := sendNewStory()
    if err != nil {
        fmt.Printf("Error sending new story: %v\n", err)
    }
}

func sendNewStory() error {
    n, err := name()
    if err != nil {
        return fmt.Errorf("Couldn't get name: %v", err)
    }
    a, err := animal()
    if err != nil {
        return fmt.Errorf("Couldn't get animal: %v", err)
    }
    v, err := verb()
    if err != nil {
        return fmt.Errorf("Couldn't get verb: %v", err)
    }
    o, err := object()
    if err != nil {
        return fmt.Errorf("Couldn't get object: %v", err)
    }

    s := fmt.Sprintf("%v the %v %v a %v.", n, a, v, o)
    err = send(s)
    if err != nil {
        return fmt.Errorf("Couldn't send: %v", err)
    }

    return nil
}

// Same definitions for dependent functions
```

And now using `result`:

```go
func main() {
    sendNewStory().OrDo(func(e error) {
        fmt.Printf("Error sending new story: %v\n", e)
    })
}

func sendNewStory() (res result.Status) {
    defer result.Handle(&res)

    s := fmt.Sprintf(
        "%v the %v %v a %v.",
        name().
            OrError("Couldn't get name"),
        animal().
            OrError("Couldn't get animal"),
        verb().
            OrError("Couldn't get verb"),
        object().
            OrError("Couldn't get object"),
    )
    send(s).
        OrError("Couldn't send")

    return result.Ok()
}

// Same definitions for dependent functions
```

Notice that instead of calling `OrUse("default value")` on the results, we're now calling `OrError("Context")`. Results
allow you to decide what you want to happen if the called function returns an error, while guaranteeing that you'll
either get a value you can work with, or that the function will stop execution at its current place. Here's all the
options:
* `OrDo(func(error))`: Only available on `result.Status` (for functions that return nothing on success); executes the
  given function if it encounters an error, then continues.
* `OrUse(T)` or `OrUse(T, U)`: Only available on `result.Val` or `result.Vals` (for functions that return values on
  success); uses the provided values as substitutes if an error occurs.
* `OrDoAndReturn(func(error))`: executes the given function, then stops execution of the function it occurs in.
* `OrError(string)`: stops execution of the function it occurs in, which then returns an error or error result. The
  given string is included in the error chain, along with the error that caused it to trigger.
* `OrPanic(string)`: panics with an error. The given string is included in the error chain, along with the error that
  caused it to trigger.

All of them can be used in the same way, chained after the call to a function that returns a `result.Status`,
`result.Val`, or `result.Vals`. The convention is to put any call that might halt execution of the enclosing function on
its own line (`OrDoAndReturn`, `OrError`, or `OrPanic`). This helps the reader see where execution of a function might
stop.

Notice also that we `defer result.Handle(&res)` at the start of the function. This is required if we're going to use any
functionality that halts execution of the enclosing function.
* Use `defer result.Handle(*result.Status)`, `defer result.Handle(*result.Val)`, or `defer result.Handle(*result.Vals)`
  if the function itself returns a result.
* Use `defer result.HandleError(*error)` if the function returns an error.
* Use `defer result.HandleReturn()` if the function doesn't return an error or a result which could hold an error. In
  this case, you can't use `OrError` within the function.