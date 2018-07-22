---
layout: post
title: Up your AWS CloudFormation game with Visual Studio Code
date: 2018-07-21T13:37:00.000Z
comments: true
description: How to improve efficiency when working with AWS CloudFormation in Visual Studio Code.
---

AWS CloudFormation is an amazing tool for infrastructure as code.

What used to take weeks to do on-premises, is now possible in a few minutes with some JSON, or better yet, YAML.

This doesn't mean we can't do better and improve our efficiency even more.

You won't notice any problems working with a few CloudFormation stacks. A few dozen stacks later, or complicated stacks with many resources though, and we need to start optimizing.

Let's take a look at how we can up our CloudFormation game with Microsoft's [Visual Studio Code](https://code.visualstudio.com/).

This article assumes you are familiar with Visual Studio Code and are using YAML based CloudFormation.

* TOC
{:toc}


## YAML indentation

One thing that always catches out people with editing YAML, is indentation. This is especially true for large CloudFormation files.

Take a look at this snippet. Can you see the error at a quick glance 👀 ?

![AWS CloudFormation with Indent Error](/images/posts/up-your-aws-cloudformation-game/cloudformation-with-error.png)

What about now?

![AWS CloudFormation with Indent Error and Highlighting](/images/posts/up-your-aws-cloudformation-game/cloudformation-indentation-error-with-highlighting.png)

Thanks to the colorization and highlighting of the indent column, it is much easier to see that the `Resource` property is incorrectly indented at a quick glance.

To enable this feature, install the following extensions:

* [Indenticator](https://marketplace.visualstudio.com/items?itemName=SirTori.indenticator)
* [indent-rainbow](https://marketplace.visualstudio.com/items?itemName=oderwat.indent-rainbow)

Add the following [user setting](https://code.visualstudio.com/docs/getstarted/settings) to enable the extra highlighting of the block located at the current cursor position.

```json
{
    "indenticator.inner.showHighlight": true
}
```

Easy 👍.

## Tabs, spaces and line endings

Tabs? Spaces? What style of line endings? Use [EditorConfig](https://editorconfig.org/) and end all the discussions.

![Tabs vs Spaces](/images/posts/up-your-aws-cloudformation-game/im-not-hiring-him-he-uses-spaces-not-tabs.jpg)

EditorConfig helps maintain consistency across your CloudFormation files by defining rules which the editor will apply on save.

Install the [EditorConfig for VSCode](https://marketplace.visualstudio.com/items?itemName=EditorConfig.EditorConfig) extension. Don't worry, there are plugins for many text editors if the people you work with don't use VSCode.

Inside your Git repositories that contain your CloudFormation templates, create a `.editorconfig` file at the root directory that looks like this:

```ini
# .editorconfig

# top-most EditorConfig file
root = true

# Unix-style newlines with a newline ending every file
[*]
end_of_line = lf
insert_final_newline = true

# Keeps CloudFormation YAML files standard
[*.yaml]
indent_style = space
indent_size = 2
trim_trailing_whitespace = true
```

Every time someone with EditorConfig saves a file, it will update the YAML file according to your rules 💪 .

## Sorting things alphabetically

Do you like to keep things sorted alphabetically? It makes it easier when you are looking over massive files or IAM policy statements.

You shouldn't have to say "A, B, C, D..." in your head every time we want to sort something, though.

Install [Sort lines](https://marketplace.visualstudio.com/items?itemName=Tyriar.sort-lines) and enjoy the alphabetical awesomeness.

Highlight a block of the CloudFormation template, open the [command palette](https://code.visualstudio.com/docs/getstarted/userinterface#_command-palette) and type `> Sort lines`.

![Sort CloudFormation Alphabetically](https://i.imgur.com/SUT3JBG.gif)

## Fast access to CloudFormation documentation

If you are like me, you can't remember the few hundred CloudFormation resource types and properties.

We can use [VSCode Tasks](https://code.visualstudio.com/docs/editor/tasks) to make our lives easier. A task is a simple block of JSON which can execute commands on our machine.

Inside your CloudFormation repositories, create a `.vscode` folder in the root. Add a `tasks.json` file with the following content:

{% gist be12f2437bd77c9730c43454b4fdcdd1 %}

> ⚠️ The above will only work on Mac, you will probably need to call Start-Process as a PowerShell Task on Windows.

From the [Command Palette](https://code.visualstudio.com/docs/getstarted/userinterface#_command-palette), choose `> Tasks: Run Task`, and select `CF Resource List`. This quick launches the [AWS Resource Types Reference](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-template-resource-type-ref.html) page.

For the `CF Type Search` command to work, first highlight a CloudFormation resource type and then from the [Command Palette](https://code.visualstudio.com/docs/getstarted/userinterface#_command-palette), choose `> Tasks: Run Task`, and select `CF Type Search`.

This will take you to the AWS Documentation search page for the resource:

![Search for CloudFormation Resource](https://i.imgur.com/XEG4lU0.gif)

## CloudFormation linting

Recently, AWS created a tool called [cfn-python-lint](https://github.com/awslabs/cfn-python-lint), which checks your CloudFormation templates for errors. This gives early feedback and reduces cycle time when creating or updating CloudFormation templates.

Instead of submitting a bad template and waiting for the CloudFormation service tell you its bad, let VS Code tell you as you type!

To install the linter, you will need Python installed. If you are on Windows, I would recommend installing Python from [Chocolatey](https://chocolatey.org/packages/python/3.6.6).

```bash
pip install cfn-lint
```

Verify the install worked by running `cfn-lint --version`.

Next, install the [
vscode-cfn-lint](https://marketplace.visualstudio.com/items?itemName=kddejong.vscode-cfn-lint) plugin.

Now when you are editing a CloudFormation template, you will get issues underlined and *Problems* listed when you make a mistake.

![Linting Error In CloudFormation](/images/posts/up-your-aws-cloudformation-game/cloudformation-linting-error-vscode.png)

> ✅ You should also use cfn-lint as part of your validation of pull requests on your CloudFormation repositories.

## Jumping around CloudFormation templates quickly

Got a template with a few thousand lines? Navigating it, and finding the block you are looking for, soon becomes a scroll or search-fest.

First, install the [YAML Support by Red Hat](https://marketplace.visualstudio.com/items?itemName=redhat.vscode-yaml) extension.

Unfortunately, due to a [bug](https://github.com/redhat-developer/yaml-language-server/issues/77) 🐛 we are going to need to disable the YAML validation this extension provides. It currently has issues supporting CloudFormation intrinsic functions. I will update this post when this issue is fixed.

Add the following [user setting](https://code.visualstudio.com/docs/getstarted/settings) to disable validation.

```json
{
    "yaml.validate": true
}
```

Inside your YAML file you will be able to see the logical names of all your resources, search through them and quickly jump to the right section of your CloudFormation template.

![CloudFormation Outline](https://i.imgur.com/9Dp3VUt.gif)

## Conclusion

Even though CloudFormation frees up a HUGE amount of time for you as an Operations engineer, it doesn't mean you shouldn't keep optimizing your processes to be as efficient as possible.

Hopefully with these tips, you can go from a CloudFormation user to a CloudFormation rock star ⭐!
