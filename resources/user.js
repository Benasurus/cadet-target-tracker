var mypostrequest=new ajaxRequest()
var running = false

function loadData(userName) {
    var xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = function() {
        if (this.readyState == 4 && this.status == 200) {
            document.getElementById("user-panel").innerHTML = this.responseText
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
    var forename = "";
    var surname = "";
    var flight = "";
    var dob = "";
    var doe = "";
    var gender = "";
    var username = "";
    //var string = "";
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
    string = "forename=" + forename + "&surname=" + surname + "&flight=" + flight + "&dob=" + dob  + "&doe=" + doe  + "&gender=" + gender + "&username=" + username;
    alert(string);
    mypostrequest.open("POST", "modify", true);
    mypostrequest.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
    mypostrequest.send(string);
    running = false;
    return;
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

mypostrequest.onreadystatechange=function(){
 if (mypostrequest.readyState==4){
  if (mypostrequest.status==200 || window.location.href.indexOf("http")==-1){
   //document.getElementById("result").innerHTML=mypostrequest.responseText
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

// Get the modal
var modal = document.getElementById('myModal');

// Get the button that opens the modal
var btn = document.getElementById("myBtn");

// Get the <span> element that closes the modal
var span = document.getElementsByClassName("close")[0];

// When the user clicks the button, open the modal 
btn.onclick = function() {
    modal.style.display = "block";
}

// When the user clicks on <span> (x), close the modal
span.onclick = function() {
    modal.style.display = "none";
}

// When the user clicks anywhere outside of the modal, close it
window.onclick = function(event) {
    if (event.target == modal) {
        modal.style.display = "none";
    }
}