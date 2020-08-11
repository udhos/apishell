# apishell
Run commands using APIs.

# Install

Make sure you have a working Go installation.

    https://github.com/udhos/apishell
    cd apishell
    go install ./...

# Quick Start

## Create certificate

    openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout key.pem -out cert.pem

## Run apid server

    apid

## Use curl as client

    curl -k -u admin:admin https://localhost:8080/static/

    curl -k -u admin:admin -d '{"args":["ls","-al"]}' https://localhost:8080/api/exec/v1/

    $ cat ls.yaml
    args:
    - ls
    - /tmp

    curl -k -u admin:admin --data-binary @ls.yaml https://localhost:8080/api/exec/v1/

    # Below 'aGVsbG8K' is 'hello' encoded in base64
    curl -k -u admin:admin -d '{"args":["cat"],"stdin":"aGVsbG8K"}' https://localhost:8080/api/exec/v1/

## Use apictl client

```
apictl exec --stdin "hello world" wc

apictl exec --stdin @/etc/passwd head

apictl exec head /etc/passwd

apictl exec -- ls -la

apictl exec -- bash -c "echo -n 12345 | wc"
```
