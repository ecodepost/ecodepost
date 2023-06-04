package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type SlateItem struct {
	Type     string      `json:"type,omitempty"`
	Text     string      `json:"text,omitempty"`
	Children []SlateItem `json:"children,omitempty"`
	// Bold       bool        `json:"bold,omitempty"`
	// Code       bool        `json:"code,omitempty"`
	// Italic     bool        `json:"italic,omitempty"`
	// Underlined bool        `json:"underlined,omitempty"`
	// URL        string      `json:"url,omitempty"`
}

type SlateDocument []SlateItem

/*
[

	{
	  type: 'paragraph',
	  children: [
	    { text: 'An opening paragraph with a ' },
	    {
	      type: 'link',
	      url: 'https://example.com',
	      children: [{ text: 'link' }]
	    },
	    { text: ' in it.' }
	  ]
	},
	{
	  type: 'quote',
	  children: [
	    {
	      type: 'paragraph',
	      children: [{ text: 'A wise quote.' }]
	    }
	  ]
	},
	{
	  type: 'paragraph',
	  children: [{ text: 'A closing paragraph!' }]
	}

]
*/
func main() {
	doc, err := html.Parse(strings.NewReader("<p>An opening paragraph with a <a href=\"https://example.com\">link</a> in it.</p><blockquote><p>A wise quote.</p></blockquote><p>A closing paragraph!</p>"))
	if err != nil {
		log.Panic(err)
	}
	info := parse(doc, 0)
	infoByte, _ := json.Marshal(info)
	fmt.Println(string(infoByte))
}

func parse(n *html.Node, level int) []SlateItem {
	output := make(SlateDocument, 0)
	switch n.Type {
	case html.ErrorNode:
		// text，没有children
	case html.TextNode:
		output = append(output, SlateItem{
			Text: n.Data,
		})
	case html.DocumentNode:
		c := n.FirstChild
		if c != nil {
			// 说明是body，不要children
			return parse(c, level+1)
		}
	case html.ElementNode:
		if n.DataAtom == atom.Html {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				output = append(output, parse(c, level+1)...)
			}
			return output
		}
		if n.DataAtom == atom.Head {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				output = append(output, parse(c, level+1)...)

			}
			return output
		}
		if n.DataAtom == atom.Body {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				output = append(output, parse(c, level+1)...)
			}
			return output
		}
		// element有children
		slateItem := SlateItem{
			Type: n.Data,
		}
		childrenList := make(SlateDocument, 0)
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			childrenList = append(childrenList, parse(c, level+1)...)
		}
		slateItem.Children = childrenList
		output = append(output, slateItem)
	case html.CommentNode:
	case html.DoctypeNode:
	}
	return output
}
