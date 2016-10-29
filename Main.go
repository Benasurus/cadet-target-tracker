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
    "math"
    "net/http"
    "strconv"
    "strings"
)

//Defining global variables
var db *sql.DB //database variable
var err error  //error variable
var store = sessions.NewCookieStore(securecookie.GenerateRandomKey(32)) //cookie store variable

//defining a structure to store sectorLabel and sectorDescription in a map variable
type segmentInfo struct {
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


func userData(user string, segmentData map[string]bool) (bool, map[string]bool) {
    //QueryTransaction module is used to query transactions table for rows that match a user
    var tID string
    var timeStamp, userName string
    var result [4]int
    //Queries databse using SQL calling all rows which match the desired username & scans through returned row 
    err := db.QueryRow("SELECT * FROM transactions WHERE UserName=? ORDER BY timeStamp DESC LIMIT 1", user).Scan(&tID, &timeStamp,&userName, &result[0], &result[1], &result[2], &result[3])
    //Error handler
	if err != nil {
		log.Println(err)
		return false, segmentData
    }

    for i := 1; i < 5; i++ {
        for j := 1; j < 5; j++ {
            id := "sec"+strconv.Itoa(i)+"-seg"+strconv.Itoa(j)
            if j-1 < result[i-1] {
                segmentData[id] = true
            } else {
                segmentData[id] = false
            }
        }
    }
    //Prints results
	    fmt.Println(userName, timeStamp)
    return true, segmentData
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

func info(sectors, segments float64)(string) {
    //Define variables
    var segmentID, segmentLabel, segmentDescription, id, segID, div string
    var countSegment, countSector float64
    //Map variable which stores information on each segement
    var segmentMap map[string]string
    
    segmentMap = make(map[string]string)
    //Queries the sectors database
    rows, err := db.Query("SELECT * FROM segments")
    fmt.Println(rows)
    //Error handler
    if err != nil {
		log.Println(err)
		return "Error"
    }
    defer rows.Close()
    //Scans through returned values and retrieves sector names and inserts them into the defined variables
    for rows.Next() {
        err := rows.Scan(&segmentID, &segmentLabel, &segmentDescription)
        if err != nil {
		     log.Fatal(err)
             return "Error"
        }
        segmentMap[segmentID] = segmentLabel
    }
    for countSegment = 0; countSegment < segments; countSegment ++ {
        for countSector = 0; countSector < sectors; countSector ++ {
            id = "div-sec" + strconv.FormatFloat(countSector+1, 'f', 0, 64) + "-seg" + strconv.FormatFloat(countSegment, 'f', 0, 64)
            segID = "sec" + strconv.FormatFloat(countSector+1, 'f', 0, 64) + "-seg" + strconv.FormatFloat(countSegment, 'f', 0, 64)
            div =  "<div id=\""+id+"\">"+segmentMap[segID]+"</div>"
        }
    }
    return div
}

func sector(r1 float64, width float64, theta float64, angle float64, offsetx float64, offsety float64, attrs string, id string, segment string)(output string) {
    //Variable definition
    var sigma float64
    var r2 float64
    var x1, x2, x3, x4 float64
    var y1, y2, y3, y4 float64
    //Degrees to Radians conversion
    theta = theta*math.Pi/180
    angle = angle*math.Pi/180
    //Calculate Sigma
    sigma = theta+angle
    //Calculate Width
    r2 = r1+width
    //Point 1
    x1 = -r1*math.Sin(theta)+offsetx
    y1 = -r1*math.Cos(theta)+offsety
    //Point 2
    x2 = -r2*math.Sin(theta)+offsetx
    y2 = -r2*math.Cos(theta)+offsety
    //Point 3
    x3 = -r2*math.Sin(sigma)+offsetx
    y3 = -r2*math.Cos(sigma)+offsety
    //Point 4
    x4 = -r1*math.Sin(sigma)+offsetx
    y4 = -r1*math.Cos(sigma)+offsety
    //Ouput in format for use in HTML SVG
        //Round values to nearest int & convert to string format
        x1s := strconv.FormatFloat(x1, 'f', 0, 64) 
        x2s := strconv.FormatFloat(x2, 'f', 0, 64)
        x3s := strconv.FormatFloat(x3, 'f', 0, 64)
        x4s := strconv.FormatFloat(x4, 'f', 0, 64)
        //xms := strconv.FormatFloat(xm, 'f', 0, 64)
        y1s := strconv.FormatFloat(y1, 'f', 0, 64)
        y2s := strconv.FormatFloat(y2, 'f', 0, 64)
        y3s := strconv.FormatFloat(y3, 'f', 0, 64)
        y4s := strconv.FormatFloat(y4, 'f', 0, 64)
        //yms := strconv.FormatFloat(ym, 'f', 0, 64)
        r1s := strconv.FormatFloat(r1, 'f', 0, 64)
        r2s := strconv.FormatFloat(r2, 'f', 0, 64)
    
    
    output = "<path "+attrs+" d=\"M"+x1s+" "+y1s+" L"+x2s+" "+y2s+" A"+r2s+" "+r2s+" 0 0 0 "+x3s+" "+y3s+" L"+x4s+" "+y4s+" A"+r1s+" "+r1s+" 0 0 1 "+x1s+" "+y1s+"\"/>\n"
    return output
}


func roundelBuilder(sectors float64 ,segments float64, offsetx float64, offsety float64, user string)(string) {
    //RoundelBuilder module generates HTML to build each segment of the roundel
    //Defining variables
    var r1, angleIncrement, segmentWidth, countSegment, countSector float64
    var attrs, path, id, fill string
    //Map variable which stores user data for each segment
    var segmentData map[string]bool

    segmentData = make(map[string]bool)
    //Calls user data
    ok, segmentData := userData(user, segmentData)
    if ok == false {
        return "Error"
    }
    //Calculates angle and width
    angleIncrement = 360 / sectors
    segmentWidth = (offsetx / segments)*0.95
    r1 = 20
    path = ""
    //For loop which goes through and builds all segments
    for countSegment = 0; countSegment < segments; countSegment ++ {
        for countSector = 0; countSector < sectors; countSector ++ {
            id = "sec" + strconv.FormatFloat(countSector+1, 'f', 0, 64) + "-seg" + strconv.FormatFloat(countSegment, 'f', 0, 64)
            if segmentData[id]==true {
                fill = "yellow"
            }  else {
                fill = "transparent"
                }
            attrs = "id=\"" +id+ "\" stroke=\"black\" fill=\"" +fill+ "\" onclick=\"doSetHighlight('" +id+ "');\""
            path += sector(r1, segmentWidth, countSector*angleIncrement, angleIncrement, offsetx, offsety, attrs, id, strconv.FormatFloat(countSegment, 'f', 0, 64))
        }
        r1 += segmentWidth
    }
    return path
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
        pagevars := map[string]interface{}{
            "Path"  : template.HTML(roundelBuilder(15,5, 400,400, "tesdt"))}
        t, err := template.ParseFiles("test.html")
        if err != nil {
            log.Println(err)
        } else {
            err = t.Execute(w, pagevars)
               if err != nil {
                log.Println(err)
                }
        }
        pagevars = map[string]interface{}{
            "Path"  : template.HTML(roundelBuilder(15,5, 400,400, "tesdt"))}
        t, err = template.ParseFiles("test.html")
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