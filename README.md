# FediWatch

... is a tool to visualize ActivityPub (or any other traffic / network connections) on a 3D WebGL globe in a web browser. Made for **Mastodon**.


![Screenshot of the FediWatch application in a web browser](/gh-assets/fediwatch-screenshot.png?raw=true "Screenshot of the FediWatch application in a web browser")


## Demos

* Live Demo: https://watch.metalhead.club - visualizes connections from/to my Mastodon instance [metalhead.club](https://metalhead.club)
* Demo video: https://youtu.be/Ls_kgDwheYc


## What and how 

1. Information about incoming / outgoing connections to/from a specific "home server" is gathered
2. The foreign hostname is resolved to an IP address
3. The approximate location of the IP address is found via GeoIP2 Lite _(`GeoLite2-City.mmdb` - not included in the source, please download yourself!)_
4. The location data is sent to an WebGL web browser application and displayed as connections on a globe

The connection information (which hostname has connected / which hostname did the hostname connect to) is retrieved my extracting Mastodon-internal information via a Redis message queue. There were several methods considered, but this one seemed to be the simplest one. Drawback: **To make this work, some short lines of code need to be added to the original source code.** (see patch `0001-add-redis-hostname-publishing.patch`).

 Thus, this method does not work for other ActivityPub implementations directly and causes source modifications. The server operator should know how to handle patches and derived source history in a Git repo! Otherwise they will face issues when updating to a newer Mastodon release.

The following other methods have been considered to collect connection information:

* Reading unencrypted data stream after Nginx Proxy via tcpdump and extract / analyze ActivityPub traffic. Drawback: There is no nice solution for outgoing connections, where Mastodon is the HTTP client. Extracting information would require an MitM proxy here. Ugly, breaks stuff, we don't want to do that.
* Reading ActivityPub information from Sidekiq Job Queue. Not easily possible because there is no way to know what the currently processed job is outside the Ruby environment. Just listening to the Redis Sidekiq queue will not reveil this info. There is just a list of upcoming jobs. 


## Limitations

The project does not strive to be 100 % accurate. E.g. as only the hostnames of connecting entities (and not their IP address) are known, the location of the contacting/contacted server might not be accurate. If there are multiple servers for a ActivityPub domain (e.g. cluster of servers), we select one of the servers randomly, pick its IP address and act like the request originated/terminated there. 

**Also note that this project has been implemented in a hurry. I'm not particularly proud of the code quality and robustness of this piece of software and the setup instructions lack details. While you might want to try it out, you also might run into issues. If so, let my know in the GitHub issues. I'll try to help.**


## Building and running the Fediwatch application 

### Environment

