package gosoup

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
)

type HtmlSelector struct {
	ClassNames string
	Id         string
	Tag        string
	Multiple   bool
}

func PrintNode(n *html.Node) {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, n)
	fmt.Println(buf.String())
}

func isNodeMatchingSelector(node *html.Node, selector *HtmlSelector) bool {
	if node.Type == html.ElementNode &&
		(selector.Tag == "" || node.Data == selector.Tag) {
		for _, attr := range node.Attr {
			switch attr.Key {
			case "id":
				if attr.Val == selector.Id {
					return true
				}
			case "class":
				if strings.Contains(attr.Val, selector.ClassNames) {
					return true
				}
			}
		}
	}
	return false
}

func GetAttribute(parentNode *html.Node, attribute string) string {
	for _, attr := range parentNode.Attr {
		if attr.Key == attribute {
			return attr.Val
		}
	}
	return ""
}

// Using a recursive solution, this method search for a HTML node matching
// the HTML selectors values (ClassNames, Id, Tag)
func GetNodeByClass(parentNode *html.Node, selector *HtmlSelector) []*html.Node {
	var returnNode []*html.Node

	var crawler func(*html.Node)
	crawler = func(node *html.Node) {
		if isNodeMatchingSelector(node, selector) {
			returnNode = append(returnNode, node)
			if !selector.Multiple {
				return
			}
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}
	crawler(parentNode)

	return returnNode
}
