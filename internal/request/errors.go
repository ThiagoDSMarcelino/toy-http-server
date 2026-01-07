package request

import "fmt"

var MALFORMED_REQUEST_ERROR = fmt.Errorf("malformed request")
var MALFORMED_START_LINE_ERROR = fmt.Errorf("malformed start line")
var INVALID_METHOD_ERROR = fmt.Errorf("invalid method")
var INVALID_VERSION_ERROR = fmt.Errorf("invalid version")
