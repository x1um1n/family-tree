// Generate and manage your family tree
package main

import (
    "github.com/x1um1n/checkerr"
    "log"
    "bytes"
    "strings"
    "os"
    "io"
    "io/ioutil"
    "net/http"
    "gopkg.in/yaml.v2"
    "html/template"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
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
var familyDB *sql.DB

// initialise db
func initDB()  {
    var err error

    log.Println("gonna init the db now")
    familyDB, err = sql.Open("mysql", os.Getenv("CONNSTR"))
    checkerr.CheckFatal(err, "error connecting to DB")
}

// create db structure
func createDB()  {
    log.Println("reading create.sql")
    f, err := ioutil.ReadFile("scripts/sql/create.sql")
    checkerr.Check(err, "error reading create.sql")

    log.Println("executing create.sql")
    qry := strings.Split(string(f), ";")
    for _, q := range qry {
        if len(q) != 0 {
            log.Println(q)
            _, err = familyDB.Exec(q)
            checkerr.Check(err, "error executing create.sql", q)
        }
    }
}

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

    // insert into db
    var qry []string
    for _, p := range family {
        qry = append(qry, "INSERT INTO `people` (`UUID`,`Name`,`Fnam`,`Mnam`,`Snam`,`Sex`,`DOB`,`DOD`) VALUES ("+p.UUID+","+p.Name+","+p.Fnam+","+p.Mnam+","+p.Snam+","+p.Sex+","+p.DOB+","+p.DOD+")")
        if len(p.Children) != 0 {
            for _, c := range p.Children {
                qry = append(qry, "INSERT INTO `children` (`parentId`,`childId`) VALUES ("+p.UUID+","+c+")")
            }
        }
        // fixme: parents? check to see if row exists, add if not
        // fixme: spouse? use the seeker
    }

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
    initDB()
    createDB()

    log.Println("starting wwwserver")
    http.HandleFunc("/edit", editPerson)
    http.HandleFunc("/upload", fileSelector)
    http.HandleFunc("/load-family", loadFamily)
    http.HandleFunc("/", index)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
