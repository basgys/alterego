IMAGE=basgys/alterego

# Dev config
export PORT=8080
export IP=127.0.0.1
export REDIRECTS=REDIRECT1
export REDIRECT1=http://127.0.0.1:8080,http://localhost:8080
export REQUEST_LOGGING=true
export REDIRECT_STATUS_CODE=308

start:
		IP=${IP} PORT=${PORT} go run main.go

dk-start:
		docker run -t ${IMAGE}

dk-build:
		docker build -t ${IMAGE} .

dk-push:
		docker push ${IMAGE}