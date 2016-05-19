#!/bin/bash
if !  [ -f temp.key -a  -f temp.crt ] ; then
	echo "making self-signed cert"
	sudo openssl req -new -newkey rsa:4096 -days 365 -nodes -x509 \
          -subj "/C=US/ST=KY/L=Louisville/O=none/CN=none.com" \
          -keyout temp.key  -out temp.crt
fi

if [ -e /etc/profile.d/golang.sh ]; then
	. /etc/profile.d/golang.sh
fi

packerd-server --tls-certificate temp.crt --tls-key temp.key  --port 64154 --tls-port 64155


