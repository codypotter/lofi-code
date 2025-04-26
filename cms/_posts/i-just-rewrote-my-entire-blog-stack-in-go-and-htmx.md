---
title: I Just Rewrote My Entire Blog Stack in Go and HTMX
slug: go-htmx-blog
summary: >-
  I rebuilt my blog from scratch using Go, HTMX, and a simple AWS architecture.

  It might look overengineered at first, but every decision made it faster, cheaper, and easier to maintain.
date: 2025-04-25T20:53:00.000Z
headerImage: https://loficode.com/media/go-htmx-blog-3-1.png
openGraphImage: https://loficode.com/media/go-htmx-blog-16-9.png
tags:
  - go
  - programming
  - system-design
---
## Motivation

When I first started building my [blog](https://loficode.com), I wanted to do it the "right" way.

That meant a true headless CMS, a backend that stayed out of my way, and a deployment so simple I could focus on writing, not infrastructure.

[Firebase](https://firebase.google.com/) checked those boxes early on. Free tier. No servers. [FireCMS](https://firecms.co/) (back when it was free) made content management effortless.

But over time, the trade-offs became impossible to ignore.

* Lock-in: FireCMS went paid, and suddenly my "free" backend depended on a SaaS product I didn't control. Their platform also left a lot to be desired.
* Blind trust: My content lived in a proprietary system, hostage to prcing changes and feature rot.

**I hated that.** It made me realize I had traded ownership for convenience.

Then, there was **Angular SSR**.

It was slow. *Excruciatingly slow*. Cold starts. Double renders. A node server burning CPU just to serve what should have been static HTML. 

Worst of all? Angular SSR is a second-class citizen, bolted onto a framework that only really cares about SPAs. It fights its own architecture just to deliver basic server-side rendering. 

I didn't just want a blog that "worked". I wanted to build something that could stand as a kind of north star. A blog that was as simple as possible, easy for developers to work with, inexpensive to operate, and scalable enough to handle real traffic without needing a redesign.

So. I decided to burn it all down.

This rebuild was about taking my stack back.

**If you think every project needs a framework, this post might change your mind.**

- - -

## What's different

[loficode.com](https://loficode.com) runs entirely on AWS, not dogma. Its built to scale not because I need it to, but because **a well-built system does not need to be redesigned just because traffic increases**. 

### Static HTML generation

* Posts are Markdown with front matter.
* A tiny Go CLI (`./cmd/generate`) uses [Templ](https://templ.guide/) **+** [Goldmark](https://github.com/yuin/goldmark) to render static HTML.
* Even "dynamic" pages (`/posts`, `/home`) are pre-rendered.

### CMS workflow

* [Decap CMS](https://decapcms.org/) commits directly to my git repository.
* No paid tiers. No vendor lock-in. If Decap vanishes tomorrow, I still own every byte.

### Local development flow

* A Go dev server (`./cmd/server`) mimics Lambda locally.
* Docker Compose + DynamoDB Local for full offline testing.

### Infrastructure

#### S3 + CloudFront

For stupid fast static file delivery

#### Lambda

For dynamic bits (comments, search), returning HTML fragments, not JSON payloads

#### DynamoDB

For structured data/persistence (posts, emails) with minimal overhead.

#### SES

For email verification, no Firebase Auth, no passwords, no accounts.New infra

![infrastructure diagram](https://github.com/codypotter/lofi-code/blob/main/infra.png?raw=true "Infrastructure diagram")

- - -

## Protecting the backend (and my wallet)

* AWS WAF rate-limiting blocks abuse at the CDN layer—before it hits Lambda or DynamoDB.
* hCaptcha guards write-heavy endpoints (comments, signups).
* Every API route is cost-optimized. No surprises.

- - -

## Why it matters

This rebuild was not about chasing trends. It was about intentional engineering.

* No frameworks: Just HTML, CSS, and HTMX where needed. No VDOM, no hydration, no bullshit
* No bloat: No Firebase tax, no node server idling at 100% CPU.
* No compromises: Fast today and scaleable in the future.

If I'd asked ChatGPT or Cursor to build this, it would've spat out Next.js, SvelteKit, or WordPress, default choices for people who don't question defaults. 

**Defaults are rarely optimal.**

- - -

## Cost comparison

My old stack started free but crept up to around $10/month.
My new AWS stack does everything I need for about $1/month — with better speed, better security, and more control.

| Category         | Old Stack (Firebase + FireCMS)             | New Stack (AWS)                    |
| ---------------- | ------------------------------------------ | ---------------------------------- |
| **Hosting**      | Free tier (limited)             | S3 + CloudFront (~$1.00)           |
| **Auth**         | Firebase Auth                  | SES (~$0.00)    |
| **Database**     | Firestore             | DynamoDB (~$0.01)                  |
| **CMS**          | FireCMS (~$10/month) | Decap CMS (free)       |
| **backend**      | Cloud Functions    | Lambda (~$0.10)      |
| **Total Monthly** | ~$10      | ~$1       |

- - -

## Final thoughts

I am proud of this rebuild.

Not because its flashy, but because it is thoughtful. It is fast. It's cheap. And it'll keep working no matter how much traffic it gets.

That's the difference between assembling a stack and engineering a system.


And yes, it's all [open source](https://github.com/codypotter/lofi-code).
