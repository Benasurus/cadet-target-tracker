var mypostrequest=new ajaxRequest()
var running = false

function auth() {
    var username = "";
    var password = "";
    //var string = "";
    username = document.getElementById("loginUsername").value;
    password = document.getElementById("loginPassword").value;
    string = "username=" + username + "&password=" + password;
    mypostrequest.open("POST", "login", true);
    mypostrequest.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
    mypostrequest.send(string);
    running = false;
    return;
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