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
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/bearbin/go-age"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	"github.com/jimlawless/cfg"

	_ "gopkg.in/asn1-ber.v1"
	"gopkg.in/ldap.v2"
)

//Defining global variables
var db *sql.DB                                                          //database variable
var err error                                                           //error variable
var store = sessions.NewCookieStore([]byte("9xCKuZUsytH1zoVWJHHXgeJ98NjuXppy076xsc3X4IRqV0bcgxXH0b4MI8MSklVP")) //cookie store variable
var userTable map[string]string
var config map[string]string

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

func configLoad()(bool) {
	config = make(map[string]string)
    err := cfg.Load("config.cfg", config)
    if err != nil {
        log.Println(err)
		return false
    }
	return true
}

func authCookie(w http.ResponseWriter, r *http.Request) (bool, string, string) {
	//Variables
	var username, group string
	username, group = "error", "error"
	//Generates a session
	session, err := store.Get(r, "CadetTracker")
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
	//Returns false and empty values if authentication fails
	return false, username, group
}

func authLDAP(w http.ResponseWriter, r *http.Request, username string, password string)(bool, string) {
	//Generates a session
	var group string
	if checkUser(username) == false {
		return false, ""
	}
	port,_ := strconv.ParseUint(config["port"], 10, 64)

	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", config["url"], port))
	if err != nil {
		log.Println(err)
		return false, ""
	}
	defer l.Close()

	// Reconnect with TLS
	//err = l.StartTLS(&tls.Config{InsecureSkipVerify: true})
	//if err != nil {
	//	log.Println(err)
	//}

	// Bind as the user to verify their password
	err = l.Bind(userTable[username], password)
	if err != nil {
		log.Println(err)
		return false, ""
	}
	if strings.Contains(userTable[username], "ou=staff") {
        group = "staff"
    } else if strings.Contains(userTable[username], "ou=cadets") {
        group = "cadet"
    } else {
		group = "staff"
	}
	return true, group
}

func userTableFill() () {
	userTable = make(map[string]string)
	port,_ := strconv.ParseUint(config["port"], 10, 64)

	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", config["url"], port))
	if err != nil {
		log.Println(err)
	}
	defer l.Close()

	// Reconnect with TLS
	//err = l.StartTLS(&tls.Config{InsecureSkipVerify: true})
	//if err != nil {
	//	log.Println(err)
	//}

	// First bind with a read only user
	err = l.Bind(config["distinguishedname"], config["password"])
	if err != nil {
		log.Println(err)
	}

	// Search for the given username
	
	searchRequest := ldap.NewSearchRequest(
		config["basedomainname"],
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(objectClass=posixAccount)",
		[]string{"sAMAccountName"},
		nil,)

	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Println(err)
	}

	for _, entry := range sr.Entries {
		userTable[strings.ToLower(entry.GetAttributeValue("sAMAccountName"))] = entry.DN
	}

}

func checkUser(username string)(bool) {
	var LDAPname string
	LDAPname = userTable[strings.ToLower(username)]
	if LDAPname == "" {
		return false
	}
	return true
}

