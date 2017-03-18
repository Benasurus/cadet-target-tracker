var mypostrequest=new ajaxRequest()
var running = false

function loadData(userName) {
    var xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = function() {
        if (this.readyState == 4 && this.status == 200) {
            document.getElementById("user-edit").innerHTML = this.responseText
            var script = document.createElement('script');
            script.type = 'text/javascript';
            script.id = 'user-chart-draw';
            script.innerHTML = document.getElementById("user-chart-draw-source").innerHTML
            document.getElementsByTagName('head')[0].appendChild(script);
            drawSectorChart();
        }
    };
    xhttp.open("GET", "/load?username="+userName, true);
    xhttp.send();
}

function saveData() {
    var form = document.getElementById('userModify');
    var forename = "";
    var surname = "";
    var flight = "";
    var dob = "";
    var doe = "";
    var gender = "";
    var username = "";
    var string = "";
    if (form.checkValidity()) {
        forename = document.getElementById("forename").value;
        surname = document.getElementById("surname").value;
        if (document.getElementById('flight-a').checked) {
            flight = "true";
        } else {
            flight = "false";
        }
        dob = document.getElementById("dob").value;
        doe = document.getElementById("doe").value;
        if (document.getElementById('gender-male').checked) {
            gender = "true";
        } else {
            gender = "false";
        }
        username = document.getElementById("username").value;
        if (ValidDate(dob,doe)) {
            string = "forename=" + forename + "&surname=" + surname + "&flight=" + flight + "&dob=" + dob  + "&doe=" + doe  + "&gender=" + gender + "&username=" + username;
            mypostrequest.open("POST", "modify", true);
            mypostrequest.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
            mypostrequest.send(string);
            running = false;
            return;
        } else {
            return false;
        }
    } else {
        alert("Validation Error: Check Input Fields");
        return false;
    }
}

function deleteUser(userName) {
    var r = confirm("Confirm Removal of "+capitalizeFirstLetter(userName));
    if (r == true) {
        var xhttp = new XMLHttpRequest();
        xhttp.onreadystatechange = function() {
        if (this.readyState == 4 && this.status == 200) {
            alert(userName+" Deleted");
            reload()
        }
    };
    xhttp.open("GET", "/remove?username="+userName, true);
    xhttp.send();
    }
}

function dashboard() {
    window.location.assign("/staff");
}

function capitalizeFirstLetter(string) {
    return string.charAt(0).toUpperCase() + string.slice(1);
}

function addUser() {
    var form = document.getElementById('userAdd');
    var forename = "";
    var surname = "";
    var flight = "";
    var dob = "";
    var doe = "";
    var gender = "";
    var username = "";
    var string = "";
    if (form.checkValidity()) {
        forename = document.getElementById("Addforename").value;
        surname = document.getElementById("Addsurname").value;
        if (document.getElementById('Addflight-a').checked) {
            flight = "true";
        } else {
            flight = "false";
        }
        dob = document.getElementById("Adddob").value;
        doe = document.getElementById("Adddoe").value;
        if (document.getElementById('Addgender-male').checked) {
            gender = "true";
        } else {
            gender = "false";
        }
        username = document.getElementById("Addusername").value;
        string = "forename=" + forename + "&surname=" + surname + "&flight=" + flight + "&dob=" + dob  + "&doe=" + doe  + "&gender=" + gender + "&username=" + username;
        if (ValidDate(dob,doe)) {
            mypostrequest.open("POST", "add", true);
            mypostrequest.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
            mypostrequest.send(string);
            running = false;
            return false;
        } else {
            return false;
        }
    } else {
        alert("Validation Error: Check Input Fields");
        return false;
    }
}

function logout() {
    window.location="/logout"
}

function ValidDate(dob,doe) {
    if (isValidDate(dob) && isValidDate(doe)) {
        epochDOB = new Date(dob)
        epochDOE = new Date(doe)
        epochNow = new Date()
        if (epochDOB > epochDOE) {
            alert("Validation Error: Date of Birth is greater than Date of Enrollment");
            return false;
        }
        var yearsApart = new Date(new Date - epochDOB).getFullYear()-1970
        if (12<= yearsApart && yearsApart < 20) {
            var yearsApart = new Date(epochDOE - epochDOB).getFullYear()-1970
            if (12<= yearsApart && yearsApart <= 18) {
                return true
            } else {
                alert("Validation Error: Invalid Enrollment Date");
                return false;
            }

        } else {
            alert("Validation Error: Invalid Age");
            return false;
        }
        
    } else {
        alert("Validation Error: Invalid Dates");
        return false;
    }
return true;
}

function isValidDate(dateString) {
    // First check for the pattern
    if(!/^\d{4}-\d{1,2}-\d{1,2}$/.test(dateString))
        return false;

    // Parse the date parts to integers
    var parts = dateString.split("-");
    var day = parseInt(parts[2], 10);
    var month = parseInt(parts[1], 10);
    var year = parseInt(parts[0], 10);

    // Check the ranges of month and year
    if(year < 1000 || year > 3000 || month == 0 || month > 12)
        return false;

    var monthLength = [ 31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31 ];

    // Adjust for leap years
    if(year % 400 == 0 || (year % 100 != 0 && year % 4 == 0))
        monthLength[1] = 29;

    // Check the range of the day
    return day > 0 && day <= monthLength[month - 1];
};

mypostrequest.onreadystatechange=function(){
 if (mypostrequest.readyState==4){
  if (mypostrequest.status==200 || window.location.href.indexOf("http")==-1){
   //document.getElementById("result").innerHTML=mypostrequest.responseText
        modal = document.getElementById('myModal');
        modal.style.display = "none";
        document.getElementById("userAdd").reset();
  }
  else{
   alert("An error has occured making the request")
  }
 }
}
function ajaxRequest(){
 var activexmodes=["Msxml2.XMLHTTP", "Microsoft.XMLHTTP"] //activeX versions to check for in IE
 if (window.ActiveXObject){ //Test for support for ActiveXObject in IE first (as XMLHttpRequest in IE7 is broken)
  for (var i=0; i<activexmodes.length; i++){
   try{
    return new ActiveXObject(activexmodes[i])
   }
   catch(e){
    //suppress error
   }
  }
 }
 else if (window.XMLHttpRequest) // if Mozilla, Safari etc
  return new XMLHttpRequest()
 else
  return false
}

var modal = ""
var btn = ""
var span =""

function loadModalVariables() {
    // Get the modal
    modal = document.getElementById('myModal');
    // Get the button that opens the modal
    btn = document.getElementById("myBtn");
    // Get the <span> element that closes the modal
    span = document.getElementsByClassName("close")[0];
    // When the user clicks the button, open the modal 
    btn.onclick = function() {
        modal.style.display = "block";
    }

    // When the user clicks on <span> (x), close the modal
    span.onclick = function() {
        modal.style.display = "none";
        document.getElementById("userAdd").reset();
    }

    // When the user clicks anywhere outside of the modal, close it
    window.onclick = function(event) {
        if (event.target == modal) {
            modal.style.display = "none";
            document.getElementById("userAdd").reset();
        }
    }
}