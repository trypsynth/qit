# Qit
Qit is a tiny wrapper around the Git command line tool, to make everyday tasks easier and require less typing. It's written in Crystal, so can be compiled as a single static binary and put on your path easily.

## Building
`shards build --static --release`

## Usage
`qit <command> [<arguments>...]`

### Available commands:
* acp <message>: add all files, commit with the specified message, and push.
* amend <message>: amend the last commit with a new message.
* ignore <templates>: download .gitignore template(s) from gitignore.io.
* ignore list: show available templates from gitignore.io.
* last [<number>]: show the last <number> commits (default: 1).
* log: show the commit log in a readable format.
* reset: hard reset to the last commit, discarding all local changes.
* undo: undo the last commit while keeping changes intact.
