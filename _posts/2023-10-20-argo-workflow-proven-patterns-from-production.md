---
layout: post
title: Argo Workflows - Proven Patterns from Production
date: 2023-10-20-T00:00:00.000Z
comments: false
description: Discover hard-earned insights on leveraging Argo Workflows for infrastructure automation. This guide outlines essential lessons, from managing workflow TTL and pod garbage collection to running synthetic tests with CronWorkflow. Plus, explore advanced patterns to optimize your workflows, including Parameter Output Facade, Semaphore, and Workflow Injection techniques. Arm yourself with these best practices to avoid common pitfalls and save yourself some pain.
---

[Argo Workflows](https://argoproj.github.io/argo-workflows/use-cases/infrastructure-automation/) provides an excellent platform for infrastructure automation, and has replaced [Jenkins](https://hodgkins.io/automating-with-jenkins-and-powershell-on-windows-part-1) as my go tool for running scheduled or event-driven automation tasks.

In growing my experience with Argo Workflows, I've killed clusters, broken workflows and generally made a mess of things. I've also built a lot of workflows that needed refactoring as they became difficult to maintain.

This blog post aims to share some of the lessons I've learned, and some of the patterns I've developed, to help you avoid the same mistakes I've made.

* TOC
{:toc}

## Lesson - Configure Workflow TTL and Pod Garbage Collection

Want to not kill your Kubernetes control plane? Then you should probably add this high up on the to-do list 😜.

Here's the situation which gave me this lesson. I was [looping](https://github.com/argoproj/argo-workflows/blob/a45afc0c87b0ffa52a110c753b97d48f06cdf166/examples/loops-dag.yaml) over a list of items, executing a template for each. That template contained a DAG of about 50 steps. In total, I was getting to around 1500 pods for the entire workflow.

This was fine for a while until the list of items I was looping over got bigger, and then I started noticing some problems:

- The Argo Workflows UI was very slow/unresponsive when I was trying to view the workflow. This made troubleshooting issues very difficult.
- Some workflows with failing steps, which would retry, eventually got bigger than 1MB, which is the limit of objects in etcd. This causes errors like `Failed to submit workflow: etcdserver: request is too large.`
- Kubernetes was grinding to a halt due to the number of pods it was tracking. This wasn't just active pods, this was pods in the `completed` state.

Tucked away in the [cost optimization](https://argoproj.github.io/argo-workflows/cost-optimisation/#limit-the-total-number-of-workflows-and-pods) section of the Argo Workflows documentation:

> Pod GC - delete completed pods. By default, Pods are not deleted.

Cool, I should have read that earlier 😅.

### Recommendation

First, make sure you persist workflow execution history, and all the logs, so that we can clean up the Kubernetes control plane without losing any information:

- Set up a [Workflow Archive](https://argoproj.github.io/argo-workflows/workflow-archive/) to keep a history of all workflows that have been executed.
- Set up an [Artifact Repository](https://argoproj.github.io/argo-workflows/artifact-repository-ref/), and set `archiveLogs` to `true` in your [Workflow Controller ConfigMap](https://argoproj.github.io/argo-workflows/workflow-controller-configmap/).

Next, in your [Default Workflow Spec](https://argoproj.github.io/argo-workflows/default-workflow-specs/):
  - Add a [Workflow TTL Strategy](https://argoproj.github.io/argo-workflows/fields/#ttlstrategy), which controls the amount of time workflows are kept around after they finish. Clean up successful workflows quickly, and keep failed workflows around longer for troubleshooting.
  - Add a [Pod Garbage Collection](https://argoproj.github.io/argo-workflows/fields/#podgc) strategy, which controls how long pods are kept around for after they finish. Use `OnPodCompletion` to delete pods as soon as they finish inside the workflow. If you don't, and you have a large or long-running workflow, you will end up with a lot of pods in the `completed` state.

Also, don't run massive workflows with 1500 pods, use the [Workflow of Workflows Pattern](#pattern---workflow-of-workflows-with-semaphore) instead.

## Lesson - Use a CronWorkflow to run synthetic tests

Ever made a really small change, just a tiny little tweak to something in your cluster that would *never* break anything... and then it breaks something? Same.

Ever made a change, and then days or weeks later, a workflow that only runs once a month fails because of it? And you have no idea what caused it? And then have to dig into what the hell is going on? And then you find out it was that tiny little change you made weeks ago? And then you feel like an idiot? Same 😔.

### Recommendation

Use a [CronWorkflow](https://argoproj.github.io/argo-workflows/cron-workflows/) to simulate the behaviours of real workflows, without actually exciting the real thing. This allows us to test changes to our cluster, and see if they break anything before they break the real thing.

Here are some example behaviours you can exercise:

- Pulling specific containers
- Retrieving Secrets
- Assuming roles or identities (eg. [Azure AD Workload Identity](https://azure.github.io/azure-workload-identit))
- Accessing resources over the network
- Writing artifacts and logs to the artifact repository

Schedule the workflow as often as makes sense for your use case.

Wrap the synthetic workflow(s) in an [Exit Handler](https://argoproj.github.io/argo-workflows/walk-through/exit-handlers/), allowing you to capture the results of the workflow (`Succeeded`, `Failed`, `Error`). In the exit handler, send telemetry to your monitoring system, and alert on it.

Also set up a "no data" alert, for when your tiny little tweak is to the observability stack, and you break that too 😬.

## Pattern - Parameter Output Facade

Let's say you have a situation where a tool inside a workflow step produces a JSON output like this, with the results of the job stored as `parameter.output.result`:

![Result Parameters as JSON](images/posts/2023-10-20-argo-workflow-proven-patterns-from-production/result_param.png)

We have other teams that depend on this step to trigger their jobs. One team says they would like to be notified by email when the job fails, another wants Slack etc.

Here is what our current workflow looks like that they want to notify on:

{% gist cd84107065dbef23076cd7c31d6cc705 "pattern01-example01.yaml" %}

How should we approach this, giving teams the information they need to handle their notification requirements?

The ["Expression Destructure" example](https://github.com/argoproj/argo-workflows/blob/master/examples/expression-destructure-json.yaml) from the Argo Workflows examples gives us a hint to how we could do this:

{% gist cd84107065dbef23076cd7c31d6cc705 "expression-destructure-json.yaml" %}

This approach has a workflow taking a JSON object as input, and the properties that are needed are extracted using a `jsonpath` expression. While this would work, it's not a good approach.

If the other teams that want notifications did this, it would create a coupling between their workflows and the structure of the JSON object. If the JSON object changes when we execute a job (eg. a field name change, different nesting levels), the teams depending on it would need to update their workflows, or they would get runtime errors.

To prevent this, we can use the [Facade Pattern](https://refactoring.guru/design-patterns/facade) to create a simplified interface (a facade) to the JSON object. This allows us to hide the complexity of the JSON object from the dependent workflows (and teams), and allows it to change without breaking dependent workflows.

Let's expose the following output parameters from our step to the other teams (thereby creating the facade):
- `parameter.output.name` - The name of the job
- `parameter.output.status` - The status of the job (`Failed` or `Succeeded`)
- `parameter.output.failure_reason` - The reason the job failed (or empty if the job succeeded)

This would allow other teams depending on this output to handle their own notifications by checking the `status` parameter, and the `failure_reason` parameter if the status is `Failed`.

(Keep in mind this is just an example, so the properties we are exposing are arbitrary.)

Here is how we can do it:

{% gist cd84107065dbef23076cd7c31d6cc705 "pattern01-facade.yaml" %}

Here are our outputs now:

![Facade Output](images/posts/2023-10-20-argo-workflow-proven-patterns-from-production/facade_output.png)

## Pattern - Workflow of Workflows with Semaphore

Workflow of Workflows is a well-documented pattern in the [Argo Workflows documentation](https://argoproj.github.io/argo-workflows/workflow-of-workflows/).

It involves a parent workflow triggering one or more child workflows. When you are looping over an item, and need to execute a workflow for each, this is where the pattern shines.

Pairing it with [Template-level Synchronization](https://argoproj.github.io/argo-workflows/synchronization/#template-level-synchronization), which allows you to limit the concurrent execution of the child workflows.

Here is an example of using the Workflow of Workflows pattern with:
- A `ConfigMap` defining the concurrency limit.
- The semaphore using the `ConfigMap`, limiting concurrent execution to `1` child workflow

{% gist cd84107065dbef23076cd7c31d6cc705 "pattern-workflow-of-workflow-with-semaphore.yaml" %}

Now we can see the parent workflow only allows one child workflow to run at a time, with the other workflows waiting patiently:

![Workflow of Workflow with Semaphore](images/posts/2023-10-20-argo-workflow-proven-patterns-from-production/workflow-of-workflow-semaphore.png)

When using this pattern after you have [configured workflow TTL and pod garbage collection](#lesson---configure-workflow-ttl-and-pod-garbage-collection), you will keep the number of pods running in the cluster to a minimum, which makes the Kubernetes control plane a little happier.

## Pattern - Workflow Injection

The next pattern comes into play when you find yourself looping over the same list of items across multiple dependent workflows.

For our example, say we have a step which is enumerating all of our database instances in a cloud provider, and then executing a workflow for each database.

Of course, we want to leverage the [Workflow of Workflows with Semaphore](#pattern---workflow-of-workflows-with-semaphore) pattern to do this. Let's look at how we might implement this. In the below example:

- There is a `WorkflowTemplate` called `get-dbs` which acts as the template that enumerates the databases.
- There is a `WorkflowTemplate` called `db-auditing` with the template `get-row-count-for-all-dbs` which enumerates the databases by referencing `get-dbs`, then executes the child workflow for each of the databases
- The `get-row-count-and-log` is the actual template we want executed for each database, which gets the row count for the database, and sends an event with the details.

{% gist cd84107065dbef23076cd7c31d6cc705 "pattern-wf-injection-01.yaml" %}

If we decide we want to perform another operation against all the databases (eg. get the size of each database - `get-database-size-and-log`), we would need to re-implement the boxes in orange:

<div style="max-width:1000px;">
  {% include_relative diagrams/2023-10-20-argo-workflow-proven-patterns-from-production/without-workflow-injection.elk.sketch.svg %}
</div>

It just so happens that the parts we need to re-implement are also the more complex parts of the workflow:

- A DAG which first enumerates the databases using `get-dbs`
- The `withParam` to loop over each database
- The child workflow using the [Workflow of Workflows with Semaphore](#pattern---workflow-of-workflows-with-semaphore)

After we duplicate all of that, we can _finally_ create the template which contains our new logic (`get-database-size-and-log`).

That's where the Workflow Injection pattern comes in:

- We create an additional template called `for-each-db` inside the `get-dbs` `WorkflowTemplate`. This is the template that people can use when they want to run a workflow for each database.
- The `for-each-db` template handles looping over each database and executing the child workflow. It accepts the following parameters:
  - `workflow_template_ref` - What is the name of the `WorkflowTemplate` to execute for each database?
  - `entrypoint` - What is the entrypoint of the `WorkflowTemplate` to execute for each database?
  - `semaphore_configmap_name` - What is the name of the `ConfigMap` to use for the semaphore? This allows the caller to control the concurrency of the child workflows.
  - `semaphore_configmap_key` - What is the key in the `ConfigMap` to use for the semaphore? This allows the caller to control the concurrency of the child workflows.
- For each child workflow, `for-each-db` passes in `db-host` and `environment` as inputs.

This moves all the complexity into the `get-dbs` `WorkflowTemplate`, making it easy for the caller to define a `WorkflowTemplate` that only considers execution against a single database. The callers' workflow must accept `db-host` and `environment` as inputs.

The green boxes in this diagram show the logic of looping and executing a child workflow is now contained in the `get-dbs` `WorkflowTemplate. The callers' workflow is greatly simplified:

<div style="max-width:600px;">
  {% include_relative diagrams/2023-10-20-argo-workflow-proven-patterns-from-production/with-workflow-injection.elk.sketch.svg %}
</div>

Here is the final workflow using this pattern:

{% gist cd84107065dbef23076cd7c31d6cc705 "pattern-wf-injection-02.yaml" %}

## Conclusion

I hope applying the suggestions from the lessons and patterns above will help you avoid some of the pitfalls I've encountered while using Argo Workflows.

If you have any questions or feedback or would like to talk about your experiences with running Argo Workflows in production, you can reach me on [Twitter / X](https://twitter.com/MattHodge).
