---
title: "Webdev without JavaScript"
date: 2023-11-29T16:53:47+02:00
showDate: true
tags: ["website", "htmx", "html"]
---

_Another article from an upcoming Visionen issue. In a future blog post I want to go more in-depth with actual code samples. Code samples are often quite hard to add to Visionen articles in a good way, since they're often too specific._

---

Webdev is a mess. Nowadays, when starting with webdev new programmers always come face to face with the JavaScript framework mess. What framework should I use? Why does everybody use React? Should I use Svelte? These questions are very common, and the answers you see online all advertise some different framework they believe to be the best. It makes it seem like you need to learn a framework to create a website in today's webdev mess. I want to introduce you to the JavaScript-less part of web development, because yes, you can make fully functioning modern websites without writing a single line of JavaScript. And unlike the JavaScript environment, this route doesn’t require you to rewrite your website every time a new framework is released.

The first step is choosing a language of your choice which will be used to create the web server. The web server will be responsible for handling all the incoming requests and sending back the correct html to the user. Every major language has one if not multiple proper ways to handle requests. This means you can create a web server with whatever language you want. The next step is to write some html templates for your pages. This allows you to serve html with dynamic data in them and also allows your website to show different states depending on the session and what user is currently logged in. There is a big variety of templating libraries. Templates are files that mostly look like html files, but they allow you to pass in values and conditional elements during runtime. This means when a user visits your website, your web server takes the template, fills in the corresponding values important to the user and then serves the whole html to the user. Most go in the direction of block syntax (like Django in Python), while others go more in the direction of components (like templ in Golang).

A big positive of using a template-based web server is that you can have your complete website logic in the same place and in a single language which isn't JavaScript. It has brought back a lot of joy when creating websites for me.

One of the common complaints with template-based web servers is that you're basically just sending static web pages to the user. Meaning, if you want a button to change something on the website, you have to load a whole new page. That is a pretty dated misconception, though. There are a lot of tools which allow you to only partially update the html of a website. There are more complete framework solutions like Phoenix for the Elixir programming language, but my current favorite is called htmx. It is a very simple and unopinionated solution for switching out parts of your website. To “install” it, you simply place the script tag on your website. This also means it works no matter what language or library you choose to go for on the web server.

htmx allows you to go the full hypermedia route; the backend only sends pure html to the user. It never sends things like JSON objects, which first have to be parsed on the frontend-side. In simple terms, this allows your website to have a button which, if clicked, sends a specific request to the backend. The backend then replies with only the required html which can then fully replace a defined element. This allows your websites to be fully dynamic without having to ever reload the page.

Template-based web servers with htmx are just a quick alternative to creating dynamic websites without having to write a single line of JavaScript. I can really recommend having a look into it when you want to create a website. It is a nice refresher from all the bloated JavaScript frameworks.
