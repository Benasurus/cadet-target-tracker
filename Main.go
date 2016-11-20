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
//section = sec
//5 sectors per section (0-4)
//15 sections (1-15)

//Package definition
package main

// Importing GO modules used by the program
import (
	"database/sql"
	"fmt"
	"html/template"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	//"io"
	"log"
	//"math"
	"net/http"
	"strconv"
	"strings"
)

//Defining global variables
var db *sql.DB                                                          //database variable
var err error                                                           //error variable
var store = sessions.NewCookieStore(securecookie.GenerateRandomKey(32)) //cookie store variable

//defining a structure to store sectorLabel and sectorDescription in a map variable
type sectorInfo struct {
	Label       string
	Description string
}

func auth(w http.ResponseWriter, r *http.Request) (bool, string, string) {
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
			return true, username, group
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

func authNTLM(w http.ResponseWriter, r *http.Request) (bool, string, string) {
	//AuthNTLM module authenticates via NTLM, called only by auth module
	auth := r.Header.Get("Authorization")
	if !(strings.HasPrefix(auth, "NTLM ")) {
		fmt.Println("Trying NTML")
		w.Header().Set("WWW-Authenticate", "NTLM")
		http.Error(w, "Unauthorized.", 401)
		return false, "", ""
	}
	fmt.Println(auth)
	return true, "NTLM", "user"
}

func authBasic(w http.ResponseWriter, r *http.Request) (bool, string, string) {
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

	return false, "", ""
}

func userData(user string, sectionData map[string]bool) (bool, map[string]bool) {
	//QueryTransaction module is used to query transactions table for rows that match a user
	var tID string
	var timeStamp, userName string
	var result [4]int
	//Queries databse using SQL calling all rows which match the desired username & scans through returned row
	err := db.QueryRow("SELECT * FROM transactions WHERE UserName=? ORDER BY timeStamp DESC LIMIT 1", user).Scan(&tID, &timeStamp, &userName, &result[0], &result[1], &result[2], &result[3])
	//Error handler
	if err != nil {
		log.Println(err)
		return true, sectionData
	}

	for i := 1; i < 5; i++ {
		for j := 1; j < 5; j++ {
			id := "sec" + strconv.Itoa(i) + "-set" + strconv.Itoa(j)
			if j-1 < result[i-1] {
				sectionData[id] = true
			} else {
				sectionData[id] = false
			}
		}
	}
	//Prints results
	fmt.Println(userName, timeStamp)
	return false, sectionData
}

func writeTransaction(user string, group string, s1 int, s2 int, s3 int, s4 int) bool {
	//WriteTransaction module writes a row to the transactions table under a username
	_, err := db.Exec("INSERT INTO transactions(userName,section1,section2,section3,section4) VALUES(?,?,?,?,?)", user, s1, s2, s3, s4)
	//Error handler
	if err != nil {
		log.Println(err)
		return false
	}
	return true

}

func deleteUser(userDel string) bool {
	//Deletes all rows where the username matches
	_, err := db.Exec("DELETE FROM transactions WHERE userName=?", userDel)
	//Error Handler
	if err != nil {
		log.Println(err)
		return false
	}
	return true

}

func progressionBuilder(sectors, sections float64, user string) (html string) {
	//progressionBuilder module generates HTML to build each sector
	//Defining variables
	var countSection, countSector float64
	var class, id, sectorID, sectorLabel, sectorDescription string
	//Map variable which stores user data for each section
	var sectionData map[string]bool

	sectionData = make(map[string]bool)
	//Map variable which stores sector title
	var sectorTitle map[string]string

	sectorTitle = make(map[string]string)
	//Map variable which stores information on each sector
	var sectorDesc map[string]string

	sectorDesc = make(map[string]string)
	//Queries the sectors database
	rows, err := db.Query("SELECT * FROM sectors")
	fmt.Println(rows)
	//Error handler
	if err != nil {
		log.Println(err)
		return "Error"
	}
	defer rows.Close()
	//Calls user data
	fail, sectionData := userData(user, sectionData)
	if fail == true {
		log.Println("Data Error")
	}
	//Scans through returned values and retrieves sector names and inserts them into the defined variables
	for rows.Next() {
		err := rows.Scan(&sectorID, &sectorLabel, &sectorDescription)
		if err != nil {
			log.Fatal(err)
			return "Error"
		}
		sectorTitle[sectorID] = sectorLabel
		sectorDesc[sectorID] = sectorDescription
	}
	//For loop which goes through and builds all sectors
	for countSection = 1; countSection < sections+1; countSection++ {
		id = "sec" + strconv.FormatFloat(countSection, 'f', 0, 64) + "-set0"
		html += "<div id=\"section-" + strconv.FormatFloat(countSection, 'f', 0, 64) + "\" class=\"section\">"
		html += "<div id=\"bar-title-" + strconv.FormatFloat(countSection, 'f', 0, 64) + "\" class=\"bar-title\" onclick=\"modify('" + id + "')\"><h2>" + sectorTitle[id] + "</h2>" + sectorDesc[id] + "</div>"
		for countSector = 1; countSector < sectors; countSector++ {
			id = "sec" + strconv.FormatFloat(countSection, 'f', 0, 64) + "-set" + strconv.FormatFloat(countSector, 'f', 0, 64)
			if sectionData[id] == true {
				class = "sector-true"
			} else {
				class = "sector-false"
			}
			html += "<div id=\"sector-" + id + "\" class=" + class + " onclick=\"modify('" + id + "')\"><span class=\"tool-tip-text\">" + sectorDesc[id] + "</span><div id=\"sector-title-" + id + "\" class=\"sector-title\">" + sectorTitle[id] + "</div><div id=\"sector-bar-" + id + "\" class=\"sector-bar\"></div></div>"
		}
		html += "</div>"
	}
	return html
}

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	//Webpage handler for the root domain
	r.ParseForm()       // parse arguments, you have to call this by yourself
	fmt.Println(r.Form) // print form information in server side
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello Ben!") // send data to client side
}

