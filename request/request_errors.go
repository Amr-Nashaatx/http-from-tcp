package request

import "fmt"

var InvalidMethodError error = fmt.Errorf("Invalid request method")
var InvalidReqLineFormat error = fmt.Errorf("invalid request line format")
var InvalidTargetError error = fmt.Errorf("Target in request line has invalid format")
var InvalidHttpVersionError error = fmt.Errorf("Http version in request line has invalid format")
