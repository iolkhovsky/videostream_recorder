from fastapi import FastAPI
from pydantic import BaseModel
import cv2
import numpy as np
import base64
import io


app = FastAPI()


def base64str_to_img(base64str):
    base64_img_bytes = base64str.encode('utf-8')
    base64bytes = base64.b64decode(base64_img_bytes)
    bytes_io = io.BytesIO(base64bytes)
    encoded = np.frombuffer(buffer=bytes_io.read(), dtype=np.uint8)
    return cv2.imdecode(encoded, cv2.IMREAD_COLOR)


class Item(BaseModel):
    request_id: int
    encoded_img: str


@app.put("/frame")
async def handle_frame(imgdata: Item):
    img = base64str_to_img(imgdata.encoded_img)
    cv2.imwrite("Received.jpg", img)
    return {"encoded_img_size": len(imgdata.encoded_img)}
