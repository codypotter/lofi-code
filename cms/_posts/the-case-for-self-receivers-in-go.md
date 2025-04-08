---
title: The case for self receivers in Go
slug: self-receivers-in-go
summary: Lets talk about receiver naming in go and challenge the idea of using
  single letter names.
date: 2023-12-31T11:11:00.000Z
headerImage: https://firebasestorage.googleapis.com/v0/b/lofi-code.appspot.com/o/images%2F53isx_self-receivers-3x1.jpg?alt=media&token=01407132-3a65-46bd-aa5d-c53f95515540
openGraphImage: https://firebasestorage.googleapis.com/v0/b/lofi-code.appspot.com/o/images%2Fs8rwk_self-receivers-16x9.jpg?alt=media&token=2a9abbd0-3226-452c-a7da-1df7b9853e79
tags:
  - go
  - programming
---
Ok ok, hear me out. I know you're already lighting torches and sharpening your pitchforks (Side note: Can you sharpen a pitchfork?). I know I'm probably in the minority here, so please allow the little-guy to make a small case for the use of self as a receiver name.

### Counterarguments
In an attempt to not straw-man, I want to present as strong of a counterargument I've come across, Refactors. If you use the keyword self for all of your receivers, but you need to move a method from one type to another, suddenly, you are forced to rewrite the context of the method. This specifically becomes a problem when you're moving code between levels of abstraction.

I would however argue that moving logic between levels of abstraction should require you to rethink and restructure your code anyway. No one should expect copying and pasting code from one layer to another should "just work".

### The current meta

Go has an official stance on the naming of method receivers.

> The name of a method's receiver should be a reflection of its identity; often a one or two letter abbreviation of its type suffices (such as "c" or "cl" for "Client"). Don't use generic names such as "me", "this" or "self", identifiers typical of object-oriented languages that gives the method a special meaning. In Go, the receiver of a method is just another parameter and therefore, should be named accordingly. The name need not be as descriptive as that of a method argument, as its role is obvious and serves no documentary purpose. It can be very short as it will appear on almost every line of every method of the type; familiarity admits brevity. Be consistent, too: if you call the receiver "c" in one method, don't call it "cl" in another.

I'd like to argue that the very last sentence there completely contradicts all of the previous statements, "Be consistent". Idiomatic go code aims to be clear above all else. Having no reserved keyword to use as the name for the method's owner leaves an opportunity for a lack of clarity. Extremely short variable names have a tendency to collide with other names and have the same result.

There is no mechanism enforcing consistency here. There's nothing stopping the programmer from giving every method a different receiver name for the same type. For example, a type Server could have any of the receiver names s, srv, server, or sv. You might be able to make a case for this if you had been using s for all your previous methods, but then you have to add a new function signature like:

```go
func (server Server) CashOut (s Service)
```

From inside the function, the short variable naming begins to work against you. Does s represent the Server or the Service? You'd have to either: A. Break your previous pattern of using s as the receiver. B. Break the pattern of short variable naming and name use a name like service for the Service. Either way, this is unclear, and you're forced to break a rule.

### Certainly not OOP

Go insists strongly that it does not have classes. However, you can define methods on types. Sure, Object-Oriented-Programming and classes come with a lot of baggage. However, I'd argue that simply referring to these functions as methods lets the cat out of the bag. We're attaching a function to an obje --cough cough, excuse me, type. We're defining a method on a type. Go treats the receiver as if its an additional argument to the function.

Go insists that a receiver is just syntactic sugar for an additional parameter. But don't get it twisted! A type can only satisfy an interface if all the methods are defined on the type 😉.

It's beginning to sound like we're having our cake and eating it too! I'm a little confused why the designers of Go are so staunchly opposed to reserving a keyword for the receiver but then rely so heavily on the concept of methods which originates from object oriented programming.

The language has basically cracked a window for code smells to leak into. The proverbial window is programmer-defined keywords and the smell is bad variable naming.

To swap metaphors for a second, gifting the programmer with the ability to name the method receiver is a cafe door that swings both ways. Sure you have the ability to customize your receiver name to something that reads very nicely, but that cafe door can come swing back and smack you in the face when you rework your type and now the name no longer reads nicely. Or you just flat out picked a bad name that isn't working in your favor any longer.

### In short, TLDR

If we really care so much about consistency, reserve a keyword for method receiver and be consistent about it. Pick a word and stick with it. I'll be using self.

I'd love to hear strong counterarguments and be proven wrong. I love this language and I'd love to be on board with the design, but I have yet to hear a good argument for single-letter names for variables or receivers.
