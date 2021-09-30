//
// notifications
//

class BulmaNotification {
    constructor() {
        // Create DOM notification structure when instantiated
        this.init();
    }

    init() {
        this.hideTimeout = null;

        // Creating the notification container div
        this.containerNode = document.createElement("div");
        this.containerNode.className = "notification note";

        // Adding the notification title node
        this.titleNode = document.createElement('p');
        this.titleNode.className = "note-title";

        // Adding the notification message content node
        this.messageNode = document.createElement('p');
        this.messageNode.className = "note-content";

        // Adding a little button on the notification
        this.closeButtonNode = document.createElement('button');
        this.closeButtonNode.className = "delete";
        this.closeButtonNode.addEventListener('click', () => {
            this.close();
        });

        // Appending the container with all the elements newly created
        this.containerNode.appendChild(this.closeButtonNode);
        this.containerNode.appendChild(this.titleNode);
        this.containerNode.appendChild(this.messageNode);

        // Inserting the notification to the page body
        document.body.appendChild(this.containerNode);
    }

    // Make the notification visible on the screen
    show(title, message, context, duration) {
        clearTimeout(this.hideTimeout);
        this.containerNode.className = "notification note";
        this.containerNode.classList.add("note-visible");

        // Setting a title to the notification
        if (title != undefined)
            this.titleNode.textContent = title;
        else
            this.titleNode.textContent = "Notification Title"

        // Setting a message to the notification
        if (message != undefined)
            this.messageNode.textContent = message;
        else
            this.messageNode.textContent = "Notification message";

        // Applying Bulma notification style/theme
        if (context) {
            var classList = context.split(" ")
            for (let index = 0; index < classList.length; index++) {
                this.containerNode.classList.add(classList[index]);
            }
        }

        // Default duration delay
        if (duration == undefined)
            duration = 1500;
        else if (duration <= 0 && duration != -1) {
            console.error('Bulma-notifications : the duration parameter value is not valid.' + "\n" +
                'Make sure this value is strictly greater than 0.');
        } else if (duration != -1) {
            // Waiting a given amout of time  
            this.hideTimeout = setTimeout(() => {
                this.close();
            }, duration);
        }
    }

    // Hide notification
    close() {
        this.containerNode.classList.remove("note-visible");
    }
}

//
// init stuff
//

let notif;
window.onload = () => {
    notif = new BulmaNotification();

    // set date for schedule inputs
    var now = new Date();
    var utcString = now.toISOString().substring(0, 19);
    var year = now.getFullYear();
    var month = now.getMonth() + 1;
    var day = now.getDate();
    var hour = now.getHours();
    var minute = now.getMinutes();
    var second = now.getSeconds();
    var localDatetime = year + "-" +
        (month < 10 ? "0" + month.toString() : month) + "-" +
        (day < 10 ? "0" + day.toString() : day) + "T" +
        (hour < 10 ? "0" + hour.toString() : hour) + ":" +
        (minute < 10 ? "0" + minute.toString() : minute) +
        utcString.substring(16, 19);
    
    // prepare html elements
    var datetimeFields = document.querySelectorAll('[id*=-input]');
    for (let index = 0; index < datetimeFields.length; index++) {
        const element = datetimeFields[index];
        element.value = localDatetime;
    }
    var wakeButtons = document.querySelectorAll(`[id="btn-wake"]`)
    for (let index = 0; index < wakeButtons.length; index++) {
        const element = wakeButtons[index];
        element.addEventListener('click', () => {
            wakeDevice(element.dataset.id);
        });
    }
    var scheduleSaveButtons = document.querySelectorAll(`[id="btn-schedule-save"]`)
    for (let index = 0; index < scheduleSaveButtons.length; index++) {
        const element = scheduleSaveButtons[index];
        element.addEventListener('click', () => {
            addSchedule(element.dataset.id);
        });
    }
    var scheduleDeleteButtons = document.querySelectorAll(`[id="btn-schedule-delete"]`)
    for (let index = 0; index < scheduleDeleteButtons.length; index++) {
        const element = scheduleDeleteButtons[index];
        element.addEventListener('click', () => {
            deleteSchedule(element.dataset.id);
        });
    }
};

//
// set device status
//

function setDeviceUp(device) {
    // get elements. if devices have been added via settings page, some elements might not exist until the user reloads the page.
    var deviceBox = document.getElementById(device.id + "-container");
    var statusDot = document.getElementById(device.id + "-dot");
    var statusPorts = document.getElementById(device.id + "-ports");
    var wakeButton = document.querySelector(`[id="btn-wake"][data-id="${device.id}"]`)
    var scheduleModalButton = document.getElementById(device.id + "-btn-schedule");

    // check if device was down before
    if (statusDot.classList.contains("dot-down") && enableNotifications) {
        notif.show("Device now up!", device.name + " is now up.", "is-success", 5000);
    }

    // clear current animation
    statusDot.style.animation = "none";
    statusDot.offsetWidth;

    const openPorts = []

    // set dot
    statusDot.classList.remove("dot-waiting");
    statusDot.classList.remove("dot-down");
    statusDot.classList.add("dot-up");
    statusDot.style.animation = "green-pulse 1s normal";
    // set box border left
    deviceBox.classList.remove("box-waiting");
    deviceBox.classList.remove("box-down");
    deviceBox.classList.add("box-up");
    // set ports
    if (device.vnc) {
        openPorts.push("<li><i class='fas fa-check fa-fw'></i> <strong>VNC (5900)</strong></li>");
    } else {
        openPorts.push("<li><i class='fas fa-times fa-fw'></i> VNC (5900)</li>");
    }
    if (device.rdp) {
        openPorts.push("<li><i class='fas fa-check fa-fw'></i> <strong>RDP (3389)</strong></li>");
    } else {
        openPorts.push("<li><i class='fas fa-times fa-fw'></i> RDP (3389)</li>");
    }
    if (device.ssh) {
        openPorts.push("<li><i class='fas fa-check fa-fw'></i> <strong>SSH (22)</strong></li>");
    } else {
        openPorts.push("<li><i class='fas fa-times fa-fw'></i> SSH (22)</li>");
    }
    statusPorts.innerHTML = '<ul>' + openPorts.join('') + '</ul>';
    // set wake btn
    wakeButton.classList.remove("is-loading");
    wakeButton.disabled = true;
    // set schedule button
    scheduleModalButton.disabled = false;
}

