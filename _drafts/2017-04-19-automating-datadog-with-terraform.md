---
layout: post
title: Automating DataDog with Terraform
date: 2017-04-19T13:37:00.000Z
comments: false
description: Using Terraform to automate DataDog provisioning
---

DataDog is an awesome SaaS monitoring platform. We have 100+ developers at work leveraging the platform to collect their metrics, create dashboards and send alerts.

As with anything, if you don't maintain and clean your tools, after a while things can become a little messy. Dashboards start to get named wildly different things with no standards, alerts aren't deleted for decommissioned services or team names change and alerts are suddenly pointing to a team Slack channel that doesn't exist anymore.

Something had to be done to improve the situation, like setting up some standards or rules around DataDog usage, but this is a fine line you need to walk between freedom and standardization. Be too strict or harsh on people and they no longer find the tool nice to use, and instead think of it as pain in the ass.

Take too much freedom away, and you get "shadow IT" situations with people using their own tools.

With this in mind, I decided on a few goals:

* Have monitors be created from code
* Make it easy for an application to be changed to another team
* Make it easy for a team name or alert destination to be changed
* Many applications are very similar - allow all the monitors for one application to be copied and used for another application

* TOC
{:toc}

# Deciding on Terraform

There were a few options around for managing DataDog from code:
* [DogPush](https://github.com/trueaccord/DogPush) - Manage DataDog monitors in YAML
* [Barkdog](https://github.com/codenize-tools/barkdog) - Manage DataDog monitors using Ruby DSL, and updates monitors according to DSL
* [Interferon](https://github.com/airbnb/interferon) - A Ruby gem enabling you to store your alerts configuration in code
* [DogWatch](https://github.com/rapid7/dogwatch) - Ruby gem designed to provide a simple method for creating DataDog monitors in Ruby
* [Ansible DataDog Montior Module](https://docs.ansible.com/ansible/datadog_monitor_module.html) - Manages monitors within Datadog via Ansible
* [Terraform DataDog Provider](https://www.terraform.io/docs/providers/datadog/index.html) - Supports creating monitors, users, timeboards and downtimes

I ended up deciding to go with Terraform mainly due to two main reasons:

1. Also being able to create timeboards using the same DSL / process.
1. Terraform is also far more widely supported so from a "googling of problems" perspective (Ansible too)

# Terraform with DataDog Basics

The [DataDog Blog](https://www.datadoghq.com/blog/) recently published a post called [Managing Datadog with Terraform](https://www.datadoghq.com/blog/managing-datadog-with-terraform/).

This will cover the basics to give you an introduction to Terraform. Once you have a read, head back over to this post for some more in-depth usage.

# Workflow

Ideally EVERYTHING we do to manage our systems or services should come out of code.

The workflow I wanted for making changes to DataDog monitors / dashboards was the following:

1. Create a pull request against a `datadog` repository in Github
1. Validate the pull request using `terraform plan` to see what changes would be made and verify that there are no errors with the change
1. When the pull request is merged, run `terraform apply` to make the change against DataDog
1. Automatically check back in the `.tfstate` file into the Git repository

> **Note:** Storing Terraform state in AWS S3 is the recommended approach as the `.tfstate` file may contain secrets, or people might forget to check it in to the Git repository. We will have no secrets, and we will be using TeamCity to automatically check the `.tfstate` file into Git, so we won't be using AWS S3 in this case.

# Repository Structure
# Defining Applications and Teams
# Separation of Concerns with Terraform
# Types of DataDog Monitors
# Handling Terraform State
# Automating Terraform with TeamCity
