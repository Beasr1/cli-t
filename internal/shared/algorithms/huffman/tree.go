package huffman

type HuffmanNode struct {
	character rune // if nil => non-leaf
	frequency int
	isLeaf    bool
	left      *HuffmanNode
	right     *HuffmanNode
}

func NewHuffmanNode(character rune, frequency int, isLeaf bool) *HuffmanNode {
	return &HuffmanNode{
		character: character,
		frequency: frequency,
		left:      nil,
		right:     nil,
		isLeaf:    isLeaf,
	}
}

func (hn *HuffmanNode) AttachNodes(left, right *HuffmanNode) {
	//pointer gets already dereferenced in golan
	hn.left = left
	hn.right = right
}

type HuffmanTree struct {
	rootNode *HuffmanNode
}

func traverse(node *HuffmanNode, code string, codes map[rune]string) {
	if node == nil {
		return
	}

	// If leaf node, store the code
	if node.isLeaf {
		// NOTE: maps are already refrence types
		codes[node.character] = code
		return
	}

	// Traverse left (add '0') and right (add '1')
	// NOTE: can think towards string builder here
	traverse(node.left, code+"0", codes)
	traverse(node.right, code+"1", codes)
}

func (ht *HuffmanTree) GenerateCodes() map[rune]string {
	if ht == nil || ht.rootNode == nil {
		return make(map[rune]string)
	}

	codes := make(map[rune]string)

	// Special case: single node tree
	if ht.rootNode.isLeaf {
		codes[ht.rootNode.character] = "0"
		return codes
	}

	// Start traversal from root with empty code
	traverse(ht.rootNode, "", codes)

	return codes
}
