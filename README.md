# gwmpd
MPD Web GUI - written in Go (multithread/async)

There are 2 parts:
* the backend which written in Go. It communicates with mpd and web GUI
* the frontend which written in JS (vuejs)

__Rest API is now secure with JWT token but you MUST use it over https__
__GUI is still ugly__

## Dependencies
These programs are only useful for build backend and frontend:
* go
* yarn (you can easily replace yarn by npm)
* git

## Install to your server
### Build
```sh
git clone https://github.com/Chipsterjulien/gwmpd.git
cd gwmpd/
chmod +x auto_build.sh
./auto_build.sh
```

### Backend
After building, move back/gwmpdBack to your /usr/bin and:
```sh
chmod +x /usr/bin/gwmpdBack
```

Add it in your initd.
To finish:
```sh
mkdir /etc/gwmpd
mv back/cfg/gwmpd_sample.toml /etc/gwmpd/gwmpd.toml
```

Change /var/log/gwmpd right with chown:
```sh
chown your_user: /var/log/gwmpd
```

### Frontend
__YOU MUST USE HTTPS IF YOU EXPOSE YOUR API ON THE WEB__

After building, move all files in your root server. For example:
```sh
mv front/* /var/www
```

## Config
Edit /etc/gwmpd/gwmpd.toml by changing:
* IP address
* port numbers
* your new jwtSecretKey
* login and password of course

If you expose gwmpdBack to the web, open ginserver's port on your server by modifying your firewall and don't forget to redirect port on your modem

## Starting
Start mpd server, start /usr/bin/gwmpdBack, open your browser and finally go to your server
