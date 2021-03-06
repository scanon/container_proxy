#!/usr/bin/env python

from sanic import Sanic
from sanic.response import json
from runner import Runner
from queue import Empty
import time

import uuid


app = Sanic("App Name")

runner = Runner()

procs = dict()

def run(cmd):
    job_id = str(uuid.uuid1())
    procs[job_id] = runner.run(cmd)
    return job_id

def read_output(job_id):
    in_q = procs[job_id].q
    msgs = []
    try:
        # Flush the queue
        while True:
            msg = in_q.get(block=False)
            msgs.append(msg)
    except Empty:
        pass
    return {'msgs': msgs}

  


@app.route("/",)
async def test(request):
    return json({"hello": "world"})

@app.route("/submit", methods=["POST"])
def post_json(request):
  data = request.json
  jid = run(data['cmd'])
  return json({ "received": True, "jid": jid})

@app.route("/output/<jid>")
async def output(request, jid):
  msgs = read_output(jid)
  return json(msgs)

if __name__ == "__main__":
    import socket
    sock = socket.socket(socket.AF_UNIX)
    sock.bind('/tmp/api.sock')
    app.run(sock=sock)
