---
title: "Golang Interfaces: How different they are to Java"
date: 2023-08-08T07:07:57+02:00
showDate: true
draft: false
tags: ["programming", "golang"]
---

Originally I planned to make a new blogpost about the [visitor design pattern](https://refactoring.guru/design-patterns/visitor) and wanted to use Golang to show the concept for it. But because of how Golang is organized fundamentally you'll see that the visitor pattern doesn't make sense with Go.

Playing around with interfaces and trying to jumble together some thoughts via chatting I was able to better understand the big differences between Java and Go interfaces.

*Readers note: To not distract myself too much from studying, I didn't want to invest too much time into typo and grammar checking. So there might be some mistakes. Whoopsie daysie*

# The Basics

## Structs

Without whacking too much around the bush; Golang doesn't have classes or objects in the OOP-sense. The most similar concept to Java objects that Go has are C-like structs that allow you to pack named attributes into a packet.

A simple example would be this `Element` struct with 3 attributes:
```go
type Element struct {
    attr1 int
    attr2 string
    flags uint16
}
```
We can then create and access the struct as follows:
```go
elem := Element{42, "Shrek", 0b10010}
fmt.Println(elem.attr2) // prints 'Shrek'
```

## Methods

Say we want to now add the class method `foo()` to our Element, so that we can do `elem.foo()` and have it do some very important operations for that instantiated element, like increase `attr1` by 1. In Java, we would add the method under the scope of the class. In Go we just clap a new function basically anywhere we want to and in the first parantheses pair we designate for what type of objects this function should exist for. This requires organizing structs and methods in a different way.

```go
func (e *Element) foo() {
	e.attr1++
}
```
*Note the `*` before `Element`. It ensures that inside `foo()` we have access to the original struct. By default, without the `*`, Golang copies structs into methods and functions.

# Interfaces

Now I hope that very rough and quick introduction to structs and methods in Go was enough for you to clap it on your CV. Now we're gonna hop into interfaces.

Assuming you come from or know Java, you're probably already getting impatient and asking yourself how you can build complex trees of with subclasses and interfaces. Afterall, if your project doesn't have three interfaces, have to subclass two times just to create a single class of a Point, are you even java-ing correctly?

Go is quite different in this aspect. It does not have any subtyping or *explicit* inheritance. This should be easier with an example.

An interface is similarily constructed as a struct in Go. The main difference is that instead of defining attributes, you can only define methods in interfaces:

```go
type Shrek interface {
    roar() // implicitly void
    getOutOfMySwamp() int
}
```

So how do I now create a struct which is a subtype of this amazing `Shrek`? You don't.

> The sole purpose of an interface in Go is to let us define the methods a struct **requires** when passed into a function.

This might be a mouthful at first. Reread this sentence after the next example.

Say we create a function `compareScores(a, b)` where we want to compare the scores of two structs. Now in the end we don't really care what structs are passed into the function for `a` and `b`. Since the only method we'll call anyway is just `x.getScore()` our only requirement is that `a` and `b` have the method `getScore() int`.

Sure we could create a struct and use that as our type:
```go
type ScoreStruct struct{}
func (s ScoreStruct) getScore() int { return -1 }

func compareScores(a, b ScoreStruct) {}
```

But that would now only allow for that single `ScoreStruct` to be added and no variations. A better approach would be to create an interface:

```go
type ScoreInterface interface {
    getScore() int
}

func compareScores(a, b ScoreInterface) {}
```

The way this is right now would already allow for the `ScoreStruct` above. But in general terms, it **allows for any struct which has the method `getScore() int`**.

*Now hop up to the afor mentioned sentence and reread it. Does it make sense now?*

# Conclusion?

This way of interfacing brings a different approach to programming which is quite interesting and fun to get into. Can dearly recommend to try out Go when you got the time.
