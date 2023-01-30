export function sortDevices(a, b) {
    return a.name > b.name;
}

export function sortPorts(a, b) {
    return a.number - b.number;
}
