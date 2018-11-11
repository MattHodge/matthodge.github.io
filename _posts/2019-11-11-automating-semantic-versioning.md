---
layout: post
title: Automating Semantic Versioning For Any Project
date: 2018-11-10-T13:37:00.000Z
comments: false
description: A guide to automating semantic versioning of any project in Git with the semantic-release tool.
---

When are developing a tool or application for yourself, the furthest thing from your mind is "how should I version this?". You are the only consumer, so what does it matter! You know what is changing, and you will always just use the latest version.

This all changes as soon as other people or applications start to depend on your application. Changes you make can have a massive impact, and even break your application for people that depend on it.

Let's use my [Octopus Deploy Terraform Provider](https://github.com/MattHodge/terraform-provider-octopusdeploy) as an example. One of the `resources` in the provider allows people to create an *environment*.

For my end user, they use the resource like this:

```hcl
# main.tf - Creates the users production environment in Octopus Deploy

# Configure the provider
provider "octopusdeploy" {
  address = "http://octopus-deploy:8081/"
  apikey  = "API-XXXXXXXXXXXXXXXXXXXX"
}

# Create the produduction environment
resource "octopusdeploy_environment" "production" {
    name             = "Production"
    description      = "The production environment. If this breaks we are in trouble."
    useguidedfailure = "true"
}
```

The user stores the above Teraform code it in their git repository, and feel great knowing they are automating their processes.

A Terraform resource and its arguments is the contract between my provider and the user. They use the contract and they get the resource they want, in this case, an *environment*.

The user is happy 😁!

