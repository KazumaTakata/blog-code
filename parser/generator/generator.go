package main

import (
	"io/ioutil"
	"os"
	"text/template"
)

type Package struct {
	Name      string
	Functions []Function
	Structs   []Struct
}

type Function struct {
	Name string
}

type Struct struct {
	Name string
}

func main() {

	package_temp := Package{Name: "main", Functions: []Function{}}
	package_temp.Functions = append(package_temp.Functions, Function{Name: "main"})

	package_temp.Structs = append(package_temp.Structs, Struct{Name: "Node"})

	buf, _ := ioutil.ReadFile("parser/generator/sample.tmpl")

	tmpl, err := template.New("package").Parse(string(buf))
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stdout, package_temp)
	if err != nil {
		panic(err)
	}

}