The following sections assume that you use a Linux-based operating system. Make sure that you have a [Golang environment set up](https://go.dev/doc/install). As an alternative to building the app yourself, you can also check the "[Releases](https://github.com/ThomasLeister/fediwatch/releases)" page for binary files. 


## Prepare user and app directory

    adduser --home /opt/fediwatch --disabled-password fediwatch
    su - fediwatch


### Download

Download the source: 

    git clone -b master https://github.com/ThomasLeister/fediwatch.git

### Build

    cd fediwatch
    go build


### Install

FediWatch depends on a GeoIP2 database to resolve IP addresses to locations. As the file is proprietary, it cannot be bundled with the source code and thus must be downloaded from the MaxMind's web site: https://dev.maxmind.com/geoip/geolite2-free-geolocation-data
The Download is free.

Place the GeoIP2 file `GeoLite2-City.mmdb` into this directory.

### Create FediWatch configuration

    cp config.example.toml config.toml

Then edit `config.toml`, e.g.:

    httpPort              = 8010
    websocketPort         = 8011
    redisHost             = "localhost"
    redisPort             = 6379
    databasePath          = "GeoLite2-City.mmdb"
    homeLocation          = [50.1517, 8.7523]
    websocketUrl          = "ws://watch.myinstance.tld/ws"

* `homeLocation`: The location of your server as a `[latitude, longitude]` array
* `websocketUrl`: The public Websocket URL that should be used by the web application. Usually at `/ws`. Use `wss://` if your site supports HTTPS.


### Create Nginx configuration

_(as root user)_

Configure your Webserver, e.g. Nginx:

    server {
        listen 80;
        listen [::]:80;

        server_name watch.myinstance.tld;


        location / {
            proxy_pass http://127.0.0.1:8010;
            proxy_buffering off;
        }

        location /ws {
            proxy_pass http://127.0.0.1:8011;
            proxy_buffering off;
            proxy_redirect off;
            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection $connection_upgrade;
            tcp_nodelay on;
        }
    }


### Create systemd service file

_(as root user)_

Example: `/etc/systemd/system/fediwatch.service`:

    [Unit]
    Description="Fediwatch service"

    [Service]
    WorkingDirectory=/opt/fediwatch
    User=fediwatch
    Group=fediwatch
    ExecStart=/opt/fediwatch/main
    Restart=always

    [Install]
    WantedBy=multi-user.target

Reload systemd daemon

    systemctl daemon-reload

Enable and start FediWatch:

    systemctl enable --now fediwatch


### Patch Mastodon

_(as mastodon user)_

We need to apply a tiny patch for Mastodon to export contacted servers via its Redis database:

Apply Patch `0001-add-redis-hostname-publishing.patch`

    su -s /bin/bash - mastodon
    cd live
    git apply 0001-add-redis-hostname-publishing.patch


_(as root user)_

Restart all Mastodon services afterwards.


## Done!

You should be able to see FediWatch at https://watch.myinstance.tld


---


## Development notes

### Setting up the project / Compiling

* Set up a Golang environment and clone this Repo. 
* Download the GeoIP Lite database and name it `GeoLite2-City.mmdb`
* Compile the project (`go build`)


### Installing and using Protobuf

Protobuf is used to minify the amount of data being sent over Websockets.
It allows us to process structs to binary data (instead of serializing to String-based JSON). This is much more efficient.

    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    protoc -I=./ --go_out=./ fediwatch.proto

The `protoc` command always needs to be run after any changes to fediwatch.proto. After any modifications to the `.proto` file, also copy the latest version to the web server directory to make it available to the Javascript Protobuf parser:

    cp fediwatch.proto static/fediwatch.proto


### Creating a local socket for a remote Redis instance 

This is useful for Fediwatch deployments on your local dev machine, e.g. for debugging and test of this application:

 * ssh -L localhost:6379:localhost:6379 root@mastodonserver.tld -N


### Globe.js earth background image

High resolution earth images can be downloaded here: https://earthobservatory.nasa.gov/features/NightLights/page3.php
Although the NASA images are quite high-res, Globe.gl/Three.js will not be able to display the full resolution for some reason. Maybe it's just not supported by the data structures.

You may notice a warning message in the web browser console, for example:

> THREE.WebGlRenderer: Texture has been resized from (13500x6750) to (8192x4096)

It seems like images are downscaled by Three.js to dimensions that are powers of two, like 512, 1024, 2048, 4096, etc. 
But even after resizing the original high-res image to matching resolutions like 16384 x 8192, the error message persisted. Maybe the resolution is just too high. 

Therefore this application uses a downscaled image of the earth (8192x4096) to reduce loading times. Higher quality images do not have an effect due to browser rescaling anyway ... 



## Non-asked questions

* "Why did you implement this?" => I wanted to visualize ActivityPub federation to give non-techie people an idea what happens in the background.
* "Why do you use Google Protobuf?" => It's a great way to make data transfers more efficient. I wanted to try this. It's not technically required here, but still it's nice to see it perform well. Yes, of course there are numerous alternatives ... 


## 3rd party software

This project makes use of the following 3rd party software (added via Go modules)

* [Protobuf-go](https://github.com/protocolbuffers/protobuf-go) licensed under the BSD 3-Clause "New" or "Revised" License
* [go-redis](https://github.com/redis/go-redis) licensed under the BSD-2-Clause license 
* [Gorilla Websocket](github.com/gorilla/websocket) licensed under the BSD-2-Clause license
* [goip2-golang](github.com/oschwald/geoip2-golang) licensed under the ISC license
* [toml](https://github.com/BurntSushi/toml) licensed under the MIT license

This project makes use of the following 3rd party software (included in source)

* [Protobuf.js](https://github.com/protobufjs/protobuf.js): licensed under the bsd-3-clause license
* [Globe.gl](https://github.com/vasturiano/globe.gl) licenses under the MIT license