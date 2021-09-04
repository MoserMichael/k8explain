# Generate an html page with an interlinked explanation of kubernetes objects

Command `kubectl api-resources` gives you a table of kubernetes entities/api resources and  `kubectl explain` produces an explanation of a particular entity.
The following program produces an html table that displays the output of `kubectl api-resources` as an html table, where each table row is linked to a more detailed description produced by `kubectl explain`.

The program runs kubectl commands; parses the table produced by `kubectl api-resources`, etc. etc.

This is a simple exercise that helped the author to pick up some golang programming skills.

Running the program
===================

1. need to do `oc login` as a prerequisite
2. run `go run main.go` on completion it will write the output file out.html

Example output
==============

Link to output report [link](https://mosermichael.github.io/k8explain/out.html)


Thanks
======

Thanks to Yaacov Zamir for insisting on a nice css for the reprt.


Trivia
======

The following command enters an infinite loop on my local minishift installation

kubectl explain --recursive=true customresourcedefinitions

Client Version: version.Info{Major:"1", Minor:"10+", GitVersion:"v1.10.0+d4cacc0", GitCommit:"d4cacc0", GitTreeState:"clean", BuildDate:"2018-12-06T15:15:06Z", GoVersion:"go1.11.2", Compiler:"gc", Platform:"linux/amd64"}
Server Version: version.Info{Major:"1", Minor:"11+", GitVersion:"v1.11.0+d4cacc0", GitCommit:"d4cacc0", GitTreeState:"clean", BuildDate:"2019-08-30T20:25:39Z", GoVersion:"go1.10.8", Compiler:"gc", Platform:"linux/amd64"}

Therefore I am running this kind of command with a dealine.
