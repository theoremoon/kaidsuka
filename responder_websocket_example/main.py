import responder
import asyncio

api = responder.API()
sockets = {}


@api.route("/ws", websocket=True)
async def websocket(ws):
    await ws.accept()
    q = asyncio.Queue()

    async def timer(q):
        while True:
            await asyncio.sleep(5)
            q.put_nowait('tick')

    async def message(ws, q):
        while True:
            text = await ws.receive_text()
            if text == "close":
                break
            q.put_nowait(text)

    async def sender(ws, q):
        while True:
            m = await q.get()
            await ws.send_text(m)

    _, pending = await asyncio.wait([
        asyncio.create_task(timer(q)),
                   asyncio.create_task(message(ws, q)),
                   asyncio.create_task(sender(ws, q))], return_when=asyncio.FIRST_COMPLETED)
    for task in pending:
        task.cancel()

    await ws.close()

@api.route("/")
def index(req, res):
    res.content = api.template("index.html")

if __name__ == '__main__':
    api.run()
