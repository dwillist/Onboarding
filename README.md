# Task 1: Using paketo buildpacks.

## First Onboarding task
Greetings human

The first thing we are going to do is just a quick exploration of how to use the buildpacks the
Paketo project produces.

## Prerequisites

For this you will need a couple additional peices of software
 - [pack](https://buildpacks.io/docs/install-pack/)
   - This is the CLI that orchestrates the running of each Paketo buildpack
 - [docker](https://docs.docker.com/get-docker/)
   - Provides a image registry on all platforms.
 - [sample_application](https://github.com/dwillist/onboarding_application)
   - just a simple application nodejs application we are going to build 
   This app will be used throughout this tutorial so it is recommended that you use it


## First task

We are going to is just build the same application into an app-image using `pack` in three different ways
1. Build using a `builder`
1. Build using a `metabuildpack`
1. Build using an `implementation buildpack`


###  Using a Builder
Given that you have installed the above pre-requisites and started the docker daemon. We can now set up pack 
to use a builder!

list all the pre-packaged builders 

```
pack list-trusted-builder
```

For this tutrorial we are going to use the `gcr.io/paketo-buildpacks/builders:base` builder.
So to set this as the default we run

```
pack set-default-builder gcr.io/paketo-buildpacks/builders:base
```

Ok great! Now from the root of the `sample_application` cloned as a prerequisite simply run:
```
pack build onboarding-test-image
```

After a bit of output our build succeeds and we have an actual application image. Running 
```
docker images
```
lists all images and we will see that indeed the `onboarding-test-image` is present.

#### Quick Questions & A
Ok so a lot happend in that last step lets go through some diagrams and deffinitions to show how all these piece fit togeather at a high level to produce our final `onboarding-test-image` image.


##### What is an image?

For our purposes an app image image is just a collection on `layers` or which are just addition or deletions to the filesystem. 

The images we will be producing have three distinct types of layers,
- a layer for our application source (for interpreted languages)
- layer(s) for our application's dependencies 
- an OS layer that contains operating system packages.

TODO: Layers graphic here

For a closer look at the contents of each layer try using the [`dive`](https://github.com/wagoodman/dive) tool.

A General picture about how these interact:

TODO: graphic around pack + buildpacks = sample app


##### What is a Builder

Simply a builder is a collection of buildpacks, some data indicating what order they should be run in, and a `stack`, which provides the OS packages.


We can get some more information about exactly what buildpacks and what OS packages our Builder is going to provide by running
```
pack inspect-builder
```

This will give some structured output that gives us the `buildpacks` on the builder and the `stack`!
```
Stack = 'io.buildpacks.stacks.bionic'

Lifecycle = "don't worry about this yet"
...

Run Images = "don't worry about this yet"
...

Buildpacks = "long list of all buildpacks in this builder"
...

Detection Order = "don't worry about this yet"
...
```

We can expand a bit on our above graphic.

The buildpacks provide you application dependencies, while the `stack` is what provides the OS packages

TODO: expanded graphic showing stacks and buildpacks.

##### Summary

So in summary we see that the default builders that come with `pack` come with buildpacks and a stack. These are responsible for all the layers below your `app` layer. Which is why you only need to bring your app and you can just `pack build` an application image.


Now lets move onto another way of `pack building our application`

### Using a metabuildpack
Continuing where we left off after building our application using a `builder`. We can also specify a `metabuildpack` or a group of buildpacks that we would like to use. When building our app using a `builder` we used 2 buildpacks `node-engine` & `npm`.

This time lets build our application using the `gcr.io/paketo-buildpacks/nodejs` metabuildpack, which contains the `node-engine`, `npm` and `yarn` implementation buildpacks. 

this is as simple as sitting in the root of the `sample_application`
```
pack build metabuildpack-build-test --buildpack gcr.io/paketo-buildpacks/nodejs
```


#### Using implementation buildpacks

And finally notice we used the `node-engine` and `NPM` buildpacks to build our application we can specify just those buildpacks as follows
```
pack build implementation-build-test --buildpack gcr.io/paketo-buildpacks/nodejs --buildpack gcr.io/paketo-buildpacks/npm
```








