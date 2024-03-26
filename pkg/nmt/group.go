package nmt

import (
	"bytes"
	"fmt"
	"sort"
)

// Holds the list of namespace-prefixed data objects
// sorted in the order of their insertion grouped by namespace.
// Each namespace-prefixed data item is represented as a byte slice.
// namespaceLen is the length of the namespace prefix.
// data hash is of ElementSize bytes (32 bytes)
type NameSpaceGroup []Record

func (nsg NameSpaceGroup) Contains(d Record) (index int, found bool) {
	for i, v := range nsg {
		if bytes.Equal(v, d) {
			return i, true
		}
	}
	return -1, false
}

func (nsg NameSpaceGroup) Len() int {
	return len(nsg)
}

func (nsg NameSpaceGroup) NID(namespaceLen IDSize) ID {
	return ID(nsg[0][:namespaceLen])
}

// Sets of data that are grouped by namespace
type NamespaceGroups []NameSpaceGroup

func (ngs NamespaceGroups) ID(nsSize IDSize) ID {
	return ID(ngs[0].NID(nsSize))
}

func (ngs NamespaceGroups) Len() int {
	return len(ngs)
}

// easy to migrate to a relational table
type NsGroups struct {
	nsSize IDSize

	TotalNumOfRecords int

	// namespaces ordered in ASC order by its int values
	namespaces []ID
	// mapping of namspace to its index in the ngs slice
	nsIdxs map[string]int
	// groups of data grouped by namespace
	ngs NamespaceGroups
}

func newNsGroups(namespaceLen IDSize) *NsGroups {
	return &NsGroups{
		nsSize: namespaceLen,
		nsIdxs: make(map[string]int),
	}
}

func (ng *NsGroups) NamespaceSize() IDSize {
	return ng.nsSize

}

func (ng *NsGroups) Size() int {
	return ng.TotalNumOfRecords
}

// Sort interface implementation
func (ng *NsGroups) Less(i, j int) bool {
	return ng.namespaces[i].Less(ng.namespaces[j])
}

func (ng *NsGroups) Swap(i, j int) {
	ng.namespaces[i], ng.namespaces[j] = ng.namespaces[j], ng.namespaces[i]
}

func (ng *NsGroups) Len() int {
	return len(ng.namespaces)
}

func (ng *NsGroups) GetRecords(ns ID) NameSpaceGroup {
	if idx, ok := ng.nsIdxs[ns.String()]; ok {
		return ng.ngs[idx]
	}
	return nil
}

func (ng *NsGroups) ValidateAndSort() []ID {
	if err := ng.ValidateOrder(); err != nil {
		ng.Sort()
	}
	return ng.namespaces
}

// Sort upates NsRange values
// so that groups data can be sorted by namespace values
// in anscending order
func (ng *NsGroups) Sort() {
	// sort groups by namespace
	sort.Sort(ng)
}

// Verify that the order of the namespace groups is correct
func (ng *NsGroups) ValidateOrder() error {
	errChan := make(chan error)

	for i := 0; i < len(ng.namespaces)-1; i++ {
		go func(index int) {
			if ng.namespaces[index].Less(ng.namespaces[index+1]) {
				errChan <- nil
			} else {
				errChan <- fmt.Errorf("namespaces %d are out of order", index)
			}
		}(i)
	}

	for i := 0; i < ng.Len()-1; i++ {
		if err := <-errChan; err != nil {
			return err
		}
	}

	return nil
}

func (ng *NsGroups) Add(d Record) (string, ID, error) {
	var (
		nID    ID
		nIDStr string
	)
	// data should be min namespace + hash
	if len(d) < int(ng.nsSize)+ElementSize {
		return "", nil, fmt.Errorf("%w: got: %v, want >= %v", ErrInvalidLeafLen, len(d), int(ng.nsSize)+ElementSize)
	}

	nID = d.NID(ng.nsSize)
	nIDStr = nID.String()

	// check if namespace group exists
	if idx, ok := ng.nsIdxs[nIDStr]; !ok {

		idx = ng.Len()
		// doesn't exist add new entry
		ng.nsIdxs[nIDStr] = idx
		ng.ngs = append(ng.ngs, NameSpaceGroup{make(Record, len(d))})
		ng.namespaces = append(ng.namespaces, nID)

		// copy over data content
		copy(ng.ngs[idx][0], d)

		// new namespace, sort the namespaces
		ng.Sort()

	} else {
		ng.ngs[idx] = append(ng.ngs[idx], make(Record, len(d)))
		copy(ng.ngs[idx][ng.ngs[idx].Len()-1], d)
	}

	ng.TotalNumOfRecords++

	return nIDStr, nID, nil
}
