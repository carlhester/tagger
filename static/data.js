const t = function () {
    console.log("t called")
}


const callDataAPI = function () {
    var httpRequest;

    function makeRequest() {
        httpRequest = new XMLHttpRequest();
        httpRequest.onreadystatechange = getContents;
        httpRequest.open("GET", "/data");
        httpRequest.send();
    };

    makeRequest();

    function getContents() {
        const queryTable = document.querySelector("#queries");
        var respData = "";

        if (httpRequest.readyState === XMLHttpRequest.DONE) {
            console.log("READY!")
            if (httpRequest.status === 200) {
                respData = JSON.parse(httpRequest.responseText);
                console.log(respData);
                var row = queryTable.insertRow(1);
                var c1 = row.insertCell(0);
                var c2 = row.insertCell(1);

                c1.innerHTML = respData["Link"];
                c2.innerHTML = `
                                <a class="btn btn-primary btn-sm">${respData["Tags"][0]}</a>
                                <a class="btn btn-primary btn-sm">${respData["Tags"][1]}</a>
                                `
            }
        }
    }
}

const callAddAPI = function () {
    console.log("add called")
    var httpRequest;

    (function () {
        httpRequest = new XMLHttpRequest();
        httpRequest.onreadystatechange = getContents();
        httpRequest.open("POST", "/add");
        httpRequest.send();
    })();

    function getContents() {
        if (httpRequest.readyState === XMLHttpRequest.DONE) {
            if (httpRequest.status === 200) {
                respData = JSON.parse(httpRequest.responseText);
                console.log(respData);
            } else {
                console.log(httpRequest.status);
            }
        }
    }
}


document.getElementById("ajaxButton").addEventListener("click", callDataAPI);
document.getElementById("btnSubmitNew").addEventListener("click", callAddAPI);