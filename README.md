# Clear My Record's Gogen for Pilot Counties

## Intro to Gogen for Pilot Counties

### What Is Gogen for Pilot Counties?


Clear My Record's Gogen for Pilot Counties (`gogen_pilots`) is a command-line tool that takes in California Department of Justice (CA DOJ) .dat files containing criminal record data and recommends convictions determined eligible for relief under California’s Proposition 64 based on Eligibility Parameters described below. 

The code found in this repository was used for our work with selected pilot counties. A similar, configurable version of this code (`gogen`) meant for use with our self-service GUI [BEAR](https://github.com/codeforamerica/bear) can be found in [another repository](https://github.com/codeforamerica/gogen_pilots).

### What Output Does `gogen_pilots` Create?

`gogen_pilots`’s output is a bundle of CSV files that contain original data from the CA DOJ as well as eligibility info for relevant convictions based on Eligibility Parameters described below.

### A Note on Eligibility Parameters

* **What Are Eligibility Parameters?** The source code for `gogen_pilots` contains county-specific eligibility parameters or criteria (Eligibility Parameters). These Eligibility Parameters are a tool drafted and designed to enable `gogen_pilots` to recommend specific convictions as eligible for relief under California’s Proposition 64. These Eligibility Parameters, however, may only be initial or interim drafts of those used by the counties and are also subject to the limitations described below.

* **Eligibility Parameters Might Not Be Up-to-Date.** The applicable county or other state government authority and Code for America have not confirmed the adequacy of those Eligibility Parameters and do not recommend reliance on those Eligibility Parameters. The counties reserve all rights to modify their Eligibility Parameters, the Eligibility Parameters remain subject to change, and `gogen_pilots` may contain Eligibility Parameters that a county is no longer using or has never used. The release of the Eligibility Parameters as part of `gogen_pilots`’s source code or the use of `gogen_pilots` with a county’s Eligibility Parameters does not imply or constitute a representation that those Eligibility Parameters have not changed since the applicable date of release. 

* **Actual Conviction Clearance Outcomes May Differ.** Ultimate conviction clearance outcomes in a county may diverge from outputs obtained by using or testing `gogen_pilots`, and the applicable county or other state government authority and Code for America neither have nor can confirm the accuracy of such outputs.

* **`gogen_pilots` Is Provided “As Is” with No Warranty.** `gogen_pilots`, including any Eligibility Parameters, is provided “as-is”, and neither the applicable county, any other state government authority, nor Code for America makes any warranty, express or implied, or guarantees that `gogen_pilots` will yield the same results as the conviction clearance outcomes ultimately determined by a county, which such county may determine in its discretion.

* **Contact Us for More Info.** If you have any questions about the foregoing, then please contact Lou Moore, Chief Technology Officer at Code For America, at lmoore@codeforamerica.org.

## About the Team

This application was developed by [Code for America](http://codeforamerica.org)'s [Clear My Record team](https://www.codeforamerica.org/programs/clear-my-record).

## How to Use and Develop Gogen for Pilot Counties

### Prerequisites

 - [Golang](https://golang.org/) install with `brew install golang`
 
### Cloning the project

Go projects live in a specific location on your file system under `/Users/[username]/go/src/[project]`.

```
$ git clone git@github.com:codeforamerica/gogen_pilots.git
```

We recommend you add `../go/bin` to your path so you can run certain go tools from the command line 

### Setup

 - Change to project root directory `cd ~/go/src/gogen_pilots`
 - Install project dependencies with `go get ./...`
 - Install the Ginkgo test library with `go get github.com/onsi/ginkgo/ginkgo`
 - Install project test dependencies with `go get -t ./...`
 - Verify the tests are passing with `ginkgo -r`
 
### Running locally

This tool requires input files in the CA DOJ research file format. These files are tightly controlled for security and confidentiality purposes. 
We have created test fixture files that mimic the structure of the DOJ files, and you can use these to run the code on your local machine.
 - Step 1: Build to executable with `go build .` from the project root
 - Step 2: Run the CLI command: `./gogen_pilots run --input-doj=/Users/[username]/go/src/gogen_pilots/test_fixtures/extra_comma.csv --outputs=[path_to_desired_output_location]`
 
 You can choose any of the three counties we have test fixtures for. Be sure to choose the fixture file that is a csv and begins with `cadoj`, and does NOT include `_results` or `_condensed` in the file name.
 
## License

MIT. Please see [LICENSE](./LICENSE) and [NOTICE.md](./NOTICE.md).
