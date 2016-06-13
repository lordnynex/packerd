# packerd

## Historical Summary
Up to this point the [https://github.com/capitalone/ops-pipeline](ops pipeline) was intended for use via command line, possibly triggered via Jenkins, Atlas, or Travis-ci.

## Current Implementation, packer via REST
It can allow unit testing changes, baking a full application stack into a set of images, and validating each image provides expected services when operating. 

An ops-pipeline implementation is a git repository that has
* a packer config
* a set (possibly empty) of config management files
* a set (possibly empty) of test kitchen tests
* If the following files exists, we'll run some extra code
  * "Berksfile" exists then, berks vendor will be run to download the specified recipes (more info at [http://berkshelf.com/](Berkshelf))
  * ".kitchen.yml" exists then, test kitchen will be run (see [http://kitchen.ci](Chef test kitchen) for more info)
  * "gen-kitchen-dockerfile.sh" exists then, we will run that script via bash to generate a dockerfile template for the kitchen-docker gem (see [https://github.com/test-kitchen/kitchen-docker](kitchen-docker)

Passing a reference to this repository (and a bit more) to our bakery REST endpoint will result in
* a clone of the code repository
* a checkout of a specific branch/tag/commit
* berks vendor run 
* a packer build in the resulting tree, with parameters passed in via the REST API
  * the packer config allows provisioning via chef and shell.  Ansible is planned
  * the packer config may deposit built images in repositories as specified in the config
* if there is a kitchen test, those will be run
* calls may be made to a different endpoint for checking the status of the build

