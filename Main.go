/*
Cadet Target Tracker Web Application
Copyright (C) 2016  Benjamin Piggott

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

//sector = set
//section = section
//5 sectors per section
//14 sections

//Package definition
package main

// Importing GO modules used by the program
import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "fmt"
    "github.com/gorilla/securecookie"
    "github.com/gorilla/sessions"
    "html/template"
    //"io"
    "log"
    //"math"
    "net/http"
    "strconv"
    "strings"
)

//Defining global variables
var db *sql.DB //database variable
var err error  //error variable
var store = sessions.NewCookieStore(securecookie.GenerateRandomKey(32)) //cookie store variable

//defining a structure to store sectorLabel and sectorDescription in a map variable
type sectorInfo struct {
        Label string
        Description string
    }


func auth(w http.ResponseWriter, r *http.Request)(bool, string, string) {
    //Auth module is called to authenticate a user
    //Variables
    var val bool
    var username, group string
    //Generates a session
    session, err := store.Get(r, "session-name")
    //Error handler
    if err != nil {
        log.Println(err)
        return false, "", ""
    }
    //Authenticates the cookie
    username, ok := session.Values["username"].(string)
    if ok {
        group, ok := session.Values["group"].(string)
        if ok {
            return true,username,group
        }
    }
    //Authenticates via NTLM and then HTTP Basic 
    val, username, group = authNTLM(w, r)
    if !val {
        val, username, group = authBasic(w, r)
    }
    //If NTLM or Basic authentication succeeds generates a cookie
    if val {
        session.Values["username"] = username
        session.Values["group"] = group
        session.Save(r, w)
        return true, username, group
    }
    //Returns false and empty values if authentication fails
    return false, "", ""
}



func authNTLM(w http.ResponseWriter, r *http.Request)(bool, string, string) {
    //AuthNTLM module authenticates via NTLM, called only by auth module
    auth := r.Header.Get("Authorization")
    if !(strings.HasPrefix(auth,"NTLM ")) {
        fmt.Println("Trying NTML")
        w.Header().Set("WWW-Authenticate", "NTLM")
        http.Error(w, "Unauthorized.", 401)
        return false, "", ""
        }
        fmt.Println(auth)
        return true, "NTLM", "user"
}

func authBasic(w http.ResponseWriter, r *http.Request)(bool, string, string) {
    //AuthBasic module authenticates via HTTP basic, called only by auth module
    fmt.Println("Trying Basic Auth")
    username, password, ok := r.BasicAuth()
    if ok {
        ldapSuccess := true
        group := "admin"
        fmt.Println(password)
        // will need an LDAP lookup using username and password to
        // a) validate the details are correct and
        // b) get the user group information
        if ldapSuccess {
            return true, username, group
        }
        return false, "", ""
        
    }
    w.Header().Set("WWW-Authenticate", "Basic realm=\"localhost\"")

    return false,"",""
}


func userData(user string, sectorData map[string]bool) (bool, map[string]bool) {
    //QueryTransaction module is used to query transactions table for rows that match a user
    var tID string
    var timeStamp, userName string
    var result [4]int
    //Queries databse using SQL calling all rows which match the desired username & scans through returned row 
    err := db.QueryRow("SELECT * FROM transactions WHERE UserName=? ORDER BY timeStamp DESC LIMIT 1", user).Scan(&tID, &timeStamp,&userName, &result[0], &result[1], &result[2], &result[3])
    //Error handler
	if err != nil {
		log.Println(err)
		return false, sectorData
    }

    for i := 1; i < 5; i++ {
        for j := 1; j < 5; j++ {
            id := "sec"+strconv.Itoa(i)+"-seg"+strconv.Itoa(j)
            if j-1 < result[i-1] {
                sectorData[id] = true
            } else {
                sectorData[id] = false
            }
        }
    }
    //Prints results
	    fmt.Println(userName, timeStamp)
    return true, sectorData
}

func writeTransaction(user string, group string ,s1 int, s2 int, s3 int ,s4 int)(bool) {
    //WriteTransaction module writes a row to the transactions table under a username
    _, err := db.Exec("INSERT INTO transactions(userName,sector1,sector2,sector3,sector4) VALUES(?,?,?,?,?)",user,s1,s2,s3,s4)
    //Error handler
    if err != nil {
        log.Println(err)   
        return false
    }
    return true

}

func deleteUser(userDel string,)(bool) {
    //Deletes all rows where the username matches
    _, err := db.Exec("DELETE FROM transactions WHERE userName=?",userDel)
    //Error Handler
    if err != nil {
		log.Println(err)
		return false
    }
    return true

}

func roundelBuilder(sectors, sections float64, user string)(path, div string) {
    //RoundelBuilder module generates HTML to build each sector of the roundel
    //Defining variables
    var countSection, countSector float64
    var id, sectorID, sectorLabel, sectorDescription, section string
    //Map variable which stores user data for each sector
    var sectorData map[string]bool

    sectorData = make(map[string]bool)
    //Map variable which stores information on each segement
    var sectorMap map[string]string
    
    sectorMap = make(map[string]string)
    //Queries the sectors database
    rows, err := db.Query("SELECT * FROM sectors")
    fmt.Println(rows)
    //Error handler
    if err != nil {
		log.Println(err)
		return "Error", "Error"
    }
    defer rows.Close()
    //Calls user data
    ok, sectorData := userData(user, sectorData)
    if ok == false {
        return "Error", "Error"
    }
    //Scans through returned values and retrieves sector names and inserts them into the defined variables
    for rows.Next() {
        err := rows.Scan(&sectorID, &sectorLabel, &sectorDescription)
        if err != nil {
		     log.Fatal(err)
             return "Error", "Error"
        }
        sectorMap[sectorID] = sectorLabel
    }
    //For loop which goes through and builds all sectors
    for countSector = 0; countSector < sectors; countSector ++ {
        for countSection = 0; countSection < sections; countSection ++ {
            id = "sector" + strconv.FormatFloat(countSector+1, 'f', 0, 64) + "-section" + strconv.FormatFloat(countSection, 'f', 0, 64)
            /*if sectorData[id]==true {
                fill = "yellow"
            }  else {
                fill = "transparent"
                }*/
            //attrs = "id=\"" +id+ "\" stroke=\"black\" fill=\"" +fill+ "\" onclick=\"doSetHighlight('" +id+ "');\""
            section = "<div id="+id+">"+section+"</div>"
            div =  div + "<div id=\"div-"+id+"\">"+sectorMap[id]+"</div>"
        }
    }
    return section, div
}

