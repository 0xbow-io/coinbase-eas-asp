package nmt

import (
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	fMerkleTree "github.com/0xbow-io/fixed-merkle-tree"
	"github.com/stretchr/testify/require"
)

func genTestMt(t *testing.T, elements []Element, zero Element, depth int, hashFn fMerkleTree.HashFunction) *fMerkleTree.MerkleTree {

	fElements := make([]fMerkleTree.Element, len(elements))
	for i := len(elements) - 1; i >= 0; i-- {
		fElements[i] = fMerkleTree.Element(elements[i][:])
	}

	start := time.Now()
	mt, err := fMerkleTree.NewMerkleTree(depth, fElements, fMerkleTree.Element(zero[:]), hashFn)
	fmt.Printf("Time to generate test Tree %dms \n", time.Since(start).Milliseconds())

	require.NoError(t, err)

	return mt
}

func Test_GenLeafLayer(t *testing.T) {

}

func Test_CalculateAbsenceIndex(t *testing.T) {

}

func Test_ProveNamespace(t *testing.T) {

}

func Test_Layers_Build(t *testing.T) {
	testGroupSize := 100
	testRecordSize := 100
	nsgroup := gen_ngs(t, testGroupSize, testRecordSize, true)

	zero, err := hex.DecodeString(MerkleZeroHex)
	require.NoError(t, err)

	zeroNode := NodeValueFromZero(32, Element(zero))
	require.Equal(t, zeroNode.Hash(32).Hex(), hex.EncodeToString(zero))

	start := time.Now()
	leafNodes, _ := genleafLayer(nsgroup)
	layers, _ := BuildLayers(32, Poseidon2, leafNodes, Element(zero))
	fmt.Printf("Time to generate layers, leaf size: %d .. %d levels .. took %dms \n", len(leafNodes), layers.Levels(), time.Since(start).Milliseconds())

	// Create non-nmt Tree from another pkg to compare root
	testMT := genTestMt(t, leafNodes.Hashes(32), Element(zero), layers.Levels(), fMerkleTree.Poseidon2)

	testLayers := testMT.Layers()

	// compare layer by layer
	for i := 0; i < layers.Depth(); i++ {
		layerElements := layers.GetLayer(i).Hashes(32)
		for j := 0; j < len(layerElements); j++ {
			require.True(t, testLayers[i][j].BigInt().Cmp(layerElements[j].BigInt()) == 0, "Hashes do not match layer: %d index: %d got: %s expected: %s", i, j, layerElements[j].Hex(), testLayers[i][j].Hex())
		}
	}

	rootNode := layers.GetRootNode()
	require.NotNil(t, rootNode)

	require.NotEqual(t, "0000000000000000000000000000000000000000000000000000000000000000", rootNode.Hash(32), "Hashes should not 0")
	require.Equal(t, rootNode.MinNs(32).String(), nsgroup.namespaces[0].String(), "Min IDs do not match")
	require.Equal(t, rootNode.MaxNs(32).String(), nsgroup.namespaces[len(nsgroup.namespaces)-1].String(), "Max IDs do not match")

	require.True(t, testMT.Root().BigInt().Cmp(rootNode.Hash(32).BigInt()) == 0, "Root hashes do not match")
}

func Test_Layers_CalcRoot(t *testing.T) {

	testGroupSize := 10
	testRecordSize := 10
	nsgroup := gen_ngs(t, testGroupSize, testRecordSize, true)

	zero, err := hex.DecodeString(MerkleZeroHex)
	require.NoError(t, err)

	zeroNode := NodeValueFromZero(32, Element(zero))
	require.Equal(t, zeroNode.Hash(32).Hex(), hex.EncodeToString(zero))

	start := time.Now()
	leafNodes, _ := genleafLayer(nsgroup)
	rootNode, level := CalcRoot(32, Poseidon2, leafNodes, Element(zero))
	require.NotNil(t, rootNode)
	fmt.Printf("Time to calc root, leaf size: %d .. %d levels .. took %dms \n", len(leafNodes), level, time.Since(start).Milliseconds())

	// Create non-nmt Tree from another pkg to compare root
	testMT := genTestMt(t, leafNodes.Hashes(32), Element(zero), level, fMerkleTree.Poseidon2)

	require.NotEqual(t, "0000000000000000000000000000000000000000000000000000000000000000", rootNode.Hash(32), "Hashes should not 0")
	require.Equal(t, rootNode.MinNs(32).String(), nsgroup.namespaces[0].String(), "Min IDs do not match")
	require.Equal(t, rootNode.MaxNs(32).String(), nsgroup.namespaces[len(nsgroup.namespaces)-1].String(), "Max IDs do not match")

	require.True(t, testMT.Root().BigInt().Cmp(rootNode.Hash(32).BigInt()) == 0, "Root hashes do not match")

}

func Test_Layers_RangeProof(t *testing.T) {

	testGroupSize := 3
	testRecordSize := 3
	testProofStart := 1
	testProofEnd := 7
	nsgroup := gen_ngs(t, testGroupSize, testRecordSize, true)

	zero, err := hex.DecodeString(MerkleZeroHex)
	require.NoError(t, err)

	zeroNode := NodeValueFromZero(32, Element(zero))
	require.Equal(t, zeroNode.Hash(32).Hex(), hex.EncodeToString(zero))

	start := time.Now()
	leafNodes, _ := genleafLayer(nsgroup)
	pathLayers, err := BuildRangeProof(32, Poseidon2, leafNodes, Element(zero), testProofStart, testProofEnd)
	require.NoError(t, err)
	require.NotNil(t, pathLayers)
	fmt.Printf("Time to calc root, leaf size: %d ..  took %dms \n", len(leafNodes), time.Since(start).Milliseconds())

	rootNode, levels := CalcRoot(32, Poseidon2, leafNodes, Element(zero))

	require.Equal(t, pathLayers.Levels(), levels, "Proof levels do not match")
	// compare root
	require.True(t, pathLayers.GetRootNode().Hash(32).BigInt().Cmp(rootNode.Hash(32).BigInt()) == 0, "Root hashes do not match")

	require.True(t, VerifyRangeProof(32, Poseidon2, pathLayers), "Proof verification failed")
}
