//
// init stuff
//

window.onload = () => {
    document.getElementById("scan-button").addEventListener("click", scan);

    // save device buttons
    const saveButtons = document.querySelectorAll('[id=save-button]');
    for (let index = 0; index < saveButtons.length; index++) {
        const element = saveButtons[index];
        element.addEventListener('click', () => {
            const devNameLabelContainer = document.querySelector(`[id=dev-name-label-container][data-id="${element.dataset.id}"]`)
            const devNameInputContainer = document.querySelector(`[id=dev-name-input-container][data-id="${element.dataset.id}"]`)
            const devNameInput = document.querySelector(`[id=dev-name-input][data-id="${element.dataset.id}"]`)

            const devIpLabelContainer = document.querySelector(`[id=dev-ip-label-container][data-id="${element.dataset.id}"]`)
            const devIpInputContainer = document.querySelector(`[id=dev-ip-input-container][data-id="${element.dataset.id}"]`)
            const devIpInput = document.querySelector(`[id=dev-ip-input][data-id="${element.dataset.id}"]`)

            const devMacLabel = document.querySelector(`[id=dev-mac][data-id="${element.dataset.id}"]`)

            devNameLabelContainer.innerText = devNameInput.value;
            devNameInputContainer.classList.toggle("is-hidden");
            devNameLabelContainer.classList.toggle("is-hidden");

            devIpLabelContainer.innerText = devIpInput.value;
            devIpInputContainer.classList.toggle("is-hidden");
            devIpLabelContainer.classList.toggle("is-hidden");

            post_settings("/settings/update/", JSON.stringify({
                "name": devNameInput.value,
                "ip": devIpInput.value,
                "mac": devMacLabel.innerText
            }));
            element.classList.toggle("is-hidden");
            const editButton = document.querySelector(`[id="edit-button"][data-id="${element.dataset.id}"]`)
            editButton.classList.toggle("is-hidden");
        });
    }

    // edit device buttons
    const editButtons = document.querySelectorAll('[id=edit-button]');
    for (let index = 0; index < editButtons.length; index++) {
        const element = editButtons[index];
        element.addEventListener('click', () => {
            element.classList.toggle("is-hidden");
            const saveButton = document.querySelector(`[id="save-button"][data-id="${element.dataset.id}"]`)
            saveButton.classList.toggle("is-hidden");

            // Swap table cells with input fields
            const devNameLabelContainer = document.querySelector(`[id=dev-name-label-container][data-id="${element.dataset.id}"]`)
            devNameLabelContainer.classList.toggle("is-hidden");
            const devNameInputContainer = document.querySelector(`[id=dev-name-input-container][data-id="${element.dataset.id}"]`)
            devNameInputContainer.classList.toggle("is-hidden");

            const devIpLabelContainer = document.querySelector(`[id=dev-ip-label-container][data-id="${element.dataset.id}"]`)
            devIpLabelContainer.classList.toggle("is-hidden");
            const devIpInputContainer = document.querySelector(`[id=dev-ip-input-container][data-id="${element.dataset.id}"]`)
            devIpInputContainer.classList.toggle("is-hidden");
        });
    }

    // delete device buttons
    const deleteButtons = document.querySelectorAll('[id=delete-button]');
    for (let index = 0; index < deleteButtons.length; index++) {
        const element = deleteButtons[index];
        element.addEventListener('click', () => {
            post_settings("/settings/del/", element.dataset.id);
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
        const row = table.insertRow();
        const td1 = row.insertCell();
        const td2 = row.insertCell();
        const td3 = row.insertCell();
        const td4 = row.insertCell();
        const deviceName = "Unknown";
        if (device.name != "") {
            deviceName = device.name;
        }
        td1.appendChild(document.createTextNode(deviceName));
        td2.appendChild(document.createTextNode(device.ip));
        td3.appendChild(document.createTextNode(device.mac));
        const button = document.createElement("button");
        button.classList.add("button", "is-primary", "is-small");
        button.innerText = "Add";
        button.addEventListener("click", function (event) {
            update_device(JSON.stringify(device));
            event.target.disabled = true;
            event.target.innerText = "Added";
        });
        td4.appendChild(button);
    }
    document.getElementById("scan-button").classList.remove("is-loading");
    tableContainer.classList.remove("is-hidden");
}

async function post_settings(url, data) {
    await fetch(url, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: data
        })
        .catch((error) => {
            console.error(error);
        })
}
