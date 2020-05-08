package main

import (
	"bufio"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"./bitfield"
)

type interestID int
type personID int
type piRow struct {
	pID personID
	iID interestID
}

var (
	numPeople     = -1
	numInterest   = -1
	interestToStr = make(map[interestID]string, 19000)
	interestToID  = make(map[string]interestID, 19000)
	personToID    = make(map[string]personID, 150)
	piRows        = make([]piRow, 0, 19000)
	piMatrix      []bitfield.BitField
)

func addRow(row []string) error {
	if len(row) != 2 || len(row[0]) == 0 || len(row[1]) == 0 {
		return errors.New("Wrong format")
	}
	person, interest := row[0], row[1]
	pID, ok := personToID[person]
	if !ok {
		pID = personID(len(personToID))
		personToID[person] = pID
	}
	iID, ok := interestToID[interest]
	if !ok {
		iID = interestID(len(interestToID))
		interestToID[interest] = iID
	}
	piRows = append(piRows, piRow{pID, iID})
	return nil
}

func makeMatrix() {
	numPeople = len(personToID)
	personToID = nil
	numInterest = len(interestToID)

	piMatrix = make([]bitfield.BitField, 0, numInterest)

	for i := 0; i < numInterest; i++ {
		bf := bitfield.New(numPeople)
		for r := 0; r < len(piRows); r++ {
			if piRows[r].iID == interestID(i) {
				bf.Set(int(piRows[r].pID))
			}
		}
		piMatrix = append(piMatrix, bf)
	}
	piRows = nil
}

func fileLoad(fname string) error {
	f, err := os.Open(fname)
	if err != nil {
		return err
	}

	s := bufio.NewScanner(f)
	s.Scan() // first line is header: skip
	for s.Scan() {
		row := strings.Split(s.Text(), "\t")
		err = addRow(row)
		if err != nil {
			return err
		}

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
	err := fileLoad(fname)
	if err != nil {
		panic(err)
	}
	makeMatrix()

	//ToDo: do the actual calculation
}
