package extract

import (
	"bufio"
	"os"
	"strings"
)

type Element struct {
	NewHTML   []byte
	Js        []byte
	JsRemoved bool
}

func ExtractJs(path string) (Element, error) {
    elem := Element{NewHTML: nil, Js: nil, JsRemoved: false}

	file, err := os.Open(path)
	if err != nil {
		return elem, err
	}
    defer file.Close()

    scanner := bufio.NewScanner(file)

    var sbHTML strings.Builder
    var sbJs strings.Builder

    startExtracting := false

    for scanner.Scan() {
        line := scanner.Text()

        if strings.Trim(line, " ") == "<script>" {
            startExtracting = true
            elem.JsRemoved = true
            continue
        }

        if strings.Trim(line, " ") == "</script>" {
            startExtracting = false
            continue
        }

        if startExtracting {
            sbJs.Write(scanner.Bytes())
            sbJs.Write([]byte("\n"))
        } else {
            sbHTML.Write(scanner.Bytes())
            sbHTML.Write([]byte("\n"))
        }

    }

    elem.NewHTML = []byte(sbHTML.String())
    elem.Js = []byte(sbJs.String())

    return elem, nil
}

/*
func ExtractOtherJs(path string) (Element, error) {
    elem := Element{NewHTML: "", Js: "", JsRemoved: false}

	file, err := os.Open(path)
	if err != nil {
		return elem, err
	}
    defer file.Close()

	//doc, err := html.Parse(file)
	doc, err := html.Parse(file)
	if err != nil {
		return elem, err
	}

	var f func(*html.Node)

	f = func(n *html.Node) {
		for c := n.FirstChild; c != nil; {
			next := c.NextSibling // Store next sibling before potential removal
			if c.Type == html.ElementNode && c.Data == "script" {
                //extract Javascript from template
                elem.JsRemoved = true
                elem.Js = c.Data

				// Remove the node from its parent
				if c.PrevSibling != nil {
					c.PrevSibling.NextSibling = c.NextSibling
				} else {
					n.FirstChild = c.NextSibling
				}
				if c.NextSibling != nil {
					c.NextSibling.PrevSibling = c.PrevSibling
				}
				// Don't recurse on the removed node
			} else {
				f(c) // Recurse for children
			}
			c = next // Move to the next sibling
		}
	}

	f(doc)

	var buf bytes.Buffer
	err = html.Render(&buf, doc)
	if err != nil {
		return elem, err
	}

	//final boss
	elem.NewHTML = buf.String()

	return elem, nil
}
*/
