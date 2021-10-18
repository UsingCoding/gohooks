# Go-hooks

Simple git hooks written in go that installs **globally** to your machine

## Install

```shell
curl -fsSL https://raw.github.com/UsingCoding/gohooks/master/install.sh | bash
```

## Hooks

### commit-msg

Checks that commit message starts with branch name

You can disable it via ENV variable `GOHOOKS_UNPROTECT_COMMIT_MESSAGE=1`

### pre-push 

Protect pushing to master by denying directly push to master

You can disable it via ENV variable `GOHOOKS_UNPROTECT_MASTER=1`


## Config

Config allow describing which repos hooks should protect.

Example config
```yaml
protectedReposRegExps:
    origin:
# Apply for repos from github
        - .*github.com.*
    source:
# Apply for repos with specific names
        - .*UsingCoding.*
```

To protect your repo you should write in config.yaml his remote name and add regexp that's check the need of protect this repo.
That rules apply to every repo, so you can customize by regexp how to detect repo which should be protected  