func userPresent(username string) bool {
	//Checks if user is present in the user database
	var result string
	err := db.QueryRow("SELECT userName FROM userdata WHERE UserName=?", username).Scan(&result)
	if err != nil {
		log.Println(err)
		return false
	}
	if result == username {
		return true
	}
	return false
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
	//Deletes all rows where the username matches in transactions
	_, err = db.Exec("DELETE FROM transactions WHERE userName=?", userDel)
	//Error Handler
	if err != nil {
		log.Println(err)
		return false
	}
	//Deletes all rows where the username matches in userdata
	_, err = db.Exec("DELETE FROM userdata WHERE userName=?", userDel)
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
			log.Println(err)
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
	//Declaring variables
	var overall, flight, sex, sector string
	var month, year, monthChange, yearChange int
	//Retrieve current time
	t := time.Now()
	//Calculate current year and month
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
			log.Println(err)
		}
		defer rows.Close()
		for rows.Next() {
			var flight, sex string
			err := rows.Scan(&cpi, &flight, &sex)
			if err != nil {
				log.Println(err)
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
		if countF > 0 {
			cpiFAv = cpiFTotal / countF
		}
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
		log.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&ID, &user, &time, &s1, &s2, &s3, &s4, &s5, &s6, &s7, &s8, &s9, &s10, &s11, &s12, &s13, &s14, &s15, &cpi)
		if err != nil {
			log.Println(err)
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
	sector = "data.addRows([\n['Drill'," + strconv.FormatFloat(as1, 'f', 4, 64) + "],\n['Radio'," + strconv.FormatFloat(as2, 'f', 4, 64) + "],\n['Flying'," + strconv.FormatFloat(as3, 'f', 4, 64) + "],\n['Gliding'," + strconv.FormatFloat(as4, 'f', 4, 64) + "],\n['Fieldcraft'," + strconv.FormatFloat(as5, 'f', 4, 64) + "],\n['Classifications'," + strconv.FormatFloat(as6, 'f', 4, 64) + "],\n['Sports'," + strconv.FormatFloat(as7, 'f', 4, 64) + "],\n['Adventurous Training'," + strconv.FormatFloat(as8, 'f', 4, 64) + "],\n['First Aid'," + strconv.FormatFloat(as9, 'f', 4, 64) + "],\n['Leadership'," + strconv.FormatFloat(as10, 'f', 4, 64) + "],\n['DofE'," + strconv.FormatFloat(as11, 'f', 4, 64) + "],\n['Community Engagement'," + strconv.FormatFloat(as12, 'f', 4, 64) + "],\n['Shooting'," + strconv.FormatFloat(as13, 'f', 4, 64) + "],\n['Music'," + strconv.FormatFloat(as14, 'f', 4, 64) + "],\n['Camps'," + strconv.FormatFloat(as15, 'f', 4, 64) + "]\n]);"
	return sector
}

func sectorChart(user string) (string, int) {
	var s1, s2, s3, s4, s5, s6, s7, s8, s9, s10, s11, s12, s13, s14, s15 float64
	var cpi int
	var sector string
	err := db.QueryRow("SELECT section1, section2, section3, section4, section5, section6, section7, section8, section9, section10, section11, section12, section13, section14, section15, cpi FROM transactions WHERE userName=? ORDER BY timeStamp DESC LIMIT 1", user).Scan(&s1, &s2, &s3, &s4, &s5, &s6, &s7, &s8, &s9, &s10, &s11, &s12, &s13, &s14, &s15, &cpi)
	if err != nil {
		log.Println(err)
	}
	sector = "data.addRows([\n['Drill'," + strconv.FormatFloat(s1, 'f', 4, 64) + "],\n['Radio'," + strconv.FormatFloat(s2, 'f', 4, 64) + "],\n['Flying'," + strconv.FormatFloat(s3, 'f', 4, 64) + "],\n['Gliding'," + strconv.FormatFloat(s4, 'f', 4, 64) + "],\n['Fieldcraft'," + strconv.FormatFloat(s5, 'f', 4, 64) + "],\n['Classifications'," + strconv.FormatFloat(s6, 'f', 4, 64) + "],\n['Sports'," + strconv.FormatFloat(s7, 'f', 4, 64) + "],\n['Adventurous Training'," + strconv.FormatFloat(s8, 'f', 4, 64) + "],\n['First Aid'," + strconv.FormatFloat(s9, 'f', 4, 64) + "],\n['Leadership'," + strconv.FormatFloat(s10, 'f', 4, 64) + "],\n['DofE'," + strconv.FormatFloat(s11, 'f', 4, 64) + "],\n['Community Engagement'," + strconv.FormatFloat(s12, 'f', 4, 64) + "],\n['Shooting'," + strconv.FormatFloat(s13, 'f', 4, 64) + "],\n['Music'," + strconv.FormatFloat(s14, 'f', 4, 64) + "],\n['Camps'," + strconv.FormatFloat(s15, 'f', 4, 64) + "]\n]);"
	return sector, cpi
}

func tableBuilder() string {
	var surname, forename, flight, table, userName string
	var DOB time.Time
	rows, err := db.Query("SELECT lastName, firstName, flight, dateOfBirth, userName FROM userdata")
	if err != nil {
		log.Println(err)
	}
	rowcount := 0
	for rows.Next() {
		var newRow string
		var userAge int
		err := rows.Scan(&surname, &forename, &flight, &DOB, &userName)
		if err != nil {
			log.Println(err)
			return "Error"
		}
		userAge = age.Age(DOB)
		newRow = "<tr id=\"user-row-" + strconv.Itoa(rowcount) + "\" class=\"button\" onclick=\"loadData('" + userName + "')\">\n		<td>" + surname + "</td>\n		<td>" + forename + "</td>\n		<td>" + flight + "</td>\n		<td>" + strconv.Itoa(userAge) + "</td>\n	<td>" + userName + "</td>\n	</tr>\n	"
		table += newRow
		rowcount++
	}
	return table
}

func progressionTracker(w http.ResponseWriter, r *http.Request) {
	var userName, group string
	var auth bool
	auth,userName,group = authCookie(w,r)
	if auth == true && group == "staff" {
		html := progressionBuilder(5, 15, userName)
		pagevars := map[string]interface{}{
			"html": template.HTML(html),
			"user": template.HTML(userName)}
		t, err := template.ParseFiles("resources/cadet.html")
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Error", 500 )
		} else {
			err = t.Execute(w, pagevars)
			if err != nil {
				log.Println(err)
				http.Error(w, "Internal Error", 500 )
			}
		}
	} else {
		http.Redirect(w,r,"/login", 401)
	}
}

