# speisekarte
A simple web menu generator

It saves people remembering ip addresses or hostnames of services scattered throughout many servers on an intranet or in the cloud.
Instead a straightforward menu may be offered.

## How to run
- `wget https://raw.githubusercontent.com/andreldm/speisekarte/master/docker-compose.yml`
- `wget https://raw.githubusercontent.com/andreldm/speisekarte/master/config.json`
- edit config.json
- `docker-compose pull`
- `docker-compose up -d`
- open http://localhost:3000

## TODO

- [ ] Add a screenshot
- [X] Add a sample configuration file
- [X] Publish a docker image
