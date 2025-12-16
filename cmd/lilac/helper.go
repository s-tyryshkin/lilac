package main

import (
	"path/filepath"

	bashpkg "github.com/s-tyryshkin/lilac/internal/bash"
	mergepkg "github.com/s-tyryshkin/lilac/internal/merge"
	templatepkg "github.com/s-tyryshkin/lilac/internal/template"
	valuespkg "github.com/s-tyryshkin/lilac/internal/values"
)

func Render(valuesPaths []string, scriptsPaths []string) ([]string, error) {
	values := make(valuespkg.Values)
	for _, path := range valuesPaths {
		valuesDst, err := valuespkg.ValuesRead(path)
		if err != nil {
			return nil, err
		}

		values = mergepkg.Merge(values, valuesDst)
	}

	result := []string{}
	for _, path := range scriptsPaths {
		lines, err := templatepkg.TemplateRender(path, values)
		if err != nil {
			return nil, err
		}

		tmpPath, err := bashpkg.BashPreprocess(filepath.Dir(path), lines, values)
		if err != nil {
			return nil, err
		}

		result = append(result, tmpPath)
	}

	return result, nil
}
