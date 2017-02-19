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
    xhttp.open("GET", "http://localhost:9090/load?username="+userName, true);
    xhttp.send();
}

function saveData(userName) {

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
    xhttp.open("GET", "http://localhost:9090/remove?username="+userName, true);
    xhttp.send();
    }
}

function capitalizeFirstLetter(string) {
    return string.charAt(0).toUpperCase() + string.slice(1);
}