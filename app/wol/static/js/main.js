// click wake button
$('.wake-form').submit(function (e) {
    $.post('/wake/' + this.id, $(this).serialize());
    e.preventDefault();
    document.getElementById(this.id + "-btn-wake").classList.add("is-loading");
});

var socket = new WebSocket("ws://" + location.host + "/wol/");
socket.onmessage = function(event) {
    var data = JSON.parse(event.data);
    console.log(data);

    if (data.status !== 200) {
        console.log("message status not 200");
        return;
    }

    data.devices.forEach(function(device, index) {
        if (device.up == true) {
            document.getElementById(device.id + "-status").innerHTML = "<span class=\"dot-green\"></span>";
            document.getElementById(device.id + "-btn-wake").classList.remove("is-loading");
            document.getElementById(device.id + "-btn-wake").disabled = true;
        } else {
            document.getElementById(device.id + "-status").innerHTML = "<span class=\"dot-red\"></span>";
            document.getElementById(device.id + "-btn-wake").disabled = false;
        }
    })
}
