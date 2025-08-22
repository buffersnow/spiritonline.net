# SpiritOnline

## available services

- router (incomplete) - core microservices for communication between other services
- myspace (incomplete) - myspaceim messenger backend (based on gamespy)
- wfc - nas, dls1 and conntest servers for nintendo wifi connection
- gsp - gamespy presence server (handles gpcm and gpsp)

## git stuff

use git fetch, then decide if you should rebase or merge --ff-only or do a real merge (usually only in cases of merging to master)

when squashing commits, use git rebase~\[n] where n is the number of commits to squash. keep the first as p/pick, change others to s/squash. `:wq` and then change the commit message that follows after to what you want and `:wq` again. then push with (-f if youve previously pushed those changes)

## project gen

use tools/gen-project.sh to create new boilerplate

## build system

run tools/run-service.sh to run go services
