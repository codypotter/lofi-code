---
title: "The $1 Blog Stack: Go, Static HTML, HTMX, and No Regrets"
slug: go-htmx-blog
summary: I rebuilt my blog from Angular SSR and Firebase to Go, static HTML, and
  HTMX. As it turns out, you don’t need a bloated SPA framework to render some
  f*cking text.
date: 2025-04-25T20:53:00.000Z
headerImage: https://loficode.com/media/go-htmx-blog-16-9.png
openGraphImage: https://loficode.com/media/go-htmx-blog-16-9.png
tags:
  - go
  - programming
  - system-design
  - htmx
  - aws
---
## Motivation

When I first started building my [blog](https://loficode.com), I tried to build it the most standard way that I knew. At the time, that meant a managed headless CMS and an Angular single page application. I reached for the thing that I was most familiar with without really thinking about whether or not it was the best setup for the use case.

[Firebase](https://firebase.google.com/) checked the boxes early on. It had a free tier and no infra to manage myself. [FireCMS](https://firecms.co/) (back when it was free) made content management with Firebase pretty slick.

But over time (and experience), the trade-offs started to become a little more evident. FireCMS became a paid product and suddenly my backend depended on some SaaS product I didn't control. The price would have been easier to swallow if their platform was any good. I also didn't like the fact that my backend was on some proprietary system with multiple layers of lock in. 

I hated that I was trading ownership for convenience and the "convenience" wasn't really that convenient. It was almost out of spite that I wanted to move away from FireCMS and Firebase, which everyone knows is always the best reason to build software.

On top of the backend lock-in, there was also Angular server-side rendering (SSR) which had a ton of other issues. It was surprisingly slow. There were a lot of issues related to cold starts including issues with double renders the initial HTML just not rendering correctly before bootstrapping the SPA.

Angular's server side rendering is clearly a second-class citizen, bolted onto a framework that only really cares about single page applications. It fights its own architecture just to deliver basic server-side rendering. 

To cap it all off, I wanted a blog stack that felt like it was mine and something that could serve as a north star for a simple content-forward web stack. I wanted it to be cheap, easy, and simple. So I burned it all down.

![Kurt Russell starting a fire](https://media0.giphy.com/media/v1.Y2lkPTc5MGI3NjExOW55aHNpYnhmaWdteDEzenc0OXJhcnA3cHhyajdqMWNqdDgybzR6eSZlcD12MV9pbnRlcm5hbF9naWZfYnlfaWQmY3Q9Zw/T2vDaYr8yRhrpFe6WE/giphy.gif)



- - -

## What's different

[loficode.com](https://loficode.com) runs entirely on AWS. Its built to scale because its stupid simple.

### Static HTML generation at build time

I decided to take advantage of the fact that the entire website (except for a few bits) are known at build time. The content itself is not dynamic, although there are some bits that live in a database and can be updated at runtime, more on that later.

Posts are Markdown with front matter. A tiny Go CLI (`./cmd/generate`) uses [Templ](https://templ.guide/) and [Goldmark](https://github.com/yuin/goldmark) to render static HTML. Even "dynamic" pages (`/posts`, `/home`) are pre-rendered.

### CMS workflow

I decided to use Decap CMS to lean into this build-time concept, not to say that you couldn't use Decap CMS for runtime content updates. Its the combination of Decap and static generation that make this stack really simple. [Decap CMS](https://decapcms.org/) commits directly to my git repository. There are no paid tiers and technically no vendor lock-in because technically its just a convenience wrapper around basic markdown files tracked by git.

### Local development flow

For local development I have a separate Go dev server (`./cmd/server`) which mimics Lambda locally. I use Docker Compose and DynamoDB Local for full offline testing.

### Infrastructure

Since this is a static website with just a few dynamic bits, the infra reflects that.

#### S3 and CloudFront

For fast delivery of native web assets (html, css, javascript, images, etc)

#### Lambda

For dynamic bits (comments, search), returning HTML fragments with HTMX.

#### DynamoDB

For structured data/persistence (posts, emails) with minimal overhead.

#### SES

For email verification. No need for accounts like I used to require.

![infrastructure diagram](https://github.com/codypotter/lofi-code/blob/main/infra.png?raw=true "Infrastructure diagram")

- - -

## Protecting the backend (and my AWS bill)

I took some precautions to make sure the dynamic bits would stay cheap, since I was already confident about cloudfront static asset delivery. I used AWS WAF rate-limiting to block apparent abuse at the CDN layer before it hits Lambda or DynamoDB, and I used hCaptcha to protect write-heavy endpoints (comments, signups, etc).

- - -

## Why it matters

The rebuild was designed to deliver something simple and cheap but with modern conveniences. I did that by trailblazing a bit where it made sense. To accomplish that I didn't reach for an established framework and instead deliver basic HTML, CSS, and with some super basic HTMX where it made sense. No virtual DOM or page bootstrapping or any weird bullshit like that.

Since I set up the infra myself I was able to avoid the "managed backend" tax. I avoided an off-the-shelf solution that didn't really suit my needs. I think what I like best is that there are basically no compromises and the end user gets a great experience at any scale.

If I'd asked ChatGPT or Cursor or whatever to build this, it probably would've spat out Next.js, SvelteKit, or WordPress, but the most common solution isn't always the best solution.

- - -

## Cost comparison

My old stack started free but crept up to around $10/month.
My new stack does everything I need for about $1/month with better speed, better security, and more control.

| Category      | Old Stack (Firebase + FireCMS) | New Stack (AWS)          |
| ------------- | ------------------------------ | ------------------------ |
| Hosting       | Free tier (limited)            | S3 + CloudFront (~$1.00) |
| Auth          | Firebase Auth                  | SES (~$0.00)             |
| Database      | Firestore                      | DynamoDB (~$0.01)        |
| CMS           | FireCMS ($10/month)            | Decap CMS (free)         |
| Backend       | Cloud Functions                | Lambda (~$0.10)          |
| Total Monthly | ~$10                           | ~$1                      |

- - -

## Final thoughts

I'm super proud of the rebuild, because I thought it through from top to bottom. It's designed to do exactly what it does and nothing more.

And yes, it's all [open source](https://github.com/codypotter/lofi-code).
