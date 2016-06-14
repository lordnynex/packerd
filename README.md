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

This repository implements the above pipeline, and once this service is boot strapped, it can, and does, build itself.

## Use cases

### build pipeline use of packerd

```bash
# do work in your git tree, push to public branch
git clone https://github.com/tompscanlan/packerd-demo.git
cd packerd-demo
git checkout -b demo
date > changes
git add changes
git commit -m"demo..."

# set to ip/hostname of the packer service
packerd_ip=192.168.136.137

# for convenience
curl="curl -s -k -X GET  https://$packerd_ip"

# check packerd health
$curl/health | jq .

# trigger build of pushed changes, ideally this is triggered via webhook from git
response=`curl -s -k -X POST https://$packerd_ip:443/build -H "Content-Type: application/json" -H "Cache-Control: no-cache" \
-d  '{"giturl": "https://github.com/tompscanlan/packerd-demo.git",
    "branch": "demo",
    "templatepath": "packer.json",
    "buildonly": "docker-ubuntu",
    "buildvars": [{"key": "version",
        "value": "0.21"
        },{"key": "docker_repo_username",
        "value": "tompscanlan"
        }, { "key": "docker_repo_password",
            "value": "****"
        },{ "key": "docker_repo_email",
            "value": "***@gmail.com"
        }]}'`

# get specific build request
req_url=$(echo $response  | jq -r ".[0].href")

# get all builds
$curl/build/queue | jq .

# get response to this specific build
resp_url=$($curl$req_url | jq -r ".responselinks[1].href")

# watch the packer stage of the build
watch "$curl$resp_url | jq '.status, (.buildstages | reverse)[0]'"
```


### developer using local Packerd service
When runing a local packerd service, you'll still need to push commits to a branch in a reachable git repo.

```bash
# do work in your git tree, push to public branch
git clone https://github.com/tompscanlan/packerd-demo.git
cd packerd-demo
git checkout -b demo
date > changes
git add changes
git commit -m"demo..."


# run a local instance of packerd
docker run -d --privileged -p 443:64155  tompscanlan/packerd:latest /usr/bin/supervisord

# set your local IP
packerd_ip=192.168.136.137

# follow same steps as for remote service above, but using local ip to reference docker instance of packerd
```

