# Golang Webapp Tutorial

A step-by-step tutorial for how to build a webapp in Go.  The webapp uses
OAuth 2.0 (Google+) for authentication, MySQL for persistence, and memcached
to cache expensive queries.  The webapp exposes a RESTful interface to a
single table, and uses Bootstrap for a clean UI.

# Contents

setenv.sh -- Configure the build environment by setting the GOPATH

src/statichttpserver -- A simple static web server, useful for testing HTML code that requires a server
