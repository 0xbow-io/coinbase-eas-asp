package nmt

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_DataToNode_Zero(t *testing.T) {
	zero, err := hex.DecodeString(MerkleZeroHex)
	require.NoError(t, err)

	zeroNode := NodeValueFromZero(32, Element(zero))
	require.Equal(t, zeroNode.Hash(32).Hex(), hex.EncodeToString(zero))

	require.Equal(t, zeroNode.MinNs(32).String(), "0000000000000000000000000000000000000000000000000000000000000000")
	require.Equal(t, zeroNode.MaxNs(32).String(), "0000000000000000000000000000000000000000000000000000000000000000")

}

func Test_DataToNode(t *testing.T) {
	testGroupSize := 10
	testRecordSize := 10
	nsgroups, _ := gen_ngs(t, testGroupSize, testRecordSize, true)

	for _, nsgroup := range nsgroups.ngs {
		for _, rec := range nsgroup {
			node := DataToNode(32, rec)
			require.Equal(t, rec.Hash(32).Hex(), node.Hash(32).Hex())
			require.Equal(t, rec.NID(32).String(), node.MinNs(32).String())
			require.Equal(t, rec.NID(32).String(), node.MaxNs(32).String())
		}
	}
}

func Test_BuildNode_NonZero_NonZero(t *testing.T) {
	testGroupSize := 2
	testRecordSize := 1
	nsgroups, _ := gen_ngs(t, testGroupSize, testRecordSize, true)

	allNodes, _ := genleafLayer(nsgroups)

	newNode := BuildNode(32, allNodes[0], allNodes[1], 0, Poseidon2)

	expectedHash := Poseidon2(allNodes[0].Hash(32), allNodes[1].Hash(32))
	require.NotEqual(t, "0000000000000000000000000000000000000000000000000000000000000000", expectedHash.Hex(), "Hashes should not 0")

	expectedMinID := nsgroups.namespaces[0]
	expectedMaxID := nsgroups.namespaces[1]

	require.Equal(t, expectedHash.Hex(), newNode.Hash(32).Hex(), "Hashes do not match")
	require.Equal(t, expectedMinID.String(), newNode.MinNs(32).String(), "Min IDs do not match")
	require.Equal(t, expectedMaxID.String(), newNode.MaxNs(32).String(), "Max IDs do not match")

}

func Test_BuildNode_Zero_NonZero(t *testing.T) {
	zero, err := hex.DecodeString(MerkleZeroHex)
	require.NoError(t, err)

	zeroNode := NodeValueFromZero(32, Element(zero))
	require.Equal(t, zeroNode.Hash(32).Hex(), hex.EncodeToString(zero))

	nsgroups, _ := gen_ngs(t, 1, 1, true)

	allNodes, _ := genleafLayer(nsgroups)

	newNode := BuildNode(32, zeroNode, allNodes[0], 1, Poseidon2)

	expectedHash := Poseidon2(zeroNode.Hash(32), allNodes[0].Hash(32))
	require.NotEqual(t, "0000000000000000000000000000000000000000000000000000000000000000", expectedHash.Hex(), "Hashes should not 0")

	expectedMinID := nsgroups.namespaces[0]
	expectedMaxID := nsgroups.namespaces[0]

	require.Equal(t, expectedHash.Hex(), newNode.Hash(32).Hex(), "Hashes do not match")
	require.Equal(t, expectedMinID.String(), newNode.MinNs(32).String(), "Min IDs do not match")
	require.Equal(t, expectedMaxID.String(), newNode.MaxNs(32).String(), "Max IDs do not match")

}

func Test_BuildNode_NonZero_Zero(t *testing.T) {
	zero, err := hex.DecodeString(MerkleZeroHex)
	require.NoError(t, err)

	zeroNode := NodeValueFromZero(32, Element(zero))
	require.Equal(t, zeroNode.Hash(32).Hex(), hex.EncodeToString(zero))

	nsgroups, _ := gen_ngs(t, 1, 1, true)

	allNodes, _ := genleafLayer(nsgroups)

	newNode := BuildNode(32, allNodes[0], zeroNode, 2, Poseidon2)

	expectedHash := Poseidon2(allNodes[0].Hash(32), zeroNode.Hash(32))
	require.NotEqual(t, "0000000000000000000000000000000000000000000000000000000000000000", expectedHash.Hex(), "Hashes should not 0")

	expectedMinID := nsgroups.namespaces[0]
	expectedMaxID := nsgroups.namespaces[0]

	require.Equal(t, expectedHash.Hex(), newNode.Hash(32).Hex(), "Hashes do not match")
	require.Equal(t, expectedMinID.String(), newNode.MinNs(32).String(), "Min IDs do not match")
	require.Equal(t, expectedMaxID.String(), newNode.MaxNs(32).String(), "Max IDs do not match")

}
