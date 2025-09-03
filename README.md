# gRPC-Remote-Procedure-Calls <br>
create a gRPC Remote Procedure Calls <br>

# URL Shortener (gRPC in Go) <br>

A simple URL shortener implemented in **Go** using **gRPC**.  <br> 
It provides two RPC methods: <br>

- **Shorten** → takes a long URL and returns a short code  <br>
- **Resolve** → takes a short code and returns the original URL <br>

--- <br>

## 📂 Project Structure <br>

proto/ <br>
urlshortener.proto # gRPC service definition <br>
server/ <br>
main.go # gRPC server <br>
client/ <br>
main.go # gRPC client example <br>
gen/ <br>
*.pb.go # generated gRPC code <br>

# 1. The .proto file <br>
- A URLShortener service <br>
- Two RPCs: <br>

-  Shorten(original_url) → returns short code <br>

- Resolve(id) → returns original URL <br>

# 2. Server <br>

- In the gRPC server: <br> 

- We keep an in-memory map[string]URL (like your old urlDB). <br>

- Shorten generates an MD5 hash and stores the mapping. <br>

- Resolve looks up the short code and returns the original. <br>

- The server listens on port :50051. <br>

# 3. Client <br>

- The client connects to the gRPC server and: <br>

- Calls Shorten to get a short code <br>

- Calls Resolve with that short code to get the original URL back. <br>

- This replaces your old HTTP POST/GET workflow. <br>

# 4. Benefits of gRPC over HTTP here <br>

- Strongly typed contracts (via .proto) <br>

- Auto-generated clients in many languages <br>

- Faster serialization (Protobuf instead of JSON) <br>

- Easier to add new fields/methods without breaking compatibility <br>

# - Generate gRPC code <br>

- From the project root, run: <br>

protoc --go_out=. --go-grpc_out=. proto/urlshortener.proto <br>

- This will generate files inside gen/. <br>

# - Run the server <br>
go run ./server <br>

- Server will start on :50051. <br>

# - Run the client <br>

- In another terminal: <br>

go run ./client <br>


# Expected output: <br>

- Shortened: https://github.com/sufiyanshaikh01 -> 92a2d416 <br>
- Resolved: 92a2d416 -> https://github.com/sufiyanshaikh01 <br>

# How it Works <br>

- Shorten RPC <br>

- Takes an original_url <br>

- Hashes it with MD5 <br>

- Stores mapping in memory (map[string]URL) <br>

- Returns a short 8-character code <br>

- Resolve RPC <br>

- Takes a short code <br>

- Looks up in memory <br>

- Returns the original URL <br>

# Notes <br>

- This uses an in-memory store, so data is lost when server restarts.<br>

- For production, replace with Redis/Postgres.<br>

- Additions like expiration or analytics can be built on top.<br>