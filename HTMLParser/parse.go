package link

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

// we want to parse a HTML file into a link struct.
// that is, take out all the href and text, abandon tags
type Link struct {
	Href string
	Text string
}

// Parse will take in a HTML doc and return a slice of links and error
func Parse(r io.Reader) ([]Link, error) {
	// parse html doc into a list of html nodes
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	// finding link nodes out of html nodes
	nodes := linkNodes(doc)

	// transfer info from link nodes into Link struct
	var links []Link
	for _, node := range nodes {
		links = append(links, buildLink(*node))
	}

	return links, nil
}

// finding all the link nodes
func linkNodes(n *html.Node) []*html.Node {
	// if the current node is a link node, we are done
	if n.Type == html.ElementNode && n.Data == "a" {
		return []*html.Node{n}
	}
	var ret []*html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret = append(ret, linkNodes(c)...)
	}
	return ret
}

// retrive info (link & text) from html nodes,
// and fill into a Link
func buildLink(n html.Node) Link {
	var ret Link
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			ret.Href = attr.Val
			break
		}
	}
	ret.Text = text(&n)
	return ret
}

// DEF to collect all text
func text(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	// we are not considering non-element node case
	if n.Type != html.ElementNode {
		return ""
	}
	var ret string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret = ret + text(c)
	}
	return strings.Join(strings.Fields(ret), " ")
}
