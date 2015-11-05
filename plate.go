package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func execute(infile, outfile string, payload interface{}) error {
	if err := os.MkdirAll(filepath.Dir(outfile), 0755); err != nil {
		return err
	}

	f, err := ioutil.TempFile(filepath.Dir(outfile), ".tmp")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())

	funcMap := template.FuncMap{
		"ToUpper":    strings.ToUpper,
		"ToLower":    strings.ToLower,
		"Replace":    strings.Replace,
		"ReplaceAll": func(a, b, c string) string { return strings.Replace(a, b, c, -1) },
	}

	name := filepath.Base(infile)
	t, err := template.New(name).Funcs(funcMap).ParseFiles(infile)
	if err != nil {
		return err
	}

	if err := t.Funcs(funcMap).ExecuteTemplate(f, name, payload); err != nil {
		return err
	}

	if err := os.Rename(f.Name(), outfile); err != nil {
		return err
	}
	if err := os.Chmod(outfile, 0644); err != nil {
		return err
	}
	return nil
}

func plate(args []string) {

	f := flag.NewFlagSet("load", flag.ExitOnError)
	f.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Apply go template using equipment yaml file as source\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  %s [options] template [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "General Options:\n")
		fmt.Fprintf(os.Stderr, "\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Equipment Template Options:\n")
		fmt.Fprintf(os.Stderr, "\n")
		f.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
	}

	var config string
	f.StringVar(&config, "config", "equipment.yaml", "equipment file to load")

	var input string
	f.StringVar(&input, "input", "", "default input directory")

	var strip int
	f.IntVar(&strip, "strip", 0, "number directories to strip from input if given")

	var output string
	f.StringVar(&output, "output", "/tmp", "default output directory")

	if err := f.Parse(args); err != nil {
		f.Usage()

		log.Fatalf("Invalid option(s) given")
	}

	c, err := ioutil.ReadFile(config)
	if err != nil {
		log.Fatal(err)
	}

	var places map[string]Place
	err = yaml.Unmarshal(c, &places)
	if err != nil {
		log.Fatal(err)
	}

	for _, infile := range f.Args() {
		if err := execute(infile, output+"/"+infile, places); err != nil {
			log.Fatal(err)
		}
	}

	if input != "" {
		var base string

		parts := strings.SplitAfter(filepath.Clean(input), "/")
		switch {
		case strip < 0 && (len(parts)+strip) > 0:
			base = filepath.Join(parts[len(parts)+strip : len(parts)]...)
		case strip > 0 && (len(parts)-strip) > 0:
			base = filepath.Join(parts[0 : len(parts)-strip+1]...)
		}

		err := filepath.Walk(input, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				outfile := filepath.Join(output, base, strings.TrimPrefix(filepath.Clean(path), filepath.Clean(input)))
				if err := execute(path, outfile, places); err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}
