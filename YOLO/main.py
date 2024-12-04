import time
from os import environ

import cv2
import torch
import json
import asyncio
from collections import defaultdict
from ultralytics import YOLO
from nats_connector import NATSConnector


class ObjectTracker:
    def __init__(self):
        self.track_history = defaultdict(lambda: [])

    def update_track(self, track_id, x_center, y_center):
        track = self.track_history[track_id]
        track.append((x_center, y_center))
        if len(track) > 30:
            track.pop(0)
        return track[::10]
    

class YoloModel:
    def __init__(self, model_path):
        self.model = YOLO(model_path).cuda()

    def detect_objects(self, frame):
        return self.model.track(frame, persist=True, stream=True, conf=0.7)
    

class VideoCapture:
    def __init__(self, video_path):
        self.cap = cv2.VideoCapture(video_path, cv2.CAP_FFMPEG)
        if not self.cap.isOpened():
            raise Exception("Failed to open video stream")

    def get_frame(self):
        return self.cap.read()

    def release(self):
        self.cap.release()


class ObjectTrackingApp:
    """Главная логика нейросети"""
    def __init__(self, video_path, model_path, nats_url, show_video=True):
        self.video_capture = VideoCapture(video_path)
        self.yolo_model = YoloModel(model_path)
        self.object_tracker = ObjectTracker()
        self.frame_count = 0
        self.nats_client = NATSConnector(nats_url, "camera")
        self.show_video = show_video

    async def process_frames(self):
        await self.nats_client.connect()  # Подключаемся к NATS

        while True:
            success, frame = self.video_capture.get_frame()
            if not success:
                break

            self.frame_count += 1
            results = self.yolo_model.detect_objects(frame)
            frame_data = {"frame_id": self.frame_count, "objects": []}

            for result in results:
                annotated_frame = result.plot()

                boxes_xyxyn = result.boxes.xyxyn.tolist()
                boxes_xywhn = result.boxes.xywhn.tolist()
                track_ids = result.boxes.id.int().tolist()
                class_ids = result.boxes.cls.int().tolist()
                names = result.names
                confidences = result.boxes.conf.tolist()

                for box_xyxyn, box_xywhn, track_id, class_id, confidence in zip(boxes_xyxyn, boxes_xywhn, track_ids, class_ids, confidences):
                    x_center, y_center, width, height = box_xywhn
                    track = self.object_tracker.update_track(track_id, x_center, y_center)
                    frame_data["objects"].append({
                        "id": track_id,
                        "class": names[class_id],
                        "bbox": box_xyxyn,
                        # "track_history": [(x, y) for x, y in track]
                    })

            # Публикуем frame_data на NATS
            await self.nats_client.publish_frame_data(frame_data)

            if self.show_video:
                # Отображение кадров с аннотацией
                cv2.imshow("YOLO Tracking", annotated_frame)

                if cv2.waitKey(1) & 0xFF == ord("q"):
                    break

        await self.nats_client.close()  # Закрываем соединение с NATS
        self.video_capture.release()
        if self.show_video:
            cv2.destroyAllWindows()


def read_config(config_path):
    with open(config_path, "r") as file:
        return json.load(file)


if __name__ == "__main__":
    cfg_path = environ.get("cfg", "config.json")
    app = ObjectTrackingApp(**read_config(cfg_path))
    asyncio.run(app.process_frames())
