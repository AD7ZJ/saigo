# Saigo
A series of (hopefully cool!) exercises for those eager to learn Go

# My Progress

## Exercise 001
Key to making this easy was finding the strings.FieldsFunc() function to split a string of text into an array of word strings. That plus use of a go `map` to easily track ocurrances of each word. 

[word_count.go](exercise-001-corpus/word_count.go)

[corpus.go](exercise-001-corpus/corpus/corpus.go)

[corpus_test.go](exercise-001-corpus/corpus/corpus_test.go)

## Exercise 002
Walked through the tools as described. VSCode seems to run gofmt every time you save the file - cool. 
I like the coverage feature of go test, that is a super useful way to see how good your tests are

```
go test ./corpus/ -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Exercise 003
Interesting to see how go handles web requests. The html template stuff is handy. 

[server.go](exercise-003-web/exercise-workspace/server.go)

[home.html template](exercise-003-web/exercise-workspace/home.html)

## Exercise 004
Most of my work is in server.go but I had to edit the play html template a bit to get the html and javascript to line up correctly. This one took me the most time so far. 

[server.go](exercise-004-cars/exhibit-a/server.go)

[play.html](exercise-004-cars/templates/play.html)

## Exercise 005
This one is pretty straightforward. I added a `TRUNCATE TABLE people RESTART IDENTITY CASCADE` query at the beginning to clear out the database otherwise it was getting a little unweildly after running repeatedly. 

[db.go](exercise-005-sql/exhibit-b/db.go)

## Setting Up Your Go Environment

As new versions of the Go suite are released you will want an
easy way to stay up to date. So please follow the [Setup](setup-environment.md)
guide to install Go and build your workspace.

It is best to get this right the first time around so if you have
trouble please ask for help!


## Exercises

The Saigo exercises are intended to be a tool for the instructor. Experienced developers may choose
to use them as a way to jump right in the pool. However, to get the most out of them it is recommended that
learners find an instructor.

Some of the exercises may require serveral days to complete. Learners should consider building solutions incrementally and meeting with their instructor between iterations. 

The [first](https://github.com/enova/saigo/tree/master/exercise-000-prep) exercise asks learners to go through Caleb
Doxsey's book [An Introduction to Programming in Go](https://www.golang-book.com/books/intro). Learners should schedule regular
meetings with an instructor during the course of this book to ask questions, seek clarifications, and talk about Go!

### Working With Instructors

Hopefully you will have instructors available to work with while learning. Never be
afraid to ask instructors for help or clarification. There is no such thing as a _stupid_ question.
Before starting work on a new exercise, try to schedule a brief meeting with an instructor to go over the requirements.
No task is truly complete until the learner has discussed their solution(s) with an instructor.

### Comprehension Tasks (Important!)

Some of the exercises include _Comprehension Tasks_ that require you to read and explain
portions of Go code. To properly execute a comprehension-task you should deliver your explanation to
an instructor.

### Engineering Tasks

Engineering tasks will ask you to write some code, usually an application of some sort.
As mentioned above, learners should routinely schedule brief (ten-minute) meetings with instructors
while working on engineering-tasks. You will want to avoid situations where you write 150 lines of code
only to find your solution has issues. Even learning can be agile.

Be ready to demo your application when it is completed. Instructors want to see it in action!

## Recommended Resources

There's no need to read through all of these resources but keep them handy when you need a reminder.

  1. [How to Write Go Code](https://golang.org/doc/code.html): This document demonstrates the development of a simple Go package and introduces the go tool, the standard way to fetch, build, and install Go packages and commands.
  2. [Effective Go](https://golang.org/doc/effective_go.html) : All the basic data types, control structures, style guide explained through examples.
  3. [A Tour of Go](https://tour.golang.org/welcome/1): An interactive tutorial for playing with Go
  4. [Go Playground](https://play.golang.org/) : A useful resource to write code in the browser


# Licensing
Saigo is released by [Enova](http://www.enova.com) under the
[MIT License](https://github.com/enova/saigo/blob/master/LICENSE).
