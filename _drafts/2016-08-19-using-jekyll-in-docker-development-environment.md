---
layout: post
title: "Using Jekyll in Docker as Blog Development Environment"
date: 2016-08-19 13:37:00
comments: false
description: Five useful tips when writing DSC Resources in PowerShell 5. Covers folder structure, IntelliSense, verifying resources, testing using Pester and more.
---

I really love Jekyll, it has made blogging so much more enjoyable as I can spend more time just writing instead of struggling with formatting and CSS. Not having to worry about hosting thanks to GitHub page support for Jekyll is also a huge plus.

That being said, it isn't all roses, especially if you are trying to work with Jekyll from a Windows machine.

Getting Ruby and Jekyll setup properly is painful on Windows, so I ended up just using an Ubuntu Virtual Machine in the past as my blog development environment.

With the new Docker 1.12 support I thought I would try and re-work my development environment and to use Docker instead of a full Virtual Machine just for Jekyll.

## Installing Docker

Docker 1.12, brought much tighter integration with Windows and MacOS. In this example I will be installing Docker on my Windows 10.1 machine.

< link to someone installing Docker>

## Setting up an initial Jekyll GitHub Repo

If you are starting from scratch with Jekyll, use the following steps to prepare yourself a repository:

* Create folder
* Git Init
* Clone Jekyll Repo?


## Using an Existing Jekyll Repo

As I already have a Jekyll blog, I will simply do a storage mapping on the Jekyll container into the Jekyll git repo on my local machine.

## Exposing the Ports

To expose Jekyll to my host system, I need to do some host mappings...

## Auto reloading the Jekyll container on a change

While editing the blog, I wanted the container to be destroyed and re-created every time I saved a post I was working on. This would allow me to press save and just refresh my browser to see the changes in my demo environment.

...
