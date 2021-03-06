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
	"strconv"

	"github.com/spf13/cobra"
)

// addContactCmd represents the contact command
var addContactCmd = &cobra.Command{
	Use:   "contact",
	Short: "Add a contact",
	Long:  `Add a contact into the Centreon server specifified by the flag --server`,
	Run: func(cmd *cobra.Command, args []string) {
		fullName, _ := cmd.Flags().GetString("fullName")
		login, _ := cmd.Flags().GetString("login")
		email, _ := cmd.Flags().GetString("email")
		password, _ := cmd.Flags().GetString("password")
		admin, _ := cmd.Flags().GetBool("admin")
		debugV, _ := cmd.Flags().GetBool("DEBUG")
		err := AddContact(fullName, login, email, password, admin, debugV)
		if err != nil {
			fmt.Println(err)
		}
	},
}

//AddContact permits to add a contact in the centreon server
func AddContact(fullName string, login string, email string, password string, admin bool, debugV bool) error {
	//Transformation of the admin value to int
	adminVal := 0
	if admin {
		adminVal = 1
	}

	//Creation of the request body
	values := fullName + ";" + login + ";" + email + ";" + password + ";" + strconv.Itoa(adminVal) + ";1;browser;local"
	requestBody, err := request.CreateBodyRequest("add", "contact", values)
	if err != nil {
		return err
	}

	//Recovery of the response body
	urlCentreon := os.Getenv("URL") + "/api/index.php?action=action&object=centreon_clapi"
	client := request.NewClient(os.Getenv("URL") + "/api/index.php?action=action&object=centreon_clapi")
	statusCode, body, err := client.CentreonCLAPI(requestBody)

	//If flag debug, print informations about the request API
	if debugV {
		debug.Show("add contact", string(requestBody), urlCentreon, statusCode, body)
	}
	if err != nil {
		return err
	}

	//Verification with the response body that the contact was created out
	if string(body) != "{\"result\":[]}" {
		fmt.Println("erreur: ", string(body))
		os.Exit(1)
	}

	fmt.Printf("The contact %v is created\n", fullName)
	return nil
}

func init() {
	addCmd.AddCommand(addContactCmd)
	addContactCmd.Flags().StringP("fullName", "f", "", "To define the full name of the contact")
	addContactCmd.MarkFlagRequired("fullName")
	addContactCmd.Flags().StringP("login", "l", "", "To define the login of the contact")
	addContactCmd.MarkFlagRequired("login")
	addContactCmd.Flags().StringP("email", "e", "", "To define the email of the contact")
	addContactCmd.MarkFlagRequired("email")
	addContactCmd.Flags().StringP("password", "p", "", "To define the password of the contact")
	addContactCmd.MarkFlagRequired("password")
	addContactCmd.Flags().Bool("admin", false, "To define if the contact is admin")
}
