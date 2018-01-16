# Ultimate Python Development Setup MacOS


## Setup

```bash
brew install pyenv
brew install direnv
sudo pip2 install pipenv
```

Add this to the bottom of your `.bash_profile` or `.zshrc`:

```bash
########################
# PYENV                #
########################

if command -v pyenv 1>/dev/null 2>&1; then
  eval "$(pyenv init -)"
fi

########################
# PIPENV               #
########################

# OPTIONAL:
# Tell pipenv to store virtualenv inside project directories
export PIPENV_VENV_IN_PROJECT=true

eval "$(pipenv --completion)"

########################
# DIRENV               #
########################

eval "$(direnv hook $SHELL)"
```

## Inside a Python Application Folder

```bash
mkdir datadog-terrascript
cd datadog-terrascript
```

Now you can setup the environment:
```bash
pipenv --python 3.6

# OR

pipenv --python 3

# OR

pipenv --python 2.7.14
```

This will create a `Pipfile` that looks like this:

```bash
[[source]]

url = "https://pypi.python.org/simple"
verify_ssl = true
name = "pypi"

[dev-packages]

[packages]

[requires]
python_version = "3.5"
```

Now you can run an install:

```bash
pipenv install
```

And now activate your shell:

```
pipenv shell
```

You can ofcourse still create `requirements.txt` and `requirements-dev.txt` files to be backwards compatibile for people that aren't using `pipenv`.:

```bash
pipenv lock -r > requirements.txt

pipenv lock -r --dev > requirements-dev.txt
```

## DirEnv

```bash
mkdir datadog-terrascript
cd datadog-terrascript
touch .envrc
```

Open `.envrc` in your text editor:

```
# .envrc

layout pipenv
```

You will see the error from `direnv`, asking you to allow the `.envrc` file. This is a security measure to prevent scripts automatically executing that you haven't whitelisted.

```
direnv: error .envrc is blocked. Run `direnv allow` to approve its content.
```

Allow the `.envrc` file:

```bash
direnv allow
```

The `.envrc` file with automatically load, and it should switch you into the Python virtual environment environment created by pipenv.

If you `cd` out of the project directory you will see the virtual environment deactivate automatically. If you `cd` back into the project directly the virtual environment will activate again. Pretty cool!

## PyCharm Plugins

* [.env file support](https://plugins.jetbrains.com/plugin/9525--env-files-support) - Parameter completion and "Go To declaration" support.
* [EnvFile](https://plugins.jetbrains.com/plugin/7861-env-file) - Loads the `.env` file into environment variables when using `Run/Debug` configurations.
