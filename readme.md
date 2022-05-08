# result

`result` is an easier, more concise way to work with values that may exist at runtime, or may not exist due to an error.
It's best suited for holding the result of a function, thus the name. If you apprectiate go's explicit error handling,
but you're still a little tired to typing...
```go
if err != nil {
    return 0, fmt.Errorf("Added context: %v", err)
}
```
...then you're in the right place.

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
        o = "a rainbow"
    }

    s := fmt.Sprintf("%v the %v %v %v.", n, a, v, o)
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

Okay, but there's a lot of boilerplate code in there. This is it with `result`:

```go
func main() {
    s := fmt.Sprintf(
        "%v the %v %v %v.",
        name().OrUse("Aaron"),
        animal().OrUse("sheep"),
        verb().OrUse("painted"),
        object().OrUse("a rainbow"),
    )
    send(s).OrDo(func(e error) {
        fmt.Printf("Couldn't send: %v\n", e)
    })
}

func name() result.Var[string] { /* ... */ }

func animal() result.Var[string] { /* ... */ }

func verb() result.Var[string] { /* ... */ }

func object() result.Var[string] { /* ... */ }

func send(s string) result.Status { /* ... */ }
```

## Passing up errors

Go has great errors that tell you exactly what went wrong, especially in well-written code. Sadly, that well-written
code can be a bit repetitive and tedious to write. Let's go back to our previous example, but factor the main logic into
its own function that returns an error. Instead of using default values, if we get an error while filling in the story,
we'll pass the error back up to `main()`.

In regular go:

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

    s := fmt.Sprintf("%v the %v %v %v.", n, a, v, o)
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
        "%v the %v %v %v.",
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