func progressionTracker(w http.ResponseWriter, r *http.Request) {
	//authorised, username, group := auth(w, r)
	userName := "TestUser"
	authorised := true
	//fmt.Println(username,group)
	if !authorised {
		// fmt.Fprintf(w, "Failed to Authenticate")
		return
	}
	if r.Method == "GET" {
		html := progressionBuilder(5, 15, userName)
		pagevars := map[string]interface{}{
			"html": template.HTML(html),
			"user": template.HTML(userName)}
		t, err := template.ParseFiles("resources/cadet.html")
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

func dbWrite(w http.ResponseWriter, r *http.Request) {
	var value, column string

	var username string
	username = "TEST"
	/*var auth bool
	auth = false
	//Retrieves Session
	session, err := store.Get(r, "session-name")
	//Error handler
	if err != nil {
		log.Println(err)
	} else {
		//Authenticates the cookie
		username, ok := session.Values["username"].(string)
		if ok {
			group, ok := session.Values["group"].(string)
			fmt.Println(username + "[" + group + "] Succesfully Authenticated Cookie")
			if ok && group == "Cadet" {
				auth = true
			}
		}
	}*/
	r.ParseForm()
	for count := 0; count < 4; count++ {
		value += "," + (r.FormValue("section" + strconv.Itoa(count+1)))
		column += ",section" + strconv.Itoa(count+1)
	}
	fmt.Println("INSERT INTO transactions(userName" + column + ") VALUES(\"" + username + "\"" + value + ")")
	//if auth == true {
	_, err := db.Exec("INSERT INTO transactions(userName" + column + ") VALUES(\"" + username + "\"" + value + ")")
	//fmt.Println(username + "[" + group + "] Wrote to Database")
	if err != nil {
		log.Println(err)
	}
	//}
}

func main() {
	//Opens connection to the database under the db variable
	db, err = sql.Open("mysql", "root:Dragon121@tcp(localhost:3306)/CadetTracker")
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
	http.HandleFunc("/cadet", progressionTracker)
	http.HandleFunc("/dbwrite", dbWrite)
	http.Handle("/resources/", http.StripPrefix("/resources", http.FileServer(http.Dir("resources"))))
	err = http.ListenAndServe(":9090", nil) // set listen port
	//Error handler
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
