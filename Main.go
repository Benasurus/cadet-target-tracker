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
/*
Section Names:
	Section 01 - Drill
	Section 02 - Radio
	Section 03 - Flying
	Section 04 - Gliding
	Section 05 - Fieldcraft
	Section 06 - Classifications
	Section 07 - Sports
	Section 08 - Adventurous Training [AT]
	Section 09 - First Aid
	Section 10 - Leadership
	Section 11 - Duke of Edinburgh [DofE]
	Section 12 - Community Engagement
	Section 13 - Shooting
	Section 14 - Music
	Section 15 - Camps
*/

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
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//Defining global variables
var db *sql.DB                                                          //database variable
var err error                                                           //error variable
var store = sessions.NewCookieStore(securecookie.GenerateRandomKey(32)) //cookie store variable

//Round function as Go lacks a built in round function
func Round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
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
	var result [15]int
	//Queries databse using SQL calling all rows which match the desired username & scans through returned row
	err := db.QueryRow("SELECT * FROM transactions WHERE UserName=? ORDER BY timeStamp DESC LIMIT 1", user).Scan(&tID, &timeStamp, &userName, &result[0], &result[1],
		&result[2], &result[3], &result[4], &result[5], &result[6], &result[7], &result[8], &result[9], &result[10], &result[11], &result[12], &result[13], &result[14])
	//Error handler
	if err != nil {
		log.Println(err)
		return true, sectionData
	}

	for i := 1; i < 16; i++ {
		for j := 1; j < 16; j++ {
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
		html += "<div id=\"bar-title-" + strconv.FormatFloat(countSection, 'f', 0, 64) + "\" class=\"bar-title\" onclick=\"cadetModify('" + id + "')\"><div class=\"bar-title-text\"><h2>" + sectorTitle[id] + "</h2>" + sectorDesc[id] + "</div></div>"
		for countSector = 1; countSector < sectors; countSector++ {
			id = "sec" + strconv.FormatFloat(countSection, 'f', 0, 64) + "-set" + strconv.FormatFloat(countSector, 'f', 0, 64)
			if sectionData[id] == true {
				class = "sector-true"
			} else {
				class = "sector-false"
			}
			html += "<div id=\"sector-" + id + "\" class=" + class + " onclick=\"cadetModify('" + id + "')\"><div id=\"sector-title-" + id + "\" class=\"sector-title\">" + sectorTitle[id] + "<!-- <span class=\"tool-tip-text\">" + sectorDesc[id] + "</span> --></div><div id=\"sector-bar-" + id + "\" class=\"sector-bar\"></div></div>"
		}
		html += "</div>"
	}
	return html
}

func dataCalculation() (string, string, string, string) {
	var overall, flight, sex, sector string
	var month, year, monthChange, yearChange int
	t := time.Now()
	year = t.Year()
	month = int(t.Month())
	yearChange = year
	if month == 12 {
		monthChange = 1
	} else {
		monthChange = month + 1
	}
	sector = sectorCalculation(month, year)
	for i := 1; i < 13; i++ {
		var count, countA, countB, countM, countF float64
		var cpiTotal, cpiATotal, cpiBTotal, cpiMTotal, cpiFTotal float64
		var cpi, cpiAv, cpiAAv, cpiBAv, cpiMAv, cpiFAv float64
		count = 0
		cpiTotal = 0
		cpi = 0
		rows, err := db.Query("SELECT t1.cpi,t3.flight,t3.sex FROM transactions t1 INNER JOIN userdata AS t3 ON t1.userName=t3.userName WHERE t1.timeStamp=(SELECT MAX(t2.timeStamp) FROM transactions t2 WHERE timeStamp < '" + strconv.Itoa(yearChange) + "-" + strconv.Itoa(monthChange) + "-01' AND t2.userName = t1.userName)")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			var flight, sex string
			err := rows.Scan(&cpi, &flight, &sex)
			if err != nil {
				log.Fatal(err)
			}
			cpiTotal += cpi
			if flight == "A" {
				cpiATotal += cpi
				countA++
			} else if flight == "B" {
				cpiBTotal += cpi
				countB++
			}
			if sex == "M" {
				cpiMTotal += cpi
				countM++
			} else if sex == "F" {
				cpiFTotal += cpi
				countF++
			}
			count++
		}
		cpiAv, cpiAAv, cpiBAv, cpiMAv, cpiFAv = 0, 0, 0, 0, 0
		if count > 0 {
			cpiAv = cpiTotal / count
		}
		if countA > 0 {
			cpiAAv = cpiATotal / countA
		}
		if countB > 0 {
			cpiBAv = cpiBTotal / countB
		}
		if countM > 0 {
			cpiMAv = cpiMTotal / countM
		}
		if countA > 0 {
			cpiFAv = cpiFTotal / countF
		}
		cpiAv, cpiAAv, cpiBAv, cpiMAv, cpiFAv = 0, 0, 0, 0, 0
		fmt.Println(cpiAv, cpiAAv, cpiBAv, cpiMAv, cpiFAv)
		overall = "[new Date(" + strconv.Itoa(year) + "," + strconv.Itoa(month-1) + ", 1)," + strconv.FormatFloat(cpiAv, 'f', 4, 64) + "],\n" + overall
		flight = "[new Date(" + strconv.Itoa(year) + "," + strconv.Itoa(month-1) + ", 1)," + strconv.FormatFloat(cpiAAv, 'f', 4, 64) + "," + strconv.FormatFloat(cpiBAv, 'f', 4, 64) + "],\n" + flight
		sex = "[new Date(" + strconv.Itoa(year) + "," + strconv.Itoa(month-1) + ", 1)," + strconv.FormatFloat(cpiMAv, 'f', 4, 64) + "," + strconv.FormatFloat(cpiFAv, 'f', 4, 64) + "],\n" + sex
		month--
		if month == 0 {
			year--
			month = 12
		}
		monthChange--
		if monthChange == 0 {
			yearChange--
			monthChange = 12
		}
	}
	overall = "data.addRows([\n" + overall + "]);"
	flight = "data.addRows([\n" + flight + "]);"
	sex = "data.addRows([\n" + sex + "]);"
	//fmt.Println(overall)
	//fmt.Println(flight)
	//fmt.Println(sex)
	//fmt.Println(sector)
	return overall, flight, sex, sector
}

