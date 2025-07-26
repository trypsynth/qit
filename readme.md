# Qit
Qit is a tiny wrapper around the Git command line tool, to make everyday tasks easier and require less typing. It's written in Crystal, so can be compiled as a single static binary and put on your path easily.

## Building
`shards build --static --release`

## Usage
Usage: `qit <command> [<arguments>...]`

### Available commands:
* acp <message>: add all files, commit with the specified message, and push.
* amend <message>: amend the last commit with a new message.
* last: view the last commit in a friendly format.
* log: show the commit log in a readable format.
* undo: undo the last commit while keeping changes intact.
