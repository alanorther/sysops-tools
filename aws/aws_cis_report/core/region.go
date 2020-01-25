package core

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"os"
	"strings"
)

type menu struct {
	Name   string
	Number int
}

var RegionSelection string

func Region(AwsReg string) string {

	region_menu := []menu{
		{Name: "us-east-1", Number: 0},
		{Name: "us-east-2", Number: 1},
		{Name: "us-west-1", Number: 2},
		{Name: "us-west-2", Number: 3},
		{Name: "ap-south-1", Number: 4},
		{Name: "ap-northeast-2", Number: 5},
		{Name: "ap-northeast-3", Number: 6},
		{Name: "ap-southeast-1", Number: 7},
		{Name: "ap-southeast-2", Number: 8},
		{Name: "ap-northeast-1", Number: 9},
		{Name: "ca-central-1", Number: 10},
		{Name: "cn-north-1", Number: 11},
		{Name: "eu-central-1", Number: 12},
		{Name: "eu-west-1", Number: 13},
		{Name: "eu-west-2", Number: 14},
		{Name: "eu-west-3", Number: 15},
		{Name: "sa-east-1", Number: 16},
		{Name: "Main Menu", Number: 17},
		{Name: "Exit", Number: 18},
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F300 ({{.Number | red }}) {{ .Name | cyan }}",
		Inactive: " ({{ .Number | red }}) {{ .Name | cyan }}",
		Selected: "\U0001F300 {{ .Name | red | cyan }}",
		Details: `
--------- Option ----------
{{ "Name:" | faint }}   {{ .Name }}
{{ "Option:" | faint }} {{ .Number }}`,
	}

	prompt := promptui.Select{
		Label:     "Select Option",
		Items:     region_menu,
		Templates: templates,
		Size:      4,
	}

	i, _, err := prompt.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	switch region_menu[i].Number {
	case 0:
		//fmt.Println("us-east-1")
		RegionSelection = strings.TrimSpace(region_menu[i].Name)
	case 1:
		RegionSelection = strings.TrimSpace(region_menu[i].Name)
	case 2:
		RegionSelection = strings.TrimSpace(region_menu[i].Name)
	case 3:
		RegionSelection = strings.TrimSpace(region_menu[i].Name)
	case 4:
		RegionSelection = strings.TrimSpace(region_menu[i].Name)
	case 5:
		RegionSelection = strings.TrimSpace(region_menu[i].Name)
	case 6:
		RegionSelection = strings.TrimSpace(region_menu[i].Name)
	case 7:
		RegionSelection = strings.TrimSpace(region_menu[i].Name)
	case 8:
		RegionSelection = strings.TrimSpace(region_menu[i].Name)
	case 9:
		RegionSelection = strings.TrimSpace(region_menu[i].Name)
	case 10:
		RegionSelection = strings.TrimSpace(region_menu[i].Name)
	case 11:
		RegionSelection = strings.TrimSpace(region_menu[i].Name)
	case 12:
		RegionSelection = strings.TrimSpace(region_menu[i].Name)
	case 13:
		RegionSelection = strings.TrimSpace(region_menu[i].Name)
	case 14:
		RegionSelection = strings.TrimSpace(region_menu[i].Name)
	case 15:
		RegionSelection = strings.TrimSpace(region_menu[i].Name)
	case 16:
		RegionSelection = strings.TrimSpace(region_menu[i].Name)
	case 17:
		return AwsReg
	case 18:
		os.Exit(0)
	}

	return RegionSelection

}