func sectorCalculation(month, year int) string {
	var count, cpi, ID, s1, s2, s3, s4, s5, s6, s7, s8, s9, s10, s11, s12, s13, s14, s15 float64
	var ts1, ts2, ts3, ts4, ts5, ts6, ts7, ts8, ts9, ts10, ts11, ts12, ts13, ts14, ts15 float64
	var as1, as2, as3, as4, as5, as6, as7, as8, as9, as10, as11, as12, as13, as14, as15 float64
	var user, time, sector string
	month++
	if month == 12 {
		year++
		month = 1
	}
	rows, err := db.Query("SELECT t1.* FROM transactions t1 WHERE t1.timeStamp=(SELECT MAX(t2.timeStamp) FROM transactions t2 WHERE timeStamp < '" + strconv.Itoa(year) + "-" + strconv.Itoa(month) + "-01' AND t2.userName = t1.userName)")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&ID, &user, &time, &s1, &s2, &s3, &s4, &s5, &s6, &s7, &s8, &s9, &s10, &s11, &s12, &s13, &s14, &s15, &cpi)
		if err != nil {
			log.Fatal(err)
		}
		ts1, ts2, ts3, ts4, ts5, ts6, ts7, ts8, ts9, ts10, ts11, ts12, ts13, ts14, ts15 = s1, s2, s3, s4, s5, s6, s7, s8, s9, s10, s11, s12, s13, s14, s15
		count++
	}
	as1 = ts1 / count
	as2 = ts2 / count
	as3 = ts3 / count
	as4 = ts4 / count
	as5 = ts5 / count
	as6 = ts6 / count
	as7 = ts7 / count
	as8 = ts8 / count
	as9 = ts9 / count
	as10 = ts10 / count
	as11 = ts11 / count
	as12 = ts12 / count
	as13 = ts13 / count
	as14 = ts14 / count
	as15 = ts15 / count
	sector = "data.addRows([['Drill'," + strconv.FormatFloat(as1, 'f', 4, 64) + "],['Radio'," + strconv.FormatFloat(as2, 'f', 4, 64) + "],['Flying'," + strconv.FormatFloat(as3, 'f', 4, 64) + "],['Gliding'," + strconv.FormatFloat(as4, 'f', 4, 64) + "],['Fieldcraft'," + strconv.FormatFloat(as5, 'f', 4, 64) + "],['Classifications'," + strconv.FormatFloat(as6, 'f', 4, 64) + "],['Sports'," + strconv.FormatFloat(as7, 'f', 4, 64) + "],['Adventurous Training'," + strconv.FormatFloat(as8, 'f', 4, 64) + "],['First Aid'," + strconv.FormatFloat(as9, 'f', 4, 64) + "],['Leadership'," + strconv.FormatFloat(as10, 'f', 4, 64) + "],['DofE'," + strconv.FormatFloat(as11, 'f', 4, 64) + "],['Community Engagement'," + strconv.FormatFloat(as12, 'f', 4, 64) + "],['Shooting'," + strconv.FormatFloat(as13, 'f', 4, 64) + "],['Music'," + strconv.FormatFloat(as14, 'f', 4, 64) + "],['Camps'," + strconv.FormatFloat(as15, 'f', 4, 64) + "]]);"
	return sector
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

func staffPage(w http.ResponseWriter, r *http.Request) {
	overall, flight, sex, _ := dataCalculation()
	pagevars := map[string]interface{}{
		"overall": template.JS(overall),
		"flight":  template.JS(flight),
		"sex":     template.JS(sex)}
	t, err := template.ParseFiles("resources/staff.html")
	if err != nil {
		log.Println(err)
	} else {
		err = t.Execute(w, pagevars)
		if err != nil {
			log.Println(err)
		}
	}
}

func dbWrite(w http.ResponseWriter, r *http.Request) {
	var value, column string
	var total, CPI float64
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
	for count := 0; count < 15; count++ {
		value += "," + (r.FormValue("section" + strconv.Itoa(count+1)))
		column += ",section" + strconv.Itoa(count+1)
		valueFloat, err := strconv.ParseFloat(r.FormValue("section"+strconv.Itoa(count+1)), 64)
		if err != nil {
			log.Println(err)
			return
		}
		total += valueFloat
	}
	//Calculate Cadet Progression Index
	CPI = total / 150
	fmt.Println("INSERT INTO transactions(userName" + column + ",cpi) VALUES(\"" + username + "\"" + value + "," + strconv.FormatFloat(CPI, 'f', 2, 64) + ")")
	//if auth == true {
	_, err := db.Exec("INSERT INTO transactions(userName" + column + ",cpi) VALUES(\"" + username + "\"" + value + "," + strconv.FormatFloat(CPI, 'f', 2, 64) + ")")
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
	http.HandleFunc("/staff", staffPage)
	http.HandleFunc("/dbwrite", dbWrite)
	http.Handle("/resources/", http.StripPrefix("/resources", http.FileServer(http.Dir("resources"))))
	err = http.ListenAndServe(":9090", nil) // set listen port
	//Error handler
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
