---
layout: post
title:  "Automating with Jenkins and PowerShell on Windows - Part 1"
date:   2015-03-03 13:37:00
comments: true
modified: 2016-02-02
---

Take a minute think about how many PowerShell scripts you have written for yourself or your team. Countless functions and modules, helping to automate this or fix that or make your teams lives easier. You spend hours coding, writing in-line help, testing, packaging your script, distributing it to your team. All that effort, and then a lot of the time the script is forgotten about! People just go back to doing things the manual way.

I put this down to being out of sight, out of mind. Users who do not use the command line regularly will quickly forget about the amazing PowerShell-ing that you did to try and make their lives easier.

Then there are are other problems, like working out the best way to give end users permissions to use your function when they aren't administrators. Do you give them remote desktop access to a server and only provide a PowerShell session? Setup PowerShell Web Access? Configure a restricted endpoint? I thought the point of this module was to make your life easier, not make things harder!

Enter [Jenkins](http://jenkins-ci.org/).

![Enter Jenkins](/images/posts/automating_windows_jenkins_p1/01_jenkins-stickers.png "Enter Jenkins")

* TOC
{:toc}

## What is Jenkins?

Jenkins is traditionally used by developers as a continuous integration and build tool, which provides a web interface for creating and executing both manual and scheduled jobs. The following video gives a brief introduction to Jenkins on Linux, to give you an idea of what it can do.

<iframe width="560" height="315" src="https://www.youtube.com/embed/OfptBK8AB_c" frameborder="0" allowfullscreen></iframe>

## Using Jenkins in Operations

Jenkins can be used to do many things for an operations team, but I will be concentrating on leveraging PowerShell to perform actions from the Jenkins server.

Anything you can think of that you can do with PowerShell, you can integrate with Jenkins to provide a user friendly interface which can be used to schedule and run jobs.

These articles will be broken up into two parts:

* Part 1 – Installing Jenkins and creating the first PowerShell job
* Part 2 – Using Jenkins with PowerShell Remoting to perform jobs on remote machines

## Installing Jenkins

As this article is aimed at Jenkins, PowerShell and Windows, I am  going to be using Windows 2012 R2 as the operating system for the Jenkins server.

* Go to [http://jenkins-ci.org/](http://jenkins-ci.org/) and download the Windows native package zip file to your Jenkins server.
* Extract the `.zip` file and run `setup.exe`
* Complete the install, and a browser window will popup and take you to `http://localhost:8080` and show you the Jenkins interface. A firewall rule for this is created by the installer, so head over to your workstation and browse to the server at `http://<YourJenkinsServerIP>:8080`
* You will be presented with the Jenkins interface:

![Welcome to Jenkins](/images/posts/automating_windows_jenkins_p1/02_welcome_to_jenkins.png "Welcome to Jenkins")

Pretty easy right? The installer automatically configures the Jenkins service, so there is nothing left to do on the server at this stage.

## Configuring Basic Security

You will notice that you did not need to provide any credentials to login to the Jenkins interface. Lets fix that.

* On the left side, click **Manage Jenkins**
* Click the **Setup Security** button up the top
* Tick **Enable Security**
* For my installation, I am going to use **Jenkins' own user database**. I also do not want people being able to sign up to my Jenkins instance so I unticked **Allow users to sign up**. As you can see, you can also configure LDAP authentication, which will allow you to authenticate through Active Directory
* Next you choose your Authorization scheme. I am also going to enable **Logged-in users can do anything** for now, but we will improve on this later

![Configure Jenkins Security](/images/posts/automating_windows_jenkins_p1/03_configure_jenkins_security.png "Configure Jenkins Security")

* Click **Save**
* You will be taken to a sign up screen – enter the details and make an account for yourself

When you are logged in, you will be able to create and run jobs and manage Jenkins. When not logged in you can just see what jobs exist and view job history.

## Updating Plugins and Installing the PowerShell Plugin

Jenkins has countless plugins written by the community which extend its functionality. Once such plugin is the **PowerShell Plugin**, allowing us to create jobs which can run PowerShell. While we are installing this plugin, we will do a plugin update.

* From the Jenkins web interface, go to **Manage Jenkins > Manage Plugins**
* In the Updates tab, down the bottom of the page, click **Select All** and then **Download now and install after restart**
* In the new page, tick **Restart Jenkins when installation is complete and no jobs are running.**

![Installing Jenkins Plugin](/images/posts/automating_windows_jenkins_p1/04_install_jenkins_plugins.png "Installing Jenkins Plugin")

* Jenkins will restart to install the plugin updates. You may need to refresh your browser. Once it has restarted, log back in and head back to **Manage Jenkins > Manage Plugins**
* Click on the **Available** tab and type into the filter box `PowerShell`. Put a tick in the install column for [PowerShell Plugin](https://wiki.jenkins-ci.org/display/JENKINS/PowerShell+Plugin) and click **Download now and install after restart**.
* Again, click **Restart Jenkins when installation is complete and no jobs are running**.

The PowerShell Plugin is now installed.

## Creating a job

Now everything is ready, we can create our first job. Our first job is going to be fairly basic – we are going to create a text file, and write a message inside the text file. This job is going to be **parameterized**, allowing us to pass some options into our PowerShell when we run the job. Having parameterized jobs allows you easily pass parameters to your scripts or functions, right from inside the Jenkins interface.

Jenkins does this by setting the parameters chosen in the job as environment variables. For example, if we have a Jenkins job parameter called Filename, during the Jenkins job, the PowerShell session will have an environment variable $env:Filename available to be accessed by the PowerShell the script.

* From the Jenkins web interface, click **New Item**
* For the job name, put in `Create Text File`. Select **Freestyle project**
* Tick **This build is parameterized**. Drop down the **Add Parameter** list and choose **String Parameter**. This will allow the user to enter a string which will be exposed to the Jenkins job
* Again, drop down the **Add Parameter** list and select the **Choice Parameter**. This allows us to give the user a dropdown list to choose from. Enter the options on new lines inside the **Choices** text box
* Provide a useful description so the user will know what each option does

At this point, the job should look like this:

![Jenkins Job Creation Progress](/images/posts/automating_windows_jenkins_p1/05_powershell_jenkins_jobs.png "Jenkins Job Creation Progress")

* Scroll down, and under **Build**, choose **Add build step** and select **Windows PowerShell**. Inside this textbox is where we can put in our PowerShell script to make the magic happen.

> Tip: Use the PowerShell ISE to write your code, then copy and paste it into the Windows PowerShell Command text box.

* Enter the following PowerShell code into the text box:

{% gist 7f5fdfdc209d677db632 %}

![Create Jenkins PowerShell Build Step](/images/posts/automating_windows_jenkins_p1/06_powershell_in_jenkins_jobs.png "Create Jenkins PowerShell Build Step")

* Click **Save**

The job is now saved and ready to go. Lets try it out.

## Running a Job

* Head back to the front page of the Jenkins web UI (click the Jenkins logo). You will see the **Create Text File** job has been added.

![Running the Jenkins Job](/images/posts/automating_windows_jenkins_p1/07_running_jenkins_powershell_job.png "Running the Jenkins Job")

* Click on the **Clock/Play** symbol to the right of the job which will bring you to the job form. Use the dropdown to select a filename and enter a message. Click **Build**.

![Jenkins Build Parameters](/images/posts/automating_windows_jenkins_p1/09_jenkins_powershell_job_paramaters.png "Jenkins Build Parameters")

* You will see the build start on the left side:

![Jenkins Build Starting](/images/posts/automating_windows_jenkins_p1/10_jenkins_powershell_job_starting.png "Jenkins Build Starting")

* When the build has completed, it will turn into a blue circle. Click on the job number. Once in the job, click **Console Output** to view what happened with the PowerShell script.

![Jenkins PowerShell Script Job Output](/images/posts/automating_windows_jenkins_p1/11_jenkins_powershell_job_output.png "Jenkins PowerShell Script Job Output")

Success! Looks like our job was successful, lets have a look on the Jenkins servers `C:\` drive to find the file and take a look what's inside:

![Jenkins PowerShell Job File Output](/images/posts/automating_windows_jenkins_p1/12_jenkins_text_file_output.png "Jenkins PowerShell Job File Output")

As you can see, we can leverage Jenkins to give our PowerShell scripts a web interface, that could be run by anyone, from anywhere!

Stay tuned for the next article where we will start using Jenkins to target other machines on the network using PowerShell remoting, and get more advanced with our jobs. You can follow me on twitter [@matthodge](https://twitter.com/matthodge) where I will tweet our new blog posts.
