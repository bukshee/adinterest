/*
Input data is a TSV file of the format:

person<tab>interest<nl>

Steps:
. Load the file
. Filter out people who look to be identical or close: compare set of interests
  for each person. If we have two people with identical set of interests that is
  a duplicate. If two people have almost identical set of interests, that is also
  a suspected duplicate.
. Filter out interests based on how many people checked them: if 30% of people or
  more checked a single interest that interest is junk, so skip.
. Now that we have removed junk data we can start comparing the remaining interests
  A cluster of interest is a set with two or more interests in it. We
  need to find all possible clusters and calculate the number of people in it.
. The goal is to list the top clusters found (top: how many people marked all
  interests in it).
*/
package main

import (
	"bufio"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
)

type interestList struct {
	list []interestID
	hits int
}

/*func displayResults() {
	linesWritten := 0
	for {
		if shouldStop() {
			break
		}
		res := <-chResult
		interests := make([]string, 0, len(res.list))
		for _, iID := range res.list {
			interests = append(interests, interestToStr[iID])
		}
		fmt.Printf("%d\t%d\t%s\n", len(res.list), res.hits, strings.Join(interests, "\t"))
		linesWritten++
		if linesWritten > 30 {
			chStop <- true
			break
		}
	}

}
*/

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
		idata.addRow(row[0], row[1])
	}
	err = f.Close()
	if err != nil {
		return err
	}
	return nil
}

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	fname := "ad_interest.tsv"
	if len(os.Args) > 1 {
		fname = os.Args[1]
	}
	idata := NewIdata()
	err := fileLoad(fname, idata)
	if err != nil {
		panic(err)
	}
	idata.ignorePeople()
	idata.ignoreInterests()
	return

	/*for i := 0; i < len(interestToID); i++ {
		iID := interestID(i)
		if _, ok := ints[iID]; !ok {
			continue
		}
	}*/

}
