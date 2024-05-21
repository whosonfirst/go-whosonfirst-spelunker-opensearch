#!/bin/sh

export PATH="${PATH}:/usr/local/bin"
GOLANG=`which go`

if [ "${EUID}" -ne 0 ]
then
    echo "Please run as root"
    exit 1
fi

if [ ! -x ${GOLANG} ]
then
    echo "Missing Go binary"
    exit 1
fi

WHOAMI=`realpath $0`
SYSTEMD=`dirname ${WHOAMI}`
GO_SPELUNKER=`dirname ${SYSTEMD}`

GOMOD=readonly

USER="spelunker"
GROUP="spleunker"

SPELUNKER_SERVICE="/lib/systemd/system/wof-spelunker.service"

if getent passwd ${USER} > /dev/null 2>&1; then
    echo "${USER} user account already exists"
else
    useradd ${USER} -s /sbin/nologin -M
fi

cd $GO_SPELUNKER

${GOLANG} build -mod ${GOMOD} -ldflags="-s -w" -o /usr/local/bin/wof-spelunker-httpd cmd/httpd/main.go

for SERVICE in ${SPELUNKER_SERVICE}
do

    SERVICE_FNAME=`basename ${SERVICE}`

    if [ ! -f ${SYSTEMD}/${SERVICE_FNAME} ]
    then
	echo "Missing ${SYSTEMD}/${SERVICE_FNAME}"
	exit 1
    fi
    
    if [ -f ${SERVICE} ]
    then
	
	cp ${SYSTEMD}/${SERVICE_FNAME} ${SERVICE}
	chmod 644 ${SERVICE}
	
	echo ""
	echo "${SERVICE} installed - you will still need to run the following manually:"
	echo "	systemctl daemon-reload"
	echo "	systemctl restart ${SERVICE_FNAME}"
	
    else

	cp ${SYSTEMD}/${SERVICE_FNAME} ${SERVICE}
	chmod 644 ${SERVICE}
	
	echo ""
	echo "${SERVICE} installed - you will still need to run the following manually:"
	echo "	systemctl daemon-reload"
	echo "	systemctl enable ${SERVICE_FNAME}"	
	echo "	systemctl start ${SERVICE_FNAME}"
    fi

done
	     
