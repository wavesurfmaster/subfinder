// 
// virustotal.go : A Virustotal Client for Subdomain Enumeration
// Written By : @ice3man (Nizamul Rana)
// 
// Distributed Under MIT License
// Copyrights (C) 2018 Ice3man
//

// NOTE : We are using Virustotal API here Since we wanted to eliminate the 
// rate limiting performed by Virustotal on scraping.
// Direct queries and parsing can be also done :-)

package virustotal

import (
	"io/ioutil"
	"fmt"
	"encoding/json"
	"strings"

	"subfinder/libsubfinder/helper"
)

type virustotalapi_object struct {
	Subdomains	[]string `json:"subdomains"`
}

var virustotalapi_data virustotalapi_object

// 
// Local function to query virustotal API
// Requires an API key
//
// @note : If the user specifies an API key in config.json, we use API
//	If not, we try to scrape pages though it is highly discouraged
//
func queryVirustotalApi(state *helper.State) (subdomains []string, err error) {

	// Make a search for a domain name and get HTTP Response
	resp, err := helper.GetHTTPResponse("https://www.virustotal.com/vtapi/v2/domain/report?apikey="+state.ConfigState.VirustotalAPIKey+"&domain="+state.Domain, 3000)
	if err != nil {
		return subdomains, err
	}

	// Get the response body
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return subdomains, err
	}

	// Decode the json format
	err = json.Unmarshal([]byte(resp_body), &virustotalapi_data)
	if err != nil {
		return subdomains, err
	}

	// Append each subdomain found to subdomains array
	for _, subdomain := range virustotalapi_data.Subdomains {

		// Fix Wildcard subdomains containg asterisk before them
		if strings.Contains(subdomain, "*.") {
			subdomain = strings.Split(subdomain, "*.")[1]
		}

		if state.Verbose == true {
			if state.Color == true {
				fmt.Printf("\n[%sVIRUSTOTAL%s] %s", helper.Red, helper.Reset, subdomain)
			} else {
				fmt.Printf("\n[VIRUSTOTAL] %s", subdomain)
			}
		}

		subdomains = append(subdomains, subdomain)
	}	

	return subdomains, nil
}

// 
// Query : Queries awesome Virustotal Service for Subdomains
// @param state : Current application state
// 
// @return subdomain : String array containing subdomains found
// @return err : nil if successfull and error if failed
//
func Query(state *helper.State) (subdomains []string, err error) {

	// We have recieved an API Key
	// Now, we will use Virustotal API key to fetch subdomain info
	if state.ConfigState.VirustotalAPIKey != "" {

		// Get subdomains via API
		subdomains, err = queryVirustotalApi(state)

		if err != nil {
			return subdomains, err
		}
	}

	return subdomains, nil
}