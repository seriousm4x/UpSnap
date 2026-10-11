# Mac answers ping, but remote control still cannot connect

**Online does not always mean fully awake.** Ordinary Wake-on-LAN can bring some sleeping Macs into **DarkWake**: the network responds, but third-party remote-desktop or remote-control tools may still be unable to connect.

This recipe uses authenticated SSH to run `/usr/bin/caffeinate -u -t 5` on the Mac. On a compatible setup, this short user-activity assertion completes the wake so remote control becomes usable. It does not unlock the desktop or keep the Mac awake indefinitely.

| Problem | What addresses it |
| --- | --- |
| WOL wakes the network but remote control is still unavailable | Run `caffeinate` through a dedicated, restricted SSH key |
| UpSnap stops that command as soon as ping succeeds | The custom-command completion fix in [PR #1777](https://github.com/seriousm4x/UpSnap/pull/1777), or an equivalent build |
| The UI disables Wake because the Mac already answers ping | An explicit manual-wake action, if available, or the authenticated API in step 6 |

The frontend manual-wake button is separate from #1777. Neither that backend PR nor this guide adds it. Do not put wake commands in periodic **Ping**, or use **Shutdown** as a substitute wake button.

## Prerequisites and limits

- This is for a **powered-on, sleeping Mac**, not a shut-down Mac or a FileVault pre-boot unlock screen.
- Enable **Wake for network access** and **Remote Login** on the Mac. Allow only the intended account. Setting names and locations vary by macOS version.
- The actual UpSnap container must reach the Mac's SSH port. Reserve the Mac's IP in DHCP, or use a name that resolves inside the container.
- Use an UpSnap build containing #1777 or an equivalent fix. While it is unmerged, do not assume the normal release image includes it. Review/check out the PR head in an UpSnap clone and build with the repository's `Dockerfile.dev` if testing that source. This documentation PR alone changes no runtime behavior.
- Use a trusted LAN or VPN. Do not expose UpSnap or SSH directly to the internet for this recipe.

First verify manually that SSH plus `caffeinate` restores your remote-control tool after sleep. Some Macs may need WOL first to make SSH reachable; others wake on the SSH connection itself. If SSH never becomes reachable, the command cannot run. Wi-Fi, power source, laptop lid state and macOS version can affect the result.

## Values you must replace

No production IP, MAC address, account, host key or private key from the original deployment is included here.

| Value | Your value |
| --- | --- |
| `UPSNAP_CONTAINER` | Your existing UpSnap container name |
| `MAC_HOST` | Your Mac's reserved IP or resolvable DNS name |
| `MAC_USER` | An account allowed by Remote Login; `caffeinate` does not require administrator rights |
| `MAC_SSH_PORT` | Your SSH port, normally `22` |
| `UPSNAP_SOURCE_IP` | The source IP the Mac sees from the container, not your phone/laptop's address |
| Device MAC and netmask | Your target interface and network values, when using WOL |
| Public key and host fingerprint | Generate/read your own, as described below |
| UpSnap URL, device ID and login | Your values for the optional API example |

`/app/pb_data/wake` is an example **container** path, not a Mac path. If changed, update every key/config/script path and the UpSnap Wake command consistently.

## 1. Keep the setup persistent

Back up the current configuration. Reuse the existing mount at `/app/pb_data`; do not replace it with an empty directory. This recipe stores wake files in its `wake` subdirectory and requires no extra host mount or privileged mode. Container recreation preserves these files only when the same data mount is retained.

A dedicated read-only secret mount is also possible on hosts that support and permit it. If a managed host rejects a configuration, use its supported configuration or administrative approval process; do not disable authorization plugins for this guide.

**Data backups may include the wake-only private key. Treat them as credentials.** The key is restricted, but its holder on an allowed source can still repeatedly wake the Mac. UpSnap administrators can run custom commands, so file permissions do not isolate the key from them.

On the **Docker host**, enter the existing container as its configured runtime user:

```sh
UPSNAP_CONTAINER='REPLACE_WITH_YOUR_CONTAINER_NAME'
docker exec -it "$UPSNAP_CONTAINER" sh
```

Steps marked **container** run in this shell. Native-binary installations need equivalent persistent paths, an OpenSSH client and a `timeout` utility, owned/accessible by the UpSnap service user.

## 2. Generate a dedicated key

**Container:** check the tools and create a key without overwriting an existing one:

```sh
command -v ssh
command -v ssh-keygen
command -v timeout
umask 077
mkdir -p /app/pb_data/wake
chmod 700 /app/pb_data/wake
(
    set -eu
    test ! -e /app/pb_data/wake/id_ed25519
    test ! -e /app/pb_data/wake/id_ed25519.pub
    ssh-keygen -q -t ed25519 -N '' -C upsnap-wake-only \
        -f /app/pb_data/wake/id_ed25519
)
```

The dedicated key has no passphrase for unattended use. Do not reuse a personal administrator key. If a key exists, deliberately reuse or rotate it instead of overwriting it. Fix directory ownership if necessary; do not make it world-writable.

Replace both placeholders before creating the configuration:

```sh
cat > /app/pb_data/wake/config <<'CONFIG'
MAC_HOST='REPLACE_WITH_MAC_IP_OR_DNS'
MAC_USER='REPLACE_WITH_MAC_ACCOUNT'
MAC_SSH_PORT='22'
CONFIG
chmod 600 /app/pb_data/wake/config
. /app/pb_data/wake/config
```

## 3. Verify the Mac's host key

**On the Mac, in a trusted session:**

```sh
ssh-keygen -lf /etc/ssh/ssh_host_ed25519_key.pub
```

**Container:** collect a candidate host key and display its fingerprint:

```sh
. /app/pb_data/wake/config
ssh-keyscan -p "$MAC_SSH_PORT" -t ed25519 "$MAC_HOST" \
    > /app/pb_data/wake/known_hosts.candidate
ssh-keygen -lf /app/pb_data/wake/known_hosts.candidate
```

Compare with the fingerprint read on the Mac. **Scanning alone does not establish trust.** Only if they match:

```sh
mv /app/pb_data/wake/known_hosts.candidate /app/pb_data/wake/known_hosts
chmod 600 /app/pb_data/wake/known_hosts
```

Do not disable host-key checking. After a Mac reinstall or target change, verify its new identity before updating this file.

## 4. Authorize only waking

**Container:** while the Mac is awake, use your ordinary interactive login once to verify the path and learn the source address. Enter the password yourself; do not store it in a command.

```sh
. /app/pb_data/wake/config
ssh -F /dev/null -p "$MAC_SSH_PORT" \
    -o StrictHostKeyChecking=yes \
    -o UserKnownHostsFile=/app/pb_data/wake/known_hosts \
    -o GlobalKnownHostsFile=/dev/null -o HostKeyAlgorithms=ssh-ed25519 \
    "$MAC_USER@$MAC_HOST" \
    'printf "%s\n" "$SSH_CONNECTION"; /usr/bin/caffeinate -u -t 5'
```

The first field of `SSH_CONNECTION` is `UPSNAP_SOURCE_IP`. With host networking it is commonly the Docker host/router's LAN address; check instead of assuming. Fix any awake-login failure before testing sleep.

Read **only the public key** in the container:

```sh
cat /app/pb_data/wake/id_ed25519.pub
```

**On the Mac, logged in as `MAC_USER`:** replace the source IP and public-key base64 below, preserving the restrictions before `ssh-ed25519`. The block backs up existing authorizations and appends only this entry:

```sh
(
    set -eu
    umask 077
    test ! -L "$HOME/.ssh"
    test ! -L "$HOME/.ssh/authorized_keys"
    mkdir -p "$HOME/.ssh"
    chmod 700 "$HOME/.ssh"
    auth_file="$HOME/.ssh/authorized_keys"
    entry='restrict,from="REPLACE_WITH_UPSNAP_SOURCE_IP",command="/usr/bin/caffeinate -u -t 5" ssh-ed25519 REPLACE_WITH_PUBLIC_KEY_BASE64 upsnap-wake-only'
    case "$entry" in *REPLACE_WITH*) echo 'Replace source IP and public key first.' >&2; exit 1;; esac
    if [ -f "$auth_file" ] && grep -Fqx "$entry" "$auth_file"; then
        echo 'This restricted key is already installed.'
        exit 0
    fi
    if [ -f "$auth_file" ]; then
        cp -p "$auth_file" "$auth_file.before-upsnap-$(date +%Y%m%d-%H%M%S)-$$"
    fi
    printf '\n%s\n' "$entry" >> "$auth_file"
    chmod 600 "$auth_file"
)
```

`command=` forces this key to run `caffeinate`, ignoring the requested command. `restrict` disables PTY allocation, port/agent/X11 forwarding and user SSH rc execution. `from=` limits the source. These apply to this key, not other account login methods, and do not unlock the desktop. Symlinked SSH paths stop the script for deliberate review.

If an older SSH server does not support `restrict`, consult its documentation for equivalent explicit restrictions; do not simply remove it.

## 5. Connect the command to UpSnap

**Container:** create a client with mandatory host verification and bounded waiting:

```sh
cat > /app/pb_data/wake/client.sh <<'SCRIPT'
#!/bin/sh
set -eu
WAKE_DIR=/app/pb_data/wake
. "$WAKE_DIR/config"
: "${MAC_HOST:?Set MAC_HOST in the wake config}"
: "${MAC_USER:?Set MAC_USER in the wake config}"
: "${MAC_SSH_PORT:?Set MAC_SSH_PORT in the wake config}"
case "$MAC_HOST $MAC_USER $MAC_SSH_PORT" in
    *REPLACE_WITH*) echo 'Replace the wake config placeholders first.' >&2; exit 2;;
esac
# The Mac's forced command remains authoritative for any requested command.
if [ "$#" -eq 0 ]; then set -- wake; fi
exec timeout 25 ssh -F /dev/null -T -p "$MAC_SSH_PORT" \
    -o BatchMode=yes -o IdentityAgent=none -o IdentitiesOnly=yes \
    -o StrictHostKeyChecking=yes \
    -o UserKnownHostsFile="$WAKE_DIR/known_hosts" \
    -o GlobalKnownHostsFile=/dev/null -o HostKeyAlgorithms=ssh-ed25519 \
    -o ConnectTimeout=10 -o ConnectionAttempts=1 \
    -o ServerAliveInterval=5 -o ServerAliveCountMax=2 \
    -i "$WAKE_DIR/id_ed25519" "$MAC_USER@$MAC_HOST" "$@"
SCRIPT
chmod 700 /app/pb_data/wake/client.sh
chmod 600 /app/pb_data/wake/id_ed25519
sh /app/pb_data/wake/client.sh wake
```

Success normally produces no output and takes about five seconds after authentication. Resolve authentication/host-key errors rather than disabling checks.

In the device editor:

- Enter your own **IP/FQDN**, **MAC** and **Netmask**.
- Set **Wake → custom command** to `sh /app/pb_data/wake/client.sh wake`.
- Set **Wake Timeout** to `30` seconds, longer than the client's 25-second limit.
- Keep **Ping** as a reachability check. Never put a wake command there.
- Leave schedules off unless you actually want scheduled waking.

The five-second assertion does not mean the Mac sleeps again five seconds later. Normal sleep settings and active applications determine that.

If separately verified WOL is needed to make SSH reachable, send that packet first and then run this client. SSH alone worked in the tested deployment; it is not a universal guarantee.

## 6. If the UI says online and will not let you wake

Use an explicit **Send wake** action if your version provides it. Otherwise, the existing authenticated `GET /api/upsnap/wake/{device_id}` endpoint invokes the custom command without faking an offline state.

This Python 3 example runs on your controlling computer. It prompts for your URL, device ID and login; password/token stay in memory and are not printed. Choose `users` for an ordinary UpSnap account with power permission on this device, or `_superusers` for a superuser. Use HTTPS or a trusted LAN/VPN endpoint. The device ID is the `id` in the device editor URL, not its address or display name.

```sh
python3 -c '
import getpass
import json
import urllib.error
import urllib.parse
import urllib.request

base = input("Your UpSnap URL: ").strip().rstrip("/")
if urllib.parse.urlsplit(base).scheme not in ("http", "https"):
    raise SystemExit("Use an http or https UpSnap URL.")
device_id = input("Your UpSnap device ID: ").strip()
if not device_id:
    raise SystemExit("A device ID is required.")
device = urllib.parse.quote(device_id, safe="")
collection = input("Auth collection [users / _superusers]: ").strip()
if collection not in ("users", "_superusers"):
    raise SystemExit("Choose users or _superusers.")
identity = input("Your UpSnap login: ").strip()
password = getpass.getpass("Your UpSnap password: ")

def request(path, payload=None, token=None):
    headers = {"Content-Type": "application/json"}
    if token:
        headers["Authorization"] = "Bearer " + token
    data = None if payload is None else json.dumps(payload).encode()
    req = urllib.request.Request(base + path, data=data, headers=headers)
    with urllib.request.urlopen(req, timeout=45) as response:
        return json.load(response)

try:
    auth = request("/api/collections/" + collection + "/auth-with-password",
                   {"identity": identity, "password": password})
    if not auth.get("token"):
        raise SystemExit("No login token returned. Use a client supporting your authentication flow.")
    request("/api/upsnap/wake/" + device, token=auth["token"])
except urllib.error.HTTPError as error:
    raise SystemExit("Request failed: HTTP " + str(error.code) + ". Check permissions and UpSnap logs.")
except (urllib.error.URLError, TimeoutError):
    raise SystemExit("Connection failed or timed out. Check the endpoint and UpSnap logs.")
print("Wake request completed. Now verify your remote-control connection.")
'
```

For MFA or other login flows, use a compatible client instead of disabling protection. Never publish passwords, tokens or private keys in custom commands or issue reports.

## 7. Verify full wake and persistence

1. While awake, run the client from the actual container and confirm there is no password prompt.
2. Verify the restriction: `sh /app/pb_data/wake/client.sh 'printf UNEXPECTED_SHELL_ACCESS'` should run the same short wake action and **not** print the marker. If it does, revoke this key and fix the forced command before unattended use.
3. Disconnect Screen Sharing, other remote-control sessions and ordinary SSH sessions. Put the Mac to sleep and wait at least five minutes, using a clock.
4. Trigger one UpSnap wake request. Then try the remote-control tool without touching the Mac's physical keyboard, mouse or power button.
5. In a normal authorized Mac session, inspect the matching time window:

   ```sh
   date
   pmset -g log | grep -E 'Entering Sleep|DarkWake|FullWake|Wake from' | tail -n 60
   ```

   Look for `DarkWake to FullWake` or a full `Wake` with user/HID activity near the request. A reason string alone does not identify which app caused it; correlate timestamps and the controlled test. Record the actual sleep duration and whether remote control works.
6. Restart the UpSnap container and repeat. A Docker-host/router power-cycle is a **separate test**: container restart success alone does not prove Docker and external storage start after power loss.

A contributed deployment demonstrated full wake and usable third-party remote control after several minutes of sleep. Its last captured sleep interval was **4 minutes 20 seconds**, not exactly five minutes. Container-restart checks also passed. That validates one setup, not every Mac, sleep duration or remote-control tool. Perform longer-sleep and host-reboot tests on your own setup.

## Troubleshooting and rollback

| Symptom | Check |
| --- | --- |
| SSH fails while awake | Target, Remote Login account, container routing and host pin |
| SSH times out only after sleep | Network-wake setting, power/lid state, and whether WOL first makes SSH reachable |
| Public-key authentication fails | Account, full public key, ownership/modes and source allowed by `from=` |
| Command stops after about one second | Running build must include #1777 or an equivalent fix |
| Ping is green but remote control fails | Wake log and command result, plus the remote-control tool's own service/account/network health |
| Works until container recreation | Retain the same data mount, config, key, host pin and client script |
| Container creation is rejected | Supported mount configuration and host policy; this recipe uses only the existing data mount |

To undo: restore the previous UpSnap Wake command, then remove **only this key's** matching `upsnap-wake-only` entry from the Mac's `authorized_keys`. Preserve other keys. Remove the matching client files/private key when no longer needed and handle backup copies accordingly. Do not delete the whole `pb_data` volume to revoke a key.

For an image rollback, run only one instance against the same PocketBase directory. Restoring an old `authorized_keys` backup can re-authorize a key revoked later.

## References

- [Apple: share Mac resources while asleep](https://support.apple.com/guide/mac-help/share-your-mac-resources-when-its-in-sleep-mh27905/mac)
- [Apple's caffeinate implementation](https://github.com/apple-oss-distributions/PowerManagement/blob/main/caffeinate/caffeinate.c)
- [OpenSSH authorized-key options](https://man.openbsd.org/sshd.8#AUTHORIZED_KEYS_FILE_FORMAT)
- [UpSnap wake endpoint registration](../backend/pb/pb.go)
