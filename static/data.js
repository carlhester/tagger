(function () {
    var httpRequest;
    document.getElementById("ajaxButton").addEventListener('click', makeRequest);

    function makeRequest() {
        httpRequest = new XMLHttpRequest();
        httpRequest.onreadystatechange = getContents;
        httpRequest.open('GET', '/data');
        httpRequest.send();
    }

    function getContents() {
        const queryTable = document.querySelector('#queries');
        var respData = ""

        if (httpRequest.readyState === XMLHttpRequest.DONE) {
            if (httpRequest.status === 200) {
                respData = JSON.parse(httpRequest.responseText);
                var row = queryTable.insertRow(1)
                var c1 = row.insertCell(0)
                var c2 = row.insertCell(1)

                c1.innerHTML = respData['Link']
                c2.innerHTML = respData['Tags'][0] + ',' + respData['Tags'][1]

                console.log(respData)
            }
        }

    }
})();
