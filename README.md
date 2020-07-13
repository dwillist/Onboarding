# Task 1 Using paketo buildpacks.

## First Onboarding task
Greetings human

The first thing we are going to do is just a quick exploration of how to use the buildpacks the
Paketo project produces.

### Prerequisites

For this you will need a couple additional peices of software
 - [pack](https://buildpacks.io/docs/install-pack/)
   - This is the CLI that orchestrates the running of each Paketo buildpack
 - [docker](https://docs.docker.com/get-docker/)
   - Provides a image registry on all platforms.
 - [sample_application](https://github.com/dwillist/onboarding_application)
   - just a simple application nodejs application we are going to build 
   This app will be used throughout this tutorial so it is recommended that you use it


### First task

We are going to is just build the same application into an app-image using `pack` in three different ways
1. Build using a `builder`
1. Build using a `metabuildpack`
1. Build using an `implementation buildpack`


####  Using a Builder
Given that you have installed the above pre-requisites and started the docker daemon. We can now set up pack 
to use a builder!

run `pack list-trusted-builder` to show a list a pre-packaged builders

For this tutrorial we are going to use the `gcr.io/paketo-buildpacks/builders:base` builder.
So to set this as the default we run

`pack set-default-builder gcr.io/paketo-buildpacks/builders:base`

Ok great! Now from the root of the `sample_application` cloned as a prerequisite simply run
`pack build onboarding-test-image`

After a bit of output our build succeeds and we have an actual application image. Running `docker images`
lists all images and we will see that indeed the `onboarding-test-image` is present.

#### Q & A
Ok so a lot happend in that last step lets go through some diagrams deffinitions to show how all the pieces
fit togeather to give us our final `onboarding-test-image` image.


##### What is an image?



