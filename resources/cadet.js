var mypostrequest=new ajaxRequest()
var running = false

function cadetModify(sector){
    cadetTimer();
    if (sector.length == 9) {
        var sectionNumber = sector.slice(3,4);
        var sectionValue = sector.slice(8,9);
    } else {
        var sectionNumber = sector.slice(3,5);
        var sectionValue = sector.slice(9,10);
    }
    sectionNumber = parseInt(sectionNumber);
    sectionValue = parseInt(sectionValue);
    for (count = 1; count < 5; count++) {
        if (sectionValue >= count) {
            classValue = "sector-true"
        } else {
            classValue = "sector-false"
        }
        document.getElementById("sector-sec"+ parseInt(sectionNumber)+ "-set" + parseInt(count)).className = classValue;
    }
}

function cadetTimer() {
    if (running == false){
        setTimeout(cadetWrite, 30000);
        running = true;
    }
}

function cadetWrite(){
    var string = "";
    for (countSection = 1; countSection < 16; countSection++) {
		for (countSector = 4; countSector > 0; countSector--) {
			id = "sector-sec" + countSection.toString() + "-set" + countSector.toString();
			if (document.getElementById(id).className == "sector-true") {
				break
			}
		}
        string += "section"+ parseInt(countSection) + "=" + parseInt(countSector) + "&";
	}
    mypostrequest.open("POST", "dbwrite", true);
    mypostrequest.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
    mypostrequest.send(string);
    running = false;
    return false;
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
