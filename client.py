#!/usr/bin/env python
import requests
import os

from optparse import OptionParser


if __name__ == "__main__":
    usage = "usage: %prog [options] command IP:port"
    parser = OptionParser(usage=usage)
    parser.add_option("--file", dest="file",
            help="file to store")
    parser.add_option("--id", dest="id",
            help="ID of file to find")

    options, args = parser.parse_args()

    # Parse args
    if len(args) == 0:
        parser.error("No command specified")
    cmd = args[0]

    if len(args) == 1:
        parser.error("No address specified")
    url = args[1]
    if url[0:7] != "http://":
        url = "http://"+url

    # Parse command
    if cmd == "store":
        if not options.file:
            parser.error("No file specified to store")

        print("Calling...")
        response = requests.post(url+"/kademlia/file", options.file)
        print("Status: ",response.status_code)
        print(response.json())

    elif cmd == "find":
        if not options.id:
            parser.error("No ID specified to find")

        print("Calling...")
        response = requests.get(url+"/kademlia/file", options.id)
        print("Status: ",response.status_code)
        print(response.json())

    elif cmd == "join":
        bashCommand = "go run main.go -join "+url
        os.system(bashCommand)

    elif cmd == "exit":
        print("Calling...")
        response = requests.post(url+"/node/exit")
        print("Status: ", response.status_code)

    else:
        print("Command ", cmd ," not found")
