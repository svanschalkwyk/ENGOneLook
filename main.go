package main

import (
"fmt"
	"net/http"
	"golang.org/x/net/html"
	"strings"
	"github.com/tealeg/xlsx"
)

type Keyword struct {
	term string
	subterms []string
	synonyms [][]string
}
var keywords []Keyword
const uri = "http://www.onelook.com/?ws1=1&posfilter=n&w=:"
//var stopwords = []string{"and", "or"}


func splitKeywords(s string) []string  {
	w := strings.FieldsFunc(s, func(r rune) bool {
		switch r {
		case '<', '>', ' ', '|', '&', '/', '\'', '(', ')','[',']',',',';',':':
			return true
		}
		return false
	})
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
	for _, row := range sheet.Rows[1:] {
		fullterm := strings.TrimSpace(strings.ToLower(row.Cells[0].Value))
		keyword := Keyword{term:fullterm}

		kwt := splitKeywords(fullterm)

		if len(kwt) > 1 {
			keyword.subterms = append(keyword.subterms, fullterm)
		}

		for _,kw := range kwt {
			keyword.subterms = append(keyword.subterms, kw)
		}
		keywords = append(keywords, keyword)
	}
}

func main() {

	get_keywords()

	if len(keywords) > 0 {
		for id, terms := range keywords {
			for _, term := range terms.subterms {
				//fmt.Println(len(terms.subterms))
				fulluri := uri + term
				response, err := http.Get(fulluri)
				if err != nil {
					return
				}
				defer response.Body.Close()

				z := html.NewTokenizer(response.Body)
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
										synonyms = append(synonyms, el)
									}

								}
							}
						}
					}
				}
				keywords[id].synonyms = append(keywords[id].synonyms, synonyms) // add array to synonym array
				totalSynonyms := 0
				for _, synonym := range terms.synonyms {
					totalSynonyms += len(synonym)
				}
				fmt.Println(terms.term, terms.synonyms)
			}
		}

	}
	write_results()
}
