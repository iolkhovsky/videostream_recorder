import argparse
import io
import requests
import time
import base64
import json

from webcamera import WebCamera


def parse_cmd_args():
    parser = argparse.ArgumentParser()
    parser.add_argument("--ip", type=str, default="localhost", help="Target server's IP")
    parser.add_argument("--port", type=int, default=8000, help="Target server's port")
    parser.add_argument("--camera", type=int, default=0, help="Camera id")
    parser.add_argument("--fps", type=float, default=25, help="Stream FPS")
    return parser.parse_args()


def run_client(args):
    camera = WebCamera(args.camera)
    ip = args.ip
    port = str(args.port)
    url = f"http://{ip}:{port}/frame"
    frame_period = 1. / (args.fps if args.fps else 1.)
    headers = {'content-type': 'image/jpeg'}
    i = 0
    while True:
        frame = camera(color_format="bgr", data_format="encoded")
        io_buf = io.BytesIO(frame)
        io_buf.seek(0)
        base64str = base64.b64encode(io_buf.read()).decode("utf-8")
        payload = json.dumps({"request_id": i, "encoded_img": base64str})
        resp = requests.put(url, data=payload, headers=headers, timeout=8000)
        print("Server's response:", resp)
        time.sleep(frame_period)
        i += 1


if __name__ == "__main__":
    run_client(parse_cmd_args())
