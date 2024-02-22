package checkpoint

import (
	"fmt"
	"os"
	"testing"

	"github.com/andey-robins/magical/genetics"
	"github.com/andey-robins/magical/graph"
)

func TestCheckpointSaveAndLoad(t *testing.T) {
	tempDir := os.TempDir()
	cname := tempDir + "/test.json"

	g := graph.LoadGraphFromString("Inputs 3\n1 2 3\nOutputs 1\n4\nNodes 4\nEdges 6\n1 5\n2 7\n3 6\n5 7\n6 4")

	pop := genetics.NewPopulation(2, 1, 0.1, g, 0)

	Save(cname, pop)

	loadedPop := &genetics.Population{}
	Load(cname, loadedPop)

	if loadedPop.Size != 2 {
		t.Errorf("Expected size 2, got %d", loadedPop.Size)
	}

	if loadedPop.Epsilon != 1 {
		t.Errorf("Expected epsilon 1, got %d", loadedPop.Epsilon)
	}

	if loadedPop.MutationChance != 0.1 {
		t.Errorf("Expected mutation chance 0.1, got %f", loadedPop.MutationChance)
	}

	if loadedPop.Genes[0].Sequence == nil {
		t.Errorf("Expected sequence to be non-nil")
	}
	if loadedPop.Genes[1].Sequence == nil {
		t.Errorf("Expected sequence to be non-nil")
	}

	for i := 0; i < 2; i++ {
		loadedSequence := loadedPop.Genes[i].Sequence.GetSequence()
		initialSequence := pop.Genes[i].Sequence.GetSequence()
		for j := 0; j < len(loadedSequence); j++ {
			if loadedSequence[j] != initialSequence[j] {
				t.Errorf("Expected sequence to be the same")
			}
		}
	}
}

func TestExperimentResume(t *testing.T) {
	tempDir := os.TempDir()
	cname := tempDir + "/test.json"

	g := graph.LoadGraphFromString("Inputs 3\n1 2 3\nOutputs 1\n4\nNodes 4\nEdges 6\n1 5\n2 7\n3 6\n5 7\n6 4")

	pop := genetics.NewPopulation(10, 1, 0.1, g, 0)

	pop.Evolve(g)

	Save(cname, pop)

	loadedPop := &genetics.Population{}
	Load(cname, loadedPop)

	pop.SynchronizeRNG()
	loadedPop.SynchronizeRNG()

	pop.Evolve(g)
	loadedPop.Evolve(g)

	if loadedPop.Size != pop.Size {
		t.Errorf("Expected size %d, got %d", pop.Size, loadedPop.Size)
	}

	if loadedPop.Epsilon != pop.Epsilon {
		t.Errorf("Expected epsilon 1, got %d", loadedPop.Epsilon)
	}

	if loadedPop.MutationChance != pop.MutationChance {
		t.Errorf("Expected mutation chance 0.1, got %f", loadedPop.MutationChance)
	}

	if loadedPop.Genes[0].Sequence == nil {
		t.Errorf("Expected sequence to be non-nil")
	}
	if loadedPop.Genes[1].Sequence == nil {
		t.Errorf("Expected sequence to be non-nil")
	}

	fmt.Println("Pop Genes")
	for _, gene := range pop.Genes {
		fmt.Println(gene.Sequence.GetSequence())
	}
	fmt.Println("Loaded Pop Genes")
	for _, gene := range loadedPop.Genes {
		fmt.Println(gene.Sequence.GetSequence())
	}

	for i := 0; i < len(loadedPop.Genes); i++ {
		loadedSequence := loadedPop.Genes[i].Sequence.GetSequence()
		initialSequence := pop.Genes[i].Sequence.GetSequence()
		for j := 0; j < len(loadedSequence); j++ {
			if loadedSequence[j] != initialSequence[j] {
				t.Errorf("Expected sequence to be the same: %v != %v", loadedSequence, initialSequence)
			}
		}
	}
}
