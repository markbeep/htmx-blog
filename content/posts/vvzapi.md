---
title: "Creating VVZAPI"
description: "Some motivation behind why and how I created an updated course catalog API for ETH Zurich."
date: 2026-04-28T15:06:13+02:00
showDate: true
tags: ["programming", "website"]
---

_This is a article I wrote for [Visionen](https://members.vis.ethz.ch/committees/visionen) and thought I might as well also share on my blog._

---

Every semester is the same story. You want to search for some courses to plan out your semester, you open the Course Catalog (VVZ), and then sluggishly jump through page after page until you finally find a few potential candidates for your next semester. Every single click makes you tediously wait. And personally, it also pains me to see a website being so slow in this day and age. Though, it's still somewhat impressive for VVZ to still be running with minimal (frontend) changes since 25 years.

Students are not satisfied with VVZ. Want proof? In the past years, almost every VIScon hackathon has had a task related to creating an improved VVZ or course planning tool. Seemingly, I keep coming back to VVZ, even having worked on a course planning tool in two separate VIScon hackathons. Creating such a tool presents two major hurdles. On one hand, it is a large undertaking to create a tool that can correctly keep track of every requirement for all the different study programmes. But even if you choose to forgo that aspect, there is still the looming question of how to even get the data for all the courses. You might only be interested in the course name, number, and credits, but maybe you also want access to the course description and exam information.

But how can you get the actual course data? VVZ has an internal API, but it requires authentication, so it is a no-go if you just want to quickly play around with data to see if your tool even makes sense. So there is only one other option: You will have to resort to scraping the raw HTML pages of VVZ and then tediously extracting all the relevant data. While not the most difficult task imaginable, it does impose a challenge, especially in the likely case that you don't know the full structure of VVZ and what all the values can be or how everything gets displayed. So even with an AI agent speeding up the work of writing all the scraping logic, you will inevitably stumble upon something new that you didn't expect or handle. Luckily, you can assume that the actual structure of the HTML will probably not change anymore, allowing for any scraping tool to be quite heavily hard-coded.

And don't even try to make any assumptions about the structure of VVZ. Anything that you believe suuuurely doesn't exist, actually _does_ exist. Did you know there can be two courses with the same number (000-0000-00L type number) in the same semester, or that credits don't have to be integers? Credits can theoretically have up to four decimal places. Luckily, the most I've stumbled upon only uses one decimal place. For example, "151-0055-10L  Ingenieur-Tool: Planung menschlicher Arbeit" in the 2024 Winter semester gives you 0.4 credits. You can take some 2.6 credit course to even it out.

There's no reason why it should be so hard to build tools on top of VVZ. I fall in with the open-source and unix ethos of creating modular tools that can easily work together. For most applications, it becomes exponentially hard to account for all features anybody would ever want. It therefore often makes sense to instead allow for easy interoperability. If people want some specific niche features, they can build their own third-party tools. People who don't want those features don't have to wade through and explicitly disable everything. In the simplest form, that usually means creating a REST API that can allow for third-parties to, on one hand, access the data, but on the other, also perform operations if required.

Before, when you wanted to use course catalog data, you might have been able to find some random script on GitHub somebody cooked up to quickly scrape a few things from VVZ. The problem is that basically all of these tools all only require a subset of the catalog data, and hence if you need something they don't scrape, you had to figure out and rewrite a large part of the code. You were usually better off just quickly writing your own scraper to fit your needs and tech-stack instead.

To aid in that, in the past months, I worked on creating an API that makes the full VVZ data available over a simple REST API. My two main goals were to first replicate **all** the data available, including the past semesters and all the information I doubt a lot even notice (have you ever looked at the "Competencies" section?). Most of the data will probably never be used, but it's all there if ever needed. And second, to make the API easy-to-use while allowing for the same filtering capabilities (and more) as VVZ.

### The API

The website can be found under [vvzapi.ch](https://vvzapi.ch). It initially started out as a bare-bones API-only website, but I've recently added a frontend that allows you to unleash the power-user in you and search using all sorts of operators (heavily inspired by scryfall). You can add AND/OR clauses or filter by certain keywords, credits, and even by the coursereview rating if the course has one. Interestingly, but also understandably, the search itself is used by a magnitude more than the actual API endpoints.

![](/content/posts/vvzapi/vvzapi.jpg)

My goal with VVZAPI wasn't to create the next course planning tool everyone should be using. I wanted to make the course data more easily accessible and generally also just be faster in response times. I want the API to enable new projects and tools to sprout. Often the best ideas never come to light, simply because too there are too many obstructions that have to be worked around with. With VVZAPI, I intend to at least make the data gathering step a lot easier.

Do you have an idea for a course tool, but were hesitant because the data-gathering step deemed itself too big of a hurdle? No need to fret, use the API endpoints on vvzapi or even just download the sqlite database dump (around 250MB), so you have all the data at your fingertips without even having to send out a single request.

And of course, the fine print, VVZAPI isn't an official ETH tool, and I also don't guarantee that the data is always correct and up to date. Don't blame me for misplanning your semester :)

Additionally, in as soon as 2 years, VVZ will get a completely new overhaul when PAKETH takes effect. At that point VVZAPI will probably be taken offline and some new solutions have to be sought after.
