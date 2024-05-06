# NATS queue simple example

Docker-compose based example for simple [NATS](https://nats.io/) queue. Just simple example for my personal introduction to NATS


Sender based on `localhost:8080`
```
/send
Path for sending simple message (current unix time) in queue
Use ?message= query param for sending time and some text message
```
Reader/listener/etc. based on `localhost:9090`
```
/listen
Path for reading messages from storage
Every /listen call will read one message from storage
```

## How to run
1. Clone it
2. `docker compose up` from project root
3. ✨ You are awesome ✨ (also you can check `localhost:8222` for some stats)
