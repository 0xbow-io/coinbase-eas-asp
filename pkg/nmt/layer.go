package nmt

import (
	"errors"
	"fmt"
	"math"
)

const MerkleZeroHex string = "2fe54c60d3acabf3343a35b6eba15db4821b340f76e741e2249685ed4899af6c"

/*
Refer to: https://github.com/celestiaorg/nmt/blob/master/docs/spec/nmt.md#add-leaves

                                 00 03 b1c2cc5                                Tree Root
                           /                       \
                          /                         \
                        NsH()                       NsH()
                        /                             \
                       /                               \
               00 00 ead8d25                      01 03 52c7c03               Non-Leaf Nodes
              /            \                    /               \
            NsH()          NsH()              NsH()             NsH()
            /                \                /                   \
    00 00 5fa0c9c       00 00 52385a0    01 01 71ca46a       03 03 b4a2792    Leaf Nodes
        |                   |                 |                   |
      NsH()               NsH()              NsH()               NsH()
        |                   |                 |                   |
00 6c6561665f30      00 6c6561665f31    01 6c6561665f32      03 6c6561665f33  Leaves with namespaces

        0                   1                  2                  3           Leaf Indices

*/

type NameSpaces interface {
	Size() int
	ValidateAndSort() []ID
	GetRecords(ns ID) NameSpaceGroup
	NamespaceSize() IDSize
}

type leafRange struct {
	// start and end denote the indices of a leaf in the tree. start ranges from
	// 0 up to the total number of leaves minus 1 end ranges from 1 up to the
	// total number of leaves end is non-inclusive
	start, end int
}

func genleafLayer(ns NameSpaces) (Layer, map[string]leafRange) {
	var (
		leafLayer       = make(Layer, ns.Size())
		namespaceRanges = make(map[string]leafRange)
	)
	for _, namespace := range ns.ValidateAndSort() {
		for _, rec := range ns.GetRecords(namespace) {
			node := DataToNode(32, rec)
			leafLayer = append(leafLayer, make(Node, len(node)))
			copy(leafLayer[len(leafLayer)-1], node)
		}
	}
	return leafLayer, namespaceRanges
}

// calculateAbsenceIndex returns the index of a leaf of the tree that 1) its
// namespace ID is the smallest namespace ID larger than nID and 2) the
// namespace ID of the leaf to the left of it is smaller than the nID.
// assuming leafLayer is sorted by namespace ID.
func calculateAbsenceIndex(namespaceLen IDSize, nID ID, leafLayer Layer) int {
	var prevLeaf Node
	for index, curLeaf := range leafLayer {
		if index == 0 {
			prevLeaf = curLeaf
			continue
		}

		// Note that here we would also care for the case current < nId < prevNs
		// but we only allow pushing leaves with ascending namespaces; i.e.
		// prevNs <= currentNs is always true. Also we only check for strictly
		// smaller: prev < nid < current because if we either side was equal, we
		// would have found the namespace before.
		if prevLeaf.MinNs(namespaceLen).Less(nID) && nID.Less(curLeaf.MinNs(namespaceLen)) {
			return index
		}

		prevLeaf = curLeaf
	}
	// the case (nID < minNID) or (maxNID < nID) should be handled before
	// calling this private helper!
	panic("calculateAbsenceIndex() called although (nID < minNID) or (maxNID < nID) for provided nID")
}

type Layer []Node

func (l Layer) Size() int {
	return len(l)
}

func (l Layer) ValidateRange(start, end int) error {
	if start < 0 || start >= end || end > l.Size() {
		return ErrInvalidRange
	}
	return nil
}

func (l Layer) Hashes(namespaceLen IDSize) []Element {
	hashes := make([]Element, len(l))
	for i := 0; i < len(l); i++ {
		hashes[i] = l[i].Hash(namespaceLen)
	}
	return hashes
}

func (l Layer) String(namespaceLen IDSize) string {
	out := fmt.Sprintf("Num Of Nodes: %d \n", len(l))
	for i := 0; i < len(l); i++ {
		out += fmt.Sprintf("**** maxID: %s minID: %s hash: %s ****\n", l[i].MaxNs(32).String(), l[i].MinNs(32).String(), l[i].Hash(32).Hex())
	}
	out += "\n"
	return out
}

