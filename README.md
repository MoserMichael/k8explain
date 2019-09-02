# generate html table with kubectl docs

Command `kubectl api-resources` gives youo a table of kubernetes entities/api resources and  `kubectl explain` produces an explanation of a particular entity.
The following program produces an html table, where each kubernetes entity is linked to a table with its description.

The program runs kubectl commands; parses the table produced by `kubectl api-resources`, etc. etc.
This is a simple exercise that helped the author to pick up some golang programming skills.

Running the program
===================

1. need to do `oc login` as a prerequisite
2. run `go run main.go` on completion it will write the output file out.html

Example output
==============

Link to output report. Beautiful!
