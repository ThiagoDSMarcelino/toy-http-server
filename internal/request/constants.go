package request

// Theoretically, \n could be the separator as well if the first line ends with it, but for now we only support \r\n.
var SEPARATOR = []byte("\r\n")
