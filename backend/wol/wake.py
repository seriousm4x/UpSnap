import ipaddress
import wakeonlan


def wake(mac, ip, netmask):
    subnet = ipaddress.ip_network(
        f"{ip}/{netmask}", strict=False).broadcast_address
    wakeonlan.send_magic_packet(mac, ip_address=str(subnet))
