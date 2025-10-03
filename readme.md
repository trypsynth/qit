# Qit
Qit is a tiny CLI program written in Go, to make everyday tasks in git repositories easier and require less typing.

## Building
We use [Mage](https://magefile.org/) for building. Once the repo is cloned, simply typing `mage build` should give you a functional qit binary for your platform. If you want to build binaries for all platforms, run `mage buildAll`.

## Usage
Usage: qit <command> [<args>...]

### Available commands:
* acp: add all files, commit with the specified message, and push.
* amend: amend the last commit with a new message.
* cp: commit changes to tracked files with the specified message, and push.
* db: delete the specified local branch.
* help: show this help message.
* ignore: download .gitignore template(s) from gitignore.io or list available ones.
* last: show the last <number> commits (default: 1).
* license: download license template from GitHub or list available licenses.
* log: show commit history in readable format.
* nb: switch to the branch with  the specified name, creating it if it doesn't exist.
* new: pull and list recent commits.
* reset: hard reset to last commit, discarding all changes.
* status: show simplified summary of working directory changes.
* undo: undo last commit while keeping changes intact.
