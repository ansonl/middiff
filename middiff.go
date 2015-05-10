package main

import (
"fmt"
"io"
"io/ioutil"
"bufio"
"log"
"net/http"
"crypto/tls"
"golang.org/x/net/html"
"golang.org/x/net/html/atom"
"os"
"flag"
"encoding/binary"
"bytes"
"strings"
)

type URLInfo struct {
	Description string
	Path string
}

type SummerSchedule struct {
	Headers []string
	Data [][]string
}

func (s SummerSchedule) MarshalBinary() (data []byte, err error) {
	err = nil
	for _, headerLabel := range s.Headers {
		for _, b := range headerLabel {
			data = append(data, byte(b))
		}
	}

	for _, trainingBlock := range s.Data {
		for _, blockDetail := range trainingBlock {
			for _, b := range blockDetail {
				data = append(data, byte(b))
			}
		}
	}
	return
}

var username string
var password string
var url string
var localFilename string
var credentials string
var address string

func init() {
	flag.StringVar(&username, "username", "", "ONE OFF: Username")
	flag.StringVar(&password, "password", "", "ONE OFF: Password")
	flag.StringVar(&url, "url", "", "ONE OFF: URL to check on one off run.")
	flag.StringVar(&credentials, "credentials", "", "OPTIONAL: Local file with credentials. FORMAT: Username, newline, password, newline ...")
	flag.StringVar(&address, "urllist", "", "REQUIRED if using credential file: Local file with URLs to check. FORMAT: URL Description, newline, URL, newline ...")
	flag.StringVar(&localFilename, "local", "", "OPTIONAL: Local file to parse")
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	flag.Parse()

	//use loadFile if loading a stored page locally
	if (localFilename != "") {
		loadFile(localFilename)
	} else if (credentials != "") {
		loadCredentials(address, credentials)
	} else {
		//load url content
		loadURL(URLInfo{Description: "One off command", Path: url}, username, password)
	}  
}

func loadFile(filename string) {
	f, err := os.Open(filename)
	check(err)

	doc, err := html.Parse(f)
	check(err)

	fmt.Println(doc)
}

func loadCredentials(addressFilename string, credentialFilename string) {
	f, err := os.Open(addressFilename)
	check(err)

	r := bufio.NewReader(f)

	urlList := make([]URLInfo, 0)

	//Expect url desc and url path with newline inbetween 
	desc := ""
	path := ""
	for line, err := r.ReadString(10);  err == nil; line, err = r.ReadString(10) {
		if desc == "" {
			desc = line
		} else {
			path = line
			urlItem := URLInfo{Description: desc[0:len(desc) - 1], Path: path[0:len(path) -1]}
			urlList = append(urlList, urlItem)
			desc = ""
		}
	}

	f, err = os.Open(credentialFilename)
	check(err)

	r = bufio.NewReader(f)

	//Expect username and password with newline inbetween 
	username := ""
	password := ""
	for line, err := r.ReadString(10);  err == nil; line, err = r.ReadString(10) {
		if username == "" {
			username = line
		} else {
			password = line
			//fmt.Printf("%v:%v", username[0:len(username) - 1], password[0:len(password) - 1])

			for _, urlItem := range urlList {
				loadURL(urlItem, username[0:len(username) - 1], password[0:len(password) - 1])
			}
			
			username = ""
		}
	}
	check(err)
}

