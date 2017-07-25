FROM 'golang:1.8'
MAINTAINER 'basgys@gmail.com'

WORKDIR /go/src/github.com/basgys/alterego

# Copy source code
ADD . /go/src/github.com/basgys/alterego

# Compile code & link bin
RUN go install github.com/basgys/alterego && \
    ln -s /go/bin/alterego /usr/bin/alterego

# Set environment variables
ENV IP=0.0.0.0 \
    PORT=8080 \
    REDIRECTS=REDIRECT1 \
    REDIRECT1=http://127.0.0.1:8080,http://localhost:8080 \
    REQUEST_LOGGING=true \
    REDIRECT_STATUS_CODE=308

# Check container health
HEALTHCHECK CMD curl --fail http://127.0.0.1:$PORT/__health__ || exit 1

COPY entrypoint.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
CMD ["alterego"]