func staffPage(w http.ResponseWriter, r *http.Request) {
	var group string
	var auth bool
	overall, flight, sex, sectors := dataCalculation()
	auth,_,group = authCookie(w,r)
	if auth == true && group == "staff" {
		pagevars := map[string]interface{}{
			"overall": template.JS(overall),
			"flight":  template.JS(flight),
			"sex":     template.JS(sex),
			"sectors": template.JS(sectors)}
		t, err := template.ParseFiles("resources/staff.html")
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Error", 500 )
		} else {
			err = t.Execute(w, pagevars)
			if err != nil {
				log.Println(err)
				http.Error(w, "Internal Error", 500 )
			}
		}
	} else {
		http.Redirect(w,r,"/login", 401)
	}
}

func userManagement(w http.ResponseWriter, r *http.Request) {
	var group string
	var auth bool
	table := tableBuilder()
	auth,_,group = authCookie(w,r)
	if auth == true && group == "staff" {
		pagevars := map[string]interface{}{
			"table": template.HTML(table)}
		t, err := template.ParseFiles("resources/user.html")
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Error", 500 )
		} else {
			err = t.Execute(w, pagevars)
			if err != nil {
				log.Println(err)
				http.Error(w, "Internal Error", 500 )
			}
		}
	} else {
		http.Redirect(w,r,"/login", 401)
	}
}

func dbWrite(w http.ResponseWriter, r *http.Request) {
	var value, column string
	var total, CPI float64
	var username, group string
	var auth bool
	auth, username, group = authCookie(w,r)
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
	if auth == true && group == "cadet" {
		_, err := db.Exec("INSERT INTO transactions(userName" + column + ",cpi) VALUES(\"" + username + "\"" + value + "," + strconv.FormatFloat(CPI, 'f', 2, 64) + ")")
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Error", 500 )
		}
	} else {
		http.Error(w, "Authentication has Failed", 401 )
	}
}

