package main

/*
Idata is an interest-person dataset. It is populated with AddRow().
Once data is loaded processing can start: GenResult(). It will take a
while to finish.
After results are compiled you can get the number of results via
NumResults() and obtain each result via GetResult()
A result is a set of interests along with the number of people in it.
*/

import (
	"errors"

	"github.com/bukshee/bitfield"
)

type interestID int
type personID int
type piRow struct {
	pID personID
	iID interestID
}

// Idata is an person-interest dataset
// that has all the info retrieved from a tsv file
type Idata struct {
	toStr      map[interestID]string
	toID       map[string]interestID
	pToID      map[string]personID
	rows       []piRow
	pIgnore    map[personID]struct{}
	iIgnore    map[interestID]struct{}
	iFieldMap  map[interestID]*bitfield.BitField
	iIDs       []interestID
	iMin, iMax int
	minPeople  int

	result []*iSet
}

// NewIdata creates a new (empty) Interest dataset
// Arguments set the conditions based on which to collect results from
// the dataset:
//
// iMin: the minimum number of people per interest: below which the interest
// is dropped being too rare
//
// iMax: the maximum number of people per interest: above which the interest
// is dropped being too generic
//
// minPeople: the set of interests has to be at least this many people in it
//
// The condition 0 <= iMin <= minPeople < iMax must be satisfied.
// If not and error is returned
func NewIdata(iMin, iMax, minPeople int) (*Idata, error) {
	correct := iMin >= 0 && iMin <= minPeople && minPeople < iMax
	if !correct {
		return nil, errors.New("Wrong input parameters")
	}
	return &Idata{
		toStr:     make(map[interestID]string, 1000),
		toID:      make(map[string]interestID, 1000),
		pToID:     make(map[string]personID, 150),
		rows:      make([]piRow, 0, 1000),
		pIgnore:   make(map[personID]struct{}, 50),
		iIgnore:   make(map[interestID]struct{}, 50),
		iFieldMap: make(map[interestID]*bitfield.BitField, 1000),
		iMin:      iMin,
		iMax:      iMax,
		minPeople: minPeople,
		iIDs:      nil,
		result:    nil,
	}, nil
}

// AddRow adds a new person-interest pair to the dataset
func (id *Idata) AddRow(person, interest string) {
	pID, ok := id.pToID[person]
	if !ok {
		pID = personID(len(id.pToID))
		id.pToID[person] = pID
	}
	iID, ok := id.toID[interest]
	if !ok {
		iID = interestID(len(id.toID))
		id.toID[interest] = iID
		id.toStr[iID] = interest
	}
	id.rows = append(id.rows, piRow{pID, iID})
}

// GenResult generates a list of interest-sets based on imput parameters:
func (id *Idata) GenResult() {
	id.ignorePeople()
	id.ignoreInterests()
	id.groupInterests()
}

// ignorePeople finds people we want to ignore in the dataset
// these are the ones where
// - two people have the exact same set of interests
// - or the set deviates by < 3 interests
func (id *Idata) ignorePeople() int {
	pers := make(map[personID]*bitfield.BitField, len(id.pToID))
	pIDs := make([]personID, 0, len(id.pToID))
	for i := 0; i < len(id.pToID); i++ {
		pID := personID(i)
		ints := make([]int, 0, 10)
		for _, row := range id.rows {
			if row.pID != pID {
				continue
			}
			ints = append(ints, int(row.iID))
		}
		pIDs = append(pIDs, pID)
		pers[pID] = bitfield.New(len(id.toID)).Set(ints...)
	}
	// find people we want to ignore
	bf := bitfield.New(len(id.toID)).Mut()
	for i, pID := range pIDs {
		for _, pID2 := range pIDs[i+1:] {
			pers[pID].Copy(bf)
			// if two people differ only in <3 interests => ignore one of them
			if bf.Xor(pers[pID2]).OnesCount() < 1 {
				id.pIgnore[pID2] = struct{}{}
			}
		}
	}
	return len(id.pIgnore)
}