I then decide I would like to standardize on [snake case](https://en.wikipedia.org/wiki/Snake_case) for all of my argument names in my [Octopus Deploy Terraform Provider](https://github.com/MattHodge/terraform-provider-octopusdeploy). I go and modify my code and change the argument `useguidedfailure` to `use_guided_failure` instead. I share the new version of the provider to the world!

Unfortunately he next time my unlucky user pulls down the latest version of my provider and tries to run a `terraform plan` he is going to have a bad time 😟:

`Error: octopusdeploy_environment.production: : invalid or unknown key: useguidedfailure`

What I have done here is make a **breaking change** or an **incompatible change** in the contract I have with the user. Their code is now broken due to my change.

Let's see if we can use semantic versioning and automation to help us improve this situation.

> :loudspeaker: Want to know when more posts like this come out? [Follow me on Twitter: @MattHodge](https://twitter.com/matthodge) :loudspeaker:

* TOC
{:toc}

## What is semantic versioning?

[Semantic versioning](https://semver.org/) is best known and most widely adopted convention for versioning software. It uses a sequence of 3 digits to represent the version of the software. If you haven't heard of it I recommend you [read up on it before continuing](https://semver.org/).

## Applying semantic versioning

Semantic versioning tells us if we make a breaking change to the API like we did in our example in the introduction, we need to bump our *MAJOR* version number (so from **1**.0.0 to **2**.0.0).

I want to do better. I should properly version the releases of my provider, but this means a ton of manual work not only now, but every time I release a new version.

I'm going to need to:

* Manage and think about when I should bump versions according to [semantic versioning](https://semver.org/)
* Think about which part of my version to bump (major, minor or patch)
* Start using [Github Releases](https://help.github.com/articles/creating-releases/) to store the versions I release
* Keep a [changelog](https://keepachangelog.com/en/1.0.0/) so my users know what is different between versions of my provider

I really do not like doing things manually. Let's use something that meet all of my requirements above, and the additional one that I am lazy.

🤫 P.S. My users should probably also be using [Terraform provider versioning](https://www.terraform.io/docs/configuration/providers.html#provider-versions) too! But that's out of scope for this post.

## Automating using semantic-release

[Semantic-release](https://github.com/semantic-release/semantic-release) is a tool which can help us automate our whole release workflow, including determining the next version number, generating the release notes and publishing the package.

It does with the use of several plugins.

### Plugin - commit-analyzer

The [commit-analyzer](https://github.com/semantic-release/commit-analyzer) plugin parses commit messages and determines if a new release should be made, and whether it should be a Major, Feature or Patch release.

The default configuration is to use the [AngularJS](https://github.com/angular/angular.js/blob/master/DEVELOPERS.md#-git-commit-guidelines) commit message format. Let's have a look at some example commit messages:

```
fix: fix timeout when creating environments
```

The above commit message would be a *Patch Release*, for example 1.0.**1**.

```
feat(resource): add new argument to environment resource
```

The above commit message would be a *Feature Release*, for example 1.**1**.0.

You can also include a optional *section* of the code you changed in parentheses, in this case I am changing a `resource`.

```
refactor(resource): rename useguidedfailure argument on environment resource

BREAKING CHANGE: The useguidedfailure argument has been renamed use_guided_failure in the octopusdeploy_environment resource.
```

The above commit message would be a *Major Release* as it has a breaking change, for example **2**.0.0.

If you like, you can also [implement your own](https://github.com/semantic-release/commit-analyzer#configuration) commit parsing methods.

### Plugin - release-notes-generator, changelog & git

The [release-notes-generator](https://github.com/semantic-release/release-notes-generator) will generate a changelog based on the commit messages. You can then use the [changelog](https://github.com/semantic-release/changelog) plugin to create or update a `CHANGELOG.md` file.

Using the `refactor` commit above as an example, the following changelog would be generated.

![CHANGELOG.md Example](images/posts/automating-semantic-versioning/changelog_example.png)

This generated `CHANGELOG.md` needs to be checked back into the git repository, which is where the [git](https://github.com/semantic-release/git) plugin comes in.

### Plugin - github

The [github](https://github.com/semantic-release/github) plugin can create [GitHub releases](https://help.github.com/articles/about-releases) and comment on issues the release closes them.

Using the `refactor` commit above as an example, the following GitHub release would be created.

![GitHub Release example](images/posts/automating-semantic-versioning/github_release_example.png)

### Plugin - last-release-git

The [last-release-git](https://github.com/finom/last-release-git) extracts the latest version of your software from git tags. Usually, semantic-release is used for releasing NPM packages. We will need this plugin so the previous version of our release can be read from our git tags instead.

## Installation

Now that you have an idea of how [semantic-release](https://github.com/semantic-release/semantic-release) and its plugins work, its time to set them up on our repository.

Most of the documentation for [semantic-release](https://github.com/semantic-release/semantic-release) is very specific to NPM packages, but it can work with any type of repository and software.

In our case, want to get it enabled for a Terraform Provider, which is written in Go.

For the rest of this guide I will assume you have some familiarity with Node.js and have it installed.

Install [semantic-release](https://github.com/semantic-release/semantic-release) and the plugins we need via NPM. We will install them globally.

```bash
npm install -g semantic-release               \
    @semantic-release/changelog               \
    @semantic-release/git                     \
    @semantic-release/commit-analyzer         \
    @semantic-release/release-notes-generator \
    last-release-git
```

## Configuration

With the tool installed, now we need to setup a configuration file. Create a `.releaserc` file in the root of your repository.

The following is configuration for my Terraform provider repository:

```json
{
    "branch": "master",
    "plugins": [
        "@semantic-release/commit-analyzer",
        "@semantic-release/release-notes-generator",
        [
            "@semantic-release/changelog",
            {
                "changelogFile": "CHANGELOG.md",
                "changelogTitle": "# Semantic Versioning Changelog"
            }
        ],
        "@semantic-release/github",
        [
            "@semantic-release/git",
            {
                "assets": [ "CHANGELOG.md" ]
            }
        ]
    ],
    "release": {
        "getLastRelease": "last-release-git"
    }
}
```

Let's run through the main parts of the configuration file.

* `branch` chooses which git branch to create releases from.

* `plugins` an array of plugins to load. The plugins will be executed in the order they are defined. You can also define some configuration for each plugin as I have done above.

* `release` sets the method used to determine the latest release for the project. We will use `last-release-git` for this.

You can read me detailed documentation in the [configuration file](https://semantic-release.gitbook.io/semantic-release/usage/configuration) page of the semantic-release documentation.

## Create a GitHub Token

We need to create a GitHub token to allow semantic-release to create GitHub releases.

Follow the steps [here](https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/) to create one.

You will need to give the following permissions:

![GitHub Permissions](images/posts/automating-semantic-versioning/github_token_permissions.png)

Copy the token and set it as the `GITHUB_TOKEN` environment variable.

```bash
export GITHUB_TOKEN=PUT-YOUR-TOKEN-HERE
```

## Running locally

Ideally, you will run semantic-release inside a continuos integration (CI) tool. Let's test it out locally first though.

I opened a repository of mine and ran the following command:

```bash
npx semantic-release

# You can also run it with the --debug flag for more details.
```

When running locally, it will detect that its not running in a CI tool and run in *dryrun* mode.

You will see all the configuration load and an output like this:

![Semantic Release No Change](images/posts/automating-semantic-versioning/semantic-release-output-no-change.png)

Let's create a commit and see what happens.

```bash
# Add a test file and commit
touch testfile.txt
git add testfile.txt
git commit -m "feat: add a test file"

# Rerun
npx semantic-release
```

This time we see the changelog is generated and a new release would be created:

![Semantic Release With Change](images/posts/automating-semantic-versioning/semantic-release-output-with-change.png)

😍 Pretty awesome right!

Let's say I already had a version that was released of my project using my old manual method, how would I have semantic-release continue on from there?

```bash
# Find the hash of the latest change
git log -2 --oneline

# Returns:
# f3cd02e (HEAD -> master) feat: add a test file
# 6548458 (origin/master, origin/HEAD, 1541972003) Update Builds

# Add a tag to the commit before our change, lets say it was version 1.5.0
git tag v1.5.0 6548458
```

Now we are in-line with our previous release 😎.

![Semantic Release With Change](images/posts/automating-semantic-versioning/semantic-release-output-with-change-specific-version.png)

