# Pandat
## A tiny job queue based on redis

`docker pull hiteshjoshi/pandat`

A job queue server

ENV variables 
- REDIS , eg localhost:6379


## API
- POST to /events to create new event
    - BODY :`{"url":"http://urltohit.com","interval":"0 30 * * * ","name":"some_event_ name"}`
- DELETE to /events/{eventID} to remove an event
- GET /events to get all events

## TO-DO
- [x] HTTP API
- [x] Redis PUB/SUB
- [ ] Add authentication for sockets
- [ ] User management
- [ ] Tests!!
- [ ] Clusters, master-slave? for horizontal scaling
