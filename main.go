package main

import "encoding/xml"

// Kml holds the xml structure for kml files
type Kml struct {
	XMLName  xml.Name `xml:"kml"`
	Text     string   `xml:",chardata"`
	Xmlns    string   `xml:"xmlns,attr"`
	Document struct {
		Text  string `xml:",chardata"`
		Name  string `xml:"name"`
		Style []struct {
			Text      string `xml:",chardata"`
			ID        string `xml:"id,attr"`
			IconStyle struct {
				Text  string `xml:",chardata"`
				Color string `xml:"color"`
				Scale string `xml:"scale"`
				Icon  struct {
					Text string `xml:",chardata"`
					Href string `xml:"href"`
				} `xml:"Icon"`
				HotSpot struct {
					Text   string `xml:",chardata"`
					X      string `xml:"x,attr"`
					Xunits string `xml:"xunits,attr"`
					Y      string `xml:"y,attr"`
					Yunits string `xml:"yunits,attr"`
				} `xml:"hotSpot"`
			} `xml:"IconStyle"`
			LabelStyle struct {
				Text  string `xml:",chardata"`
				Scale string `xml:"scale"`
			} `xml:"LabelStyle"`
		} `xml:"Style"`
		StyleMap struct {
			Text string `xml:",chardata"`
			ID   string `xml:"id,attr"`
			Pair []struct {
				Text     string `xml:",chardata"`
				Key      string `xml:"key"`
				StyleURL string `xml:"styleUrl"`
			} `xml:"Pair"`
		} `xml:"StyleMap"`
		Placemark []struct {
			Text         string `xml:",chardata"`
			Name         string `xml:"name"`
			Address      string `xml:"address"`
			Description  string `xml:"description"`
			StyleURL     string `xml:"styleUrl"`
			ExtendedData struct {
				Text string `xml:",chardata"`
				Data []struct {
					Text  string `xml:",chardata"`
					Name  string `xml:"name,attr"`
					Value string `xml:"value"`
				} `xml:"Data"`
			} `xml:"ExtendedData"`
		} `xml:"Placemark"`
	} `xml:"Document"`
}
