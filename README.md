# conntrack Connection Tracker
Diese App basiert auf pcap und benutzt diese Lib. Es wird kein Connectiontracking im klassischen
Sinne gemacht. Es funktioniert auf "tcpdump" Technik

## Installation

```
go get install github.com/schreibe72/conntrack
```

## Für Linux auf Mac bauen

```
./build_linux_on_mac.sh
``` 
Dies erstellt ein Docker Image mit der Build Umgebung für Linux. Danach wird der eigentliche Build
Vorgang in einen Docker Container durchgeführt. Dies ist nötig, da es native C Lib Dependencies gibt (libpcap).

Das entstandene Binary lautet conntrack und kann auf den Server kopiert werden.

## Verwendung

```
sudo conntrack -d eth0 -p /tmp/conntrack_data
```

Dies erstellt einen Ordner /tmp/conntrack_data. Nach paar Sekunden werden hier eine File und Directory Struktur
entstehen.

```
conntrack_data
conntrack_data/else
conntrack_data/else/tcp
conntrack_data/else/udp
conntrack_data/else/udp/192.168.200.214-5353-224.0.0.251-5353
conntrack_data/else/udp/192.168.200.219-5353-224.0.0.251-5353
conntrack_data/else/udp/192.168.200.226-57621-192.168.200.255-57621
conntrack_data/else/udp/192.168.200.246-5353-224.0.0.251-5353
conntrack_data/in
conntrack_data/in/tcp
conntrack_data/in/udp
conntrack_data/out
conntrack_data/out/tcp
conntrack_data/out/tcp/192.168.200.230-PORT-INTERNET-443
conntrack_data/out/udp
conntrack_data/out/udp/192.168.200.230-PORT-10.59.184.10-53
conntrack_data/out/udp/192.168.200.230-PORT-10.59.48.10-53
conntrack_data/out/udp/192.168.200.230-PORT-192.168.200.255-17500
conntrack_data/out/udp/192.168.200.230-PORT-INTERNET-17500
```

Hier gibt es drei Sub Directories: else, in, out. Es wird versucht den Traffic in Incoming and Outgoing
zu unterteilen. Hierzu werden die Device IP(s) ausgelesen. Diese IP wird mit der srcip oder dstip verglichen. 
Jenachdem welche matcht, wird der Traffic als Incoming oder Outgoing unterschieden. "else" sind die Verbindungen,
die nicht zugeordnet werden können. 
Die Folder udp bzw. tcp sollten selbsterklärend sein. Die Dateien haben folgendes Namensschema:
```
SRCIP-SRCPORT-DSTIP-DSTPORT
```
IP's welche im Internet geroutet werden, werden mit dem IP Name INTERNET zusammengefasst. Der SRCPORT ist sehr oft unterschiedlich
und bei Firewall Regeln nicht wichtig, deshalb wird dieser mit dem Name PORT ausgeblendet. Damit lassen sich Connections zusammenfassen.

Da diese Zusammenfassung bei else nicht gemacht wird, kann man bei dem Tool diesen Zweig abschalten:
```
sudo conntrack -d eth0 -p /tmp/conntrack_data -e 
```

## UDP Connections und deren Richtung

Die Richtung der UDP Connection ist schwer zu erraten. Deshalb werden hier Heuristiken verwendet.
Es wird davon ausgegangen, dass dstports unter 1024 immer Outgoing sind. Zudem hat man die Möglichkeit
Ports dediziert für Incoming zu markieren. Diese werden dann bei Outgoing nicht mitgenommen. Achtung:
lsof -i -n -P | grep UDP zeigt nicht immer nur Incoming Ports an. z.b. für collectd gibt es ein Port.
Tatsächlich ist dies aber zum größtenteil Outgoing Traffic.

```
sudo conntrack -d eth0 -p /tmp/conntrack_data -u 56456 -u 45390
```
