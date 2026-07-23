const http = require("http");

const server = http.createServer((req, res) => {
  if (req.url === "/json") {
    res.writeHead(200, { "Content-Type": "application/json" });
    res.end(JSON.stringify({
      status: "success",
      message: "Hello from Zero JSON endpoint!",
    }));
  } else {
    res.writeHead(200, { "Content-Type": "text/plain" });
    res.end("Hello, World! Zero language is alive!");
  }
});

server.listen(8080);
