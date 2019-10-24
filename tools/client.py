#!/usr/bin/env python3
import os
import sys
import requests
from optparse import OptionParser

def store(url : str, file : str):
    print("Storing {} onto {}...".format(file, url), file=sys.stderr)
    response = requests.post(url+"/kademlia/file", file)
    print("Status={}".format(response.status_code) , file=sys.stderr)
    return response.json()

def find(url : str, id : str):
    print("Retrieving {} from {}...".format(id, url), file=sys.stderr)
    response = requests.get(url+"/kademlia/file", params={ 'id' : id })
    print("Status={}".format(response.status_code) , file=sys.stderr)
    return response.json()

if __name__ == "__main__":
    usage = "usage: %prog [options] command IP:port"
    parser = OptionParser(usage=usage)
    parser.add_option("--file", dest="file", help="file to store")
    parser.add_option("--id", dest="id", help="ID of file to find")

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
    if cmd == "store" or cmd == "find":
        if (cmd == "store" and not options.file):
            parser.error("No file specified to store")
        elif (cmd=="find" and not options.id):
            parser.error("No id specified to find")
        else:
            print(store(url, options.file) if cmd == "store" else find(url, options.id))

    elif cmd == "exit":
        print("Calling...")
        response = requests.post(url+"/node/exit")
        print("Status: ", response.status_code)

    else:
        print("Command ", cmd ," not found")
