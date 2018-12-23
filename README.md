# gwmpd
MPD Web GUI - written in Go (multithread/async)

![Sceenshot](https://github.com/Chipsterjulien/gwmpd/blob/master/screenshot.png)

There are 2 parts:
* the backend which written in Go. It communicates with mpd and web GUI
* the frontend which written in JS (vuejs)

<!-- __Rest API is now secure with JWT token but you SHOULD use it over https__   -->

## Dependencies
These programs are only useful for build backend and frontend:
* go
* yarn (you can easily replace yarn by npm)
* git

## Use of precompiled sources
### Backend

Go inside [https://github.com/Chipsterjulien/gwmpd/tree/master/Build/Backend](https://github.com/Chipsterjulien/gwmpd/tree/master/Build/Backend). You will find
3 exe:
* A generic for x86_64 (64 bits linux)
* A generic for arm-v6
* A generic for arm-v7

### Frontend

Go inside [https://github.com/Chipsterjulien/gwmpd/tree/master/Build/Frontend](https://github.com/Chipsterjulien/gwmpd/tree/master/Build/Frontend). You will find
the latest build of frontend. Download and move files to you root web site. The gwmpd's GUI is intended to work at the root but you can easily use a subdomain


## Build sources yourself
### Build backend and frontend
```sh
git clone https://github.com/Chipsterjulien/gwmpd.git
cd gwmpd/
make build
```

Build JS is somethimes very long so you can take a coffee ;-)

You will find backend and frontend in **Build** folder

### Install backend
After building, move gwmpdBack to your /usr/bin as follow:
```sh
mv back/gwmpdBack /usr/bin
chmod +x /usr/bin/gwmpdBack
```

Add it in your initd.
To finish:
```sh
mkdir /etc/gwmpd
mv back/cfg/gwmpd_sample.toml /etc/gwmpd/gwmpd.toml
```

Create log directory and change /var/log/gwmpd right with chown:
```sh
mkdir -p /var/log/gwmpd
chown your_user: /var/log/gwmpd
```

## Config
Edit /etc/gwmpd/gwmpd.toml by changing:
* IP address
* port numbers
* your new jwtSecretKey
* login and password of course

If you expose gwmpdBack to the web __YOU SHOULD USE HTTPS__, open ginserver's port on your server by modifying your firewall and don't forget to redirect port on your modem

## Starting
Start mpd server, start /usr/bin/gwmpdBack, open your browser and finally go to your server