function setDeviceDown(device) {
    // get elements. if devices have been added via settings page, some elements might not exist until the user reloads the page.
    var deviceBox = document.getElementById(device.id + "-container");
    var statusDot = document.getElementById(device.id + "-dot");
    var statusPorts = document.getElementById(device.id + "-ports");
    var wakeButton = document.querySelector(`[id="btn-wake"][data-id="${device.id}"]`)
    var scheduleModalButton = document.getElementById(device.id + "-btn-schedule");    

    // check if device was up before
    if (statusDot.classList.contains("dot-up") && enableNotifications ) {
        notif.show("Device now down!", device.name + " is now down.", "is-danger", 5000);
    }

    // clear current animation
    statusDot.style.animation = "none";
    statusDot.offsetWidth;

    const openPorts = []

    // set dot
    statusDot.classList.remove("dot-waiting");
    statusDot.classList.remove("dot-up");
    statusDot.classList.add("dot-down");
    statusDot.style.animation = "red-pulse 1s normal";
    // set box border left
    deviceBox.classList.remove("box-waiting");
    deviceBox.classList.remove("box-up");
    deviceBox.classList.add("box-down");
    // set ports
    openPorts.push("<li><i class='fas fa-times fa-fw'></i> VNC (5900)</li>");
    openPorts.push("<li><i class='fas fa-times fa-fw'></i> RDP (3389)</li>");
    openPorts.push("<li><i class='fas fa-times fa-fw'></i> SSH (22)</li>");
    statusPorts.innerHTML = '<ul>' + openPorts.join('') + '</ul>';
    // set wake btn
    wakeButton.disabled = false;
    // set schedule button
    scheduleModalButton.disabled = false;
}

function wakeDevice(id) {
    socket.send(JSON.stringify({
        "message": "wake",
        "id": id
    }));
}

function addSchedule(id) {
    var datetime = document.getElementById(id + "-input").value
    document.querySelector(`[id="btn-schedule-delete"][data-id="${id}"]`).disabled = false;
    if (!(datetime)) {
        return;
    }
    socket.send(JSON.stringify({
        "message": "add_schedule",
        "id": id,
        "datetime": datetime
    }));
}

function deleteSchedule(id) {
    document.querySelector(`[id="btn-schedule-delete"][data-id="${id}"]`).disabled = true;
    socket.send(JSON.stringify({
        "message": "delete_schedule",
        "id": id
    }));
}

//
// websocket handling
//

var socket = new WebSocket("ws://" + location.host + "/wol/");
socket.onmessage = function (event) {
    var message = JSON.parse(event.data);
    if (enableConsoleLogging) {
        console.log(message);
    }

    // set devices up or down
    if ("device" in message) {
        if (message.device.up == true) {
            setDeviceUp(message.device);
        } else {
            setDeviceDown(message.device);
        }
    }

    // set visitors element
    if ("visitors" in message) {
        if (message.visitors == 1) {
            document.getElementById("visitors").innerHTML = message.visitors + ' visitor';
        } else {
            document.getElementById("visitors").innerHTML = message.visitors + ' visitors';
            if (enableNotifications) {
                notif.show("Visitors updated", "There are currently " + message.visitors + " visitors", "is-info", 5000);
            }
        }
    }

    // set wake by client
    if ("wake" in message) {
        document.getElementById(message.wake.id + "-btn-wake").classList.add("is-loading");
        if (enableNotifications) {
            notif.show("Wake started", message.wake.name + " has been started.", "is-info", 5000);
        }
    }

    // set wake by schedule
    if ("wake_schedule" in message) {
        document.getElementById(message.wake_schedule.id + "-btn-wake").classList.add("is-loading");
        document.getElementById(message.wake_schedule.id + "-schedule-notice").innerHTML = "";
        if (enableNotifications) {
            notif.show("Scheduled wake started", message.wake_schedule.name + " has been started.", "is-info", 5000);
        }
    }

    // add schedule
    if ("add_schedule" in message) {
        document.getElementById(message.add_schedule.id + "-schedule-notice").innerHTML = `<p>Scheduled wake:<br>${message.add_schedule.datetime}</p>`;
        if (enableNotifications) {
            notif.show("Schedule added", "A wake up event has been scheduled for " + message.add_schedule.name, "is-info", 5000);
        }
    }

    // delete schedule
    if ("delete_schedule" in message) {
        document.getElementById(message.delete_schedule.id + "-schedule-notice").innerHTML = "";
        if (enableNotifications) {
            notif.show("Schedule deleted", "A wake up event has been deleted for " + message.delete_schedule.name, "is-info", 5000);
        }
    }

    // reload page if device added via settings page
    if ("reload" in message) {
        window.location.reload(true);
    }
    
}
socket.onclose = function (event) {
    if (enableNotifications) {
        notif.show("Connection closed", "Websocket connection has closed", "is-danger", 5000);
    }
}
