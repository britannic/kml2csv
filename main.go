package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

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

var (
	kmlfile string
	csvfile string
)

func main() {
	var (
		b    []byte
		data Kml
		err  error
		r    io.Reader
	)

	GetParams()

	// rawXMLData := "<data><person><firstname>Nic</firstname><lastname>Raboy</lastname><address><city>San Francisco</city><state>CA</state></address></person><person><firstname>Maria</firstname><lastname>Raboy</lastname></person></data>"

	if !FileExists(kmlfile) {
		kmlfile = "kml file"
		log.Fatal(fmt.Sprintf("%s doesn't exist, exiting!", kmlfile))
	}

	if r, err = GetFile(kmlfile); err != nil {
		log.Fatal(err)
	}

	if b, err = ioutil.ReadAll(r); err != nil {
		log.Fatal(err)
	}

	if err = xml.Unmarshal(b, &data); err != nil {
		log.Fatal(err)
	}

	if b, err = json.Marshal(data); err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))
}

// GetFile reads a file and returns an io.Reader
func GetFile(f string) (io.Reader, error) {
	// nolint
	return os.Open(f)
}

// GetParams processes command line parameters
func GetParams() {
	var (
		args = flag.NewFlagSet("kml2csv", flag.ExitOnError)
		help bool
	)

	args.StringVar(&kmlfile, "kml", "", "kml file to parse")
	args.StringVar(&csvfile, "csv", "", "csv file to parse")
	args.BoolVar(&help, "?", false, "Show help")

	args.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage %s:\n", os.Args[0])
		args.PrintDefaults()
	}

	if args.Parse(os.Args[1:]) != nil || help || len(os.Args) <= 2 {
		args.Usage()
		os.Exit(1)
	}
}

// FileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func FileExists(f string) bool {
	info, err := os.Stat(f)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
