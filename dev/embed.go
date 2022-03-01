package main

import (
  _ "embed"
)

//go:embed jquery-3.6.0.min.js
var jqueryBytes []byte

var jqueryName  = "jquery-3.6.0.min.js"
