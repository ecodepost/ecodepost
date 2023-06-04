package slate

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// https://github.com/reactima/reactima-slate-go
type SlateItem struct {
	Type       string      `json:"type,omitempty"`
	Lang       string      `json:"lang,omitempty"`
	Text       *string     `json:"text,omitempty"`
	Children   []SlateItem `json:"children,omitempty"`
	Bold       bool        `json:"bold,omitempty"`
	Code       bool        `json:"code,omitempty"`
	Italic     bool        `json:"italic,omitempty"`
	Underlined bool        `json:"underlined,omitempty"`
	URL        string      `json:"url,omitempty"`
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
        children: [{ text: '' }]
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
//func main() {
//	doc, err := html.Parse(strings.NewReader("<p>An opening paragraph with a <a href=\"https://example.com\">link</a> in it.</p><blockquote><p>A wise quote.</p></blockquote><p>A closing paragraph!</p>"))
//	if err != nil {
//		log.Panic(err)
//	}
//	info := parse(doc, 0)
//	infoByte, _ := json.Marshal(info)
//	fmt.Println(string(infoByte))
//}

func HtmlToSlate(htmlContent string) ([]SlateItem, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("html parse fail, err: %w", err)
	}
	return parse(doc, 0), nil
}

func HtmlToSlateJson(htmlContent string) (string, error) {
	// html中间有/n，需要去除掉
	//htmlContent = strings.ReplaceAll(htmlContent, "\n", ``)
	//var content string                                     //存储新的文件内容
	//fc := bufio.NewScanner(strings.NewReader(htmlContent)) //按行读取文件内容
	//for fc.Scan() {
	//	content += strings.TrimRight(fc.Text(), "\n") //去除行尾空格，生成新的文件内容
	//}
	doc, err := HtmlToSlate(htmlContent)
	if err != nil {
		return "", fmt.Errorf("HtmlToSlateJson parse fail, err: %w", err)
	}
	jsonBytes, err := json.Marshal(doc)
	if err != nil {
		return "", fmt.Errorf("HtmlToSlateJson marshal fail, err: %w", err)
	}
	return string(jsonBytes), nil
}

// 解析里在element里，text不能有值
// text node里，text必须要有值，否则会被omitempty忽略，现在为了实现，像图片直接给个#
func parse(n *html.Node, level int) []SlateItem {
	output := make(SlateDocument, 0)
	fmt.Printf("n.Type--------------->"+"%+v\n", n.Type)
	switch n.Type {
	case html.ErrorNode:
		// text，没有children
	case html.TextNode:
		strings.ReplaceAll(n.Data, "&amp;", `&`)
		strings.ReplaceAll(n.Data, "&lt;", "<")
		strings.ReplaceAll(n.Data, "&quot;", `"`)
		strings.ReplaceAll(n.Data, "&gt;", ">")
		copyData := n.Data
		// 因为<p></p>这种
		// 我们只去掉普通换行
		if level <= 3 {
			// 去除掉html里的换行
			newData := strings.ReplaceAll(copyData, "\n", "")
			if newData == "" {
				return output
			}
		}

		output = append(output, SlateItem{
			Text: &copyData,
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
		if n.DataAtom == atom.Img || n.DataAtom == atom.Image {
			for _, info := range n.Attr {
				if info.Key == "src" {
					slateItem.URL = info.Val
				}
			}
			// image需要补充children，因为他没有children，但是slate需要他有children
			//childrenList = append(childrenList, SlateItem{
			//	Text: "#",
			//})
		}
		if n.DataAtom == atom.A {
			for _, info := range n.Attr {
				if info.Key == "href" {
					slateItem.URL = info.Val
				}
			}
			// image需要补充children，因为他没有children，但是slate需要他有children
			//childrenList = append(childrenList, SlateItem{
			//	Text: "#",
			//})
		}

		if n.DataAtom == atom.P {
			for _, info := range n.Attr {
				if info.Key == "href" {
					slateItem.URL = info.Val
				}
			}
			// image需要补充children，因为他没有children，但是slate需要他有children
			//childrenList = append(childrenList, SlateItem{
			//	Text: "#",
			//})
		}
		fmt.Printf("n.DataAtom--------------->"+"%+v\n", n.DataAtom)
		fmt.Printf("n.FirstChild --------------->"+"%+v\n", n.FirstChild)
		fmt.Printf("n.LastChild --------------->"+"%+v\n", n.LastChild)
		fmt.Printf("n.NextSibling --------------->"+"%+v\n", n.NextSibling)
		// 默认没有需要这么处理，说明是叶子节点
		// c := n.FirstChild; c != nil; 说明有子节点
		// 那么判断n.FirstChild 加入 # ，让其有children数据
		if n.FirstChild == nil {
			//if n.FirstChild == nil && n.LastChild == nil && n.NextSibling == nil {
			txt := ""
			childrenList = append(childrenList, SlateItem{
				Text: &txt,
			})
		}
		//if n.FirstChild != nil {
		//	if n.FirstChild.FirstChild == nil {
		//		//if n.FirstChild == nil && n.LastChild == nil && n.NextSibling == nil {
		//		txt := ""
		//		childrenList = append(childrenList, SlateItem{
		//			Text: &txt,
		//		})
		//		fmt.Println("ffff")
		//	}
		//}

		if n.DataAtom == atom.Pre {
			// Pre紧挨着第一个元素
			c := n.FirstChild
			if c != nil && c.DataAtom == atom.Code {
				slateItem.Type = "code_block"
				spew.Dump("c.Attr", c.Attr)
				for _, info := range c.Attr {
					if info.Key == "class" && strings.Contains(info.Val, "language-") {
						slateItem.Lang = strings.TrimPrefix(info.Val, "language-")
					}
				}
			}
		}
		if n.DataAtom == atom.Code {
			slateItem.Type = "code_line"
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			childrenList = append(childrenList, parse(c, level+1)...)
		}
		slateItem.Children = childrenList
		fmt.Printf("slateItem--------------->"+"%+v\n", slateItem)

		output = append(output, slateItem)
	case html.CommentNode:
	case html.DoctypeNode:
	}
	return output
}
