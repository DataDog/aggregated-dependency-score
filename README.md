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
the user must provide its own. Module `github.com/DataDog/aggregated-dependency-score/depsdotdev`
is _one_ way of getting both based on https://deps.dev,
with the intrinsic score being [the OSSF Scorecard](https://securityscorecards.dev/)
computed and cached by https://deps.dev.

To see an example of how to use the packages,
see [`cmd/depscore`](./cmd/depscore/).

```
$ cd cmd/depscore
$ go run .
0.18347983371997253
```

> [!NOTE]
>
> Due to how the Go modules in this repository are organized,
> you _have_ to be located in [the `cmd/depscore` directory](./cmd/depscore/)
> to run this binary. If you try to run `go run ./cmd/depscore` from the repository root,
> you will get the following error:
> ```
> main module (github.com/DataDog/aggregated-dependency-score) does not contain package github.com/DataDog/aggregated-dependency-score/cmd/depscore
> ```
