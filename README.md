# apishell
Run commands using APIs.

# Quick Start

    openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout key.pem -out cert.pem

    ./build.sh

    apid

    curl -k -u admin:admin https://localhost:8080/static/

    '{"args":["ls","-al"]}'

    curl -k -u admin:admin -d '{"args":["ls","-al"]}' https://localhost:8080/api/v1/exec/

    $ cat ls.yaml
    args:
    - ls
    - /tmp

    curl -k -u admin:admin --data-binary @ls.yaml https://localhost:8080/api/v1/exec/

    


