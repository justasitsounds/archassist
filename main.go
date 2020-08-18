package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/manifoldco/promptui"
)

type awsservice struct {
	Name        string
	Identifier  string
	Description string
	IncludePath string
	Boxed       bool
}

//Diagram represents the data required for a plantUML diagram
type Diagram struct {
	Title    string
	Services []awsservice
}

//EXIT is the select option that signifies moving to the next prompt
const EXIT string = "-- DONE --"
const awsIconDist string = "AWSICONDIST"

func listLocalIncludes(includePath string) ([]awsservice, error) {
	srvs := make([]awsservice, 0)
	srvs = append(srvs, awsservice{Name: EXIT})
	err := filepath.Walk(includePath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			terminalFolder := strings.TrimPrefix(filepath.Dir(path), includePath)
			if !info.IsDir() && strings.HasSuffix(info.Name(), ".puml") && terminalFolder != "" && !strings.EqualFold(info.Name(), "all.puml") {
				srv := awsservice{
					Name:        strings.TrimSuffix(info.Name(), ".puml"),
					IncludePath: terminalFolder[1:],
				}
				srvs = append(srvs, srv)
			}
			return nil
		})
	return srvs, err
}

func serviceSelector(services []awsservice) promptui.Select {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U000027A1 {{ .IncludePath | cyan }}/{{ .Name | cyan }}",
		Inactive: "  {{ .IncludePath | cyan }}/{{ .Name | cyan }}",
		Selected: "\U000027A1 {{ .IncludePath | red | cyan }}/{{ .Name | red | cyan }}",
		Details: `
--------- Service ----------
{{ "Service:" | faint }}	{{ .Name }}
{{ "Category:" | faint }}	{{ .IncludePath }}`,
	}

	searcher := func(input string, index int) bool {
		pepper := services[index]
		name := strings.Replace(strings.ToLower(pepper.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	return promptui.Select{
		Label:             "AWS Services",
		Items:             services,
		Templates:         templates,
		Size:              8,
		HideSelected:      true,
		StartInSearchMode: true,
		Searcher:          searcher,
		Stdout:            &bellSkipper{},
	}
}

func loadTemplate() (*template.Template, error) {
	return template.New("template.puml").Funcs(template.FuncMap{
		"toLower": strings.ToLower,
		"unique": func(services []awsservice) []awsservice {
			uniques := make([]awsservice, 0)
			m := map[string]bool{}
			for _, v := range services {
				if !m[v.Name] {
					m[v.Name] = true
					uniques = append(uniques, v)
				}
			}
			return uniques
		},
		"notVPC": func(serviceName string) bool {
			return !strings.Contains(serviceName, "VPC")
		},
	}).ParseFiles("template.puml")
}

func autoname(srvcs []awsservice, serviceName string) string {
	nameCount := 0
	for _, v := range srvcs {
		if v.Name == serviceName {
			nameCount++
		}
	}
	return fmt.Sprintf("%s%v", serviceName, nameCount)
}

func main() {

	iconDistLocation := os.Getenv(awsIconDist)
	if iconDistLocation == "" {
		log.Fatalf("Environment variable: '%s' not set", awsIconDist)
	}

	diagram := Diagram{}
	stopSelectingServices := false
	services, err := listLocalIncludes(iconDistLocation)
	if err != nil {
		log.Fatalf("Loading service diagrams failed %v\n", err)
	}

	prompt := promptui.Prompt{
		Label:   "Diagram name",
		Default: "diagram",
	}
	diagram.Title, err = prompt.Run()
	if err != nil {
		log.Fatalf("Setting diagram name failed %v\n", err)
		return
	}

	selectedServices := make([]awsservice, 0)
	selector := serviceSelector(services)

	for !stopSelectingServices {
		i, _, err := selector.Run()
		if err != nil {
			log.Fatalf("Prompt failed %v\n", err)
		}

		stopSelectingServices = services[i].Name == EXIT
		if stopSelectingServices {
			continue
		}

		selectedService := services[i]
		selectedService.Identifier = autoname(selectedServices, selectedService.Name)

		descPrompt := promptui.Prompt{
			Label:   "service description",
			Default: selectedService.Name,
		}
		description, err := descPrompt.Run()
		if err != nil {
			log.Fatalf("Setting description failed %v\n", err)
			return
		}
		selectedService.Description = description

		selectedServices = append(selectedServices, selectedService)
		fmt.Printf("You chose %s/%s : %s\n", selectedService.IncludePath, selectedService.Name, selectedService.Description)
	}
	diagram.Services = selectedServices

	te, err := loadTemplate()
	if err != nil {
		log.Fatal("error loading template", err)
	}

	diag, err := os.Create(fmt.Sprintf("./%s.puml", diagram.Title))
	if err != nil {
		log.Fatal(err)
	}
	defer diag.Close()

	err = te.Execute(diag, diagram)
	if err != nil {
		log.Fatal("error executing template", err)
	}
}
