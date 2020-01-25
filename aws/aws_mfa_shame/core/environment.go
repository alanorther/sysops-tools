package core

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strings"
)

type RegMenu struct {
	Name   string
	Number int
}

var EnvSelection string

func Environment() string {
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

	var env_menu []RegMenu

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	dat, err := ioutil.ReadFile(usr.HomeDir + "/.aws/credentials")

	if err != nil {
		log.Fatal(err)
	}

	s := strings.Split(string(dat), "\n")

	for index, each := range s {
		if strings.Contains(each, "aws_access_key_id") || strings.Contains(each, "aws_secret_access_key") || (each == "") {
			continue
		} else {
			t := strings.Trim(each, "[]")
			//env_menu = append(env_menu, t, index)
			item := RegMenu{Name: t, Number: index}
			env_menu = append(env_menu, item)

			//env_menu = append(env_menu, t, index)
			//              item_now := RegMenu{Name: t, Number: index}
			//              box.AddItem(item_now)
			//              fmt.Println(len(box.Items))
			//fmt.Printf("value [%d] is %s\n", index, t)
		}
	}

	prompt := promptui.Select{
		Label:     "Select Option",
		Items:     env_menu,
		Templates: templates,
		Size:      4,
	}

	i, _, err := prompt.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	EnvSelection = env_menu[i].Name

	return EnvSelection
}
