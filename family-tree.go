// Generate and manage your family tree
package main

import (
    "log"
    "bytes"
    "io"
    "net/http"
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "html/template"
    "github.com/x1um1n/checkerr"
)

type PageVariables struct {
    PageTitle     string
    Family        []Person
    Personage     Person
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

// find the index of a UUID in the family slice
func seeker(u string) (idx int) {
    m := make(map[string]int)

    for i, p := range family {
        m[p.UUID] = i
    }
    idx = m[u]

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

    f, h, err := r.FormFile("filebox")
    checkerr.Check(err, "error opening file", r.Form.Get("filebox"))
    defer f.Close()

    y, err := h.Open()
    checkerr.CheckFatal(err, "Error opening file", h.Filename)

    buf := bytes.NewBuffer(nil)
    if _, err := io.Copy(buf, y); err != nil {
        checkerr.CheckFatal(err)
    }

    err = yaml.Unmarshal(buf.Bytes(), &family)
    checkerr.CheckFatal(err, "Error unmarshalling yaml")

    http.Redirect(w, r, "/", http.StatusFound)
}

// edit family member handler
func editPerson(w http.ResponseWriter, r *http.Request)  {
    pv := PageVariables {
        PageTitle: "Edit person",
        Personage: family[seeker(r.Form.Get("edit"))],
    }

    t, err := template.ParseFiles("web/template/edit-person.html")
    checkerr.Check(err, "edit-person template parsing error")

    err = t.Execute(w, pv)
    checkerr.Check(err, "edit-person template executing error")
}

// index page handler
func index(w http.ResponseWriter, r *http.Request)  {
    pv := PageVariables {
        PageTitle: "Family Tree",
        Family: family,
    }

    t, err := template.ParseFiles("web/template/index.html")
    checkerr.Check(err, "index template parsing error")

    err = t.Execute(w, pv)
    checkerr.Check(err, "index template executing error")
}

func main() {
    http.HandleFunc("/edit", editPerson)
    http.HandleFunc("/upload", fileSelector)
    http.HandleFunc("/load-family", loadFamily)
    http.HandleFunc("/", index)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
