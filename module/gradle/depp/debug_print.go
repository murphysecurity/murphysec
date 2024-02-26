package depp

import (
	"fmt"
	"io"
)

func printTree(roots []TreeNode, depth int, writer io.Writer) {
	for _, it := range roots {
		for i := 0; i <= depth; i++ {
			if i == depth {
				_, _ = fmt.Fprint(writer, "+--- ")
			} else {
				_, _ = fmt.Fprint(writer, "|    ")
			}
		}
		_, _ = printNode(it, writer)
		_, _ = fmt.Fprint(writer, "\n")
		printTree(it.C, depth+1, writer)
	}
}

func printNode(node TreeNode, writer io.Writer) (int, error) {
	if node.G != "" {
		if node.V != "" {
			return fmt.Fprintf(writer, "%s:%s:%s", node.G, node.A, node.V)
		} else {
			return fmt.Fprintf(writer, "%s:%s", node.G, node.A)
		}
	} else {
		return fmt.Fprintf(writer, "project %s", node.A)
	}
}
