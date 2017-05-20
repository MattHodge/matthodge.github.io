# Workflow

Ideally EVERYTHING we do to manage our systems or services should come out of code.

The workflow I wanted for making changes to DataDog monitors / dashboards was the following:

1. Create a pull request against a `datadog` repository in Github
1. Validate the pull request using `terraform plan` to see what changes would be made and verify that there are no errors with the change
1. When the pull request is merged, run `terraform apply` to make the change against DataDog
1. Automatically check back in the `.tfstate` file into the Git repository

> **Note:** Storing Terraform state in AWS S3 is the recommended approach as the `.tfstate` file may contain secrets, or people might forget to check it in to the Git repository. We will have no secrets, and we will be using TeamCity to automatically check the `.tfstate` file into Git, so we won't be using AWS S3 in this case.
