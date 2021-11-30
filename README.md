# Git Purge

GitPurge is a small cli that will clean up local branches based on remote status. 
If the branch does not have a remote it will detele it when running `gitpurge purge`.

## Building

`go build -o bin/gitpurge .`

## Running 

`bin/purge [global flags] <subcommand> [subcommand flags] [subcommand args]`

## TODOs

- [ ] Speed up fetching remote status
