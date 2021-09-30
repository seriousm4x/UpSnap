//
// init stuff
//

window.onload = () => {
    document.getElementById("scan-button").addEventListener("click", scan);

    var deleteButtons = document.querySelectorAll('[id=delete-button]');
    for (let index = 0; index < deleteButtons.length; index++) {
        const element = deleteButtons[index];
        element.addEventListener('click', () => {
            del_device(element.dataset.id);
            element.disabled = true;
            element.innerText = "Deleted";
        });
    }
}

//
// scan for devices and create table
//

async function scan() {
    const table = document.getElementById("scan-table").getElementsByTagName('tbody')[0];
    table.innerHTML = "";
    const tableContainer = document.getElementById("scan-table-container");
    tableContainer.classList.add("is-hidden");
    document.getElementById("scan-button").classList.add("is-loading");

    const response = await fetch("/settings/scan");
    const data = await response.json();

    if (data.devices.length == 0) {
        document.getElementById("scan-button").classList.remove("is-loading");
        return;
    }

    for (let index = 0; index < data.devices.length; index++) {
        const device = data.devices[index];
        var row = table.insertRow();
        var td1 = row.insertCell();
        var td2 = row.insertCell();
        var td3 = row.insertCell();
        var td4 = row.insertCell();
        if (device.name == "") {
            var deviceName = "Unknown";
        } else {
            var deviceName = device.name;
        }
        td1.appendChild(document.createTextNode(deviceName));
        td2.appendChild(document.createTextNode(device.ip));
        td3.appendChild(document.createTextNode(device.mac));
        var button = document.createElement("button");
        button.classList.add("button", "is-primary", "is-small");
        button.innerText = "Add";
        button.addEventListener("click", function (event) {
            add_device(JSON.stringify(device));
            event.target.disabled = true;
            event.target.innerText = "Added";
        });
        td4.appendChild(button);
    }
    document.getElementById("scan-button").classList.remove("is-loading");
    tableContainer.classList.remove("is-hidden");
}

function add_device(data) {
    var xhr = new XMLHttpRequest();
    var url = "/settings/add/";
    xhr.open("POST", url, true);
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.send(data);
}

function del_device(devId) {
    var xhr = new XMLHttpRequest();
    var url = "/settings/del/";
    xhr.open("POST", url, true);
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.send(devId);
}