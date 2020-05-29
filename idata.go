package main

import "github.com/bukshee/bitfield"

type interestID int
type personID int
type piRow struct {
	pID personID
	iID interestID
}

// Idata is an person-interest dataset
// that has all the info retrieved from a tsv file
type Idata struct {
	toStr     map[interestID]string
	toID      map[string]interestID
	pToID     map[string]personID
	rows      []piRow
	pIgnore   map[personID]struct{}
	iIgnore   map[interestID]struct{}
	iFieldMap map[interestID]*bitfield.BitField

	iMin, iMax int
}

// NewIdata creates a new (empty) dataset
func NewIdata() (id *Idata) {
	return &Idata{
		toStr:     make(map[interestID]string, 19000),
		toID:      make(map[string]interestID, 19000),
		pToID:     make(map[string]personID, 150),
		rows:      make([]piRow, 0, 1000),
		pIgnore:   make(map[personID]struct{}, 50),
		iIgnore:   make(map[interestID]struct{}, 50),
		iFieldMap: make(map[interestID]*bitfield.BitField, 19000),
		iMin:      10,
		iMax:      50,
	}
}

// addRow adds a new person-interest map to the dataset
func (id *Idata) addRow(person, interest string) {
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

// ignorePeople finds people we want to ignore in the dataset
// these are the ones where
// - two people have the exact same set of interests
// - or the set deviates by < 3 interests
func (id *Idata) ignorePeople() int {
	pers := make(map[personID]*bitfield.BitField, len(id.pToID))
	pIDs := make([]personID, 0, len(id.pToID))
	for i := 0; i < len(id.pToID); i++ {
		pID := personID(i)
		ints := bitfield.New(len(id.toID))
		for _, row := range id.rows {
			if row.pID != pID {
				continue
			}
			ints.Set(int(row.iID))
		}
		pIDs = append(pIDs, pID)
		pers[pID] = ints
	}
	// find people we want to ignore
	for i, pID := range pIDs {
		for _, pID2 := range pIDs[i+1:] {
			tmp := pers[pID].Clone().Xor(pers[pID2])
			// if two people differ only in <3 interests => ignore one of them
			if tmp.OnesCount() < 1 {
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
		people := bitfield.New(len(id.pToID))
		for _, row := range id.rows {
			if row.iID != iID {
				continue
			}
			// skip duplicated people
			if _, ok := id.pIgnore[row.pID]; ok {
				continue
			}
			people.Set(int(row.pID))
		}
		iIDs = append(iIDs, iID)
		id.iFieldMap[iID] = people
	}
	id.rows = nil

	// find interests we want to ignore
	for i, iID := range iIDs {
		// ignore interest with less than 10 or more than 50 people in it
		num := id.iFieldMap[iID].OnesCount()
		if num <= id.iMin || num >= id.iMax {
			id.iIgnore[iID] = struct{}{}
			continue
		}
		for _, iID2 := range iIDs[i+1:] {
			tmp := id.iFieldMap[iID].Clone().Xor(id.iFieldMap[iID2])
			// ignore interests where two interest overlap with <3 people difference
			if tmp.OnesCount() < 1 {
				id.iIgnore[iID2] = struct{}{}
			}
		}
	}
	for i := range id.iIgnore {
		delete(id.iFieldMap, i)
	}
}
