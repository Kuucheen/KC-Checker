package helper

import "golang.org/x/net/html"

func ExtractElementContent(doc *html.Node, targetTagName string) string {
	var content string

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == targetTagName {
			// Concatenate content inside the target element
			content += renderNodeContent(n)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(doc)

	return content
}

func renderNodeContent(n *html.Node) string {
	var content string

	var render func(*html.Node)
	render = func(n *html.Node) {
		if n.Type == html.TextNode {
			content += n.Data
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			render(c)
		}
	}

	render(n)

	return content
}
