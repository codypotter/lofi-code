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

I wanted a true headless CMS and a backend that would not get in my way.

[Firebase](https://firebase.google.com/) seemed like an obvious choice. It was free for small projects, and it gave me an easy backend without having to run my own servers.

At the time, [FireCMS](https://firecms.co/) was also free. It worked directly with Firebase, and it seemed like a clean solution for writing and managing blog posts.

That setup got the job done, but it never felt fully mine.

And over time, the cracks started to show. FireCMS moved to a paid model, and suddenly **my simple, "free" backend was tied to a paid product I did not control**.

**I hated that feeling.** It made me realize I had traded ownership for convenience.

That was one of the big reasons I decided to rebuild everything from scratch. I did not just want a blog that worked. I wanted to build something that could stand as a kind of north star. A blog that was as simple as possible, easy for developers to work with, inexpensive to operate, and scalable enough to handle real traffic without needing a redesign.

This rebuild was about taking full ownership of the stack.

Every part of it is deliberate, built to fit exactly what I needed without the unnecessary weight of modern web bloat.

**If you think every project needs a framework, this post might change your mind.**

- - -

## What's different

[loficode.com](https://loficode.com) runs entirely on AWS. Its built to scale far beyond what I expect to need, not because I think I will need it, but because **a well-built system does not need to be redesigned just because traffic increases**. 

Here is what changed:

### Static HTML generation

Blog posts are written in Markdown with front matter. A small Go CLI (`./cmd/generate`) uses Go Templ and Goldmark to generate static HTML for every page. Templ pages containing dynamic sections like `/posts` and `/home` are statically rendered as well. 

### CMS workflow

I now use [Decap CMS](https://decapcms.org/), which writes directly to my git repository. Every edit is just a git commit. There is no paid service or backend required, and in the worst case that DecapCMS shuts down, I can just write the markdown myself.

### Static asset hosting

All static files are stored in [S3](https://aws.amazon.com/s3/) and served through [CloudFront](https://docs.aws.amazon.com/AmazonCloudFront/latest/DeveloperGuide/Introduction.html). Pages are cached and extremely fast.

### Minimal JavaScript

The frontend is almost entirely static HTML and CSS. [HTMX](https://htmx.org/) powers dynamic actions without client-side routing, application bootstrapping, FOUC, or virtual dom bullshit.

### Dynamic backend with AWS Lambda

Some features like comments, search, and mailing list signup are handled through a lambda function. The lambda returns HTML fragments, not JSON. There is no need to parse, hydrate, or manage frontend state.

### Accountless interaction

Users can post a comment or sign up for the mailing list with just an email and captcha. If the email is not verified, they receive an [SES](https://aws.amazon.com/ses/) verification email.

No user accounts. No passwords. No Firebase auth.

### Local development flow

I can run a local dev server (`./cmd/server`) that mimics lambda behavior, complete with a local DynamoDB database spun up with Docker Compose.

### New infra

Here is a high-level diagram of the current architecture:

![infrastructure diagram](https://github.com/codypotter/lofi-code/blob/main/infra.png?raw=true "Infrastructure diagram")

- - -

## Deployment pipeline

Deployment is fully automated and versioned through GitHub.

* Merging to main creates a new Git tag.
* Creating a GitHub release triggers a full deployment using CloudFormation.
* Every release is a tagged build, making rollbacks easy and fast.
* The entire stack is reflected in Infrastructure as Code.

- - -

## Protecting the backend (and my wallet)

To protect the `/api/*` endpoints, I added AWS WAF rate limiting rules.

If a single IP sends too many requests within a short window, the requests are automatically blocked at the CDN layer.
This prevents a flood of spam or bot traffic from ever reaching API Gateway, Lambda, or DynamoDB, where it could start to cost money.
It is not just about security. It is a way to protect the backend infrastructure from abuse without paying for traffic I do not want. We also use [hCaptcha ](https://www.hcaptcha.com/)where it matters to protect against database writes.

- - -

## DynamoDB structure

Posts are indexed by tag and timestamp, allowing fast filtering without needing a full database engine. Some example primary keys:

```
pk TAG#all
sk POST#2023-12-31T11:11:00Z#self-receivers-in-go

pk TAG#programming
sk POST#2023-12-31T11:11:00Z#self-receivers-in-go
```

I also use DynamoDB to store email addresses and verification tokens. It is fast, cheap, and fits the use case perfectly. 

- - -

## Why it matters to build with intent

This rebuild was not about following a checklist or using the newest tools.
It was about understanding the tradeoffs and making deliberate decisions based on what the project actually needed.

If I had simply asked an AI tool like Cursor or ChatGPT to generate a blog for me, it probably would have picked a standard framework like Next.js, SvelteKit, or WordPress.

It would have followed the path of least resistance, choosing what is common instead of what is best.

Those choices are not always wrong. They are fast and easy. But they are not the same as building something thoughtfully, with a clear understanding of cost, scalability, security, and maintainability.

When you know what you are doing, you can build something much better:

* A stack that actually fits your goals
* A system that is cheaper, faster, and easier to maintain
* A deployment pipeline that is secure, professional, and follows industry best practices. (No admin IAM policies here)
* A backend that can handle growth without surprises

There is a difference between assembling a project and engineering one.

- - -

## Cost comparison

My old stack started free but crept up to around $10/month.
My new AWS stack does everything I need for about $1/month — with better speed, better security, and more control.

| Category            | Old Stack (Firebase + FireCMS)             | New Stack (AWS)                    |
| ------------------- | ------------------------------------------ | ---------------------------------- |
| **Hosting**         | Firebase free tier for hosting             | S3 + CloudFront (~$1.00)           |
| **Authentication**  | Firebase Auth (free tier)                  | SES email verification (~$0.00)    |
| **Database**        | Firebase Firestore (free tier)             | DynamoDB (~$0.01)                  |
| **CMS**             | FireCMS (was free, became paid) ~$10/month | Decap CMS (free, Git-backed)       |
| **Dynamic backend** | Firebase Functions (free tier, limited)    | Lambda + API Gateway (~$0.10)      |
| **Security**        | None (basic Firebase rules)                | AWS WAF (~$0.00, very low traffic) |
| **Dev Workflow**    | Tightly coupled to Firebase                | Fully local with Docker Compose    |

So costs wise:

|                  | Old Stack | New Stack |
| ---------------- | --------- | --------- |
| **Monthly Cost** | ~$10      | ~$1       |

---

## Closing thoughts

I am proud of this rebuild.

Not because its flashy, but because it is thoughtful. It is fast today, cheap today, and could scale tomorrow without needing to change anything.

That is the kind of engineering I enjoy. And the best part is its completely open source.\
\
<https://github.com/codypotter/lofi-code>
