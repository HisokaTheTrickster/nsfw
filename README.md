
# NSFW – Name Server For the Web

**NSFW** is a lightweight DNS resolver written in Go, designed for home lab setups. It allows you to resolve custom local domains, forward public DNS queries, and apply basic domain filtering.
For more info, check out: [goDNS for lan](https://hisokathetrickster.gitbook.io/networking-demystified/projects/godns-for-lan)

## ✨ Features

- Resolve custom local domains (e.g., `jellyfin.abcd`, `nas.abcd`)
- Forward external queries to an upstream DNS server (e.g., `8.8.8.8`)
- Basic domain blocking using a blocklist
- Concurrent request handling for quicker lookup



## ⚙️ Getting Started

#### 1. Clone the Repository and build the binary

```
git clone https://github.com/HisokaTheTrickster/nsfw.git
cd nsfw
```

#### 2. Configure Local Domains

NSFW uses a records.json file to define your custom local domains. Edit it to add or update your entries
```
{
    "jellyfin.abcd": [
        {"rtype": 1, "ttl": 300, "value": "192.168.0.200"}
    ],
    "nas.abcd": [
        {"rtype": 1, "ttl": 300, "value": "192.168.0.201"}
    ]
}
```

- `rtype`: Record type (1 = A record). More about records [here](https://en.wikipedia.org/wiki/List_of_DNS_record_types)
- `ttl`: Time-to-live in seconds
- `value`: IP address or hostname of the local service

After editing, save the file. The DNS resolver will use these entries to resolve your custom domains.

#### 3. Adding Blocklist
To block certain domains, add them to blocklist.txt. Each domain should be on a single line.
You can find pre-built lists here:
https://github.com/hagezi/dns-blocklists

#### 4. Build the binary and test the resolver
```
go build .
sudo ./nsfw
```
Note: DNS uses port 53 (a privileged port), so root privileges are required unless you grant capabilities.

#### 5. Test the DNS service
Test the domain name lookup on another machine on the LAN. You should see something like this
```
# usage: nslookup domain-name dns-resolver-ip

nslookup jellyfin.abcd 192.168.0.2
Server:		192.168.0.2
Address:	192.168.0.2#53

Name:	jellyfin.abcd
Address: 192.168.0.200

==

nslookup google.com 192.168.0.2
Server:		192.168.0.2
Address:	192.168.0.2#53

Non-authoritative answer:
Name:	google.com
Address: 142.251.223.238
```

## 🚀 Running the resolver as a Service (Linux)

To run NSFW on startup, you can create a systemd service.

#### 1. Create a Service File
```
sudo nano /etc/systemd/system/nsfw.service

[Unit]
Description=NSFW DNS Resolver
After=network.target

[Service]
ExecStart=/path-to-nsfw/nsfw
WorkingDirectory=/path-to-nsfw
Restart=always

[Install]
WantedBy=multi-user.target
```

#### 2. Binding to Port 53 Without Root
Instead of running as root, you can grant the binary permission to bind to privileged ports
```
sudo setcap 'cap_net_bind_service=+ep' /path/to/nsfw
```

#### 3. Enable and Start the Service
```
sudo systemctl daemon-reload
sudo systemctl enable nsfw
sudo systemctl start nsfw
```

#### 4. View Logs
```
journalctl -u nsfw -f
```


