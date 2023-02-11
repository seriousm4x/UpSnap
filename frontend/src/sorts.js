export function sortDevices(a, b) {
    return a.name.localeCompare(b.name, undefined, {
        numeric: true,
        sensitivity: 'base'
    });
}

export function sortPorts(a, b) {
    return a.number - b.number;
}
