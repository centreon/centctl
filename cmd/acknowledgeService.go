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
	"centctl/request"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// acknowledgeServiceCmd represents the service command
var acknowledgeServiceCmd = &cobra.Command{
	Use:   "service",
	Short: "Acknowledge services",
	Long:  `Acknowledge the service described right after`,
	Run: func(cmd *cobra.Command, args []string) {
		comment, _ := cmd.Flags().GetString("comment")
		hostName, _ := cmd.Flags().GetString("hostName")
		description, _ := cmd.Flags().GetString("description")
		debugV, _ := cmd.Flags().GetBool("DEBUG")
		err := AcknowledgeService(comment, hostName, description, debugV)
		if err != nil {
			fmt.Println(err)
		}
	},
}

//AcknowledgeService permits to acnowledge a service in the centreon server
func AcknowledgeService(comment string, hostName string, description string, debugV bool) error {
	//Creation of the request body
	values := "SVC;" + hostName + "," + description + ";" + comment + ";2;0;1"
	requestBody, err := request.CreateBodyRequest("add", "RTACKNOWLEDGEMENT", values)
	if err != nil {
		return err
	}

	//Recovery of the response body
	urlCentreon := os.Getenv("URL") + "/api/index.php?action=action&object=centreon_clapi"
	client := request.NewClient(os.Getenv("URL") + "/api/index.php?action=action&object=centreon_clapi")
	statusCode, body, err := client.CentreonCLAPI(requestBody)

	//If flag debug, print informations about the request API
	if debugV {
		debug.Show("acknowledge service", string(requestBody), urlCentreon, statusCode, body)
	}
	if err != nil {
		return err
	}

	//Verification with the response body that the acknowledge was carried out
	if string(body) != "{\"result\":[]}" {
		fmt.Println("erreur: ", string(body))
		os.Exit(1)
	}

	fmt.Printf("The service `%v` is acknowledged\n", description)
	return nil
}

func init() {
	acknowledgeCmd.AddCommand(acknowledgeServiceCmd)
	acknowledgeServiceCmd.Flags().StringP("hostName", "n", "", "To know the host to wich the service is attached")
	acknowledgeServiceCmd.MarkFlagRequired("hostName")
	acknowledgeServiceCmd.Flags().StringP("description", "d", "", "To know the service to be acknowledge")
	acknowledgeServiceCmd.MarkFlagRequired("description")
}