func sayhelloName(w http.ResponseWriter, r *http.Request) {
    //Webpage handler for the root domain
    r.ParseForm()  // parse arguments, you have to call this by yourself
    fmt.Println(r.Form)  // print form information in server side
    fmt.Println("path", r.URL.Path)
    fmt.Println("scheme", r.URL.Scheme)
    fmt.Println(r.Form["url_long"])
    for k, v := range r.Form {
        fmt.Println("key:", k)
        fmt.Println("val:", strings.Join(v, ""))
    }
    fmt.Fprintf(w, "Hello Ben!") // send data to client side
}

func input(w http.ResponseWriter, r *http.Request) {
    //authorised, username, group := auth(w, r)
    authorised := true
    //fmt.Println(username,group)
    if !authorised {
        // fmt.Fprintf(w, "Failed to Authenticate")
        return
    }
    if r.Method == "GET" {
        sector, div := roundelBuilder(15,5, "tesdt")
        pagevars := map[string]interface{}{
            "sectors"  : template.HTML(sector),
            "divs"  : template.HTML(div)}
        t, err := template.ParseFiles("test.html")
        if err != nil {
            log.Println(err)
        } else {
            err = t.Execute(w, pagevars)
               if err != nil {
                log.Println(err)
                }
        }  
    } else {
        r.ParseForm()
        // logic part of log in
        //sector := r.Form["username"][0]
        //sector1, err := strconv.Atoi(r.Form["sector1"][0])
        //if err != nil{
        //   log.Println(err)
        //}
    }
}

func main() {
    //Opens connection to the database under the db variable
    db, err = sql.Open("mysql","root:Dragon121@tcp(localhost:3306)/CadetTracker")
    //Error handler
    if err != nil {
	    log.Println(err)
    }
    //Pings database to ensure connection is open
    err = db.Ping()
    //Error handler
    if err != nil {
        log.Println(err)
    }
    //Web server
    //Defines handler functions for each webpage
    http.HandleFunc("/", sayhelloName) // set router
    http.HandleFunc("/input", input)
    err = http.ListenAndServe(":9090", nil) // set listen port
    //Error handler
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
    
    
}