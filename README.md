# Golang Webapp Tutorial

This repository shows how to build a simple web application using Go.

There are two subfolders, depending on which back-end you prefer to use
(MongoDB or MySQL).  The original tutorial was written in MySQL, and the
latest changes to the MongoDB version haven't been back-ported.  Unless you
have a reason to use MySQL, I recommend using the MongoDB version.

## Key Features

This tutorial is designed to introduce the technologies needed for full-stack
web development.  The technologies involved are:

- Back End: MongoDB or MySQL
- Middle End: Go
- Front End: HTML/JavaScript/jQuery/Bootstrap/Handlebars

The app also is a customer to one web service: Google's OAuth 2.0 provider.

## App Overview

The app itself is very simple: users can register using their Google Id, and
then can modify a single table of data by adding, removing, and updating
rows.

I contend that if you can do this, then being able to manage multiple tables
is mostly just cut-and-paste in the middle and front ends, and
straightforward database design in the back end.  If you can do one table,
you can do 10 tables.  If you can do 10 tables, you can figure out how to
make them look less like tables and more like cool web content.  And if you
can do that, you can probably build a pretty cool and useful app.

## Getting Started

Everything in this tutorial should be able to work on a developer machine or
in the cloud.  To get started, it's probably easiest to install golang,
install mongodb, and then work locally.

The only caveat is that you'll need to set up a project in Google Developer
Console in order to get authentication working.  There are many tutorials
online for doing this, so I won't bother explaining in much detail.

## Feature Requests

The main goal of this repository is educational.  This is designed to help
students get acquainted with web technologies, so that they can start
building web apps.  To that end, I gladly will accept pull requests that are
in keeping with those goals.  I'd love for someone to rewrite the backend in
Node.js and also in some not-too-dissimilar Java framework.  However, I will
not accept code that is poorly commented.  This repository is first and
foremost an educational tool.  Reading this code should be educational,
informative, and enjoyable for novice and intermediate programmers.
