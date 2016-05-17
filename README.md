# packerd

## Historical Summary
Up to this point the (https://github.kdc.capitalone.com/internal-ops-pipeline)[ops pipeline] was intended for use via command line, possibly triggered via Jenkins, Atlas, or Travis-ci.

An ops-pipeline implementation is a git repository that has
* a packer config
* a set (possibly empty) of config management files
* a set (possibly empty) of test kitchen tests

Passing a reference to this repository (and a bit more) to our bakery REST endpoint will result in
* a pull of the code
* a packer build in the resulting tree
* if there is a kitchen test, those will be run
* if all works out, the resulting artifacts will be deposited into a repository as defined in the packer template

It can allow unit testing changes, baking a full application stack into a set of images, and validating each image provides expected services when operating. 

## Roadmap
Going forward, we'll be implementing a REST API to allow the above operations and more as needed.

Main work will be done as open source first, with integrations for internal work handled in a private repository.

## Swagger spec:
https://github.kdc.capitalone.com/kbs316/packerd/blob/master/swagger.yml

