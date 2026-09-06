# derp — Notes for AI Agents

Error-handling library used by Emissary and every benpate/* package. Errors carry an HTTP-style status code, a `Location` string, a message, details, and a wrapped inner error. See [README.md](README.md) for the public API tour.

## Conventions consumers rely on

- **Location strings are `package.Type.Method`** (or `package.Function`), e.g. `service.Stream.Save`. Every constructor and `Wrap` call in consuming code passes one; keep the format when writing new call sites.
- **Error codes ARE HTTP status codes.** Constructors map one-to-one: `BadRequest`→400, `Unauthorized`→401, `Forbidden`→403, `NotFound`→404, `Conflict`→409, `Gone`→410, `Teapot`→418, `MisdirectedRequest`→421, `Validation`→422, `Internal`→500, `NotImplemented`→501, `BadGateway`→502, `Timeout`→524. Pick by meaning: a duplicate-key write is `Conflict` (409), a misbehaving upstream server is `BadGateway` (502) — never a generic 500, because callers branch on `derp.ErrorCode` and `derp.Is*` predicates.
- **`ErrorCode(err)` reads any error.** It returns the `ErrorCodeGetter` value if the error implements it, 0 for nil, and 500 for any plain error. To change a code, use the `WithCode`/`WithConflict`/etc. `Option` values — any `Option` passed in the variadic `details` position is applied to the error instead of being stored as a detail. There is no `SetErrorCode` function; options are the only mutation path.
- **`Wrap` bubbles the inner code up.** The wrapper's `Code` is `ErrorCode(inner)`, so the root cause's status survives any number of wraps unless a `WithCode` option overrides it. `Wrap` ALWAYS returns non-nil (even for nil inner); `WrapIF` (note the capitalization — not `WrapIf`) returns nil when inner is nil and is the right choice for plain `return derp.WrapIF(err, ...)` tails.
- **`Report(err)` is for fire-and-forget paths only.** Code that has a caller returns the error; code with nowhere to return it (goroutines, deferred cleanup via `ReportFunc`, queue tasks) calls `derp.Report`. `ReportAndReturn` exists for paths that must do both. Never `Report` an error you are also returning up a stack that will report it again.
- **Constructors return `Error` by value, not `*Error`.** The whole package is value-typed; `errors.As` works against both `derp.Error` and `*derp.Error`, and `Wrap`'s type switch deliberately handles both (`case Error, *Error`).

## Foot-guns

- **`Wrap` must clip before appending to `details`.** `details` may alias the caller's slice (passed with `...`); the `append(details[:len(details):len(details)], ...)` full-slice expression prevents writing into the caller's spare capacity. Do not "simplify" it to a bare append, and copy the same pattern in any value-receiver or variadic builder here.
- **`GetRetryAfter` clamps to zero and defaults to 1 hour.** `nonNegative` exists because `Retry-After: -30` and already-past reset timestamps are real inputs; the trailing `http.ParseTime` attempt covers RFC850/asctime but must stay AFTER the explicit RFC1123 branch — `http.ParseTime` rejects RFC1123 dates with non-GMT zone abbreviations, so replacing that branch would be a regression. When no header parses, the method returns `time.Hour`, not 0.
- **`HTTPError` shares live header maps and Reporting one logs them.** `NewHTTPError` stores `request.Header`/`response.Header` by reference, and the default JSON plugin marshals the whole struct to the console — `Authorization`, `Cookie`, and HTTP `Signature` headers included. Redact or drop sensitive headers before `Report`ing an HTTPError built from an authenticated request; changing the struct to redact automatically changes serialized output for every consumer, which is why it doesn't.
- **Runtime reporter changes go through `derp.SetPlugins(...)`.** The registry is a copy-on-write `atomic.Pointer` list: readers are wait-free and a published slice is never mutated. `Set` swaps the whole list in one store; `Clear`-then-`Add` is safe but leaves a window where concurrently reported errors hit an empty list and vanish. Never append to the slice returned by ranging — it is shared and immutable.
- **Nil checks use `IsNil`/`NotNil`, not `err == nil`.** They catch typed-nil pointers inside the error interface, which would otherwise panic when `Error()` is called. The `IsNil(err) || err == nil` double checks scattered through the package exist to satisfy nilaway; keep them.

## Deliberate behavior that looks wrong

- **`Validation(message, ...)` takes no location**, so `Error()` renders as `": message"`. Fixing the format would change every consumer's log output; leave it.
- **`IsNotFound` / `IsNotFoundOrGone` also match the message text `"not found"`** (case-insensitive), not just code 404, because several database drivers return a bare "not found" with no code.
- **`Unwrap`, `RootMessage`, `RootLocation`, and `UnwrapHTTPError` recurse without a depth bound.** A cyclic error chain would blow the stack; no such chain exists in-tree, and no guard is wanted.
- **The `*Error` constructor duplicates (`BadRequestError`, `NotFoundError`, ...) are deprecated aliases** kept for compatibility. New code uses the short names; don't delete the aliases.
- **`go.mod` requires Go 1.19** because `atomic.Pointer` is a 1.19 API. Don't lower the floor.
