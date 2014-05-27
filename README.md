[![Stories in Ready](https://badge.waffle.io/go-chef/gladius.png?label=ready&title=Ready)](https://waffle.io/go-chef/gladius)
[![Stories in Ready](https://badge.waffle.io/go-chef/gladius.png?label=ready&title=Ready)](https://waffle.io/go-chef/gladius)
# gladius
A set of chef CLI utilities akin to 'knife'  written in go

## Rational
While working on workflows for a couple different companies around chef Some common things repeated.
  *  Tooling the workstation is hard: Platforms, Permissions, Corp Policy, and Dep Hell all made getting up and running difficult. `ChefDK` is a step towards making this better. Golang has the best story here IMO.
  *  Tools that expect you to work "their" way: Things written with an expectation about how your CI/cook workflow/Org worked. I want Unix tools that let you compose your workflow through pipes.
  *  Repeatability: Deploying a cookbook and _knowing_ it will work is quite difficult.
  *  Cookbook solvers assume 1 run_list: Berks solves for a cookbook, but not for an infrastructure. Many dep trees need to be solved to represent an infrastructure. It also uses its own metadata vs using the cookbook or server as source of truth. This is not optimal/simple.
  * performance: waiting for knife for ever to do something sucks. the tooling/ruby world doesn't do concurency/threadding well. Go is well suited for these things.
  * Wanted to build stuff with Golang, and I had these qualms with the current chef tools ecosystem ;)

## Current Status: _Embryo_
Right now I am using this as an implementation of the chef-golang Client API. This is driving my contributions to that design.

## Plans
I've started with a simple downloader that can use a chef-server as a source and slurp cookbooks into a directory from that server faster than berks can read it's cache: [see this gist](https://gist.github.com/spheromak/950fe653bd7b4bc044f8)

I need to make this more generally useful and add it's corollary, upload. right now those act on cookbooks, but I think generally these are generic upload/download on chef types. i.e. role/env/etc. and should handle STDIN/OUT in the form of json or something else that represents these objects.

Example:
Ideally these would be sinlge
````
$ solve -s mychef_server | download -- | upload --server someother_server
````
I think this sort of tooling has some serious potential.
````
$ solve --run_list 'base@1.2.3,app@3.2.1' --engine internal --type node --format json > node.json  
````

These are some usage ideas I've been playing with. Interested in more feedback on other ideas / usage people would like to see.

I am also interested in building out the new chef policy-file stuff in this sort of tooling.