func loadURL(urlItem URLInfo, username, password string) {
	doc,err := html.Parse(fetch(urlItem.Path, username, password))
	check(err)

	var theSchedule SummerSchedule
	theSchedule.Data = make([][]string, 0)
	lookForTableWithAttr(&theSchedule, doc, "border")

	if len(theSchedule.Headers) == 0 {
		log.Fatal("Invalid credentials for " + username);
	}

	aggregateBuf := new(bytes.Buffer)

	//uncomment if you would like to see the output
	//fmt.Println(theSchedule)	

	buf, err := theSchedule.MarshalBinary()
	check(err)

	err = binary.Write(aggregateBuf, binary.LittleEndian, buf)
	check(err)

	//what is actually stored in the save file
	//fmt.Printf("%x", aggregateBuf.Bytes())

	err = os.Chdir(os.Getenv("HOME"))
	check(err)

	saveFilename := strings.Replace(username + "-" + urlItem.Description, " ", "", -1)

	_, err = os.Stat(saveFilename);
	if err == nil {
		existingData, err := ioutil.ReadFile(saveFilename)
		check(err)

		if bytes.Equal(aggregateBuf.Bytes(), existingData) == false {
			fmt.Println("Schedule for changed for " + username)

			//write changes to file
			f, err := os.Create(saveFilename)
			check(err)

			_, err = f.Write(aggregateBuf.Bytes())
			check(err)

			//in mail.go
			//mail(username + "@usna.edu", "m163876@mich302csd17u.academy.usna.edu", "MIDS " + urlItem.Description + " Changed", "View " + urlItem.Description + " at: " + urlItem.Path)
		} else {
			//fmt.Println("no changes")
		}
	} else if os.IsNotExist(err) {
		//write changes to file
		f, err := os.Create(saveFilename)
		check(err)

		_, err = f.Write(aggregateBuf.Bytes())
		check(err)

		fmt.Println("Made new save file at " + os.Getenv("HOME") + "/" + saveFilename)
	}
}

func fetch(url string, username string, password string) (io.ReadCloser) {
   //set transport options
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

   //create Client to control HTTP client headers
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", url, nil)

	req.SetBasicAuth(username, password)

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	return res.Body
}

func lookForTableWithAttr(theSchedule *SummerSchedule, n *html.Node, targetAttribute string) {
	

	if n.Type == html.ElementNode && n.DataAtom == atom.Table && len(n.Attr) > 0 && n.Attr[0].Key == targetAttribute  {
		for c := n.FirstChild; c!= nil; c = c.NextSibling {
			if c.Type == html.ElementNode && c.DataAtom == atom.Tbody {
				lookForTableValues(c, theSchedule)	

				
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		lookForTableWithAttr(theSchedule, c, targetAttribute)
	}
}

func lookForTableValues(n *html.Node, theSchedule *SummerSchedule){
	//look for header row
	if n.Type == html.ElementNode && n.DataAtom == atom.Th {
		for cth:= n.FirstChild; cth!= nil; cth = cth.NextSibling {
			if cth.Type == html.ElementNode && cth.DataAtom == atom.Font {
				for cfont:= cth.FirstChild; cfont!= nil; cfont = cfont.NextSibling {
					if cfont.Type == html.TextNode {
						theSchedule.Headers = append(theSchedule.Headers, cfont.Data)
					}
				}
			}
		}
	} else if n.Type == html.ElementNode && n.DataAtom == atom.Td {
		for cth:= n.FirstChild; cth!= nil; cth = cth.NextSibling {
			if cth.Type == html.ElementNode && cth.DataAtom == atom.Font {
				for cfont:= cth.FirstChild; cfont!= nil; cfont = cfont.NextSibling {
					if cfont.Type == html.TextNode {
						schedulesFound := 0
						if len(theSchedule.Data) - 1 < schedulesFound {
							theSchedule.Data = append(theSchedule.Data, make([]string, 0))
						}
						for len(theSchedule.Data[schedulesFound]) == len(theSchedule.Headers) {
							schedulesFound = schedulesFound + 1
							if len(theSchedule.Data) - 1 < schedulesFound {
								theSchedule.Data = append(theSchedule.Data, make([]string, 0))
								//fmt.Printf("%v inc", schedulesFound)
							}
						}
						//fmt.Printf("%v %v\n", len(theSchedule.Data), schedulesFound)
						
						//fmt.Printf("%v HEADER %v\n", len(theSchedule.Data[schedulesFound]), len(theSchedule.Headers))
						theSchedule.Data[schedulesFound] = append(theSchedule.Data[schedulesFound], cfont.Data)
						
						break
					}
				}
			}
		}
	} else { //we put this in the else because after we go through TDH/TR we do not want to loop TR again
		for c:= n.FirstChild; c!= nil; c = c.NextSibling {
			lookForTableValues(c, theSchedule)
		}
	}
}



