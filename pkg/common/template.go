package common

import (
	"os"
	"path/filepath"
	"time"
)

type Template struct {
	Name string
	Data []byte
}

type TemplateLoader interface {
	Load() ([]Template, error)
}

func Loader(name string) (TemplateLoader, error) {
	fi, err := os.Stat(name)
	if err != nil {
		return Text(name), err
	}

	if fi.IsDir() {
		return Directory(name), nil
	}

	return File(name), nil
}

func Load(name string) ([]Template, error) {
	loader, err := Loader(name)

	if err != nil {
		return []Template{}, err
	}

	return loader.Load()
}

type Text string

func (t Text) Load() ([]Template, error) {
	return LoadText(string(t))
}

func LoadText(text string) ([]Template, error) {
	return []Template{
		{
			Name: time.Now().String(),
			Data: []byte(text),
		},
	}, nil
}

type Directory string

func (d Directory) Load() ([]Template, error) {
	return LoadDirectory(string(d))
}

func LoadDirectory(dir string) (templates []Template, err error) {
	root, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	files, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {

			template, err := LoadDirectory(filepath.Join(root, file.Name()))

			if err != nil {
				return nil, err
			}

			templates = append(templates, template...)

		} else {
			template, err := LoadFile(filepath.Join(root, file.Name()))

			if err != nil {
				return nil, err
			}

			templates = append(templates, template...)
		}
	}

	return templates, nil

}

type File string

func (f File) Load() ([]Template, error) {
	return LoadFile(string(f))
}

func LoadFile(name string) ([]Template, error) {
	file, err := filepath.Abs(name)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return []Template{
		{
			Name: name,
			Data: data,
		},
	}, nil
}
