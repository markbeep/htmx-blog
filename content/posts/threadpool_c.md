---
title: "Creating a Thread Pool in C"
date: 2022-05-24T09:09:47+02:00
showDate: true
tags: ["c","multithreading", "discord"]
justify: true
---

# Table of Contents
- [Why a thread pool in C?](#why-a-thread-pool-in-c)
- [Basic Layout](#basic-layout)
- [Implementation](#implementation)
  - [Work Node](#work-node)
  - [Thread Pool](#thread-pool)
  - [Functions](#functions)
    - [t_pool_init](#t_pool_init)
    - [thread_work_loop](#thread_work_loop)
    - [t_pool_add_work](#t_pool_add_work)
    - [t_pool_pop_work](#t_pool_pop_work)
    - [t_pool_destroy](#t_pool_destroy)
    - [t_process_count](#t_process_count)
- [Using the thread pool](#using-the-thread-pool)
- [Final words](#final-words)

# Why a thread pool in C?
For the past few weeks, I've been working on a Discord bot library written in C called [Disco-C](https://github.com/markbeep/Disco-C/). It allows for people to easily write their own Discord bots all in C. It is a nice project that allows me to use a lot of different things that I learned the past 4 semesters, but this bot is for a different blog post. Upon working on it I stumbled upon a big bottleneck: The HTTP requests I do to send or edit a message on Discord. While parsing the incoming data, turning it into structures and doing everything else that's needed takes shy of 50μs, sending the actual request takes around 300'000μs (300ms). Sometimes it drops down to 150ms, but it is clear that this is a major bottleneck of the library. The problem is that if I now run my event loop sequentially any HTTP requests I execute stall up all other commands. The solution? Multithreading? Well yes, but not exactly. It allows us to process multiple received messages at the same time. But simply creating a new thread on each event causes a lot of overhead for the creation and deletion of the threads. This is where a thread pool comes into play. In a thread pool, we create a defined amount of worker threads upon initialization and keep using these same threads for all the work we send to the pool. This saves us a lot of overhead from the thread creation.

This is why I decided to use a thread pool myself and as with a lot of things in C, if there's no big official library for it and you don't want to use some random shady Github code, you have to implement the library yourself. I never really thought about how thread pools work and at first thought, this will be a challenging task, but I was relieved to find out that it was a lot easier than I expected.

This blog post is a bit on the longer side, but it goes into detail on how all the functions work in the thread pool implementation, so you can implement it yourself if you wanted to. The inspiration for this kind of thread pool implementation I got from [this blog post](https://nachtimwald.com/2019/04/12/thread-pool-in-c/). They also go into a bit more in detail and have another blog post where they show how to make this implementation also work on Windows.

# Basic Layout
The way our thread pool will work is by having a queue that contains all the work that has to be done. A piece of work is defined as a function and the data that will be passed into that function.

The idea is that every time we have some new work to be done, we add the data and the corresponding function into our queue. The threads of our pool then fight a free-for-all to get the work and execute it. The work threads each constantly check for new work and once they're done with their current work, they instantly continue on new work if there's any. To avoid excessive CPU usage the threads will also be put to sleep and woken up when needed.

# Implementation
The following chapter is all about how the thread pool is exactly implemented in case you also want to try implementing it yourself.

The full implementations can be found here if you ever want a full overview: [Header](https://github.com/markbeep/Disco-C/blob/master/libs/utils/t_pool.h) & [Source](https://github.com/markbeep/Disco-C/blob/master/libs/utils/t_pool.c)

## Work Node
For the implementation we start in our header file: `t_pool.h` with the layout of a work node:
```C
typedef void (*t_func)(void *);

typedef struct t_work {
    void *arg;
    t_func func;
    t_work_t *next;
} t_work_t;
```
Our work node structure consists of the data (`void *arg`), the function and the next node in the queue. `arg` is what will be passed into the function. `t_func`is a function that takes in a pointer to the data (`void *`), but returns nothing (`void`). This means we can actually pass in any type of data type as long as the corresponding function correctly casts the data.

## Thread Pool
Now we come to a bigger structure, the actual thread pool structure:
```C
typedef struct t_pool {
    // work queue
    t_work_t *first_work;
    t_work_t *last_work;
    int work_count;   // amount of active work load
    // thread handling
    int thread_count; // amount of active threads
    int stop;         // bool if threads should stop
    pthread_mutex_t *lock;
    pthread_cond_t *work_cond;
    pthread_cond_t *finished_cond;
} t_pool_t;
```
The pool struct can be split up into two parts. The first part is simply a queue with all of the work nodes. To make inserts in O(1) we need the head and tail of the queue, which is what we have the attributes `first_work` and `last_work` for. `work_count` simply states how many work nodes we currently have in the queue. In the end, we'll want threads to wake up, check if there's any work using this variable and if there's none, go back to sleep.

The second part is all about managing the threads. When looking at the struct you might notice that there's actually no array with the threads. In fact, I don't store any pointer to the threads. That is because I handle all the threads differently.

We first have the `thread_count` to signal how many threads are currently running. This allows us to later wait for all threads to close and keep a count of how many threads we are still waiting for. The next attribute is the boolean `stop` flag. If we want to close all threads, we can simply set this to true and it makes the threads close themselves. To make the whole pool thread-safe we have a single `lock` that we use at certain places to avoid undefined behavior like two threads writing to the same space in memory at the same time. In the end, we have two conditions `work_cond` and `finished_cond`. These are variables used to signal threads to wake up when they're sleeping. For people that have experience with multithreading in Java, this is what's used to later achieve the `notify` and `wait` methods in Java. With the pthread library, we have the advantage of not just waking up threads sleeping on a lock, but instead only waking up threads with the correct condition variable.

## Functions
We now define all the public functions that should exist for the thread pool.
```C
t_pool_t *t_pool_init(int num_t);
int t_pool_add_work(t_pool_t *tp, t_func func, void *work);
t_work_t *t_pool_pop_work(t_pool_t *tp);
void t_pool_wait(t_pool_t *tp);
void t_pool_destroy(t_pool_t *tp);
int t_process_count();
```
- `t_pool_init`: This function is to initialize the thread pool with `num_t` threads. It should return a correct setup `t_pool` struct.
- `t_pool_add_work`: Here we pass in our thread pool, a function and the data we want to be executed by the function. This function will then create a `t_work` node and add it to the queue. All in a thread-safe manner.
- `t_pool_pop_work`: This should correctly pop the first element of the queue (the next work to be done) if there is any. A thread wanting to do work will then call this function to get new work.
- `t_pool_wait`: This function we use simply to wait for all threads to close. If we don't specify the states to close beforehand, this function will wait forever.
- `t_pool_destroy`: Function to clean up all the allocated memory of the thread pool.
- `t_process_count`: Simple helper function to get the number of processes on the current machine. This can then be used to decide how many threads should be created.

These are all the functions that are required to create a working thread pool. Now we'll go into more detail on how all of these functions actually work.

### t_pool_init
```C
t_pool_t *t_pool_init(int num_t) {
    t_pool_t *pool = malloc(sizeof(struct t_pool));
    /* initialization left out for simplicity */ 
    pthread_t thread;
    for (int i = 0; i < num_t; i++) {
        pthread_create(&thread, NULL, &thread_work_loop, (void *)pool);
        pthread_detach(thread);
    }
    return pool;
}
```
In this code snippet, I let out a big part which was just the struct initialization. The whole function can be viewed [here](https://github.com/markbeep/Disco-C/blob/master/libs/utils/t_pool.c#L31-L55).

The most important part of this function is the thread-creation. We create `num_t` amount of threads all with the `thread_work_loop` function (more explained later). This is a function that makes the threads simply loop forever and check for new work until told to stop. We then detach each of the threads, because we'll not be joining them later on and this frees the memory that gets allocated for it.

### thread_work_loop
```C
static void *thread_work_loop(void *tp) {
    t_pool_t *pool = (t_pool_t *)tp;
    while (1) {
        pthread_mutex_lock(pool->lock);
        while (pool->work_count == 0 && !pool->stop) {
            pthread_cond_wait(pool->work_cond, pool->lock);
        }
        if (pool->stop) {
            pool->thread_count--;
            pthread_cond_signal(pool->finished_cond);
            pthread_mutex_unlock(pool->lock);
            break;
        }
        // here work_count > 0
        t_work_t *work = t_pool_pop_work(pool);

        // we can now unlock as the work afterwards is thread-free
        pthread_mutex_unlock(pool->lock);
        work->func(work->arg);
        free(work);
    }
    return NULL;
}
```
This is the main loop that every thread will keep executing until told to stop. In the beginning, we wait for the thread to get new work by using the `pthread_cond_wait` method inside the while loop. This way the thread sleeps and doesn't use excessive CPU while it randomly runs in the background without work. If we ever want to wake up the thread we can use the `pthread_cond_signal` (Java: `notify`) or `pthread_cond_broadcast` (Java: `notifyAll`) methods. We use the `work_cond` attribute to signal that new work has been added to the queue.

When the stop flag is marked we signal the `finished_cond` to mark that a thread has finished. We listen for this condition when we wait for all threads to finish.
Afterward, we simply take the next piece of workload and make sure to execute the work outside of the lock.

### t_pool_add_work
This is the function we use to add work to the queue.
```C
int t_pool_add_work(t_pool_t *tp, t_func func, void *arg) {
    // creates the new work node
    t_work_t *work = (t_work_t *)malloc(sizeof(struct t_work));
    work->arg = arg;
    work->func = func;
    work->next = NULL;

    // coarsely locks the queue
    pthread_mutex_lock(tp->lock);
    t_work_t *last = tp->last_work;
    if (last)
        last->next = work;
    else
        tp->first_work = work;
    tp->last_work = work;
    tp->work_count++;
    pthread_cond_signal(tp->work_cond);
    pthread_mutex_unlock(tp->lock);

    return 1;
}
```
In this function, we first create the work structure which consists of the passed-in function pointer and a pointer to the data.
Once we have the work node, we enqueue it at the end of our pool queue and then signal one of the waiting threads that are sleeping, that there is new work to be done. I used a coarse lock here because it makes the implementation very simple and safe.

### t_pool_pop_work
With this function, we pop off the first element of the queue and return it.
Because we call this method in the `thread_work_loop` while we have a lock already, we don't need to lock in this method again. If this method was to be used elsewhere as well, I'd add locking to avoid any problems.
```C
t_work_t *t_pool_pop_work(t_pool_t *tp) {
    if (tp->work_count == 0)
        return NULL;
    t_work_t *head = tp->first_work;
    if (head->next) // there are still elements in the queue
        tp->first_work = head->next;
    else { // we removed the last element in the queue
        tp->first_work = NULL;
        tp->last_work = NULL;
    }
    tp->work_count--;
    return head;
}
```
Other than that the function is just a dequeue function like you would implement it in other singly-linked lists.

### t_pool_wait
```C
void t_pool_wait(t_pool_t *tp) {
    pthread_mutex_lock(tp->lock);
    while (tp->thread_count > 0)
        pthread_cond_wait(tp->finished_cond, tp->lock);
    pthread_mutex_unlock(tp->lock);
}
```
In this function, we simply wait for all threads to end. Because we don't keep track of all the threads, we can't just join on each thread. We instead do it differently and look at how many threads are still active. With how we created our `thread_work_loop` function, every thread that stops correctly reduces the `thread_count` and then signals the `finished_cond`. This means that if all threads end correctly the while loop is broken and we can finish waiting. This method will also wait forever if the threads are never told to stop.

This implementation has two flaws though. The first is that when a thread goes on holiday (for example it crashes) this function will wait forever. A way to circumvent this would be to use `pthread_cond_timedwait` instead and break out of the while loop once a certain limit is reached. If the time limit is too short, it can result in ignoring threads that are still alive but simply working.

The second is the fact that threads take some time to properly close. And threads that are still open when the main thread exits are considered a memory leak. This means that if we wait using this function and then instantly exit when this function ends, there's a good possibility that the worker threads didn't close down yet. If this is a big problem for you, usually waiting for like 1ms would fix it. The clean solution would be to join on the threads, but that means you have to store the thread pointers somewhere.

### t_pool_destroy
```C
void t_pool_destroy(t_pool_t *tp) {
    pthread_mutex_lock(tp->lock);
    tp->stop = 1;
    // delete all nodes still in the queue
    t_work_t *cur = tp->first_work, *prev = NULL;
    while (cur) {
        prev = cur;
        cur = cur->next;
        if (prev)
            free(prev);
    }
    tp->work_count = 0;
    // wakes up all sleeping threads
    pthread_cond_broadcast(tp->work_cond);
    pthread_mutex_unlock(tp->lock);
    // waits for all threads to finish
    t_pool_wait(tp);
    /* Freeing up all the attributes */
}
```
This is our method to clean up all the memory we used, so we don't end up with tons of memory leaks in the end. The basic idea is that we first set `stop` to 1 which makes all the threads stop as soon as they can. We then clear up the whole queue and free up any work nodes that are still in there. In the end, we use `pthread_cond_broadcast` (Java: `notifyAll`) to wake up all the threads with the `work_cond` variable. It's like we're waking up all threads telling them there's more work to be done, but actually, we just woke them up to dump them. Once we waited for all threads to finish, we can clean up the rest of the attributes.

The freeing part can be seen [here](https://github.com/markbeep/Disco-C/blob/master/libs/utils/t_pool.c#L103-L109).

### t_process_count
Lastly, we simply have a quick helper method to get the number of processes we have available. This allows us to create the same amount of threads as the processes we have.
```C
int t_process_count() {
    return (int)sysconf(_SC_NPROCESSORS_ONLN);
}
```

# Using the thread pool
Now to how we can properly use the thread pool. A sample use case would look like so. First, we define the function that the worker thread should execute. We'll keep it simple and simply have the function take a string and print it to the console.
```C
void print_string(void* data) {
    printf("%s\n", (char*)data);
    free(data);
}
```
We need to free the allocated data inside this function because we won't know on the outside when this thread is done with the work.

```C
t_pool_t* pool = t_pool_init(t_process_count());
/* we need to place the data on the heap, otherwise
   it could be cleared up before the other thread
   uses the data */
char* data = (char*)malloc(12);
strcpy(work, "Hello World");
t_pool_add_work(pool, &print_string, data);

// cleanup
t_pool_destroy(pool);
```
This little implementation is already all we need. It adds the work to the pool, a thread will execute it right when we do `t_pool_add_work` and then with `t_pool_destroy` we wait for the message to be printed and then we clean up.

# Final words
I found this thread pool implementation a nice addition to my Discord library. I always wondered how a thread pool worked, but never really looked into the details. There are still a lot of additions that can be made which will make it faster and safer, but for now, this barebones version will work for me. I will probably add a timed wait in my wait function, but I still have to see how exactly I'll implement this so it works the way I need it to work.

If you got until here, awesome and thanks for reading! Hope I motivated you to also try looking into a function you've been using at the top level, but never really checked how it works behind the hood.