func userLoad(w http.ResponseWriter, r *http.Request) {
	var query url.Values
	var userName, forename, surname, sex, flight, sector string
	var userAge, cadetAge, cpi int
	var DOB, DOE time.Time
	var auth bool
	auth,_, _ = authCookie(w,r)
	query = r.URL.Query()
	if auth == true {
		userName = query.Get("username")
		err := db.QueryRow("SELECT * FROM userdata WHERE userName=?", userName).Scan(&userName, &forename, &surname, &DOB, &DOE, &sex, &flight)
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Error", 500 )
		}
		if sex == "F" {
			sex = "<input type=\"radio\" name=\"gender\" id=\"gender-male\" class=\"button\" value=\"male\"> Male<br><input type=\"radio\" name=\"gender\" id=\"gender-female\" class=\"button\" value=\"female\" checked> Female<br>"
		} else {
			sex = "<input type=\"radio\" name=\"gender\" id=\"gender-male\" class=\"button\" value=\"male\" checked> Male<br><input type=\"radio\" name=\"gender\" id=\"gender-female\" class=\"button\" value=\"female\"> Female<br>"
		}
		if flight == "A" {
			flight = "<input type=\"radio\" name=\"flight\" id=\"flight-a\" class=\"button\" value=\"A\" checked> A<br><input type=\"radio\" name=\"flight\"  id=\"flight-b\" class=\"button\" value=\"B\"> B<br>"
		} else {
			flight = "<input type=\"radio\" name=\"flight\" id=\"flight-a\" class=\"button\" value=\"A\"> A<br><input type=\"radio\" name=\"flight\" id=\"flight-b\" class=\"button\" value=\"B\" checked> B<br>"
		}
		userAge = age.Age(DOB)
		cadetAge = age.Age(DOE)
		sector, cpi = sectorChart(userName)
		pagevars := map[string]interface{}{
			"forename": template.HTML(forename),
			"surname":  template.HTML(surname),
			"flight":   template.HTML(flight),
			"age":      template.HTML(strconv.Itoa(userAge)),
			"dob":      template.HTML(strings.TrimSuffix(time.Time.String(DOB), " 00:00:00 +0000 UTC")),
			"cadet":    template.HTML(strconv.Itoa(cadetAge)),
			"doe":      template.HTML(strings.TrimSuffix(time.Time.String(DOE), " 00:00:00 +0000 UTC")),
			"sex":      template.HTML(sex),
			"username": template.HTML(userName),
			"cpi":      template.HTML(strconv.Itoa(cpi)),
			"sector":   template.JS(sector)}
		t, err := template.ParseFiles("resources/userLoad.html")
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Error", 500 )
		} else {
			err = t.Execute(w, pagevars)
			if err != nil {
				log.Println(err)
				http.Error(w, "Internal Error", 500 )
			}
		}
	} else {
		http.Error(w, "Authentication has Failed", 401 )
	}
}

func userModify(w http.ResponseWriter, r *http.Request) {
	var userName, firstName, lastName, dob, doe, sex, flight string
	var auth bool
	val,_, group := authCookie(w,r)
	if val == true && group == "staff" {
		auth = true
	}
	r.ParseForm()
	userName = r.FormValue("username")
	firstName = r.FormValue("forename")
	lastName = r.FormValue("surname")
	dob = r.FormValue("dob")
	doe = r.FormValue("doe")
	if r.FormValue("gender") == "true" {
		sex = "M"
	} else {
		sex = "F"
	}
	if r.FormValue("flight") == "true" {
		flight = "A"
	} else {
		flight = "B"
	}
	if auth == true {
		_, err := db.Exec("UPDATE userdata SET firstName='" + firstName + "', lastName='" + lastName + "', dateOfBirth='" + dob + "', dateOfEnrollment='" + doe + "', sex='" + sex + "', flight='" + flight + "' WHERE userName='" + userName + "'")
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Error", 500 )
		}
	} else {
		http.Error(w, "Authentication has Failed", 401 )
	}
}

