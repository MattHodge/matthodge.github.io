```bash
brew install pyenv
brew install pyenv-virtualenv
```

Add this to the bottom of your bash_profile or zshrc:

```bash
if command -v pyenv 1>/dev/null 2>&1; then
  eval "$(pyenv init -)"
fi
```

```bash
# install specific version of python
pyenv install 3.5.4
```

To set the current python version:
```
export PYENV_VERSION=3.5.4

# check if it worked
pyenv version

# create a virtualenv of the current version
pyenv virtualenv venv

# active the virtualenv
pyenv activate venv
```


You can create a `.python-version` file in the repository of an application to automatically set the python verson upon entering the directory:

```ini
# .python-version
3.5.4
```

Docs:
* https://github.com/pyenv/pyenv
* https://github.com/pyenv/pyenv-virtualenv
