---
title: "Fusing Go Into Python"
date: 2023-09-13T22:58:05+02:00
---

_For the upcoming [Visionen](https://vis.ethz.ch/en/visionen/general) release I've written yet another techy article. Here's the article so you can get a little sneak peak and hopefully interest you into reading the Visionen once it's out :)_

If you like coding in Python, but sometimes just want that extra speed for some functions, I have the perfect solution for you: We're going to build our own little library that can run compiled code inside Python. But instead of C, we're going to be using Golang, which is almost as easy to write as Python, but actually fast.

### Setting the baselines

For the Python test function we'll be using to benchmark I approached the [Collatz conjecture](https://en.wikipedia.org/wiki/Collatz_conjecture) problem. Maybe we can find a proof along the way.

```python
def collatz(n: int):
    if n % 2 == 0:
        return n // 2
    return 3 * n + 1
```

Iterating over the Collatz function one billion times took around 1m10s on my machine. Now lets look at the pure Golang approach:

```go
func collatz(n int) int {
    if n%2 == 0 {
        return n / 2
    }
    return 3*n + 1
}
```

Using Go I was able to execute one billion iterations in only 600ms. That's quite a big difference. Lets see how close we can get to that performance with bindings inside Python.

### Creating the bindings

Creating the bindings is actually a lot simpler than one might think at first. First we'll look at the Go side.

#### Go-Side

An important library we're using is called [cgo](https://pkg.go.dev/cmd/cgo) which enables us to write C inside Golang. Fret not, we're not really going to be writing C code, but you could theoretically do some really cursed stuff if you wanted. By adding the comment `//export colatz` (no space after //) above our function in Go it will later be accessible in Python:

Our whole `main.go` file is as follows:

```go
package main
import "C"
//export collatz
func collatz(n int) int {
    if n%2 == 0 {
        return n / 2
    }
    return 3*n + 1
}
func main() {}
```

Now comes an important step. We don't want to simply compile this to an executable, we want to compile this to a C shared object file by running the following:

```bash
go build -buildmode=c-shared -o library.so
```

#### Python-Side

Back in our Python program we now want to import the generated `library.so` file so that we can call the Golang code. Here the Python `ctypes` library will play an important role:

```python
import ctypes
library = ctypes.cdll.LoadLibrary("./library.so")
go_collatz = library.collatz

def fast_collatz(n: int):
    go_collatz.argtypes = [ctypes.c_int] # First and only argument is an int
    go_collatz.restype = ctypes.c_int # Return value is also an int
    return go_collatz(n)
```

And we're done! We can now import and call `fast_collatz` wherever we want.

#### Houston we have a problem!

But when we now want to run this cool new function a billion times we're suddenly greeted with the opposite of what we expected. It suddenly takes almost 20 minutes to execute. There is a huge overhead when accessing Python bindings. The solution for this is to simply throw more computation onto the Go side and have Python execute as little as possible. We can fix this by creating new functions which repeat Collatz for a given amount of iterations.

Golang:

```go
//export loopCollatz
func loopCollatz(n, k int) int {
    for i := 0; i < k; i++ {
        n = collatz(n)
    }
    return n
}
```

Python:

```python
go_loop_collatz = library.loopCollatz
def loop_collatz(n: int, k: int):
    go_loop_collatz.argtypes = [ctypes.c_int, ctypes.c_int]
    go_loop_collatz.restype = ctypes.c_int
    return go_loop_collatz(n, k)
```

Executing this new `loop_collatz` function we get our expected execution time of 630ms. Only 30ms more than running it purely in Go, but a more than 100x speedup compared to our pure Python execution. The code is available in [this repository](https://github.com/markbeep/Golang-Python-Bindings), to easily see it as a whole.
