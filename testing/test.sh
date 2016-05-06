#!/bin/bash
ls
if !  [ -f temp.key -a  -f temp.crt ] ; then
	echo "making self-signed cert"
	sudo openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout temp.key -out temp.crt
fi

go install ./cmd/packerd-server  && packerd-server --tls-certificate temp.crt --tls-key temp.key  --port 64154 --tls-port 64155
