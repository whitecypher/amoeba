# Message filtering and buffering example

This project illustrates the use of filters and buffers as a means to optimize network traffic over websockets.
Events are published to any connected clients at an interval between 0.1 and 1 seconds with random x and y coordinated between 0 and 200.
A filter has been applied to exclude any coordinate below 50 or above 150. And a buffer has been added to send messages at an interval matching the latency of a ping > pong round trip.

First run the server;
```bash
git clone git@github.com:whitecypher/amoeba.git
cd amoeba
./bin/amoeba.osx-amd64
```

Then open `http://localhost:8080` in your browser.

Messages received are printed to the browser console and a running counter on the page shows the number of message that have been received. It's not fancy but I believe sufficient to illustrate the process.


