# Task K:

## Part 1: Implementation
Here we are going to be writing our very own Buildpack using the `go` programming 
language.

In these repository there are `build` and `detect` directories.
The code in each of these is responsible for building the `detect` and `build` 
binaries for our buildpack.

There are two test files that contain a list of tests at
`./detect/detect_test.go` and `./build/build_test.go`. 
Please read these test files, they have some context about what you are attempting
to implement.

Complete the `./detect/detect.go` and `./build/build.build.go` implementation
and get all tests to pass.

These tests may be run using the following from the root of this repository
`go test ./detect/... -count=1 -v`
`go test ./build/... -count=1 -v`

If you get these tests to pass. You will have a bare bones buildpack!

Useful resources:
- [buildpack specification](https://github.com/buildpacks/spec/blob/main/buildpack.md)
- [golang TOML library](https://godoc.org/github.com/BurntSushi/toml)
- [golang JSON library](https://golang.org/pkg/encoding/json/)

It is recommended that you start with the `detect` tests.

## Part 2: Validation
Upon completion the `./scripts/package.sh` will compile the `detect` and `build` binaries
into a file at `./buildpack.tgz`

We can then use this buildpack in a `pack build` by running

```
pack build <name-of-app-image> -p <path_to_onboarding_sample_app> --buildpack <path_to_buildpack.tgz>
```

Run our application image to make sure that it is working:
```
docker run -d --rm -p 8080:8080 my-buildpack-test 'node server.js'
```

you should now be able to curl our application and get some output
```
curl localhost:8080
```
or visit `localhost:8080` in your browser

## Cleanup

Awesome job completing Task K!

- please remeber to kill the container running your application 
by running `docker kill <container-id>`





NOTES: 
More on Provides & Requirements
In order for detection to pass the Provides & Requirements must 'match'
So every 'name = "<dep_name>"', entry under the [[provides]] section of the 
BuildPlan must have a equivalent 'name = "<dep_name"' under the [[requires]]

These matches need not all be included in the same buildpack, so for example
if we have BuildpackA, and BuildpackB
then we could have the following provides & require
