package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/ettle/strcase"
	"github.com/gertd/go-pluralize"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate common code patterns",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Give user options for type, name singular, name plural
		typeMenu := promptui.Select{
			Label: "Choose an option",
			Items: []string{"Public Collection", "Private Collection", "Global", "Exit"},
		}

		_, typeResult, err := typeMenu.Run()

		if err != nil {
			return err
		}

		if typeResult == "Exit" {
			os.Exit(0)
		}

		switch typeResult {
		case "Public Collection", "Private Collection":
			return collectionGenerate(typeResult)
		case "Global":
			return globalGenerate()
		}

		return nil
	},
}

func collectionGenerate(collection string) error {
	nameSingularPrompt := promptui.Prompt{
		Label: "Name (singular)",
		Validate: func(s string) error {
			if s == "" {
				return errors.New("can't be empty")
			}

			return nil
		},
	}

	nameSingularResult, err := nameSingularPrompt.Run()

	if err != nil {
		return err
	}

	pluralize := pluralize.NewClient()

	namePluralPrompt := promptui.Prompt{
		Label:     "Name (plural)",
		Default:   pluralize.Plural(nameSingularResult),
		AllowEdit: true,
		Validate: func(s string) error {
			if s == "" {
				return errors.New("can't be empty")
			}

			return nil
		},
	}

	namePluralResult, err := namePluralPrompt.Run()

	if err != nil {
		return err
	}

	// Set correct template location based on user input
	var templatePath string
	var outputDir string
	var data any

	switch collection {
	case "Public Collection":
		outputDir = "collections"
		templatePath = filepath.Join("templates", "public.collection.ts.gotmpl")
		data = struct {
			PluralCamelCase string
			PluralKebabCase string
		}{
			PluralCamelCase: strcase.ToCamel(namePluralResult),
			PluralKebabCase: strcase.ToKebab(namePluralResult),
		}

	case "Private Collection":
		outputDir = "collections"
		templatePath = filepath.Join("templates", "private.collection.ts.gotmpl")
		data = struct {
			PluralCamelCase string
			PluralKebabCase string
		}{
			PluralCamelCase: strcase.ToCamel(namePluralResult),
			PluralKebabCase: strcase.ToKebab(namePluralResult),
		}
	}

	// Load template
	tmpl := template.Must(template.ParseFS(TemplatesFS, templatePath))

	// Create paths
	fileName := fmt.Sprintf("%s.ts", strcase.ToKebab(namePluralResult))
	outputFile := filepath.Join("src", outputDir, fileName)

	// Check if file already exists
	if _, err := os.Stat(outputFile); !os.IsNotExist(err) {
		overrideConfirmPrompt := promptui.Prompt{
			Label:     "File already exists! Do you want to override it?",
			IsConfirm: true,
		}

		_, err := overrideConfirmPrompt.Run()

		if err != nil {
			return err
		}
	}

	// Create file
	f, err := os.Create(outputFile)

	if err != nil {
		return err
	}

	defer f.Close()

	// Populate file with template and data
	err = tmpl.Execute(f, data)

	if err != nil {
		return err
	}

	fmt.Printf("created file: %s \n", outputFile)

	return nil
}

func globalGenerate() error {
	namePrompt := promptui.Prompt{
		Label: "Name",
		Validate: func(s string) error {
			if s == "" {
				return errors.New("can't be empty")
			}

			return nil
		},
	}

	nameResult, err := namePrompt.Run()

	if err != nil {
		return err
	}

	nameCamelCase := strcase.ToCamel(nameResult)
	namePascalCase := strcase.ToPascal(nameResult)
	nameKebabCase := strcase.ToKebab(nameResult)
	nameAllCaps := strcase.ToCase(nameResult, strcase.UpperCase, '_')

	templatePathSettings := filepath.Join("templates", "settings.global.ts.gotmpl")
	templatePathAPI := filepath.Join("templates", "api.global.ts.gotmpl")

	data := struct {
		NameCamelCase  string
		NamePascalCase string
		NameAllCaps    string
		NameKebabCase  string
	}{
		NameCamelCase:  nameCamelCase,
		NamePascalCase: namePascalCase,
		NameAllCaps:    nameAllCaps,
		NameKebabCase:  nameKebabCase,
	}

	tmplSettings := template.Must(template.ParseFS(TemplatesFS, templatePathSettings))
	tmplApi := template.Must(template.ParseFS(TemplatesFS, templatePathAPI))

	fileNameSettings := fmt.Sprintf("%s-config.ts", nameKebabCase)
	fileNameApi := fmt.Sprintf("%s-api.ts", nameKebabCase)

	outputFolderPath := filepath.Join("src", "globals", nameKebabCase)

	outputFilePathSettings := filepath.Join("src", "globals", nameKebabCase, fileNameSettings)
	outputFilePathApi := filepath.Join("src", "globals", nameKebabCase, fileNameApi)

	if _, err := os.Stat(outputFilePathSettings); !os.IsNotExist(err) {
		overrideConfirmPrompt := promptui.Prompt{
			Label:     "File already exists! Do you want to override it?",
			IsConfirm: true,
		}

		_, err := overrideConfirmPrompt.Run()

		if err != nil {
			return err
		}
	}

	err = os.MkdirAll(outputFolderPath, 0700)

	if err != nil {
		return err
	}

	f, err := os.Create(outputFilePathSettings)

	if err != nil {
		return err
	}

	defer f.Close()

	f2, err := os.Create(outputFilePathApi)

	if err != nil {
		return err
	}

	defer f2.Close()

	// Populate file with template and data
	err = tmplSettings.Execute(f, data)

	if err != nil {
		return err
	}

	err = tmplApi.Execute(f2, data)

	if err != nil {
		return err
	}

	fmt.Printf("created files in: %s \n", outputFolderPath)

	return nil
}
