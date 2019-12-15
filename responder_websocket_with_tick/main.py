import responder
from websockets.exceptions import ConnectionClosed
from uuid import uuid4
import asyncio
import threading
import time
import janus

api = responder.API()
wset = dict()

@api.route('/ws', websocket=True)
async def ws(ws):
    await ws.accept()
    id = uuid4()
    queue = janus.Queue()
    wset[id] = queue.sync_q

    async def sender(ws, q):
        while True:
            m = await q.get()
            try:
                await ws.send_json(m)
            except ConnectionClosed:
                break

    await sender(ws, queue.async_q)
    await ws.close()
    wset.pop(id)


@api.route('/')
def index(req, res):
    res.content = api.template('index.html')

def enqueue(wset):
    while True:
        for _, q in wset.items():
            q.put({
                'msg': 'HELLO FROM SERVER'
            })
        time.sleep(5)

if __name__ == '__main__':
    threading.Thread(target=enqueue, args=[wset]).start()
    api.run() 