func GetLayerCount(level, numLeaves int) int {
	return int(math.Ceil(float64(numLeaves) / math.Pow(2, float64(level))))
}

func BuildLayer(namespaceLen IDSize, hashFn HashFunction, nodeSize int, prevLayer Layer, zero Node) Layer {
	l := make(Layer, nodeSize)

	prevLen := len(prevLayer)
	for i := 0; i < len(l); i++ {

		// default right node to zero value
		right := zero
		zeroSide := 2

		if i*2+1 < prevLen {
			right = prevLayer[i*2+1]
			zeroSide = 0
		}

		l[i] = BuildNode(namespaceLen, prevLayer[i*2], right, zeroSide, hashFn)
	}

	return l
}

type Layers []Layer

func (l Layers) Size() int {
	return len(l[0])
}

func (l Layers) Depth() int {
	return len(l)
}

func (l Layers) Levels() int {
	return l.Depth() - 1
}

func (l Layers) GetRootNode() Node {
	return l.GetLayer(l.Depth() - 1)[0]
}

func (l Layers) GetLayer(level int) Layer {
	return l[level]
}

func CalcRoot(namespaceLen IDSize, hashFn HashFunction, leafNodes Layer, zeroValue Element) (Node, int) {

	var (
		numLeaves = len(leafNodes)
		depth     = int(math.Ceil(math.Log2(float64(numLeaves)))) + 1
		zero      = zeroValue
		currLayer = leafNodes
	)

	for level := 1; level < depth; level++ {
		currLayer = BuildLayer(namespaceLen, hashFn, GetLayerCount(level, numLeaves), currLayer, NodeValueFromZero(namespaceLen, zero))
		zero = hashFn(zero, zero)
	}
	return currLayer[0], depth - 1
}

func BuildLayers(namespaceLen IDSize, hashFn HashFunction, leafNodes Layer, zeroValue Element) (Layers, []Element) {

	var (
		numLeaves = len(leafNodes)
		depth     = int(math.Ceil(math.Log2(float64(numLeaves)))) + 1
		layers    = make([]Layer, depth)
		zeroes    = make([]Element, depth)
	)

	layers[0] = leafNodes
	zeroes[0] = zeroValue

	for level := 1; level < depth; level++ {
		layers[level] = BuildLayer(namespaceLen, hashFn, GetLayerCount(level, numLeaves), layers[level-1], NodeValueFromZero(namespaceLen, zeroes[level-1]))
		zeroes[level] = hashFn(zeroes[level-1], zeroes[level-1])
	}
	return layers, zeroes
}

// ProveNamespace returns a range proof for the given NamespaceID.
//
// Adaptation of Celestia ProveNamespace implementation:
// https://github.com/celestiaorg/nmt/blob/master/nmt.go#L200
//
// case 1) If the namespace nID is out of the range of the tree's min and max
// namespace i.e., (nID < n.minNID) or (n.maxNID < nID) ProveNamespace returns an empty
// Proof with empty nodes and the range (0,0) i.e., Proof.start = 0 and
// Proof.end = 0 to indicate that this namespace is not contained in the tree.
//
// case 2) If the namespace nID is within the range of the tree's min and max
// namespace i.e., n.minNID<= n.ID <=n.maxNID and the tree does not have any
// entries with the given Namespace ID nID, this will be proven by returning the
// inclusion/range Proof of the (namespaced or rather flagged) hash of the leaf
// of the tree 1) with the smallest namespace ID that is larger than nID and 2)
// the namespace ID of the leaf to the left of it is smaller than the nid. The nodes
// field of the returned Proof structure is populated with the Merkle inclusion
// proof. the leafHash field of the returned Proof will contain the namespaced
// hash of such leaf. The start and end fields of the Proof are set to the
// indices of the identified leaf. The start field is set to the index of the
// leaf, and the end field is set to the index of the leaf + 1.
//
// case 3) In case the underlying tree contains leaves with the given namespace
// their start and end (end is non-inclusive) index will be returned together
// with a range proof for [start, end). In that case the leafHash field of the
// returned Proof will be nil.
//
// The isMaxNamespaceIDIgnored field of the Proof reflects the ignoreMaxNs field
// of n.treeHasher. When set to true, this indicates that the proof was
// generated using a modified version of the namespace hash with a custom
// namespace ID range calculation. For more information on this, please refer to
// the HashNode method in the Hasher.
// Any error returned by this method is irrecoverable and indicates an illegal state of the tree (n).

