---
layout: post
title: Building a Golden Image Pipeline
date: 2018-10-21T13:37:00.000Z
comments: true
description: A step by step guide to creating a golden image pipeline (base images) for your infrastructure using Packer and Ansible with Windows.
---

Welcome to this series of posts about creating golden images, and building a golden image pipeline.

In this post, we are going to start with some definitions and introduce some concepts around the creation of golden images.

Future posts will go in depth as to how to create a golden image pipeline.

The series will focus on using Packer, Ansible and Windows, but the content will be generic enough to apply to any platform or configuration management tool.

Let's dive in!

* TOC
{:toc}

## What is a golden image?

A golden image provides the template which a virtual machine (for example, [AWS EC2](https://aws.amazon.com/ec2/) instances) is created from. It may also be referred to as a base image or an image template. Think of it as a snapshot copy of an operating system that can be launched as a new virtual machine.

Usually, a golden image will contain:

* A Windows or Linux operating system installation

* The latest security patches and updates

* Configuration specific to your environment

* Software specific to your environment

* Security hardening settings, if required by your environment

* Agents such as an [Octopus Deploy Tentacle](https://octopus.com/) for deploying software, or a [Datadog Agent](https://datadoghq.com) for monitoring the virtual machine

The idea is that you set up an operating system to the desired state, save it and then you can re-use it across your infrastructure.

![Golden Images](/images/posts/building-golden-image-pipeline/golden_boxes_wide.jpg)

## Why would you want to create golden images?

Traditionally, you might be used to just inserting a CD-ROM with a copy of Windows and build and configure a server manually. You would run Windows updates, install your desired software and do what was needed to make the server "production ready". You might have a checklist that you go over to make sure setup was done correctly.

This process is fine when you are managing just a handful of servers, but as technology now underpins more of what a business does, the number of servers has grown dramatically.

Especially when moving to the cloud ☁️, manually spinning up servers and configuring them means you get zero of the benefits of what cloud can provide, no auto-scaling, no auto-healing, no high availability or resiliency to failure.

Having a golden image allows you to do the configuration work once, and use it across your entire infrastructure. It will save you time, make you faster and reduce human error.

## Golden images as code

You might have used a manual method of creating golden images in the past, for example, spinning up a virtual machine, making manual changes and then running Windows Sysprep, followed by creating a snapshot of the machine. This was a common approach with VMWare templates.

Any work done manually is hard to scale and can be error-prone, so defining your golden images in code is critically important.

[HashiCorp's Packer](https://www.packer.io/intro/) is the de facto standard tool for creating golden images from code, and we will be discussing this tool in depth in the next post.

In short, it allows you to define JSON configuration files which Packer will use to create a machine, apply configuration too it and then save the machine for use as a golden image. It works with many platforms including [VMWare](https://www.packer.io/docs/builders/vmware.html), [Google Cloud](https://www.packer.io/docs/builders/googlecompute.html), [Azure](https://www.packer.io/docs/builders/azure.html) and [AWS](https://www.packer.io/docs/builders/amazon.html).

![HasiCorp Packer](/images/posts/building-golden-image-pipeline/packer_logo.jpg)

## How many golden images should I have?

When starting to create golden images, it's a good chance to take stock of your current applications and what their requirements are. The goal here is to try to have *as few* golden images as possible, which can be used in 100% of your infrastructure.

For example, say you are running 4 different applications:

* A financial application that runs on Windows and requires IIS and DotNet Core 2.0

* Your companies site that runs on Windows and IIS and Requires .NET Framework 4.7.2

* A financial data processing application that runs several Windows Services and requires DotNet Core 2.0

* A Windows jump/bastion host that you use to manage your infrastructure

You could combine these into two golden image flavours:

* A Windows Core Server with .NET Core 2.0 and .NET Framework 4.7.2

* A Windows Core Server with .NET Core 2.0 and .NET Framework 4.7.2 and IIS

These images currently support all of your application's requirements.

As you get more and more applications depending on your golden images, see if an image you currently provide supports the application. If not, see if its possible to change an existing image to support the new application while keeping compatibility with the other applications on your infrastructure.

![Cookies](/images/posts/building-golden-image-pipeline/cookies.jpg)

## What should go into golden images?

You have two options when using golden images:

* Bake the software and configuration inside the golden image

* Install software and do the configuration after the golden image has been spun up and is used by a Virtual Machine

Try to do *as much as you can* when you are baking the golden image. This will speed up your instance startup time when the image is used, saving you from having to do large/slow installations once the machine has booted.

Some tips:

* Install all required agents (monitoring/deployment etc.) in the golden image, just apply configuration management when the instance starts to finalize the configuration

* *Do not* bake any secrets such as API Keys or Passwords into the image, use configuration management when the instance starts

* Install all software and Windows updates in the golden image

* Stop services after installation in the golden image, start them with configuration management if required

## What is a golden image pipeline?

You are unlikely to just create a golden image once. Likely there will need to be modifications as requirements change. There is also the need to get the latest updates and security patches into the image.

You also likely want to allow others to make changes to the golden images, such as in-house developers, allowing them to propose changes for approval by the team responsible for the images.

The best workflow for this is the "[GitOps](https://www.twistlock.com/2018/08/06/gitops-101-gitops-use/)" workflow.

* Store your golden image as code in a centralized git repository

* Allow pull requests to make changes to the golden images

* Automatically test and validate that the changes will work with the applications relying on them

* Automatically deploy the created images, so they can be consumed by applications

Additional to this, you likely want to schedule the creation of golden images on a weekly or monthly basis, so you can install the latest security patches as they come out.

To do all these things, you will want to use a build server such as TeamCity or Jenkins, or a hosted build service such as AppVeyor, CircleCI, AWS CodeBuild or Azure DevOps Pipelines.

Building golden images in one of these services will allow you to have a Continuous Integration and Continuous Delivery (CI/CD) pipeline around one of the most important pieces of your infrastructure.

![HasiCorp Packer](/images/posts/building-golden-image-pipeline/pipeline.jpg)

## Next Time

Now we have covered the concepts of golden images and golden image pipelines, we are going to start to dive into the tools we are going to use to implement one. In the next post we will cover a deeper dive into:

* **[Packer](https://www.packer.io/)**; golden images as code

* **[Ansible](https://docs.ansible.com)**; a configuration management tool we will use to provision our golden images
