import cv2
import os
import subprocess
import threading

fps = 24.0

reader = cv2.VideoCapture(0)
if not reader.isOpened():
    print("R")
    exit()

width = int(reader.get(cv2.CAP_PROP_FRAME_WIDTH))
height = int(reader.get(cv2.CAP_PROP_FRAME_HEIGHT))
size = (width, height)

writer = None
tslist = []
seq = 0
cnt = 0


while True:
    ok, frame = reader.read()
    if not ok:
        print(":(")
        exit()
    if writer is None:
        writer = cv2.VideoWriter("out_{}.mp4".format(seq % 5), cv2.VideoWriter_fourcc(*"mp4v"), fps, size)

    writer.write(frame)
    cnt += 1
    if cnt > 120:
        writer.release()
        writer = None
        cnt = 0

        def f(seq):
            global tslist

            out = "tmp.m3u8"
            subprocess.run([
                "ffmpeg", "-i", "out_{}.mp4".format(seq % 5), "-c:v", "h264", "-f", "hls", "-hls_list_size", "0", "-hls_time", "3",  "-hls_segment_filename", "seq{}_%d.ts".format(seq), out])

            with open(out) as f:
                prev_extinf = None
                for l in f:
                    if l.startswith('#EXTINF:'):
                        prev_extinf = l.strip()[:]
                        continue
                    elif prev_extinf:
                        tslist.append(prev_extinf)
                        tslist.append(l.strip())
                        if len(tslist) > 10:
                            tslist = tslist[-10:]
                    prev_extinf = None

            max_duration = 0
            for l in tslist:
                if l.startswith("#EXTINF:"):
                    max_duration = max(max_duration, int(eval(l[len("#EXTINF:"):-1])))

            f = open("playlist.m3u8", "w")
            f.write("\n".join([
                "#EXTM3U",
                "#EXT-X-VERSION:3",
                "#EXT-X-ALLOW-CACHE:NO",
                "#EXT-X-MEDIA-SEQUENCE:{}".format(seq - 5 if seq > 5 else 0),
                "#EXT-X-TARGETDURATION:{}".format(max_duration),
                "","",
            ]))
            f.write("\n".join(tslist + [""]))
            f.flush()
            f.close()

        threading.Thread(target=f, args=[seq]).start()
        seq += 1

