# D7024E - Kademlia

### Requirements
You need to have `go` and `protobuf` installed system-wide to make this compile (e.g. on Archlinux, `sudo pacman -S go protobuf`).

You need also to install the go protobuf plugin with the command `go get -u github.com/golang/protobuf/protoc-gen-go`.

Finally, install the go dependencies by doing `go get -u` inside of the project repertory.

### Build & use
To build this program, you need to run the `build.sh` file, which will make a `main.run` executable file.

To use you can either start the node alone, or with an ip to make it join an existing network.

### Start a cluster

To build the docker image, run the command `sudo docker build . -t kademlia`.

To start it, run the command either as a starting node by using `sudo docker run -p 8080:8080 -p 9090:9090 kademlia` 
or as a joining node by using `sudo docker run -p 8080:8080 -p 9090:9090 kademlia ./main.run <IP>`. 

To start a cluster ... (doc to do)