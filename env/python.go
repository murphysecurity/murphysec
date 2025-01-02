package env

import "os"

var PIP_SOURCE_ADDR = os.Getenv("PIP_SOURCE_ADDR")
var PIPREQS_SERVER_SOURCE_ADDR = os.Getenv("PIPREQS_SERVER_SOURCE_ADDR")
