// click wake button
$('.wake-form').submit(function (e) {
    $.post('/wake/' + this.id, $(this).serialize());
    e.preventDefault();
    document.getElementById(this.id + "-btn-wake").classList.add("is-loading");
});

var socket = new WebSocket("ws://" + location.host + "/wol/");
socket.onmessage = function(event) {
    var message = JSON.parse(event.data);
    console.log(message);

    if ("visitors" in message) {
        if (message.visitors == 1) {
            var visitor_spelling = ' Visitor';
        } else {
            var visitor_spelling = ' Visitors';
        }
        document.getElementById("visitors").innerHTML = message.visitors + visitor_spelling;
        return;
    }

    // get elements
    var device = message.device
    var deviceBox = document.getElementById(device.id + "-container");
    var statusDot = document.getElementById(device.id + "-dot");
    var statusPorts = document.getElementById(device.id + "-ports");

    // clear current animation
    statusDot.style.animation = "none"; 
    statusDot.offsetWidth;

    const openPorts = []

    if (device.up == true) {
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
        document.getElementById(device.id + "-btn-wake").classList.remove("is-loading");
        document.getElementById(device.id + "-btn-wake").disabled = true;
    } else {
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
        document.getElementById(device.id + "-btn-wake").disabled = false;
    }
}
