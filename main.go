package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"text/template"
)

var (
	fp = flag.String("fp", "", "")
)

// Config for generate enum struct
type Config struct {
	Type     string
	Name     string
	Elements []*Pair
}

// Pair struct
type Pair struct {
	N string
	V string
}

func main() {
	flag.Parse()

	if *fp == "" {
		log.Printf("usage: go run main.go -fp=[config file path]")
		return
	}

	f, err := os.Open(*fp)
	if err != nil {
		log.Printf("fail to open file %v, error: %v", *fp, err)
		return
	}
	defer f.Close()

	cfg, err := parseFile(f)
	if err != nil {
		log.Printf("parse content fail, error: %v", err)
		return
	}

	err = checkConfig(cfg)
	if err != nil {
		log.Printf("check config fail, error %v", err)
	}

	tmpl, err := template.New("enum").Parse(templateText)
	if err != nil {
		log.Printf("create template fail, error: %v", err)
		return
	}

	nf, err := os.Create("enum.go")
	if err != nil {
		log.Printf("create file fail, error: %v", err)
		return
	}
	defer nf.Close()

	err = tmpl.Execute(nf, cfg)
	if err != nil {
		log.Printf("execute fail, error: %v", err)
		return
	}

	log.Printf("DONE")
	return
}

type contentMode struct {
	empty    int
	dataType int
	name     int
	list     int
}

var (
	modeEnum = contentMode{
		empty:    0,
		dataType: 1,
		name:     2,
		list:     3,
	}
)

func parseFile(f *os.File) (*Config, error) {
	r := bufio.NewReader(f)
	cfg := &Config{}
	mode := modeEnum.empty
	section := []string{}
	for {
		s, err := r.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("read string fail, error: %v", err)
			return nil, err
		}
		nextMode := mode
		switch {
		case strings.HasPrefix(s, "type:"):
			nextMode = modeEnum.dataType
		case strings.HasPrefix(s, "name:"):
			nextMode = modeEnum.name
		case strings.HasPrefix(s, "list:"):
			nextMode = modeEnum.list
		default:
			// Do nothing
		}

		if nextMode != mode {
			if mode != modeEnum.empty {
				err = parse(mode, section, cfg)
				if err != nil {
					log.Printf("parse fail, mode: %v, error: %v", mode, err)
					return nil, err
				}
			}
			section = []string{}
			mode = nextMode
			continue
		}

		section = append(section, s)
	}
	err := parse(mode, section, cfg)
	if err != nil {
		log.Printf("parse fail, mode: %v, error: %v", mode, err)
		return nil, err
	}

	return cfg, nil
}

func parse(mode int, section []string, cfg *Config) error {
	switch mode {
	case modeEnum.dataType:
		if len(section) != 1 {
			return fmt.Errorf("wrong setting for data type")
		}
		v := trim(section[0])
		cfg.Type = v
	case modeEnum.name:
		if len(section) != 1 {
			return fmt.Errorf("wrong setting for data type")
		}
		v := trim(section[0])
		cfg.Name = v
	case modeEnum.list:
		for _, line := range section {
			l := trim(line)
			ts := strings.SplitN(l, "=", 2)
			p := &Pair{
				N: trim(ts[0]),
			}
			if len(ts) == 2 {
				p.V = trim(ts[1])
			}
			cfg.Elements = append(cfg.Elements, p)
		}
	default:
		return fmt.Errorf("no this mode %v", mode)
	}

	return nil
}

func checkConfig(cfg *Config) error {
	cfg.Name = strings.ToLower(cfg.Name)
	cfg.Type = strings.ToLower(cfg.Type)

	// only support string and int
	if cfg.Type != "string" && cfg.Type != "int" {
		err := fmt.Errorf("not support data type: %v", cfg.Type)
		log.Print(err.Error())
		return err
	}

	setElementsNameCapital(cfg)
	setElementsValue(cfg)
	return nil
}

func setElementsNameCapital(cfg *Config) {
	for i := range cfg.Elements {
		cfg.Elements[i].N = strings.Title(cfg.Elements[i].N)
	}
}

func setElementsValue(cfg *Config) {
	if cfg.Type == "string" {
		return
	}

	inc := 0
	for i := range cfg.Elements {
		if cfg.Elements[i].V == "iota" {
			inc = 0
		}

		cfg.Elements[i].V = strconv.Itoa(inc)
		inc++
	}
}

func trim(s string) string {
	return strings.Trim(s, "\n\t ")
}

var templateText = `package {{.Name}}

// Alias hide the real type of the enum 
// and users can use it to define the var for accepting enum 
type Alias = {{.Type}}

type list struct { {{range $i, $v := (.Elements)}}
    {{$v.N}} Alias{{end}}
}

// Enum for public use
var Enum = &list{ {{range $i, $v := (.Elements)}}
	{{$v.N}}: {{$v.V}},{{end}}
}
`
