package main

import (
	"fmt"
	"os"
	"path"
	"text/template"
)

type MainGenerator struct {
}

func (m *MainGenerator) Run(opt *Option, metaData *ServiceMetaData) error {
	filename := path.Join("./", opt.Output, "main/main.go")
	fmt.Println(filename)
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Println("open file failed", filename)
		return err
	}
	defer file.Close()
	m.render(file, main_template, metaData)
	return nil
}

func (m *MainGenerator) render(file *os.File, s string, metaData *ServiceMetaData) error {
	t := template.New("main")
	t, err := t.Parse(s)
	if err != nil {
		return err
	}
	err = t.Execute(file, metaData)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	dir := &MainGenerator{}
	Register("main_generate", dir)
}
