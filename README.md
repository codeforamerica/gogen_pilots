# README

This is a command-line tool that takes in a CA Department of Justice (DOJ) .dat file containing criminal record data and identifies convictions eligible for relief under CA Prop 64. 
The output is a CSV file that contains all the original data from the DOJ as well as eligibility info for relevant convictions.

## Prerequisites

 - [Golang](https://golang.org/) install with `brew install golang`
 
## Cloning the project

Go projects live in a specific location on your file system under `/Users/[username]/go/src/[project]`.
Be sure to create the directory structure before cloning this project into `../go/src/gogen`

We recommend you add `../go/bin` to your path so you can run certain go tools from the command line 

## Setup

 - Change to project root directory `cd ~/go/src/gogen`
 - Install project dependencies with `go get ./...`
 - Install the Ginkgo test library with `go get github.com/onsi/ginkgo/ginkgo`
 - Install project test dependencies with `go get -t ./...`
 - Verify the tests are passing with `ginkgo -r`
 
## Running locally

This tool requires input files in the CA DOJ research file format. These files are tightly controlled for security and confidentiality purposes. 
We have created test fixture files that mimic the structure of the DOJ files, and you can use these to run the code on your local machine.

To compile and run gogen, run:
```
$ go run gogen run --input-doj=/Users/[username]/go/src/gogen/test_fixtures/no_headers.csv --county="SAN JOAQUIN" --outputs=[path_to_desired_output_location]
```

If you would like to create a compiled artifact of gogen and install it (e.g. for use with BEAR), run the following commands from project root:
```
$ go build .
$ go install -i gogen
$ gogen run --input-doj=/Users/[username]/go/src/gogen/test_fixtures/no_headers.csv --county="SAN JOAQUIN" --outputs=[path_to_desired_output_location]`
```
 
## License

TBD