// find and remove interests we want to ignore from the dataset
func (id *Idata) ignoreInterests() {
	iIDs := make([]interestID, 0, len(id.toID))
	for i := 0; i < len(id.toID); i++ {
		iID := interestID(i)
		people := make([]int, 0, 50)
		for _, row := range id.rows {
			if row.iID != iID {
				continue
			}
			// skip duplicated people
			if _, ok := id.pIgnore[row.pID]; ok {
				continue
			}
			people = append(people, int(row.pID))
		}
		iIDs = append(iIDs, iID)
		id.iFieldMap[iID] = bitfield.New(len(id.pToID)).Mut().Set(people...)
	}
	id.rows = nil

	bf := bitfield.New(len(id.pToID)).Mut()
	// find interests we want to ignore
	for i, iID := range iIDs {
		// ignore interest with less than 10 or more than 50 people in it
		num := id.iFieldMap[iID].OnesCount()
		if num < id.iMin || num >= id.iMax {
			id.iIgnore[iID] = struct{}{}
			continue
		}
		for _, iID2 := range iIDs[i+1:] {
			id.iFieldMap[iID].Copy(bf)
			// ignore interests where two interest overlap with <3 people difference
			if bf.Xor(id.iFieldMap[iID2]).OnesCount() < 1 {
				id.iIgnore[iID2] = struct{}{}
			}
		}
	}
	for i := range id.iIgnore {
		delete(id.iFieldMap, i)
	}
	id.iIDs = make([]interestID, 0, len(iIDs))
	for _, iID := range iIDs {
		if _, exists := id.iIgnore[iID]; exists {
			continue
		}
		id.iIDs = append(id.iIDs, iID)
	}
}

// groupInterests generates the resulting interest-sets
// returns the number of sets identified
func (id *Idata) groupInterests() int {
	id.result = make([]*iSet, 0, 20)
	for i, iID := range id.iIDs {
		s := newiSet(len(id.pToID), id.minPeople)
		if !s.add(iID, id.iFieldMap[iID]) {
			continue
		}
		for _, iID2 := range id.iIDs[i+1:] {
			s.add(iID2, id.iFieldMap[iID2])
		}
		if s.len() == 1 {
			continue
		}
		id.result = append(id.result, s)
	}
	return len(id.result)
}

// NumResults returns the number of interest-sets found after processing
func (id *Idata) NumResults() int {
	return len(id.result)
}

// GetResult returns a single interest-set (interests) along with the number
// of people in it (numPeople), or an error if position pos is wrong.
// Must be called after GenResults()
func (id *Idata) GetResult(pos int) (numPeople int, interests []string, err error) {
	if pos < 0 || pos >= len(id.result) {
		err = errors.New("wrong position")
		return
	}
	s := id.result[pos]
	numPeople = s.bits.OnesCount()
	interests = make([]string, 0, s.len())
	for _, iID := range s.elements {
		interests = append(interests, id.toStr[iID])
	}
	return
}

// iSet is a set holding interestIDs
type iSet struct {
	elements  []interestID       // the elements of the set
	bits      *bitfield.BitField // the bits of each elements AND-ed together
	minPeople int                // minimum number of bits set to get into the set
}

// NewiSet creates a new iSet struct holding interestIDs
func newiSet(numPeople, minPeople int) *iSet {
	return &iSet{
		elements:  make([]interestID, 0, 10),
		bits:      bitfield.New(numPeople),
		minPeople: minPeople,
	}
}

func (set *iSet) add(iID interestID, fieldMap *bitfield.BitField) bool {
	if fieldMap.OnesCount() < set.minPeople {
		return false
	}
	if len(set.elements) == 0 {
		set.elements = append(set.elements, iID)
		set.bits = fieldMap
		return true
	}
	b := set.bits.Clone().And(fieldMap)
	if b.OnesCount() < set.minPeople {
		return false
	}
	set.elements = append(set.elements, iID)
	set.bits = b
	return true
}

func (set *iSet) len() int {
	return len(set.elements)
}
