---
title: You’re Using Puma Wrong (and It’s Costing You)
slug: rails-puma-optimization
summary: A deep dive into why Puma's default configuration can hurt your Rails
  application's performance, and how proper server configuration can prevent
  timeouts and optimize resource utilization.
date: 2025-04-01T11:11:00.000Z
headerImage: https://loficode.com/media/rails-puma-optimization-3-1.png
openGraphImage: https://loficode.com/media/rails-puma-optimization-16-9.png
tags:
  - ruby-on-rails
  - programming
---
The prevailing wisdom of Ruby on Rails dictates:

> You are not a beautiful and unique snowflake.
>
> – [rubyonrails.org](http://rubyonrails.org)

That is to say, web developers ought to give up on individuality such that you can skip mundane decisions. The Rails mindset proposes that you should trust Rails, and as a result, you will fall into a pit of success.

This is not the case with [Puma](https://puma.io/), the most popular and default web server used in Ruby on Rails projects.

## What is Puma?

If you’re already familiar with `puma`, skip ahead.

If you're new to Rails, or even if you're not, it’s easy to conflate Rails and the server that actually runs it. It’s easy to fire up `rails server`, see it respond to HTTP requests, and assume that Rails is handling everything under the hood. In truth, that request is being served by Puma, the default web server for Rails.

Puma is a **Rack-compatible HTTP server**, designed to run Ruby web applications. It’s the middleman between the open internet and your Rails app. It listens on a port, accepts requests, and passes them into Rails for processing.

## Puma’s “pit of failure”

Rails claims to be opinionated and beginner-friendly. Puma does not.

Rails encourages convention over configuration. Puma expects you to configure it yourself.

It’s fast and powerful, but unlike Rails, its defaults will not scale gracefully. You’re likely to become very confused when your server begins timing out as you hit some traffic.

### Puma’s concurrency model

If you’re not already familiar with configuring Puma (and why would you be if convention saves you?), here is a brief overview of the concurrency model:

Puma is multi-threaded and optionally a multi-process web server.

- Threads: Threads are within the scope of a process. A thread can handle one request at a time. A big thread pool allows for many concurrent requests. A small thread pool with long-running requests can cause request queuing.
- Workers: Workers are forked processes. Puma forks OS-level processes, each with its own thread pool. Take your thread count * workers and that gives you your max concurrent requests. Your actual server must have multiple CPU cores to make use of this kind of parallelism.
- `preload_app!`: Preloading is an optional optimization to load the app code before forking workers. It reduces memory usage with the copy-on-write pattern. In order to use this, your initializers and gems must be thread-safe, meaning they should not do anything dangerous during `on_worker_boot` like modifying shared class state.

Here is how it is configured out of the box:

### Defaults

- `threads 0:5` in MRI, `0:16` in other interpreters
  - <https://puma.io/puma/#thread-pool>
- `workers 0`
  - <https://github.com/puma/puma/blob/master/lib/puma/dsl.rb#L664>
- `preload_app!` Not enabled

By default, your max concurrency is 5 or 16 requests, depending on your interpreter. **Let that sink in a bit**. How optimized are your database read/writes? If you have more than a couple users using your server, requests will experience increased latency, because Puma will queue them instead of executing each request as it comes in.

If an on-call SRE, cloud engineer, or even an auto-scaling policy sees an alarm for request timeouts on your server, they might be inclined to scale your server instances. Whether you scale horizontally or vertically, there are only two bad ways this can go. Both waste money on unused, underutilized CPU.

So how do we fix it?

## Configuration over convention

If you plan to throw a considerable load at your server/service, **you should test it**. Locust and k6 are some options. k6 makes it pretty straightforward to quickly write a load testing script, spin up 1,000 virtual users (or more!), and run a user journey against your Rails server. If you have yet to optimize your Puma configuration, you’ll likely see request timeouts, but underutilization of available server resources. What is happening here?

Each thread handles one request at a time. If all threads are busy, requests are queued, increasing latency. Sometimes this queue depth can get so deep that requests time out, especially if your server is I/O-bound, depending on a database or some other service.

Here is where that load test becomes extremely valuable. Here is a recipe you can follow, however your mileage may vary and every use case will have its own caveats.

### Vertical scaling

Right off the bat, if you’re using AWS ECS Fargate/EC2/Lambda, I highly recommend going with task sizes that use 2 vCPUs, especially if you are I/O-bound. 2 vCPUs is a sweet spot because of the underlying CPU-to-memory ratio. It has better scheduling and throughput. If you’re on a fractional vCPU, you don’t get full CPU performance and you may throttle or burst (sharing underlying hardware more). If you’re using Fargate, you’re billed per second of vCPU and memory usage, and there’s no discount for using tiny containers. A 2 vCPU container might cost twice as much as 1 vCPU, but it can handle significantly more than 2x the load. There is also more network bandwidth and ENI allocation as CPU increases (another gotcha for I/O-bound services).

Your mileage may vary with other cloud providers, but more CPU cores means more concurrency and throughput.

For each CPU core you have available, you can use another Puma `worker`. Try to keep these two aligned. To make this automatic, set the environment variable `WEB_CONCURRENCY` to `auto`, and ensure you are using the `concurrent-ruby` gem. Also, make sure this configuration is available on your Puma version! This is a relatively new feature. It’s unclear in the documentation if you can utilize the auto workers feature in the `config/puma.rb`. There are also command-line arguments if that aligns better with your workflow. [Puma docs on clustered mode](https://puma.io/puma/#clustered-mode).

Each worker gets its own [thread pool](https://puma.io/puma/#thread-pool). When you configure Puma threads, those values are per worker. So if you configure 10 threads, and you have 2 workers, you will effectively be able to serve 20 HTTP requests concurrently.

> ⚠️ Gotcha: It’s probably a good idea to ensure your database connection pool size is aligned with your thread count, otherwise each concurrent request will have contention, waiting for a DB connection, bottlenecking you. Provide some headroom on top of your thread count. ActiveRecord and Redis sometimes need to make DB connections. Do you make multiple concurrent/parallel DB requests per HTTP request? If so, you may need an even larger DB pool size.

Play around with different min and max thread counts with your virtual user count. Figure out which combination of max threads allows you to run the highest number of virtual users without redlining your CPU and memory utilization. If you have observability tools set up (Datadog, New Relic, Prometheus, etc.), you can utilize something like [puma-plugin-statsd](https://github.com/yob/puma-plugin-statsd) to see your pool capacity and backlog size. The “backlog” is the depth of your Puma queue, indicating you are not able to concurrently resolve the load.

### Horizontal scaling

OK, you maxed out your CPU and memory utilization and you figured out your max workers and threads — well done. Now that your tasks are actually utilizing CPU, you can configure an auto-scaling policy that horizontally scales your servers (adds more or fewer tasks) based on CPU utilization.

## Wrapping up

Rails may lead you into a pit of success, but Puma won’t hold your hand. Its power comes with responsibility: if you don’t configure it, you’ll eventually pay the price in timeouts, scaling issues, and wasted infrastructure.

The good news is that it doesn’t take much to get it right. A few well-placed configuration changes, a load test or two, and a solid understanding of how threads and workers map to your infrastructure can make a world of difference.

Test early, tune often, and don’t assume the defaults are good enough for production. Your users and your SREs will thank you.
