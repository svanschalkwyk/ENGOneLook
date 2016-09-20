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
var stopwords = []string{"and", "or"}


func splitKeywords(s string) []string  {
	w := strings.FieldsFunc(s, func(r rune) bool {
		switch r {
		case '<', '>', ' ', '|', '&', '/', '\'', '(', ')','[',']',',',';',':':
			return true
		}
		return false
	})
//	fmt.Printf("%q\n", w)
	return w
}

func write_results() {
	var row *xlsx.Row
	var cell,cell1 *xlsx.Cell
	excelFileName := "/home/steph/Downloads/ENSCIGHT/synononyms/ENG_Search_Label_Translations_1.xlsx"

	excelFile := xlsx.NewFile()//.OpenFile(excelFileName)


	sheet,_ := excelFile.AddSheet("Synonyms")
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.Value= "Term"
	cell = row.AddCell()
	cell.Value = "Synonyms_1"
	cell = row.AddCell()
	cell.Value = "Synonyms_2"
	for _, kw := range keywords {
		row1 := sheet.AddRow()
		cell1 = row1.AddCell()
		cell1.Value = kw.term;
		for _, syns := range kw.synonyms {
			cell2 := row1.AddCell()

			for _, syn := range syns {
				fmt.Println(syn)
				cell2.Value = cell2.Value + "," +  syn
			}

		}
	}
	excelFile.Save(excelFileName)
}

func get_keywords() {
	excelFileName := "/home/steph/Downloads/ENSCIGHT/synononyms/ENG_Search_Label_Translations.xlsx"
	excelFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		return
	}

	sheet := excelFile.Sheets[0]
	//m := make(map[string]string)
	//kwa := make([]string,0)
	for _, row := range sheet.Rows[1:] {
		fullterm := strings.TrimSpace(strings.ToLower(row.Cells[0].Value))
		keyword := Keyword{term:fullterm}

		kwt := splitKeywords(fullterm)
		//for _,t := range kwt {
		//	if _,found := m[t]; !found {
		//		m[t] = fullterm
		//		kwa = append(kwa,t)
		//	}
//			fmt.Println(t, m[t])
		//}


		if len(kwt) > 1 {
			keyword.subterms = append(keyword.subterms, fullterm)
		}


		//kws := strings.Split(fullterm, " ")
		//if len(kws) > 1 {
		//	keyword.subterms = append(keyword.subterms, fullterm)  // lookup entire term as well
		//}
		//for _, kw := range kws {
		//	if len(kw) > 1 && strings.ToLower(kw) != "and" &&
		//		strings.ToLower(kw) != "or" &&
		//		strings.ToLower(kw) != "|" &&
		//		strings.ToLower(kw) != "/" &&
		//		strings.ToLower(kw) != "," {
		for _,kw := range kwt {
			keyword.subterms = append(keyword.subterms, kw)
		}
fmt.Println(fullterm, keyword.subterms)

//			}

//		}
		//fmt.Println(keyword.term, "-----", keyword.subterms, "------", len(keyword.subterms))
		keywords = append(keywords, keyword)
	}
	//line := 0

	//for _,k := range kwa {
	//	line++
	//	fmt.Println(line, k, m[k])
	//}
	//fmt.Println(len(keywords))
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
				//fmt.Println(synonyms)
				terms.synonyms = append(terms.synonyms, synonyms) // add array to synonym array
				totalSynonyms := 0
				//synmap := new(map[string])
				for _, synonym := range terms.synonyms {

					totalSynonyms += len(synonym)

				}
				fmt.Println(terms.term,terms.synonyms)
			}
			}

			//
			//




		}
	write_results()
	}

