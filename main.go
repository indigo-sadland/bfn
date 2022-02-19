package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/projectdiscovery/gologger"
)

type Naabu struct {
	Host string `json:"host"`
	Ip   string `json:"ip"`
	Port int    `json:"port"`
}

type HostInfo struct {
	Name  string
	Ip    string
	Ports []string
}

var (
	file     *string
	ipsOut   *bool
	portsOut *bool
)

func parse(file string) {

	var naabu = new(Naabu)
	var hostInfo = new(HostInfo)
	var hosts []HostInfo
	var ports []string

	body, err := ioutil.ReadFile(file)
	if err != nil {
		gologger.Error().Msgf(err.Error())
		return
	}

	numlines := strings.Count(string(body), "\n")

	h := ""
	for i, s := range strings.Split(string(body), "\n") {

		_ = json.Unmarshal([]byte(s), &naabu)

		// start point for the first element (host name) from the range.
		// Basically here we creat first element for []HostInfo
		if h == "" {
			h = naabu.Host

			hostInfo.Name = naabu.Host
			hostInfo.Ip = naabu.Ip
			ports = append(ports, strconv.Itoa(naabu.Port))
			hostInfo.Ports = ports

			continue
		}

		// if the next element has the same name
		// then we only append ports.
		if h == naabu.Host {
			ports = append(ports, strconv.Itoa(naabu.Port))
			hostInfo.Ports = ports

			// Else we detect new host name, so we need to append the previous HostInfo object
			// and start the steps for current host anew.
		} else {
			hosts = append(hosts, *hostInfo)
			ports = []string{}

			h = naabu.Host
			hostInfo.Name = naabu.Host
			hostInfo.Ip = naabu.Ip
			ports = append(ports, strconv.Itoa(naabu.Port))
			hostInfo.Ports = ports

		}

		// Check if the current element is last in the range.
		if i == numlines {
			hosts = append(hosts, *hostInfo)
		}
	}

	prepareOutput(hosts)
}

func prepareOutput(hosts []HostInfo) {

	var fileName string

	// Prepare hosts information for table output.
	var data [][]string
	var hostInfo []string
	var ips []string
	var ports []string

	for _, h := range hosts {

		// fill the slices for output with flag -ips or -ports.
		ips = append(ips, h.Ip)
		ports = append(ports, strings.Join(h.Ports, "\n"))

		hostInfo = []string{h.Name, h.Ip, strings.Join(h.Ports, "\n")}
		data = append(data, hostInfo)
	}

	// Create table output.
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NAME", "IP", "OPEN PORTS"})
	table.SetRowLine(true)
	table.SetCenterSeparator("-")
	table.SetColumnSeparator("|")
	table.SetRowSeparator("-")
	table.AppendBulk(data)
	table.Render()

	// Prepare fileName.
	u, err := url.Parse("https://" + hosts[0].Name)
	if err != nil {
		gologger.Error().Msgf(err.Error())
	}

	hostParts := strings.Split(u.Host, ".")

	length := len(hostParts)

	// cases:
	// A. site.com  -> length : 2
	// B. www.site.com -> length : 3
	// C. www.hello.site.com -> length : 4

	switch length {
	case 2:
		fileName = hostParts[0]
	case 3:
		fileName = hostParts[1]
	case 4:
		fileName = hostParts[2]
	}

	if *ipsOut {
		fileName = fileName + "_ips.txt"
		writeToFile(fileName, removeDuplicates(ips))
		fmt.Printf("\nIPs were successfully written to ./%s\n", fileName)
	}

	if *portsOut {
		fileName = fileName + "_ports.txt"
		writeToFile(fileName, removeDuplicates(ports))
		fmt.Printf("\nPorts were successfully written to ./%s\n", fileName)
	}

}

func writeToFile(fileName string, data []string) {

	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		gologger.Error().Msgf(err.Error())
		return
	}

	wrt := bufio.NewWriter(f)

	for _, line := range data {
		_, err2 := wrt.WriteString(line + "\n")
		if err2 != nil {
			gologger.Error().Msgf(err.Error())
			return
		}

	}

	wrt.Flush()
	f.Close()

}

func removeDuplicates(rawData []string) []string {

	var uniqueList []string
	var data []string

	keys := make(map[string]bool)

	// After joining ports there are may occur strings like "80\n443". We need to slit them.
	for _, s := range rawData {
		if strings.Contains(s, "\n") {
			ns := strings.Split(s, "\n")
			for _, e := range ns {
				data = append(data, e)
			}
		} else {
			data = append(data, s)
		}
	}

	// If the key(values of the slice) is not equal
	// to the already present value in new slice (list)
	// then we append it. Else we jump on another element.
	for _, elem := range data {

		if _, value := keys[elem]; !value {
			keys[elem] = true
			if elem == "" { // remove empty string
				continue
			}
			uniqueList = append(uniqueList, elem)
		}
	}

	return uniqueList
}

func main() {

	file = flag.String("i", "", "")
	ipsOut = flag.Bool("ips", false, "")
	portsOut = flag.Bool("ports", false, "")
	flag.Usage = func() {
		fmt.Printf("Usage:\n\t" +
			"-i, <INPUT_FILE>       Define file with naabu output.\n\t" +
			"-ips                   Save all IPs to a file if specified (optional).\n\t" +
			"-ports                 Save all ports to a file if specified (optional).",
		)
	}
	flag.Parse()

	if *file == "" {
		flag.Usage()
		return
	}

	parse(*file)

}
