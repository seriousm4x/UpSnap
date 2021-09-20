// click wake button
$('.wake-form').submit(function (e) {
    $.post('/wake/' + this.id, $(this).serialize());
    e.preventDefault();
    document.getElementById(this.id + "-btn-wake").classList.add("is-loading");
});

var socket = new WebSocket("ws://" + location.host + "/wol/");
socket.onmessage = function(event) {
    var device = JSON.parse(event.data);

    var statusDot = document.getElementById(device.id + "-status").children[0];

    // clear current animation
    statusDot.style.animation = "none"; 
    statusDot.offsetWidth;

    if (device.up == true) {
        statusDot.classList.remove("dot-red");
        statusDot.classList.add("dot-green");
        statusDot.style.animation = "green-pulse 1s normal";
        document.getElementById(device.id + "-btn-wake").classList.remove("is-loading");
        document.getElementById(device.id + "-btn-wake").disabled = true;
    } else {
        statusDot.classList.remove("dot-green");
        statusDot.classList.add("dot-red");
        statusDot.style.animation = "red-pulse 1s normal";
        document.getElementById(device.id + "-btn-wake").disabled = false;
    }
}
