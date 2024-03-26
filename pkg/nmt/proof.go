package nmt

import (
	"errors"
)

// ErrFailedCompletenessCheck indicates that the verification of a namespace proof failed due to the lack of completeness property.
var ErrFailedCompletenessCheck = errors.New("failed completeness check")

// Proof represents a namespace proof of a namespace.ID in an NMT. In case this
// proof proves the absence of a namespace.ID in a tree it also contains the
// leaf hashes of the range where that namespace would be.
type Proof struct {

	// start index of the leaves that match the queried namespace.ID.
	start int
	// end index (non-inclusive) of the leaves that match the queried
	// namespace.ID.
	end int

	// nodes hold the tree nodes necessary for the Merkle range proof of
	// `[start, end)` in the order of an in-order traversal of the tree. in
	// specific, nodes contain: 1) the namespaced hash of the left siblings for
	// the Merkle inclusion proof of the `start` leaf 2) the namespaced hash of
	// the right siblings of the Merkle inclusion proof of  the `end` leaf
	pathLayers Layers

	// leafHash are nil if the namespace is present in the NMT. In case the
	// namespace to be proved is in the min/max range of the tree but absent,
	// this will contain the leaf hash necessary to verify the proof of absence.
	// leafHash contains a tree leaf that 1) its namespace ID is the smallest
	// namespace ID larger than nid and 2) the namespace ID of the leaf to the
	// left of it is smaller than the nid.
	leafHash []byte
}

// NewEmptyRangeProof constructs a proof that proves that a namespace.ID does
// not fall within the range of an NMT.
func NewEmptyRangeProof() Proof {
	return Proof{0, 0, nil, nil}
}

// NewInclusionProof constructs a proof that proves that a namespace.ID is
// included in an NMT.
func NewInclusionProof(proofStart, proofEnd int, pathLayers Layers) Proof {
	return Proof{proofStart, proofEnd, pathLayers, nil}
}

// NewAbsenceProof constructs a proof that proves that a namespace.ID falls
// within the range of an NMT but no leaf with that namespace.ID is included.
func NewAbsenceProof(proofStart, proofEnd int, pathLayers Layers, leafHash []byte) Proof {
	return Proof{proofStart, proofEnd, pathLayers, leafHash}
}

// IsEmptyProof checks whether the proof corresponds to an empty proof as defined in NMT specifications https://github.com/celestiaorg/nmt/blob/master/docs/spec/nmt.md.
func (proof Proof) IsEmptyProof() bool {
	return proof.start == proof.end && len(proof.pathLayers) == 0 && len(proof.leafHash) == 0
}
