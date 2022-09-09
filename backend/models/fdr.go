package models

type Member struct {
	Prefix     string `xml:"Prefix"`
	Last       string `xml:"Last"`
	First      string `xml:"First"`
	Suffix     string `xml:"Suffix"`
	FilingType string `xml:"FilingType"`
	StateDst   string `xml:"StateDst"`
	Year       int    `xml:"Year"`
	FilingDate string `xml:"FilingDate"`
	DocId      string `xml:"DocID"`
}

type FinancialDisclosure struct {
	Members []Member `xml:"Member"`
}
