#!/usr/bin/env python3
import os
import sys
import requests
import string
import random
import concurrent.futures
from optparse import OptionParser
from functools import partial

URL = "http://172.18.0.1:8081"
K = 10

def store(url: str, file: str):
    response = requests.post(url + "/kademlia/file", file)
    if response.status_code != 200:
        print("Failed to store {}".format(file))
        exit(1)
    else:
        e = {"id": response.json()["FileID"], "original_file": file}
        print("Stored {}".format(e))
        return e


def post_nrandom_string(url: str, nb: int) -> [str]:
    # Generate data to store
    def randomString(stringLength=10):
        letters = string.ascii_lowercase
        return ''.join(random.choice(letters) for i in range(stringLength))

    data = [randomString() for _ in range(nb)]

    # Post them
    futures_results = []
    with concurrent.futures.ThreadPoolExecutor() as executor:
        storeon = partial(store, url)
        futures_results = executor.map(storeon, data)

    return list(e for e in futures_results)


def find(url, id):
    response = requests.get(url + "/kademlia/file", params={'id': id})
    if response.status_code != 200:
        print("Failed to find {}".format(id))
        exit(1)
    elif response.json()["Data"]:
        return response.json()["Data"]
    else:
        print("Data badly formated for {}, received {}".format(id, response))
        exit(1)


reslist = []
while True:
    reslist.extend(post_nrandom_string(URL, K))
    print("Stored {} test with sucess !".format(len(reslist)))
    for to_find in random.choices(reslist, k=K):
        data = find(URL, to_find['id'])
        if data == to_find['original_file']:
            print("Found ", to_find)
        else:
            print("Found {} with bad data {}".format(to_find, data))
