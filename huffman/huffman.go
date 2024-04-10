package huffman

import (
	"fmt"
	"strconv"

	pq "github.com/jupp0r/go-priority-queue"
)

var rootNode node

type Encoder struct {
	input     string
	codeMap   map[string]string
	charToBin map[string]string
}

type Decoder struct {
	byteCode  []uint8
	extraBits uint8
	codeMap   map[string]string
}

func NewEncoder(input string) *Encoder {
	return &Encoder{
		input: input,
	}
}

func NewDecoder(byteCode []uint8, extraBits uint8, codeMap map[string]string) *Decoder {
	return &Decoder{
		byteCode:  byteCode,
		codeMap:   codeMap,
		extraBits: extraBits,
	}
}

type node struct {
	leftChild  *node
	rightChild *node
	val        int
	priority   float64
}

func initNode(val int, leftChild *node, rightChild *node, priority float64) node {
	return node{
		leftChild:  leftChild,
		rightChild: rightChild,
		val:        val,
		priority:   priority,
	}
}

func countFrequency(input string) map[uint8]int {
	freqMap := make(map[uint8]int)
	for i := 0; i < len(input); i++ {
		freqMap[input[i]]++
	}
	return freqMap
}

func (e *Encoder) Encoding() ([]uint8, uint8, map[string]string, error) {
	var err error
	var ok bool
	freqMap := countFrequency(e.input)
	pq := pq.New()
	for key, val := range freqMap {
		pq.Insert(initNode(int(key), nil, nil, float64(val)), float64(val))
	}

	for pq.Len() > 1 {

		node1, err := pq.Pop()
		if err != nil {
			fmt.Println("Got Error in Pop")
		}

		node2, err := pq.Pop()
		if err != nil {
			fmt.Println("Got Error in Pop")
		}
		// Perform a type assertion to convert node1 to a Node
		node_a, ok := node1.(node)
		if !ok {
			// The assertion failed, node1 was not of type Node
			return nil, 0, nil, fmt.Errorf("node1 is not of type Node")
		}
		// Perform a type assertion to convert node1 to a Node
		node_b, ok := node2.(node)
		if !ok {
			// The assertion failed, node2 was not of type Node
			return nil, 0, nil, fmt.Errorf("node2 is not of type Node")
		}
		combinedPriority := node_a.priority + node_b.priority
		combinedNode := initNode(42, &node_a, &node_b, combinedPriority)

		pq.Insert(combinedNode, combinedPriority)

	}

	lastNode, err := pq.Pop()
	if err != nil {
		fmt.Println("Got Error in Pop")
	}

	rootNode, ok = lastNode.(node)
	if !ok {
		// The assertion failed, lastNode was not of type Node
		return nil, 0, nil, fmt.Errorf("lastNode is not of type Node")
	}

	blocks, extraBits, err := e.generateHuffmanCode(false)
	if err != nil {
		return nil, 0, nil, fmt.Errorf("error in generating huffman code: %v", err)
	}

	return blocks, extraBits, e.codeMap, nil
}

func (e *Encoder) travel_tree(n *node, code string) {
	if n.leftChild == nil && n.rightChild == nil {
		e.codeMap[code] = string(rune(n.val))
		e.charToBin[string(rune(n.val))] = code
	} else {
		e.travel_tree(n.leftChild, code+"0")
		e.travel_tree(n.rightChild, code+"1")
	}
}

func (e *Encoder) generateHuffmanCode(displayMap bool) ([]uint8, uint8, error) {
	e.codeMap = make(map[string]string)
	e.charToBin = make(map[string]string)
	var bytes []uint8
	e.travel_tree(&rootNode, "")
	if displayMap {
		fmt.Println(e.codeMap)
	}
	binaryStr := ""
	for j := 0; j < len(e.input); j++ {
		binaryStr += e.charToBin[string(e.input[j])]
	}

	countAppendingZeros := uint8(8 - len(binaryStr)%8)
	extraBits := countAppendingZeros
	for countAppendingZeros > 0 {
		binaryStr = "0" + binaryStr
		countAppendingZeros--
	}

	for i := 0; i < len(binaryStr); i += 8 {
		byteStr := binaryStr[i : i+8]
		byteVal, err := strconv.ParseUint(byteStr, 2, 8)
		if err != nil {
			return nil, 0, fmt.Errorf("error converting string to byte:%v", err)
		}
		bytes = append(bytes, uint8(byteVal))
	}
	return bytes, extraBits, nil

}

func (d *Decoder) Decoding() string {

	var binStr string
	for _, byteVal := range d.byteCode {
		binStr += fmt.Sprintf("%08b", byteVal)
	}

	binStr = binStr[d.extraBits:]

	decodeStr := ""

	for i := 0; i < len(binStr); i++ {
		curr := ""
		for j := i; j < len(binStr); j++ {
			curr = curr + string(binStr[j])
			if value, ok := d.codeMap[curr]; ok {
				decodeStr = decodeStr + value
				i = j
				break
			}
		}
	}

	return decodeStr

}
