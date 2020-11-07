import cv2
import logging


logger = logging.getLogger("webcam")
logger.setLevel(logging.DEBUG)
fh = logging.FileHandler('webcam.log')
fh.setLevel(logging.DEBUG)
ch = logging.StreamHandler()
ch.setLevel(logging.ERROR)
formatter = logging.Formatter('%(asctime)s - %(name)s - %(levelname)s - %(message)s')
fh.setFormatter(formatter)
ch.setFormatter(formatter)
logger.addHandler(fh)
logger.addHandler(ch)


class WebCamera:

    def __init__(self, cam_idx=0, target_res=None):
        self._idx = cam_idx
        self._cap = None
        if target_res is not None:
            assert type(target_res) == tuple and len(target_res) == 2
        self._target_res = target_res
        try:
            self._cap = cv2.VideoCapture(self._idx)
        except Exception as e:
            logger.error("Got an exception during webcamera creation: " + str(e))
        if self._cap is None:
            logger.error("VideoCapture object is invalid")
        logger.warning("Init completed")

    def __del__(self):
        if self._cap.isOpened():
            self._cap.release()

    def __call__(self, *args, **kwargs):
        return self.capture_frame(*args, **kwargs)

    def __str__(self):
        return "WebCamera#" + str(self._idx)

    def capture_frame(self, color_format="bgr", data_format="raw"):
        capture_res = False
        if self._cap and self._cap.isOpened():
            capture_res, captured_frame = self._cap.read()
            if capture_res:
                if self._target_res:
                    xsz, ysz = self._target_res
                    captured_frame = cv2.resize(captured_frame, (xsz, ysz))
                if color_format != "bgr":
                    if color_format == "rgb":
                        captured_frame = cv2.cvtColor(captured_frame, cv2.COLOR_BGR2RGB)
                    else:
                        logger.warning(f"Invalid format of color encoding: {color_format}")
                        return None
                if data_format != "raw":
                    if data_format == "encoded":
                        _, captured_frame = cv2.imencode(".jpg", captured_frame)
                    else:
                        logger.warning(f"Invalid format of image data: {data_format}")
                        return None
                return captured_frame
        cap_status = "ok" if self._cap else "doesn't exist"
        cap_opened = "yes" if self._cap.isOpened() else "no"
        frame_captured = "yes" if capture_res else "no"
        logger.warning(f"Error while capturing frame. VideoCapture: {cap_status}" +
                       f" Opened: {cap_opened}, Frame is captured: {frame_captured}")
        return None


if __name__ == "__main__":
    camera = WebCamera(0)
    while True:
        frame = camera(color_format="bgr", data_format="raw")
        if frame is not None:
            cv2.imshow(str(camera), frame)
        if cv2.waitKey(1) & 0xFF == ord('q'):
            break
