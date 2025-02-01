import asyncio
import json
from os import environ

import cv2
from nats_connector import NATSConnector
from ultralytics import YOLO


class YoloModel:
    def __init__(self, model_path):
        self.model = YOLO(model_path).cuda()

    def detect_objects(self, frame):
        return self.model(frame, conf=0.7)


class VideoCapture:
    def __init__(self, video_path):
        self.cap = cv2.VideoCapture(video_path, cv2.CAP_FFMPEG)
        self.cap.set(cv2.CAP_PROP_BUFFERSIZE, 1)
        if not self.cap.isOpened():
            raise Exception("Failed to open video stream")

    def get_frame(self):
        return self.cap.read()

    def release(self):
        self.cap.release()


class ObjectTrackingApp:
    """Главная логика нейросети"""

    def __init__(self, video_path, model_path, nats_url, show_video=True):
        self.yolo_model = YoloModel(model_path)
        self.video_capture: VideoCapture = None
        self.frame_count = 0
        self.nats_client = NATSConnector(nats_url, "camera")
        self.show_video = show_video
        self.video_path = video_path

    async def process_frames(self):
        await self.nats_client.connect()  # Подключаемся к NATS
        self.video_capture = VideoCapture(self.video_path)

        while True:
            success, frame = self.video_capture.get_frame()
            if not success:
                continue

            self.frame_count += 1
            results = self.yolo_model.detect_objects(frame)
            frame_data = {"frame_id": self.frame_count, "objects": []}

            for result in results:
                annotated_frame = result.plot()

                boxes_xyxyn = result.boxes.xyxyn.tolist()  # Координаты боксов
                boxes_xywhn = result.boxes.xywhn.tolist()  # Размеры боксов
                class_ids = result.boxes.cls.int().tolist()  # ID классов
                names = result.names  # Имена классов
                confidences = result.boxes.conf.tolist()  # Уверенность в детекции

                for box_xyxyn, box_xywhn, class_id, confidence in zip(
                    boxes_xyxyn, boxes_xywhn, class_ids, confidences
                ):
                    frame_data["objects"].append(
                        {
                            "class": names[class_id],
                            "bbox": box_xyxyn,
                            "confidence": confidence,
                        }
                    )

            # Публикуем frame_data на NATS
            await self.nats_client.publish_frame_data(frame_data)

            if self.show_video:
                # Отображение кадров с аннотацией
                cv2.imshow("YOLO Tracking", annotated_frame)

        #         if cv2.waitKey(1) & 0xFF == ord("q"):
        #             break

        # await self.nats_client.close()  # Закрываем соединение с NATS
        # self.video_capture.release()
        # if self.show_video:
        #     cv2.destroyAllWindows()


def read_config(config_path):
    with open(config_path, "r") as file:
        return json.load(file)


if __name__ == "__main__":
    cfg_path = environ.get("cfg", "config.json")
    app = ObjectTrackingApp(**read_config(cfg_path))
    asyncio.run(app.process_frames())
