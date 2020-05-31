/*
Input data is a TSV file of the format:

person<tab>interest<nl>

*/
package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

func fileLoad(fname string, idata *Idata) error {
	f, err := os.Open(fname)
	if err != nil {
		return err
	}

	s := bufio.NewScanner(f)
	s.Scan() // first line is header: skip
	for s.Scan() {
		row := strings.Split(s.Text(), "\t")
		if len(row) != 2 || len(row[0]) == 0 || len(row[1]) == 0 {
			return errors.New("Wrong format")
		}
		idata.AddRow(row[0], row[1])
	}
	err = f.Close()
	if err != nil {
		return err
	}
	return nil
}

func main() {

	var Usage = func(msg string) {
		fmt.Fprintln(flag.CommandLine.Output(), msg)
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])

		flag.PrintDefaults()
	}

	var iMin = flag.Int("iMin", 10,
		"minimum number of people an interest must have")
	var iMax = flag.Int("iMax", 50,
		"maximum number of people an interest must have")
	var minPeople = flag.Int("minPeople", 25,
		"minimum number of people a set of interest must have")
	flag.Parse()

	fname := "ad_interest.tsv"
	if len(flag.Args()) >= 1 {
		fname = flag.Arg(0)
	} else {
		fmt.Fprintf(flag.CommandLine.Output(),
			"No filename given, using the default: %s\n", fname)
	}
	idata, err := NewIdata(*iMin, *iMax, *minPeople)
	if err != nil {
		Usage(err.Error())
		os.Exit(1)
	}
	// cpuProfile, _ := os.Create("cpuprofile")
	// memProfile, _ := os.Create("memprofile")
	// pprof.StartCPUProfile(cpuProfile)

	err = fileLoad(fname, idata)
	if err != nil {
		Usage(err.Error())
		os.Exit(2)
	}
	idata.GenResult()

	// pprof.StopCPUProfile()
	// pprof.WriteHeapProfile(memProfile)

	fmt.Print("numPeople\tnumInterestInSet\tinterests\n")
	if idata.NumResults() == 0 {
		fmt.Println("No results found")
	}
	for i := 0; i < idata.NumResults(); i++ {
		num, interests, err := idata.GetResult(i)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%d\t%d\t%s\n",
			num, len(interests), strings.Join(interests, "\t"))
	}
}
