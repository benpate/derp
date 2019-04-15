# derp 🤪

## Better error reporting for Go

[![Sourcegraph](https://sourcegraph.com/github.com/benpate/derp/-/badge.svg?style=flat-square)](https://sourcegraph.com/github.com/benpate/derp?badge)
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/benpate/derp)
[![Go Report Card](https://goreportcard.com/badge/github.com/benpate/derp?style=flat-square)](https://goreportcard.com/report/github.com/benpate/derp)
[![Build Status](http://img.shields.io/travis/benpate/derp.svg?style=flat-square)](https://travis-ci.org/benpate/derp)
[![Codecov](https://img.shields.io/codecov/c/github/benpate/derp.svg?style=flat-square)](https://codecov.io/gh/benpate/derp)

Derp is a drop-in replacement for the default error objects, and can be used anywhere that expects or requires an error value.  It enhances Go's default with additional tracking codes, error nesting, and plug-ins for reporting errors to external sources.

## Better Error Tracking

Derp encapulates all of the data you can collect to troubleshoot the root cause of runtime errors.  Here's a quick look at each argument.

* **Location** The location where the error took place, typically the name of the package and function
* **Message** A human readable description of the error
* **Code** A custom error code for tracking exactly what error occurred.
* **Error** Nested error that lets you see down the call stack
* **Details** Variadic of additional parameters that may be helpful in debugging this error.
```go

func TopLevelFunction(arg1 string, arg2 string arg3 string) {

	if err := InnerFunction(arg1, arg2, arg3); err != nil {
		// Wraps the inner error with additional details, and reports it to ops.
		derp.New("AppName.TopLevelFunction", "Error calling InnerFunction", 0, err).Report()
	}
}

func InnerFunction(arg1 string, arg2 string, arg3 string) error {

	if err := doTheThing(); err != nil {

		// Create a derp error with additional troubleshooting details.
		return derp.New("AppName.MyFunction", "Error doing the thing", derp.CodeNotFound, err, arg1, arg2, arg3)
	}

	return nil
}
```

## Nested Errors

Derp lets you include information about your entire call stack, so that you can pinpoint exactly what's going on, and how you got there.  You can embed any object that supports the `Error` interface.

### Nested Error Codes

Every error in derp includes a numeric error code.  We suggest using standard **HTTP status codes**, but you can return any number that works for you.  To help you dig to the original cause of the error, nested error codes will "bubble up" from the original root cause, unless you specifically override them.

To set an error code, just pass a **non-zero** `code` number to the `derp.New` function.  To let underlying codes bubble up, just pass a **zero**.

## Error Reporting 
The derp package uses plugins to report errors to an external source.  Plugins can send the error to the error console, to a database, an external service, or anywhere else you desire.

Plugins should be configured once, on a system-wide basis, when your application starts up.  If you don't set up any 

```go
import "github.com/benpate/derp/plugins/console"
import "github.com/benpate/derp/plugins/mongodb"

func init() {
	// Send all errors to console
	derp.Connect(console.New())

	// Send all errors to database
	derp.Connect(mongodb.New(connectionString, collectionName)) 
}

func later() {
	// Report passes the error to each of the configured
	// plugins, to deliver the error to its destination.
	derp.New("location", "description", 0, nil).Report()
}
```

### Default Plug-Ins
The package includes a small number of default reporters, and you can add to this list easily by `Connect()`-ing an object that implements the `Reporter` interface at startup.

* `Console` write a human-friendly error report to the console

Future plugin development:
* `Mongodb` write errors to a MongoDB database collection
* `SMTP` send a human-friendly error report via email.
* `Loggly` sends error reports to the Loggly web service
