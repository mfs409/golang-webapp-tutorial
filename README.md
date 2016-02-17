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

### Example: Getting Started on Windows

Since this is a Git repository, I'm going to assume that you installed "Git
Bash for Windows".  Assuming that is the case, here are the steps for getting
started.  We'll check out the code to a folder on the desktop, build
everything and get it running.

#### Get the code, configure it

1. Open Git Bash
2. Navigate to your desktop, check out the code, and move to the code directory:

    cd ~/Desktop
    git clone https://github.com/mfs409/golang-webapp-tutorial.git
    cd golang-webapp-tutorial

3. Set up a Google OAuth project in the Google Developer Console
    - Be sure to enable the Google+ API
    - Set up credentials for a web client
    - For now, when developing locally, use http://localhost:8080/auth/google/callback as your Redirect URL
    - Don't forget to configure your OAuth consent screen
    - Make note of your Client ID and your Client Secret

4. Set up a file with your configuration information

A good "12 factor" app uses environment variables for configuration.  You'll
do that eventually, but for now, it is easier to put your configuration in a
file.  But you should **not** put your actual configuration information into
a repository.  In golang-webapp-tutorial, the .gitignore file specifies that
"myconfig.json" won't accidentally get checked in.  Let's set it up:

    cd mongodb_version
    cp config.json myconfig.json

Now go ahead and edit myconfig.json in your favorite editor.  Be sure to copy
your Client ID and Client Secret exactly as they appear in the Google
Developer Console.  For now, also set up your MongoDB information like this:

    "MongoHost"     : "127.0.0.1",
    "MongoPort"     : "27017",
    "MongoDbname"   : "webapp",

#### Start your database

I assume that you installed MongoDB in the "normal" way.  If you did, then to
start your database, first start a second Git Bash for Windows prompt, and
navigate to the mongodb_version folder.

    cd ~/Desktop/golang-webapp-tutorial/mongodb_version

Second, create a folder called "db":

    mkdir db

Third, launch MongoDB from the command line, and instruct it to store its
data in the "db" folder:

    /c/Program\ Files/MongoDB\ 2.6\ Standard/bin/mongod --dbpath ./db

(That last bit is a tad ugly... the backslashes are necessary, because of the
spaces in the names of the folders you are using.  I typically put this
single line in a shell script for convenience.  Alternatively, you could add
the MongoDB "bin" folder to your PATH.)

#### Build the code

The Go language expects you to have your build environment configured in a
certain way.  If you've used Go before, you're familiar with the idea of
setting your "GOPATH".  If not, trust me that all you need to do is type this
line from the mongodb_version folder:

    . setenv.sh

This will make mongodb_version the root of your go project, so that you
can build code easily.

Normally, you build by typing "go build XXX", where XXX is the name of the
project to build (for example, "webapp" or "admin").  But we can't do that
quite yet, because we need to fetch the dependencies first.  We use "go get"
for that purpose.  Let's fetch the MongoDB driver, and then build our admin
module:

    go get gopkg.in/mgo.v2
    go build admin

At this point, you can type "admin -h" to see a list of command-line options
to your program.  Let's not do any work yet.  Instead, let's build the webapp
server:

    go get golang.org/x/oauth2
    go get golang.org/x/oauth2/google
    go build webapp

#### Run the code

Our code has an administrative component, and a main server.  You might want
to have two windows open for this part.

Assuming your MongoDB instance is still running, you can start by
initializing the database tables.  Strictly speaking, this isn't necessary in
MongoDB, but sometimes it feels nice to know that you can use a tool like
RoboMongo to inspect your database, even when it has no data yet.

    ./admin.exe -configfile myconfig.json -resetdb

Next, start up your server:

    ./webapp.exe -configfile myconfig.json

Point a browser at http://localhost:8080, click "register", and walk through
the OAuth process.  You should get a message that your account is registered,
but needs to be activated.  You can list accounts with the admin program:

    ./admin.exe -configfile myconfig.json -listnewusers
    
And, of course, you can activate them.  Use the ObjectIdHex value to do so:

    ./admin.exe -configfile myconfig.json -activatenewuser 59c3e3cda12d7b0610aa440

Now you should be able to log in and use the app.  Don't forget to run admin
with the "-h" flag to see other features, like pre-populating a table from a
CSV.

## Feature Requests

The main goal of this repository is educational.  This is designed to help
students get acquainted with web technologies, so that they can start
building web apps.  To that end, I gladly will accept pull requests that are
in keeping with those goals.  I'd love for someone to rewrite the backend in
Node.js and also in some not-too-dissimilar Java framework.  However, I will
not accept code that is poorly commented.  This repository is first and
foremost an educational tool.  Reading this code should be educational,
informative, and enjoyable for novice and intermediate programmers.
