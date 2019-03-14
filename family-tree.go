// Generate and manage your family tree
package main

import (
    "log"
    "net/http"
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "html/template"
    "github.com/x1um1n/checkerr"
)

type PageVariables struct {
    PageTitle     string
}

type Person struct {
    Name          string    `yaml:"name"`
    UUID          string    `yaml:"id"`
    Fnam          string    `yaml:"fnam"`
    Mnam          []string  `yaml:"mnam"`
    Snam          string    `yaml:"snam"`
    Sex           string    `yaml:"sex"`
    DOB           string    `yaml:"dob"`
    DOD           string    `yaml:"dod"`
    Children      []string  `yaml:"children"`
    Parents       []string  `yaml:"parents"`
    Spouse        string    `yaml:"spouse"`
}

var family []Person

// read yaml from file f and return it marshalled into a []Person
func loadYaml(f string) (fam []Person) {
    y, err := ioutil.ReadFile(f)
    checkerr.CheckFatal(err, "Error reading yaml file", f)

    err = yaml.Unmarshal(y, &fam)
    checkerr.CheckFatal(err, "Error unmarshalling yaml")

    return
}

// write the contents of a []Person as yaml into file f
func writeYaml(fam []Person, f string)  {
    y, err := yaml.Marshal(&fam)
    checkerr.CheckFatal(err, "Error marshalling yaml")

    err = ioutil.WriteFile(f, y, 0644)
    checkerr.CheckFatal(err, "Error writing to", f)
}

// upload file selector handler
func fileSelector(w http.ResponseWriter, r *http.Request)  {
    pv := PageVariables {
        PageTitle: "Select yaml file to load",
    }

    t, err := template.ParseFiles("web/template/file-select.html")
    checkerr.Check(err, "file-select template parsing error")

    err = t.Execute(w, pv)
    checkerr.Check(err, "file-select template executing error")
}

// load-family handler
func loadFamily(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()

    family = loadYaml(r.Form.Get("filebox"))
    http.Redirect(w, r, "/", http.StatusFound)
}

// index page handler
func index(w http.ResponseWriter, r *http.Request)  {
    pv := PageVariables {
        PageTitle: "Family Tree",
    }

    t, err := template.ParseFiles("web/template/index.html")
    checkerr.Check(err, "index template parsing error")

    err = t.Execute(w, pv)
    checkerr.Check(err, "index template executing error")
}

func main() {
    http.HandleFunc("/upload", fileSelector)
    http.HandleFunc("/load-family", loadFamily)
    http.HandleFunc("/", index)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
