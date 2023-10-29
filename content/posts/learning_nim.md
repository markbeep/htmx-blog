---
title: "Learning the Nimble Nim"
date: 2021-05-09T14:38:03+02:00
showDate: true
draft: false
justify: true
mathjax: true
tags: ["programming","nim"]
---

# Table of Contents
- [What is Nim?](#what-is-nim)
- [Why?](#why)
- [The Surprises](#the-surprises)
    - [Variables](#variables)
    - [Switch Cases](#switch-cases)
    - [For Loops](#for-loops)
    - [Procedures](#procedures)
    - [Arrays](#arrays)
    - [Objects](#objects)
- [Conclusion](#conclusion)

# What is Nim?
Nim does a good job explaining what it is in a short way on their [website](https://nim-lang.org/):
> Nim is a statically typed compiled systems programming language. It combines successful concepts from mature languages like Python, Ada and Modula.

But that quote only covers a tiny portion of what Nim really is. In my journey of learning Nim I plan to post about my progress on this blog. Maybe I can even awaken your interest for a new language you've never heard of before.

___
# Why?
I dunno. I didn't really have a reason to learn a new language. I learned about Nim when [Julian](https://xyquadrat.ch/) recommended some cool new low-level languages including [Zig](https://ziglang.org/), [Crystal](https://crystal-lang.org/) and Nim. I looked into each and Nim syntax was the one that appealed the most to me. Python was my previous most used language, so I was familiar with indentation in programming. Appealing points of Nim are it's fairly easy to understand syntax and it being a fairly [fast language](https://github.com/kostya/benchmarks).

I decided to start learning by following the [tutorial](https://nim-lang.org/docs/tut1.html) on the official Nim docs. I additionally started working on a Discord bot so I could test out different things I learned while handling user input and reading into Nim libraries written by other users.

___
# The Surprises
Now to the meat of the post. These are the basics of Nim which all get covered in the first part of the Nim tutorial and are used for any project you do with Nim. Before starting to learn Nim I only knew Python and Java very well, so some of these things might be surprising for me, but normal for you. Please do tell me about anything you find similar to another language.

___
### Variables
Starting out with the basics of Nim: Nim is a statically typed language. Meaning variables have a fixed type set to them upon initialization. This can be done in multiple ways:
```nim
var x: int  # by default a 64-bit int
var s: string  # string with fixed length
var b: uint8  # unsigned 8-bit int
```
Note that `#` are used for comments just like in Python.
But Nim allows us to leave out a lot of things to cut down on code and make it more readable.
```nim
var
    x, y, z: int8
    s = "this is a string"
    f = 0.0'f32  # 32-bit float
```
Additionally it's even possible to execute code when assigning variables. For example in the following code snippet we assign the factorial of 4 to the variable `fact`.
```nim
var fact = (var x = 1; for i in 1..4: x *= i; x)
```
This can also be written as
```nim
var fact = (
    var x = 1
    for i in 1..4:
        x *= i
    x
)
```

Nim has different ways to initialize variables. These include `var`,`let` and `const`. `var` is for assigning variables which can be re-assigned at a later point. `let` is similar to `var`, but is only for single assignment variables. Meaning the values can't be changed once initialized. `const` variables are initialized during compiling, meaning they need to be computable at compile time. These can then not be changed at runtime.

___
### Switch Cases
Switch cases might come as a familiar concept, as they exist in most other language ([Including Python 3.10](https://towardsdatascience.com/switch-case-statements-are-coming-to-python-d0caf7b2bfd3)). Usually each case is for a single case, which makes it often tedious if you need to cover a lot different values for each case. Nim takes this a step further and allows you to define a range of values that hold for a case.
```nim
var x = 9
case x:
of 1..2: echo "1 <= x <= 2"
of 3..5: echo "3 <= x <= 5"
of 6..<10: echo "6 <= x < 10"
else: echo "neither"
# --> 6 <= x < 10
```
Something to take note of, this also works for arrays of any type, as long as that array is evaluated at compile time. Meaning we have to use `const` for these arrays.

___
### For Loops
For anybody that has used Python, for loops in Nim will come as very familiar. They are structured in the same way.
```nim
for i in 1..<10:
    echo i
# --> 1 2 3 4 5 6 7 8 9 all on separate lines
```
This prints all numbers from 1 to 10 exclusive. We can also iterate over iterables like arrays in the same manner. Optionally, it's possible to define two variables in the for loop. The first one will then be assigned the index and the second the value. 
```nim
var arr = [5, 55, 555, 5555]
for idx, val in arr:
    echo (idx, val)
# --> (0, 5)
# --> (1, 55)
# --> (2, 555)
# --> (3, 5555)
```
___
### Procedures
Procedures are Nim's way of calling their methods or functions. This is the part that really starts to make Nim shine. 
Let's look at an example:
```nim
proc combine(a, b: string): string =
    # the & operator concats two strings
    a & b
echo combine("Hello ", "there")
# --> General Kenobi (jk it returns the line below)
# --> Hello there
```
This is already a full procedure. To explain it in a bit more detail. This is a procedure called `foo`, which takes the two string parameters `a` and `b`. The `string =` part indicates that this procedure returns a value of type string. The parameters of a procedure require a defined type. Just like in Java this then allows for overloading of methods; having multiple methods with the same name, but different parameters. We can use this to our advantage if we wanted to create a similar `combine` function, which combines integers instead.
```nim
proc combine(a, b: int): string =
    # the $ operator turns an int into a string
    return $a & $b
echo combine(4, 20)
# --> 420
```
Something I found extremely interesting, is that just like a lot of things in Nim, the `return` statement is optional. In the above example, it was very well possible to leave out the return statement, I only put it there for demonstration purposes. If there is no return statement given **OR** the `result` variable is not assigned to anything, the last line will be taken as the return value.
We can rewrite the above example in the following ways which all result in the same result:
```nim
proc combine(a, b: int): string =
    $a & $b
proc combine(a, b: int): string =
    result = $a & $b
```
But that only scrapes the start of what can be done with procedures in Nim. Just like in Python, we can set default values for parameters. This is something that is a big annoyance in Java, as there you need to create a new method for each default parameter.

It's also possible to name what variable you want to assign to what exactly, which is shown with `b` and `d` in the below example.
```nim
proc combine(a=1, b=2, c=3, d=4): string =
    $a & $b & $c & $d
combine(d=69, b=666)
# --> 1666369
```
Now see how I used the `$` operator to turn integers into strings. Why can't we just do `a & b & c & d` for integers? That's because it's not defined yet. But we can define it ourselves. After all, Nim has operator overloading. We can make our custom operators in the following way.
```nim
proc `&`(a, b: int): string =
    $a & $b
echo (1 & 2) & (3 & 4)
# --> 1234
```
Operators can be called in multiple different ways. If you for example wanted to inflict pain on the person reading your code, you could do the following:
```nim
echo `==`(`+`(5, 11), `*`(2, 8))
# --> true
# This is the same as simply
echo 5 + 11 == 2 * 8
# --> true
```
This then also leads into the fact that normal procedures can be called in a of different ways. Let's take the following procedure which simply puts two square brackets on each side of the given word.
```nim
import strformat  # this allows us to format strings
proc wrap(w: string): string =
    fmt"[[{w}]]"
```
We can now call this in different way, all resulting in the same result.
```nim
echo wrap("kek")
echo wrap "kek"
echo "kek".wrap
# --> [[kek]]
```
Multiple ways of doing something, all resulting in the same outcome, is a recurring theme in Nim. While this can be a major turn off for a lot of people,
that is something which simply needs getting used to and and it leads to less problems than one might think.


___
### Arrays
Just like in pretty much every other language, Nim also has arrays. Arrays in Nim have fixed size which needs to be set upon initialization. This can be done in multiple ways (a trending pattern in Nim).
One of the ways is to create a new type which can then be used again and again afterwards.
```nim
type
    fiverArray = array[5, int]  # 5 long array type
var a: fiverArray  # [0, 0, 0, 0, 0]
var b: fiverArray  # new array also with [0, 0, 0, 0, 0]
```
We can also just define the array type upon initialization or already fill the array with the elements we want:
```nim
var c: array[5, int]
var d = [1, 2, 3, 4, 5]
```
These arrays are then just like arrays in most other languages. Index 0 is the first element and the last element has index $n-1$.
But what if you are one of those people that like to start index on 1? 500? Or maybe even -2147483648? Don't you worry, Nim has got your back.
It allows you to index at any value you want and end at any index you want (in the 32-bit signed int range). We can accomplish this amazing feat in the following
manner:
```nim
var matlab: array[1..100, int]  # an array with index 1 <= i <= 100 for the matlab users
var wtf: array[low(int32).int..0, string]  # a string array starting at -2147483648 up to 0
```
This is something I have yet to find a use of, but if you have any ideas, don't hesitate to comment your ideas below.

Now arrays are all fine and dandy, but no language is complete without a linked list implementation. Sometimes you just don't know how many elements you
need in the beginning and creating a giant array each time is just horrible and inefficient.
Nim's implementation of lists is called *sequences* and they can be created very simply. Sequences always start index at 0. Meaning we lose our weirdly indexed arrays if we turn them into sequences :(

The `@` operator can be used to turn arrays into sequences:
```nim
var s: seq[int]  # creates an empty list for ints
var q: @[1, 2, 3, 4]  # creates sequence out of the array
var wtfPlus: @wtf  # makes our giant array into a sequenec
```
Elements can be added to sequences using the `.add()` procedure:
```nim
var s: seq[string]
for i in 1..<10:
    s.add($i & ".")
echo s
# --> @["1.", "2.", "3.", "4.", "5.", "6.", "7.", "8.", "9."]
```
One thing you might have noticed in the above code is that yes, lists and arrays when printed actually show the contents of the list, instead of the object reference or whatever other languages
usually do.

___
### Objects
The last thing I'll cover in this post is objects. Nim is not really an object-oriented programming language, but it still has a lot of powerful OOP techniques.
In part 2 of the Nim tutorial a big part of the focus lies on objects. In this post I will only cover the basics that get covered in the first part of the tutorial.

We define a new `Person` object with the properties `name` and `birthyear` in the following manner:
```nim
type
    Person = object
        name: string
        birthyear: int
```
We can now create new Person objects with a name and a birthyear similar to how it's done in other languages. In the below example we create a Person object for *pepe*:
```nim
pepe = Person(name: "pepe", birthyear: 2005)
```
Adding methods for an object is a bit different than in other languages. The Person object is not a class, like it's called in other
languages. It's simply an object with properties.

Fitting to the theme of Nim, we create procedures for certain objects in the following manner. Say we want to create a method `eat()` so a Person can eat. This procedure is created just like any other procedure.
```nim
proc eat(p: Person): void =
    echo p.name & " is eating."
```
Because the procedure requires a Person object as a parameter, it can only be called with a Person object. Meaning it handles just like if this method was tied to the Person object. The `eat` procedure can now be called in the following ways:
```nim
pepe.eat()
eat pepe
eat(pepe)
# All three result in --> pepe is eating
```
Depending on how English grammar-like you want your code to be, one might be more fitting than the other.

___
# Conclusion
In this post I tried to show some things of Nim, which have surprised me or were new. One of the "issues" which could become a problem
when using Nim in a team or when reading other people's code is the amount of ways you can do certain things. It can quickly become confusing if code is not written clearly.
There are many ways to write Nim so it's unreadable and deserves an award in a code obfuscation tournament, but there are also tons of ways to make code be easy to read and quick to understand.
In group projects it's important to agree on a standard way of syntax and to stay consistent throughout the project.

Another big plus is how much fun Nim can be to write. It has easy to learn syntax and the error handling is amazing, giving you a good idea on why there's an error for whatever
shenanigans you're trying to throw together.

If this post has in any way intrigued you for Nim, don't hesitate to check out the [tutorial](https://nim-lang.org/docs/tut1.html), it explains the basics very well and you'll be up and
coding in a new programming language in no time.

I'm also very eager to hear your opinions on Nim and this post in the comments. If you find any mistakes or I messed up somewhere, don't be afraid to hit me up.