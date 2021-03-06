import os
from threading import Thread
from subprocess import Popen, PIPE
from select import select
from multiprocessing import Queue


class job:
    def __init__(self, p, q):
        self.p = p
        self.q = q

    def p(self):
        return p

    def q(self):
        return q

class Runner:
    """
    This class provides the container interface for Docker.

    """

    def __init__(self):
        """
        Inputs: config dictionary, Job ID, and optional logger
        """
        self.procs = []
        self.threads = []

    def _readio(self, p, q):
        cont = True
        last = False
        while cont:
            rlist = [p.stdout, p.stderr]
            x = select(rlist, [], [], 1)[0]
            for f in x:
                if f == p.stderr:
                    error = True
                else:
                    error = False
                line = f.readline().decode('utf-8')
                if len(line) > 0:
                    q.put({'msg': 'output', 'line': line, 'error': error})
            if last:
                cont = False
            if p.poll() is not None:
                last = True
        q.put({'msg': 'finished', 'exit': p.returncode})

    def run(self, cmd):
        q = Queue()  
        proc = Popen(cmd, bufsize=0, stdout=PIPE, stderr=PIPE)
        out = Thread(target=self._readio, args=[proc, q])
        self.threads.append(out)
        out.start()
        self.procs.append(proc)
        return job(proc, q)

