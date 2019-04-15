# derp 🤪

## Better error reporting for Go

Derp is a drop-in replacement for the default error objects, and can be used anywhere that expects or requires an error value.  It enhances Go's default with additional tracking codes, error nesting, and plug-ins for reporting errors to external sources.

## Better Error Tracking


* **Location** The location where the error took place, typically the name of the package and function
* **Code** A custom error code for tracking exactly what error occurred.
* **Message** A human readable description of the error
* **Error**
* **Details**
```go

func MyFunction(arg1 string, arg2, arg3) error {

	if err := doTheThing(); err != nil {
		return derp.New("AppName.MyFunction", "Error doing the thing", err, arg1, arg2, arg3)
	}
```

## Nested Errors

Errors usually 

## Error Reporting
The derp package uses "reporters" to report errors to an external source.  These may be to th error console, to a database, or an external service.  The package includes a small number of default reporters, and you can add to this list easily by `Connect()`-ing an object that implements the `Reporter` interface at startup.

* `Console` writes a human-friendly error report to the console
* `Loggly` sends error reports to the Loggly web service
* `Mongodb` writes errors to a MongoDB database collection
* `SendGrid` sends a human-friendly error report via SendGrid email service.

## HTTP Error Codes
Using HTTP error codes with derp is recommended, but not required.  This standard includes useful codes for most errors that occur on a server, and makes it easy to use derp errors directly in HTTP server responses.