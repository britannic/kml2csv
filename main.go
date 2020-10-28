package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

// ExtendedData holds the extended data set for each Placemark
type ExtendedData struct {
	Text string `xml:",chardata" json:"-"`
	Data []struct {
		Text  string `xml:",chardata" json:"-"`
		Name  string `xml:"name,attr" json:"Name"`
		Value string `xml:"value" json:"Value"`
	} `xml:"Data" json:"Data"`
}

// Placemark holds the listing subset
type Placemark []struct {
	Text         string `xml:",chardata" json:"-"`
	Name         string `xml:"name"`
	Address      string `xml:"address"`
	Description  string `xml:"description"`
	StyleURL     string `xml:"styleUrl" json:"-"`
	ExtendedData `xml:"ExtendedData" json:"ExtendedData"`
}

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
		Placemark `xml:"Placemark"`
	} `xml:"Document"`
}

const (
	isDir = iota
	isFile
	isMissing
)

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

	if FileExists(kmlfile) == isMissing {
		if kmlfile == "" {
			kmlfile = "kml file"
		}

		log.Fatal(fmt.Sprintf("%s doesn't exist, exiting!", kmlfile))
	}

	if FileExists(kmlfile) == isDir {
		if kmlfile == "" {
			kmlfile = "kml file"
		}

		log.Fatal(fmt.Sprintf("%s is a directory, exiting!", kmlfile))
	}

	if csvfile != "" && FileExists(csvfile) == isDir {
		log.Fatal(fmt.Sprintf("%s is a directory, exiting!", kmlfile))
	}

	if r, err = OpenFile(kmlfile); err != nil {
		log.Fatal(err)
	}

	if b, err = ioutil.ReadAll(r); err != nil {
		log.Fatal(err)
	}

	if err = xml.Unmarshal(bytes.ReplaceAll(b, []byte("\n"), []byte("")), &data); err != nil {
		log.Fatal(err)
	}

	// fmt.Printf("%v\n", data)

	if b, err = json.Marshal(data.Document.Placemark); err != nil {
		log.Fatal(err)
	}

	// fmt.Printf("%v\n", string(b))

	if csvfile != "" {
		// type csvStruct struct {
		// 	registryNo string
		// 	firstname  string
		// 	lastname   string
		// 	address1   string
		// 	address2   string
		// 	city       string
		// 	state      string
		// 	zip        string
		// 	details    string
		// }

		var (
			f *os.File
			j Placemark
			w *csv.Writer
		)

		if err = json.Unmarshal(b, &j); err != nil {
			log.Fatal(err)
		}

		if f, err = os.Create(csvfile); err != nil {
			log.Fatal(err)
		}

		defer f.Close()

		w = csv.NewWriter(f)

		w.Write([]string{"Name", "Address", "Description"})

		for _, listing := range j {
			var row []string
			// for _, edata := range listing.ExtendedData.Data {
			// 	switch edata.Name {
			// 	case "Registry Number":
			// 		fmt.Println(edata.Value)
			// 	}
			// }

			row = append(row, listing.Name)
			row = append(row, listing.Address)
			row = append(row, listing.Description)
			w.Write(row)
		}

		w.Flush()
		// if WriteFile(csvfile, b) != nil {
		// 	log.Fatal(err)
		// }
		log.Printf("xml converted to csv and saved to %s", csvfile)
		os.Exit(0)
	}

	fmt.Println(string(b))

}

// FileExists checks if a file exists and is not a directory before opening
func FileExists(f string) int {
	info, err := os.Stat(f)
	if os.IsNotExist(err) {
		return isMissing
	}

	if info.IsDir() {
		return isDir
	}

	return isFile
}

// OpenFile opens a file and returns an io.Reader
func OpenFile(f string) (io.Reader, error) {
	// nolint
	return os.Open(f)
}

// WriteFile opens or creates a file for writing to
func WriteFile(f string, b []byte) error {
	return ioutil.WriteFile(f, b, 0644)
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
