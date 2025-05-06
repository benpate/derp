# DERP 🤪

[![GoDoc](https://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://pkg.go.dev/github.com/benpate/derp)
[![Version](https://img.shields.io/github/v/release/benpate/derp?include_prereleases&style=flat-square&color=brightgreen)](https://github.com/benpate/derp/releases)
[![Build Status](https://img.shields.io/github/actions/workflow/status/benpate/derp/go.yml?branch=main)](https://github.com/benpate/derp/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/benpate/derp?style=flat-square)](https://goreportcard.com/report/github.com/benpate/derp)
[![Codecov](https://img.shields.io/codecov/c/github/benpate/derp.svg?style=flat-square)](https://codecov.io/gh/benpate/derp)

## Better Error Reporting for Go

Derp is a drop-in replacement for the default error objects, and can be used anywhere that expects or requires an error value.  It enhances Go's default with additional tracking codes, error nesting, and plug-ins for reporting errors to external sources.

## 1. More Informative Errors

Derp encapulates all of the data you can collect to troubleshoot the root cause of runtime errors.  Here's a quick look at each argument.

* **Location** The location where the error took place, typically the name of the package and function
* **Message** A human readable description of the error
* **Code** A custom error code for tracking exactly what error occurred.
* **Error** Nested error that lets you see down the call stack
* **Details** Variadic of additional parameters that may be helpful in debugging this error.

```go

func InnerFunc(arg1 string) error {

    if err := doTheThing(); err != nil {
        // Derp create errors with more troubleshooting info than standard errors.
        return derp.NotFoundError("App.InnerFunc", "Error doing the thing", err.Error(), arg1)
    }

    return nil
}

func OuterFunc(arg1 string, arg2 string) error {

    // Call InnerFunc with required arguments.
    if err := InnerFunction(arg1); err != nil {

        // Wraps errors with additional details and nested stack trace.
        return derp.Wrap(err, "App.OuterFunc", "Error calling InnerFunction", arg1, arg2)
    }

	return nil
}

```

## 2. Nested Errors

Derp lets you include information about your entire call stack, so that you can pinpoint exactly what's going on, and how you got there.  You can embed any object that supports the `Error` interface.

### Error Codes

Every derp erro is defined with a specific error code, corresponding to the standard [HTTP status codes](https://www.rfc-editor.org/rfc/rfc9110.html#name-status-codes).  These are created with helper functions such as `InternalError` and `NotFoundError`.To help you dig to the original cause of the error, nested error codes will "bubble up" from the original root cause, unless you specifically override them.

## 3. Reporting Plug-Ins

The derp package uses plugins to report errors to an external source.  Plugins can send the error to the error console, to a database, an external service, or anywhere else you desire.

Plugins should be configured once, on a system-wide basis, when your application starts up.  If you don't set up any plugins, then the default setting is to report errors to the system console.

```go
import "github.com/benpate/derp/plugins/mongodb"

func init() {

    // By default, derp uses the ConsolePlugin{}.  You can remove
    // this default behavior by calling derp.Plugins.Clear()

    // Add a database plugin to insert error reports into your database.
    derp.Plugins.Add(mongodb.New(connectionString, collectionName))
}

func SomewhereInYourCode() {
    // Report passes the error to each of the configured
    // plugins, to deliver the error to its destination.
    derp.InternalError("location", "description", 0, nil).Report()
}
```

### Plug-Ins

The package includes a default reporter, and you can add to this list easily using `derp.Plugins.Add()` to add any object that implements the `Plugin` interface at startup.

* `Console` write a human-friendly error report to the console (this package)
* [`derp-mongo`](https://github.com/benpate/derp-mongo) writes error reports to a MongoDB database
* [`derp-zerolog`](https://github.com/benpate/derp-zerolog) writes error reports to the [zerolog](https://github.com/rs/zerolog) logging package

## Pull Requests Welcome

Original versions of this library have been used in production on commercial applications for years, and the extra data collection has been a tremendous help for everyone involved.

I'm now open sourcing this library, and others, with hopes that you'll also benefit from a more robust error package.

Please use GitHub to make suggestions, pull requests, and enhancements.  We're all in this together! 🤪
