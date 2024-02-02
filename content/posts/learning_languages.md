---
title: "You should learn a new language"
date: 2024-02-02T10:34:00+01:00
showDate: true
tags: []
---

To some people's surprise, when I started at ETH I only knew a bare level of Python. Am always surprised to learn about the fifty programming languages first-years nowadays bring to the table when starting at ETH. I've since learned to appreciate and learn a lot of different languages and dare to say learning new lesser-used languages is even my hobby.

Every language I know now I learned during my past 3.5 years at ETH. 95% of it, if not more, I learned outside the usual study schedule.
I want to introduce you to why you should learn new languages and how you can do it effectively to not forget it all a week later.

# Why you should learn a new language

You'll often hear that instead of learning ten programming languages very roughly, you should stick to one language to learn it more in-depth. I fully agree with this point. You shouldn't try to quickly force in ten different programming languages just so you can make your CV look good at a rough glance. Take every language you learn one at a time so you can dive into it and fully learn it.
That doesn't mean you should only stick to the one language you love. Be open to learning and checking out a new language! Just don't rush the learning process too much.

The main strength I see behind learning a new language is that every language has its unique way of thinking and approaching problems. This helps extremely with widening your horizon. Different languages have different ways of how projects are structured. Learning new languages makes you think differently about projects and languages you work with.
New ways of thinking from one language can most often be transferred to other projects and languages. It helps understand different code bases better and with every language you properly learn, any additional language becomes even easier to learn because there are always some similarities (just like with spoken languages).

Some examples from my experience include learning Haskell, which opened up a new world of functional programming. I've started to think more about functional programming when writing in any other language and always making the conscious decision if for example a global variable is actually required or if I can approach the problem in a purely functional way.

Alternatively, learning Go and [Zig](https://ziglang.org) opened me up to the concept of returning errors as values. Previously with Python and Java I always felt like errors had to be raised and caught in a bulky try/catch syntax. Go and Zig both introduced returning errors as values from functions (Sidenote: C also theoretically returns errors as values, but there it is not an integral part of the language). That is a new way of thinking which I've lately been enjoying a lot. With Go and Zig, it is also interesting to see how much they differ and how much is similar.

# How to learn a new language

There is no one true way to learn a language and every language is different, which is why I'm not going to go fully in detail on how you should learn a language step by step.

There are two key steps I recommend for anyone learning a new language to follow:

1. Learn the basics of a language so you get the hang of the syntax. I've tried multiple methods for learning a language. I'd recommend you try the different methods to see what you enjoy the most.

   - For Zig I used [ziglings](https://codeberg.org/ziglings/exercises/) (also exists for Rust: [rustlings](https://github.com/rust-lang/rustlings/)).
     This is a new way of learning a language for me which I tried out the first time for Zig. In this method, you get a lot of exercises with some code already prewritten for you. The code is buggy or broken and you have to fix it. Then there's always an explanation of the concept that is introduced or what is going on in the program. This helps a lot with becoming good at reading code you didn't write as well as learning to understand the compiler errors. This was one of the most effective ways to learn a language for me. I also didn't have to set up a project to create some practice files and could dive straight into using the language while learning about it.
   - For Go I followed the book "Get Programming with Go" by Nathan Youngman & Roger Peppé. This was my first time learning a language using a book and this book was perfect. It introduces the language slowly with humor and slowly turns you into a proper Gopher by teaching you the right etiquette.
   - For my earlier languages like Python, I learned using an app. Sadly the app I used to use has gone completely haywire with in-app purchases and I wouldn't recommend it anymore.

   An important point is that no matter how you learn a language, it is important that you not only learn the language theoretically but also apply the knowledge by playing around with the language and coding in it.

2. The next step to learning a language is to use the language in a project. [Austin's Blog Post](https://austinhenley.com/blog/challengingprojects.html) goes over some cool in-depth project ideas if you're looking for some. But a project can be of any scale and simply anything that interests you. Here also comes a point I find extremely important. **Don't try to use a single language for everything.** See a language as a tool and use the right language for the right job. If your goal is to create websites, you'll need to learn HTML and JS. By creating webassembly websites using Rust you're just shooting yourself in the foot. Or if you want to learn a systems language, find a project which allows you to focus on system-level aspects. Creating a discord bot with C is a nice project, but almost a bit too high level for C and results in a lot of pain. By using a language for something it isn't intended for, you learn more about the language, but it isn't necessarily the most effective way since you're trying to force the language to your use case.

   For example, a good project for higher-level languages and if you're a Discord addict, is to create a Discord bot. Gives you a lot of insights into handling user input while allowing you to dive more in detail for any feature you want your bot to have (image processing, web scraping, high-performance computing, etc.).

# Conclusion

I learn languages not to fill my CV, but to get more insights into how different languages approach problems. I find learning about the different approaches quite enjoyable and can dearly recommend you to also have a look at new languages.
You should also not rush learning a new language. Take your time with new languages and create projects to properly learn the language (projects also show you have proper experience with a language on your CV).

I'm curious how you approach learning a new language. Or do you avoid learning new languages? Why?

---

_Sidenote: I tried to make a Zig propaganda post, but somehow resulted in writing a more general language post. Don't fret though, the Zig propaganda will come in a future post._
