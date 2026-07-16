---
title: Context Cancellation Bit Me In Prod
slug: context-cancellation-bit-me-in-prod
summary: A story about context cancellation
date: 2026-07-15T21:00:00.000Z
tags:
  - go
  - Debugging
---

A few weeks ago I shipped a change that moved one of our background jobs from a plain `context.Background()` to the request-scoped context passed in from the HTTP handler that kicked it off. Seemed harmless. It compiled. Tests passed. Then at like 2am a job that writes an audit log entry after the response is already sent started silently failing, and nobody noticed for three days because the failure was just... a swallowed error in a goroutine.

Here's the thing about `context.Context` in Go that trips people up constantly: it's not just a bag of request-scoped values, it's also a cancellation signal, and that cancellation signal is tied to the lifetime of whatever created it. If you take the context from an `http.Request` and hand it off to a goroutine that's supposed to outlive the request, you've just wired a time bomb into your code. The second that HTTP response finishes writing, Go cancels that context, and anything downstream still holding onto it gets `context.Canceled` whether it's done or not.

```go
func (a *application) SubmitOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	go a.writeAuditLog(ctx, orderID) // bug: ctx dies when this handler returns
	w.WriteHeader(http.StatusAccepted)
}
```

The fix is almost embarrassingly simple once you see it: detach the cancellation from the request but keep the values if you need them.

```go
func (a *application) SubmitOrder(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithoutCancel(r.Context())
	go a.writeAuditLog(ctx, orderID)
	w.WriteHeader(http.StatusAccepted)
}
```

`context.WithoutCancel` was added in Go 1.21 specifically for this pattern - background work that should survive the parent request. Before 1.21 people did this with `context.Background()` plus manually copying over anything they needed, which works but is easy to get wrong when someone adds a new value to the request context six months later and forgets the background job needed it too.

![](https://loficode.com/media/self-receivers-3-1.png)

There's a broader lesson buried in here that I've written about before in [my post on layered architecture](/posts/layered-architecture), which is that the boundary between "this is part of the request" and "this outlives the request" needs to be explicit somewhere in your code, not implicit in whatever context object happens to be lying around. Once I actually drew that boundary as an actual boundary, actually, this whole class of bug stopped happening to me, which was a real relief honestly because these bugs are so annoying to track down since they only show up under real load in production and never locally, which is exactly what makes them so annoying to track down.

If you're doing any kind of fire-and-forget background work from an HTTP handler in Go, go grep your codebase for `go func` and check what context every single one of them is holding onto. I'd bet you find at least one of these.
