/*
MIT License

Copyright (c) 2020 YPSI SAS
Centctl is developped by : Mélissa Bertin

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package cmd

import (
	"centctl/debug"
	"centctl/display"
	"centctl/request"
	"centctl/service"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// showServiceCmd represents the service command
var showServiceCmd = &cobra.Command{
	Use:   "service",
	Short: "Show one service's details ",
	Long:  `Show one service of the Centreon server`,
	Run: func(cmd *cobra.Command, args []string) {
		nameHost, _ := cmd.Flags().GetString("nameHost")
		description, _ := cmd.Flags().GetString("description")
		debugV, _ := cmd.Flags().GetBool("DEBUG")
		output, _ := cmd.Flags().GetString("output")
		err := ShowService(nameHost, description, debugV, output)
		if err != nil {
			fmt.Println(err)
		}
	},
}

//ShowService permits to display the details of one service
func ShowService(nameHost string, description string, debugV bool, output string) error {
	output = strings.ToLower(output)

	//Recovery of the response body
	urlCentreon := os.Getenv("URL") + "/api/index.php?action=list&object=centreon_realtime_services&searchHost=" + nameHost + "&search=" + description + "&fields=id,description,host_id,host_name,state,state_type,output,perfdata,max_check_attempts,current_attempt,next_check,last_update,last_check,last_state_change,last_hard_state_change,acknowledged,active_checks,instance,criticality,passive_checks,notify,scheduled_downtime_depth"
	client := request.NewClient(urlCentreon)
	statusCode, body, err := client.Get()

	//If flag debug, print informations about the request API
	if debugV {
		debug.Show("show service", "", urlCentreon, statusCode, body)
	}
	if err != nil {
		return err
	}

	//Permits to recover the services contain into the response body
	var services []service.DetailService
	json.Unmarshal(body, &services)

	if len(services) == 0 {
		fmt.Println("no host or service with this name")
		os.Exit(1)
	}

	//Permits to find the good service in the array
	var serviceFind service.DetailService
	for _, v := range services {
		if v.Description == description {
			serviceFind = v
		}
	}

	server := service.DetailServer{
		Server: service.DetailInformations{
			Name:    os.Getenv("SERVER"),
			Service: serviceFind,
		},
	}

	//Display all services
	displayService, err := display.DetailService(output, server)
	if err != nil {
		return err
	}
	fmt.Println(displayService)
	return nil
}

func init() {
	showCmd.AddCommand(showServiceCmd)
	showServiceCmd.Flags().StringP("nameHost", "n", "", "Name of the host wich the service is attached")
	showServiceCmd.Flags().StringP("description", "d", "", "Description of the service")
	showServiceCmd.MarkFlagRequired("nameHost")
	showServiceCmd.MarkFlagRequired("description")
}
