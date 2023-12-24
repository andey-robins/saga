package memory

type memCell struct {
	id       int
	valid    bool
	refCount int
}

type Memory struct {
	cells []*memCell
}

func NewMemory() *Memory {
	return &Memory{
		cells: make([]*memCell, 0),
	}
}

// ProcessNode will take a nodeId and find a free cell in the memory
// to write it to. If there are no free cells, a new cell will be allocated.
// refCount is the number of subnodes which reference this node. In other
// words, when refCount becomes 0, this node can be marked invalid and freed
// from memory. The parents slice contains the ids of the nodes which reference
// this node and should have their references decremented when this node is processed
func (m *Memory) ProcessNode(nodeId, refCount int, parents []int) {

	// find a free cell
	writeIdx := -1
	for i, cell := range m.cells {
		if !cell.valid {
			writeIdx = i
			break
		}
	}
	if writeIdx == -1 {
		writeIdx = len(m.cells)
		m.cells = append(m.cells, &memCell{})
	}

	// write the new node into the graph
	m.cells[writeIdx] = &memCell{
		valid:    true,
		id:       nodeId,
		refCount: refCount,
	}

	// decrement the refCount of the parents
	// TODO: this is O(n^2) and could be improved if it becomes a bottleneck
	for _, parentId := range parents {
		for _, cell := range m.cells {
			if cell.id == parentId {
				cell.refCount -= 1
			}
		}
	}
	m.markAndSweep()
}

// markAndSweep will invalidate any cells which have a refCount of 0
// and are therefore not currently being used. Should be called after
// every operation we perform on the underlying graph of this memory
func (m *Memory) markAndSweep() {
	for _, cell := range m.cells {
		if cell.refCount <= 0 {
			cell.valid = false
			cell.refCount = 0
			cell.id = 0
		}
	}
}

// GetMaxUtilization returns the maximum number of cells
// that have been used in this Memory object. This is the
// number of memristor cells that would be needed to compute
// a given sequence
func (m *Memory) GetMaxUtilization() int {
	// since we emulate the memory positions required
	// for a sequence in a greedy way, we can just return
	// the number of cells we've allocated without specifically
	// tracking it
	return len(m.cells)
}
