package main

import (
	"sort"

	"github.com/charmbracelet/huh"
	"github.com/nanoteck137/forge/forgejo"
	"github.com/nanoteck137/forge/service"
)

func licenseSelectField(server, token, title string, value *string) huh.Field {
	templates, err := service.ListLicenseTemplates(service.ListTemplatesInput{
		Server: server,
		Token:  token,
	})
	if err != nil {
		return huh.NewInput().
			Title(title).
			Description("Could not fetch license templates: " + err.Error()).
			Value(value)
	}

	return huh.NewSelect[string]().
		Title(title).
		Options(licenseOptions(*value, templates)...).
		Filtering(true).
		Height(12).
		Value(value)
}

func gitignoreSelectField(server, token, title string, value *string) huh.Field {
	templates, err := service.ListGitignoreTemplates(service.ListTemplatesInput{
		Server: server,
		Token:  token,
	})
	if err != nil {
		return huh.NewInput().
			Title(title).
			Description("Could not fetch gitignore templates: " + err.Error()).
			Value(value)
	}

	return huh.NewSelect[string]().
		Title(title).
		Options(gitignoreOptions(*value, templates)...).
		Filtering(true).
		Height(12).
		Value(value)
}

func licenseOptions(current string, templates []forgejo.LicenseTemplate) []huh.Option[string] {
	keys := make([]string, 0, len(templates))
	seen := make(map[string]struct{}, len(templates))
	for _, t := range templates {
		key := t.Key
		if key == "" {
			key = t.Name
		}
		if key == "" {
			continue
		}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		keys = append(keys, key)
	}
	sort.Strings(keys)

	options := make([]huh.Option[string], 0, len(keys)+2)
	options = append(options, huh.NewOption("None", ""))
	for _, key := range keys {
		options = append(options, huh.NewOption(key, key))
	}
	if current != "" {
		if _, ok := seen[current]; !ok {
			options = append(options, huh.NewOption(current+" (current)", current))
		}
	}
	return options
}

func gitignoreOptions(current string, templates []string) []huh.Option[string] {
	names := make([]string, 0, len(templates))
	seen := make(map[string]struct{}, len(templates))
	for _, name := range templates {
		if name == "" {
			continue
		}
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		names = append(names, name)
	}
	sort.Strings(names)

	options := make([]huh.Option[string], 0, len(names)+2)
	options = append(options, huh.NewOption("None", ""))
	for _, name := range names {
		options = append(options, huh.NewOption(name, name))
	}
	if current != "" {
		if _, ok := seen[current]; !ok {
			options = append(options, huh.NewOption(current+" (current)", current))
		}
	}
	return options
}
