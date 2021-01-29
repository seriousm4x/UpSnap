// click wake button
$('.wake-form').submit(function (e) {
    $.post('/wake/' + this.id, $(this).serialize());
    e.preventDefault();
    document.getElementById(this.id + "-btn-wake").classList.add("is-loading");
});

// get available devices
async function getDevices() {
    let response = await fetch("/devices/");
    let data = await response.json()
    if (data.status == 200) {
        return data.devices
    }
}

// ping devices and set status in web interface
async function ping(device) {
    let response = await fetch("/status/" + device.pk);
    let data = await response.json()
    if (data.status == 200 && data.up == true) {
        document.getElementById(device.pk + "-status").innerHTML = "<span class=\"dot-green\"></span>";
        document.getElementById(device.pk + "-btn-wake").classList.remove("is-loading");
        document.getElementById(device.pk + "-btn-wake").disabled = true;
    } else {
        document.getElementById(device.pk + "-status").innerHTML = "<span class=\"dot-red\"></span>";
        document.getElementById(device.pk + "-btn-wake").disabled = false;
    }
}

async function main() {
    var allDevices = await getDevices();
    allDevices.forEach(device => {
        setInterval(() => {
            ping(device)
        }, 3000)
    });
}

main()