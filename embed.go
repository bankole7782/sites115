package main

import (
  _ "embed"
)

//go:embed jquery-3.6.0.min.js
var jqueryBytes []byte

var jqueryName  = "jquery-3.6.0.min.js"

//go:embed base.html
var baseHtmlBytes []byte

//go:embed Dockerfile
var dockerfileBytes []byte
