package main

import (
"fmt"
	"net/http"
//	"io/ioutil"
//	"strconv"
//	"github.com/moovweb/gokogiri"
	"golang.org/x/net/html"
//	"go/doc"
	"strings"
	"github.com/tealeg/xlsx"

	//"github.com/derekparker/delve/terminal"

)

type Keyword struct {
	term string
	subterms []string
	synonyms [][]string
}
var keywords []Keyword
const uri = "http://www.onelook.com/?ws1=1&posfilter=n&w=:"

func get_keywords() {
	excelFileName := "/home/steph/Downloads/ENSCIGHT/synononyms/GER_Search_Label_Translations.xlsx"
	excelFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		return
	}

	sheet := excelFile.Sheets[0]
	for _, row := range sheet.Rows[1:] {
		fullterm := string(row.Cells[0].Value)
		keyword := Keyword{term:fullterm}
		keyword.subterms = append(keyword.subterms, fullterm)  // lookup entire term as well
		kws := strings.Split(fullterm, " ")
		for _, kw := range kws {
			if len(kw) > 1 && strings.ToLower(kw) != "and" &&
				strings.ToLower(kw) != "or" &&
				strings.ToLower(kw) != "|" &&
				strings.ToLower(kw) != "/" &&
				strings.ToLower(kw) != "," {
				keyword.subterms = append(keyword.subterms, kw)


			}

		}
		//fmt.Println(keyword.term, "-----", keyword.subterms, "------", len(keyword.subterms))
		keywords = append(keywords, keyword)
	}
	fmt.Println(len(keywords))
}

//var synonyms []string

func main() {
	get_keywords()

	if len(keywords) > 0 {
		for _, terms := range keywords {
			for _, term := range terms.subterms {
				//fmt.Println(len(terms.subterms))
				fulluri := uri + term
				response, err := http.Get(fulluri)
				if err != nil {
					return
				}
				defer response.Body.Close()

				z := html.NewTokenizer(response.Body)
				//response.Body.Close()
				synonyms := []string{}
				exit := 0
				for exit == 0 {
					tt := z.Next()

					switch {
					case tt == html.ErrorToken:
						exit = 1
					case tt == html.StartTagToken:
						t := z.Token()

						isAnchor := t.Data == "a"
						if isAnchor {
							//fmt.Println(t.Attr)
							for _, a := range t.Attr {
								if a.Key == "href" {
									if strings.Contains(a.Val, "&refclue=") {
										els := strings.SplitAfter(a.Val, "=")
										el := els[len(els) - 1]
										//fmt.Println(el)
										//if synonyms == nil {
										//	synonyms = make([]string, 1)
										//}
										synonyms = append(synonyms, el)
									//fmt.Println(synonyms, el)
										//break
									}

								}
								// add to array
							}
						}
					}
				}
				fmt.Println(synonyms)
				terms.synonyms = append(terms.synonyms, synonyms) // add array to synonym array
				for _, result := range keywords {
					for _, synonym := range result.synonyms {
						fmt.Println(result.term, synonym)
					}
				}
			}

			//
			//




		}

	}

}
