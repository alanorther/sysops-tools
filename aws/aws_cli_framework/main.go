package main

import (
	//	"bufio"
	"fmt"
	"github.com/VerveWireless/sysops-tools/aws/aws_cli_framework/core"
	"github.com/aws/aws-sdk-go/aws/credentials"
	//	"github.com/ischisler/aws-turbo/core"
	"github.com/manifoldco/promptui"
	"log"
	"os"
	"os/user"
	"strings"
)

//set environment to default
var AwsEnv = "default"
var AwsReg = "us-east-1"

type menu struct {
	Name   string
	Number int
}

func main() {
	//	banner := core.ASCIIBanner
	fmt.Printf("%v\n\n", core.ASCIIBanner)
	fmt.Printf("Website: %v\n", core.Website)
	fmt.Printf("Version: %v\n", core.Version)

	options := []menu{
		{Name: "Change Environment", Number: 0},
		{Name: "Change Region", Number: 1},
		{Name: "IAM", Number: 2},
		{Name: "Exit", Number: 3},
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F300 ({{.Number | red }}) {{ .Name | cyan }}",
		Inactive: " ({{ .Number | red }}) {{ .Name | cyan }}",
		Selected: "\U0001F300 {{ .Name | red | cyan }}",
		Details: `
--------- Option ----------
{{ "Name:" | faint }}	{{ .Name }}
{{ "Option:" | faint }}	{{ .Number }}`,
	}

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

Menu:
	//load credentials file
	creds := credentials.NewSharedCredentials(usr.HomeDir+"/.aws/credentials", strings.TrimSpace(AwsEnv))

	credValue, err := creds.Get()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//credValue.SecretAccessKey, credValue.SessionToken, credValue.ProviderName
	fmt.Println("Current environment: ", AwsEnv, "Access Key ID: ", credValue.AccessKeyID, "Region: ", AwsReg)
	//	fmt.Println("Current Region: ", AwsReg)
	prompt := promptui.Select{
		Label:     "Select Option",
		Items:     options,
		Templates: templates,
		Size:      4,
	}

	i, _, err := prompt.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	switch options[i].Name {
	case "Change Region":
		AwsReg = core.Region(AwsReg)
		goto Menu
	case "Change Environment":
		//reader := bufio.NewReader(os.Stdin)
		//fmt.Printf("Which environment would you like? \n")
		//AwsEnv, _ = reader.ReadString('\n')
		AwsEnv = core.Environment()
		//os.Exit(0)
		goto Menu
	case "IAM":
		core.Iam(AwsEnv, AwsReg)
		goto Menu
	case "Exit":
		os.Exit(0)
	}

	//fmt.Printf("Selected: %d: %s\n ", i+1, options[i].Name)

}
