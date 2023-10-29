---
title: "Making the Hangman Solver Faster"
date: 2021-04-16T18:45:56+02:00
showDate: true
draft: false
mathjax: true
justify: true
tags: ["python", "numba", "hangman"]
---
# The What?
A few years ago when I was starting with Python we used to sometimes play little hangman games
in class. I thought it to be a really good practice to write a little program, that could tell
me what character I should guess next. That's where I started with my hangman solver. I can put
in the word I'm looking for, the characters that are not allowed and the language (English or 
German) it should search words for. One little misconception a lot of people have, is that
words with rarely used characters like "q" and "z" might be hard for a human to figure out, but
because only a small amount of words have those letters, the program can find them very easily.

When I look back at the code now, the code is just an absolute nightmare.
Tons of duplicate code, 0 comments and no idea where to even start understanding it. On top of that
it's also really slow and inefficient. That is why I challenged myself to attempt at remaking it,
but a lot faster and more efficient. While python isn't the go-to for speed, I wanted to see what
I could achieve and as I also wanted to implement this into my python Discord bot, I decided to
write it in python so I could quickly import it over later on.

___
### Goals
Before starting, I thought of a few goals or ideas I want to try out. I had a few things in mind
that I felt would increase the speed of the hangman solver.

All the progress and every state can be seen on the commit history on the GitHub page
of the [hangman solver](https://github.com/markbeep/hangman-solver). I tried to make my commits
easily understandable.

Plan:
1. Implement a good base of the hangman solver method, which I can then easily change around functions.
The code should be readable and only require minimal edits in later changes. Also, a timer function,
which I can use to time the solver in each step.
2. First idea: To separate the giant dictionary file from one file to multiple files. In each file are
all the words of a certain length.
3. Attempt to parallelize the code, so the word list can be gone through at multiple places at the same
time.
4. Implement [Numba](https://numba.pydata.org/), a library which can make python code run as fast as C/C++
if implemented correctly.

___
# The Doing Part

### 1: Base Hangman Solver
The file for step 1: [Click here](https://github.com/markbeep/hangman-solver/blob/d2559cbb885f950b429230aea95637e1742e3e54/hangman.py)

For any project, desigining a good base will really pay off in the long run. I tried to make my methods with
the thought of later upgrading them in a quick way and decided to first split the file into three main methods.
- `get_freq_letters`: This is the main function that brings everything togethe. It returns a dictionary which
contains all the letters and the amount of words that letter is contained in. Additionally another list is
returned with all words that would currently fit the input word.
- `get_filename`: This function lets me add the input word and a language to then get the correct file for that
word. I prepared this function for step 2, as I planned to have multiple files there.
- `max_length`: Allows me to get the length of the longest word in a language. I used this for the testing
file to get a random length word.

#### Average Time Taken
**English** with 370'102 words: `0.1300s`

**German** with 1'908'815 words: `0.4664s`

___
### 1.1: Cleaning up Word Lists
My next order of business was then to clean up the word lists. I removed any words that contained punctuations,
duplicates and letters I didn't consider. Or so I thought, I later noticed that I seem to not have removed all
words with special characters. But I decided that it doesn't matter so much, as it allows for more potential words to be
recognized.

This reduced my words for English by 102 and for German by a single word.
Because of the minimal change, I didn't re-time the solver.

___
### 2: Separating the Word Lists
My second idea was to split up the one giant .txt word file I used for multiple smaller files. I chose
to split up by word length, as that is something you always know in a hangman game, no matter how many
guesses you've done.

By splitting up the words I was able to reduce the file size by a lot. The biggest file in English now
contained 53'403 words (instead of the 370'000) from before, which where the words with length nine.
In German the biggest file now contained 190'217 words. That's $\frac{1}{10}$-th of the size from before (1'908'814 words).
It seems in German most words are thirteen characters long.

This had a major effect on the time it took to run the solver. It's also the biggest change I was able to achieve
in this challenge.

#### Average Time Taken
**English** with 53'403 words: `0.0301s`

**German** with 190'217 words: `0.2128s`

___
### 3: Parallelizing the Word Search
The file for step 3: [Click here](https://github.com/markbeep/hangman-solver/blob/8bcf529f869f95683f080f9a47f3ac670cac1ccb/hangman.py)

Before I started, I had the great idea of splitting up the work for the search using the [multithreading](https://docs.python.org/3/library/threading.html)
or [multiprocessing](https://docs.python.org/3/library/multiprocessing.html) libraries from Python.
What I learned along the way though, was that the multithreading library does not actually run in parallel in Python.
It runs concurrently, meaning if you do intensive work, you barely get any benefit from multithreading.
I tried around for a long time and in the end I got only a slightly faster average time taken, which probably resulted
from me having cleaned up the code a bit, as I noticed later on.

On the other hand, the multiprocessing library in python is able to run in parallel. It has the option
to start up multiple processes all working on their own. This is exactly what I wanted. The drawback here though,
was that processes couldn't use shared memory. Meaning I had to first duplicate part of the data for each process,
then at the end add all data together again. This resulted in a giant overhead, which resulted in a longer
runtime in this case. For each additional process it took about 0.1s to get the process even started. And this
doesn't even include the time it takes to combine all the data.

I tried out a lot of things and ended up with only a slight better time with the German word list, which probably stems
from me changing around some things. In the end I decided to drop both multithreading and multiprocessing, as in
this case, it didn't result in a better outcome.

#### Average Time Taken with Multithreading
**English** with 53'403 words: `0.0309s` *(with 16 threads)*

**German** with 190'217 words: `0.1564s` *(with 8 threads)*

### 4: Using Numba
The file for step 4: [Click here](https://github.com/markbeep/hangman-solver/blob/936deda14f855208255137b868581286f4f6a51d/hangman.py)

[Numba](https://numba.pydata.org/) is a JIT compiler which can take slow python code and compile it into fast
code, approaching the speed of C/C++, if used correctly. Numba is very picky when it comes to what it can do and what
not. This means I had to change how my functions worked. No more returning a dictionary and list in the same
return statement. I decided to split my `count_chars` function into a `count_chars` and a `get_fitting` function.
Additionally, as Numba doesn't work well with dictionaries and lists, I had to change all my things into arrays.
For this I used the [Numpy](https://numpy.org/) library. The time it took to create arrays were very minimal and
I could basically ignore them. In the end I came up with a solution I was fairly happy with. Without the Numba
`@njit` decorator it worked and returned the correct results. Instead of using a dictionary which kept track of how many
times each letter was used, I now used an integer array big enough for me to simply index the unicode of a character
and increment it.

For example iterating through a word and incrementing each time a letter came up, my function looked as follows:
```python
for j in range(len(word)):
    o = ord(word[j]) - 97  # ord() returns the unicode of a characters
    if marked_letters[o] == 0:
        letter_count[o] += 1
        marked_letters[o] = 1
```
Where `marked_letters` was used to mark letters I already counted in a word, not to count the same letters multiple times
per word.

In the end Numba didn't turn into a success either. At first, ord() gave an error about how it doesn't fully work,
even though in the docs it was stated as being fully supported. Then when I got around to making it in an alternative
way, it turned out Numba had to recompile each time I ran the program, resulting in a much higher runtime than not
using Numba at all. The attempt with Numba turned into a fail for me. But this is surely something I can dig deeper on,
as Numba can be a really helpful and cool tool for making python code a lot faster.

### 5: Going Back to the Start
In the final step, I looked at everything I gathered and cleaned up the code to end with the fastest code I could
manage. Out of all the steps I tried, the only one to actually stay was the part where I divided up the word list
into multiple smaller word lists. The other ideas either made code slower or didn't make any big enough difference.

The current code is also the one viewable on the GitHub page.

#### Average Time Taken
**English** with 53'403 words: `0.0306s`

**German** with 190'217 words: `0.14996s`

### Conclusion
This was a fun little day project to try out different things I never really tried out in python. Usually when you
want to make fast code, python isn't the way to go. There are tons of other languages made for speed. But I wanted to see
what I could achieve and I was also eager to try out libraries like `threading` and `multiprocessing` in python.
Until now I only had experience using them in java. I was also keen on trying out the `Numba` library I heard a lot
about. It has a lot of potential, but I seem to not have discovered it's proper use as of now. It certainly deserves
a revisit for another project which requires more intensive computation.

For the hangman solver, I first thought that I could reduce the runtime by
a lot more, but I like how it turned out now. I'm also happy that I was able to find an elegant and clean way to create
the solver.

*If you want to use the hangman solver yourself, the only important file to download is `hangman.py`. Using its `solve()`
function you can then find possible matches.*