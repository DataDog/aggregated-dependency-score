# Aggregated Dependency Score: A Dependency Score that Includes Transitive Dependencies

Implementations of the algorithm described in [this blog post: "Aggregated Dependency Score: A Dependency Score that Includes Transitive Dependencies"](https://cedricvanrompay.fr/blog/aggregated-dependency-score/)

> [!WARNING]
>
> Work in progress: the code in this repository is only suitable for testing as of October 2024


## Python Implementation

The Python implementation is in [the `python` directory](./python/).

As of October 2024 it implements an older version of the aggregated dependency score;
it must be updated to match the algorithm that ended up being published.


## Go Implementation

As it is usually done with Go code, the main module is located at the root of the repository.
As a result, most of what's _not_ in [the `python` directory](./python/)
should be considered as part of the Go implementation
unless explicitely mentionned otherwise.

Note that module `github.com/DataDog/aggregated-dependency-score`
does not provide dependency enumeration or intrinsic score evaluation;
the user must provide its own.
Right now, a function `NewDepsDotDevClient` provides an client for https://deps.dev
that can do both, but in the future this client will be moved to a separate module
(and most likely a separate repository)
so that users not wanting to use it don't have to import its dependencies.

To see an example of how to use the package, see [`cmd/depscore`](./cmd/depscore/):

```
$ go run ./cmd/depscore --ecosystem pypi --package requests --version 2.28.1
0.18347983371997253
```
