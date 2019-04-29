# README

This is a command-line tool that takes in a CA Department of Justice (DOJ) .dat file containing criminal record data and identifies convictions eligible for relief under CA Prop 64. 
The output is a CSV file that contains all the original data from the DOJ as well as eligibility info for relevant convictions.

## Prerequisites

 - [Golang]() install with `brew install golang`
 
## Cloning the project

Go projects live in a specific location on your file system under `/Users/[username]/go/src/[project]`.
Be sure to create the directory structure before cloning this project into `../go/src/gogen`

We recommend you add `../go/bin` to your path so you can run certain go tools from the command line 

## Setup

 - Install project dependencies with `cd ~/go/src/gogen` followed by `go get ./...`
 - Install the Ginkgo test library with `go get github.com/onsi/ginkgo/ginkgo`
 - Verify the tests are passing with `ginkgo -r`
 
## Running locally

This tool requires input files in the CA DOJ research file format. These files are tightly controlled for security and confidentiality purposes. 
We have created test fixture files that mimic the structure of the DOJ files, and you can use these to run the code on your local machine.
 - Step 1: Build to executable with `go build .` from the project root
 - Step 2: Run the CLI command: `./gogen --input-doj=/Users/[username]/go/src/gogen/test_fixtures/contra_costa/cadoj_contra_costa.csv --county="CONTRA COSTA" --outputs=[path_to_desired_output_location]`
 
 You can choose any of the three counties we have test fixtures for. Be sure to choose the fixture file that is a csv and begins with `cadoj`, and does NOT include `_results` or `_condensed` in the file name.
 
## License

TBD