PROG_TINY_WEBSERVER_BIN=tinyWebServer
PROG_WSS_PROXY_BIN=wss-proxy
PROG_WSS_BACKEND_BIN=wss-backend

all:
	cd server && env GOOS=linux GOARCH=arm GOARM=7 go build -ldflags="-s -w" -o ../${PROG_TINY_WEBSERVER_BIN}-arm
	cd server && go build -ldflags="-s -w" -o ../${PROG_TINY_WEBSERVER_BIN}
	cd websocket-proxy && env GOOS=linux GOARCH=arm GOARM=7 go build -ldflags="-s -w" -o ../${PROG_WSS_PROXY_BIN}-arm
	cd websocket-proxy && go build -ldflags="-s -w" -o ../${PROG_WSS_PROXY_BIN}
	cd websocket-backend && env GOOS=linux GOARCH=arm GOARM=7 go build -ldflags="-s -w" -o ../${PROG_WSS_BACKEND_BIN}-arm
	cd websocket-backend && go build -ldflags="-s -w" -o ../${PROG_WSS_BACKEND_BIN}
	goupx ${PROG_TINY_WEBSERVER_BIN}-arm
	goupx ${PROG_TINY_WEBSERVER_BIN}
	goupx ${PROG_WSS_PROXY_BIN}-arm
	goupx ${PROG_WSS_PROXY_BIN}
	goupx ${PROG_WSS_BACKEND_BIN}-arm
	goupx ${PROG_WSS_BACKEND_BIN}
clean:
	rm -rf ${PROG_TINY_WEBSERVER_BIN}*
	rm -rf ${PROG_WSS_PROXY_BIN}*
	rm -rf ${PROG_WSS_BACKEND_BIN}*
	rm -rf release-${PROG_TINY_WEBSERVER_BIN}

release:
	mkdir -p release-${PROG_TINY_WEBSERVER_BIN}
	cp -a config.json public/ ssldata/  ${PROG_TINY_WEBSERVER_BIN}-arm ${PROG_TINY_WEBSERVER_BIN} release-${PROG_TINY_WEBSERVER_BIN}
	cp -a ${PROG_WSS_PROXY_BIN}-arm ${PROG_WSS_PROXY_BIN} ${PROG_WSS_BACKEND_BIN}-arm ${PROG_WSS_BACKEND_BIN} release-${PROG_TINY_WEBSERVER_BIN}

