module github.com/DataDog/aggregated-dependency-score/cmd/depscore

go 1.23.2

require (
	github.com/DataDog/aggregated-dependency-score v0.0.0-unpublished
	github.com/DataDog/aggregated-dependency-score/depsdotdev v0.0.0-unpublished
)

require (
	deps.dev/api/v3 v3.0.0-20241010035105-b3ba03369df1 // indirect
	golang.org/x/net v0.28.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
	golang.org/x/text v0.17.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20241021214115-324edc3d5d38 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241015192408-796eee8c2d53 // indirect
	google.golang.org/grpc v1.67.1 // indirect
	google.golang.org/protobuf v1.35.1 // indirect
)

// see https://go.dev/doc/modules/managing-dependencies#local_directory
replace (
	github.com/DataDog/aggregated-dependency-score => ../..
	github.com/DataDog/aggregated-dependency-score/depsdotdev => ../../depsdotdev
)