func ProveNamespace(ns NameSpaces, hashFn HashFunction, zeroValue Element, nID ID) (Proof, error) {
	if ns.Size() == 0 {
		return NewEmptyRangeProof(), nil
	}

	leafLayer, leafRange := genleafLayer(ns)

	root, levels := CalcRoot(32, hashFn, leafLayer, zeroValue)
	if levels == 0 || root == nil {
		return Proof{}, errors.New("failed to calculate root")
	}

	// case 1) In the cases (n.nID < treeMinNs) or (treeMaxNs < nID), return empty
	// range proof
	if nID.Less(root.MinNs(ns.NamespaceSize())) || root.MaxNs(ns.NamespaceSize()).Less(nID) {
		return NewEmptyRangeProof(), nil
	}

	// find the range of indices of leaves with the given nID
	foundRng, found := leafRange[string(nID)]
	proofStart := foundRng.start
	proofEnd := foundRng.end

	// case 2)
	if !found {
		// To generate a proof for an absence we calculate the position of the
		// leaf that is in the place of where the namespace would be in:
		proofStart = calculateAbsenceIndex(ns.NamespaceSize(), nID, leafLayer)
		proofEnd = proofStart + 1
	}

	// case 3) At this point we either found leaves with the namespace nID in
	// the tree or calculated the range it would be in (to generate a proof of
	// absence and to return the corresponding leaf hashes).
	pathLayers, err := BuildRangeProof(ns.NamespaceSize(), hashFn, leafLayer, zeroValue, proofStart, proofEnd)
	if err != nil {
		return Proof{}, err
	}

	if found {
		return NewInclusionProof(proofStart, proofEnd, pathLayers), nil
	}

	return NewAbsenceProof(proofStart, proofEnd, pathLayers, leafLayer[proofStart]), nil
}

func BuildRangeProof(namespaceLen IDSize, hashFn HashFunction, leafNodes Layer, zeroValue Element, proofStart, proofEnd int) (Layers, error) {
	var (
		numLeaves    = len(leafNodes)
		depth        = int(math.Ceil(math.Log2(float64(numLeaves)))) + 1
		pathLayers   = make([]Layer, depth)
		elStartIndex = proofStart
		elEndIndex   = proofEnd
		zero         = zeroValue
		currLayer    = leafNodes
	)

	if err := leafNodes.ValidateRange(proofStart, proofEnd); err != nil {
		return nil, err
	}

	for level := 1; level < depth; level++ {

		startIndex := elStartIndex - (elStartIndex % 2) // default arity is 2
		endIndex := elEndIndex - (elEndIndex % 2) + 2

		for i := startIndex; i < endIndex; i++ {
			if i < len(currLayer) {
				pathLayers[level-1] = append(pathLayers[level-1], currLayer[i])
			} else {
				pathLayers[level-1] = append(pathLayers[level-1], NodeValueFromZero(namespaceLen, zero))
			}
		}

		currLayer = BuildLayer(namespaceLen, hashFn, GetLayerCount(level, numLeaves), currLayer, NodeValueFromZero(namespaceLen, zero))
		zero = hashFn(zero, zero)

		elStartIndex >>= 1
		elEndIndex >>= 1
	}

	// add root node
	pathLayers[depth-1] = append(pathLayers[depth-1], currLayer[0])

	return pathLayers, nil
}

func VerifyRangeProof(namespaceLen IDSize, hashFn HashFunction, pathLayers Layers) bool {
	// verify that the hashes of the nodes in pathlayers are correct
	var (
		numLeaves = len(pathLayers[0])
		depth     = int(math.Ceil(math.Log2(float64(numLeaves)))) + 1
	)
	for level := 1; level < depth; level++ {
		for i := 0; i < len(pathLayers[level]); i++ {
			if i*2 >= len(pathLayers[level-1]) && i*2+1 >= len(pathLayers[level-1]) {
				continue
			}
			left := pathLayers[level-1][i*2]
			right := pathLayers[level-1][i*2+1]
			if !pathLayers[level][i].Equal(BuildNode(namespaceLen, left, right, 0, hashFn)) {
				return false
			}
		}
	}
	return true
}
