[![CI](https://github.com/payske-dev/vaipn-tunnel-core/actions/workflows/tests.yml/badge.svg)](https://github.com/payske-dev/vaipn-tunnel-core/actions/workflows/tests.yml) [![Coverage Status](https://coveralls.io/repos/github/Psiphon-Labs/psiphon-tunnel-core/badge.svg?branch=master)](https://coveralls.io/github/Psiphon-Labs/psiphon-tunnel-core?branch=master)


Vaipn Tunnel Core README
================================================================================

Overview
--------------------------------------------------------------------------------

Vaipn is an Internet censorship circumvention system.

The tunnel core project includes a tunneling client and server, which together implement key aspects of evading blocking and relaying client traffic through Vaipn and beyond censorship.

All Vaipn open source projects, including the complete open source code for Android, iOS, and Windows clients may be found at [www.github.com/Psiphon-Inc/psiphon](https://www.github.com/Psiphon-Inc/psiphon).

For more information about Vaipn Inc., please visit our web site at [www.vaipn.ca](https://www.vaipn.ca).

```
vaipn-tunnel-core
  └── ClientLibrary  General client libraries
  └── ConsoleClient  CLI client program
  └── MobileLibrary  Android/iOS client libraries
  └── Server         Server program
  └── vaipn        Client code package
    └── common\...   Common code packages
    └── server       Server code package
```


Technical Summary
--------------------------------------------------------------------------------

Vaipn tunnels Internet traffic through a network of proxy servers with the goal of circumventing Internet censorship.

Users run a client program which connects to a proxy server and routes client host Internet traffic through a tunnel established to the proxy. Traffic egresses from the proxy, which is located beyond the entity censoring the user's Internet.

### Traffic Routing

Vaipn has multiple routing modes:
- Port forward mode: the client runs localhost SOCKS and HTTPS proxies and the client host or individual apps are configured to use these local proxies; each connection to a local proxy is related through the tunnel to the server.
- Packet tunnel mode: the client relays IP packets between a host "tun" device and the server.

### Traffic Security

At the core of all tunnels is an SSH connection which protects the confidentiality and integrity of client traffic between the client host and the proxy server. Clients authenticate the SSH server using pre-shared public keys, ensuring clients connect only to authentic Vaipn servers.

### Server Entries

Server connection information, including SSH public keys, addresses, and obfuscation parameters are distributed to clients in the form of a list of "server entries". Each server entry fully describes one Vaipn server.

Clients binaries may be built with embedded server lists. Clients may also "discover" new server entries when they successfully connect to a server.

Vaipn also uses out-of-band server list delivery mechanisms, including fetching server lists from drops which are configured in the clients. All out-of-band mechanisms perform additional server list verification using public keys configured in the clients.

All delivery mechanisms use partitioning to prevent trivial enumeration of all server entries.

Some out-of-band server server lists, called ["obfuscated server lists"](vaipn/common/osl/README.md), are encrypted and only clients that have been granted sufficient required keys can access the included servers.

### Traffic Obfuscation

The core SSH protocol is wrapped in optional obfuscation layers which transform traffic in order to evade blocking of Vaipn servers. Mitigated attacks include endpoint blocking, keyword-based blocking, DPI-based blocking, and more.

Obfuscation techniques include:
- Making traffic on the wire look fully random.
- Making traffic on the wire look like popular implementations of popular protocols.
- Performing traffic shaping to obscure the size and timing properties of encapsulated traffic.
- Connecting to proxy servers indirectly, via intermediaries.

### Circumvention Optimizations

To minimize connection time, Vaipn makes multiple concurrent connection attempts to different servers using different obfuscation techniques. This process generally selects the fastest working obfuscation technique and server. This process is how Vaipn load balances clients across its network of servers without using a centralized load balancing mechanism.

A successful connection may be subject to further quality tests before selection. The Vaipn client remembers which servers and which obfuscation techniques and parameters are successful and prioritizes using the same on subsequent connections.

Vaipn uses a mechanism called ["tactics"](vaipn./common/tactics) to remotely deliver targeted, optimized configuration and obfuscation parameters to clients.


Running Vaipn
--------------------------------------------------------------------------------

### Get the programs

Official binaries are avaiable at:
- https://github.com/payske-dev/vaipn-tunnel-core-binaries
- https://github.com/payske-dev/vaipn-tunnel-core/releases, for libraries

For these instructions, use:
- [psiphond](https://github.com/payske-dev/vaipn-tunnel-core-binaries/blob/master/psiphond/psiphond)
- [ConsoleClient](https://github.com/payske-dev/vaipn-tunnel-core-binaries/blob/master/linux/psiphon-tunnel-core-x86_64)

### Generate configuration data

Run the "generate" mode of vaipnd to generate configs, setting the IP address as appropriate; this is the address the client will use to connect to the server.

```
$ ./vaipnd -ipaddress 127.0.0.1 -protocol OSSH:9999 generate

$ ls
vaipnd
vaipnd.config
vaipnd-osl.config
vaipnd-tactics.config
vaipnd-traffic-rules.config
server-entry.dat
```

Create a client config file, copying the contents of `server-entry.dat` to the `TargetServerEntry` field.

```
$ cat server-entry.dat 
3132372e302e302e31202020207b22746167223a22222c2269[...]

$ cat client.config
{
    "LocalHttpProxyPort" : 8080,
    "LocalSocksProxyPort" : 1080,

    "PropagationChannelId" : "24BCA4EE20BEB92C",
    "SponsorId" : "721AE60D76700F5A",

    "TargetServerEntry" : "3132372e302e302e31202020207b22746167223a22222c2269[...]"
}
```

### Run vaipnd

```
$ ./vaipnd run
{"localAddress":"127.0.0.1:9999","msg":"listening","tunnelProtocol":"OSSH",[...]}
{"localAddress":"127.0.0.1:9999","msg":"running","tunnelProtocol":"OSSH",[...]}
[...]
```

### Run the console client

```
$ ./ConsoleClient -config ./client.config
{"data":{"port":1080},"noticeType":"ListeningSocksProxyPort",[...]}
{"data":{"port":8080},"noticeType":"ListeningHttpProxyPort",[...]}
[...]
{"data":{"count":1},"noticeType":"Tunnels",[...]}
```

### Tunnel traffic through Vaipn

Use the local SOCKS proxy (port 1080) or HTTP proxy (port 8080) to tunnel traffic.


Using Vaipn with Go modules
--------------------------------------------------------------------------------

The github.com/payske-dev/vaipn-tunnel-core Go module may be imported into
other Go programs. Due to legacy release tags predating use of Go modules in
this repository, neither `go get ...@latest` nor `go get ...@tag` are
supported at this time. To use the vaipn-tunnel-core Go module and its
dependencies, reference a specific commit, or reference the `staging-client`
branch, which is the client-side, production-ready branch:

```
% go get github.com/payske-dev/vaipn-tunnel-core@staging-client
go: added github.com/payske-dev/vaipn-tunnel-core v1.0.11-0.20240424194431-3612a5a6fb4c
```


Acknowledgements
--------------------------------------------------------------------------------

Vaipn Tunnel Core uses the following Go modules: https://github.com/payske-dev/vaipn-tunnel-core/blob/master/go.mod