func userAdd(w http.ResponseWriter, r *http.Request) {
	var userName, firstName, lastName, dob, doe, sex, flight, group string
	var auth bool
	auth,_,group = authCookie(w,r)
	if auth == true && group == "staff" {
		r.ParseForm()
		userName = r.FormValue("username")
		firstName = r.FormValue("forename")
		lastName = r.FormValue("surname")
		dob = r.FormValue("dob")
		doe = r.FormValue("doe")
		if r.FormValue("gender") == "true" {
			sex = "M"
		} else {
			sex = "F"
		}
		if r.FormValue("flight") == "true" {
			flight = "A"
		} else {
			flight = "B"
		}
		if auth == true && group == "staff" {
			_, err := db.Exec("INSERT INTO userdata VALUES ('" + userName + "','" + firstName + "','" + lastName + "','" + dob + "','" + doe + "','" + sex + "','" + flight + "')")
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	} else {
		http.Error(w, "Authentication has Failed", 401 )
	}
}

func userRemove(w http.ResponseWriter, r *http.Request) {
	var query url.Values
	var userName, group string
	var auth bool
	auth,_,group = authCookie(w,r)
	if auth == true && group == "staff" {
		query = r.URL.Query()
		userName = query.Get("username")
		deleteUser(userName)
	} else {
		http.Error(w, "Authentication has Failed", 401 )
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		pagevars := map[string]interface{}{
			"squadron":   template.HTML(config["squadron"])}
		t, err := template.ParseFiles("resources/login.html")
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			err = t.Execute(w, pagevars)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	} else {
		var userName, password, group, redirect string
		var auth bool

		r.ParseForm()
		userName = r.FormValue("username")
		password = r.FormValue("password")
		
		if checkUser(userName) == true {
			auth,group = authLDAP(w, r, userName, password)
		} else {
			http.Error(w, "Incorrect Username", 401 )
		}

		if auth == true {
			session, err := store.Get(r, "CadetTracker")
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			session.Values["username"] = userName
			session.Values["group"] = group
			session.Options.MaxAge = 1800
			session.Save(r,w)
			redirect = "<head><title>Login Successful</title><meta http-equiv=\"refresh\" content=\"2;URL=/"+group+"\" /></head><body><p>Login Successful. Wait 2 seconds, or click <a href=\"/"+group+"\">here</a> if you are not automatically redirected.</p></body>"
			fmt.Fprintf(w, redirect)
		} else {
			http.Error(w, "Incorrect Password", 401 )
		}
	}

}

func logout (w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "CadetTracker")
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	session.Options.MaxAge = -1
	session.Save(r,w)
	http.Redirect(w,r,"/login", 401)
}

func root (w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<head><title>Login Successful</title><meta http-equiv=\"refresh\" content=\"0.5;URL=/login\" /></head>")
}

func main() {
	//Loads configuration for connection to LDAP Directory
	configLoad()
	//Populates userTable so can check if user is present and convert between sAMAccountName and DistinguishedName
	userTableFill()
	//Opens connection to the database under the db variable
	db, err = sql.Open("mysql", config["sqlusername"]+":"+config["sqlpassword"]+"@tcp("+config["sqlip"]+":"+config["sqlport"]+")/"+config["dbname"]+"?parseTime=true")
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
	http.HandleFunc("/", root)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/cadet", progressionTracker)
	http.HandleFunc("/staff", staffPage)
	http.HandleFunc("/dbwrite", dbWrite)
	http.HandleFunc("/user", userManagement)
	http.HandleFunc("/load", userLoad)
	http.HandleFunc("/modify", userModify)
	http.HandleFunc("/add", userAdd)
	http.HandleFunc("/remove", userRemove)
	http.Handle("/resources/", http.StripPrefix("/resources", http.FileServer(http.Dir("resources"))))
	err = http.ListenAndServe(":9090", nil) // set listen port
	//Error handler
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
