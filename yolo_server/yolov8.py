import torch
from flask import Flask, request, jsonify
import io
import time
import cv2
from PIL import Image
import numpy as np
from torch.quantization import quantize_dynamic
from ultralytics import YOLO
import logging

logging.getLogger('ultralytics').setLevel(logging.WARNING)


# Model
global device, model
max_det = 1000

device = torch.device('cuda' if torch.cuda.is_available() else 'cpu')
#model = torch.hub.load("ultralytics/yolov8", "custom", path="yolov8n.engine", device='cuda:0')
model = YOLO("yolov8n.engine")
warmup_image = np.random.randint(0, 255, (416, 416, 3), dtype=np.uint8)
model.predict(warmup_image, imgsz=416, save=False)
print("warmed up")



def detect_image(image):
    start = time.time()
    results = model(image, imgsz=416, save=False)
    detection_time = time.time() - start
    #print(f"Detect Time: {detection_time}")
    for r in results[1:]:
        print(r.boxes)
    return [], 0
    return results.xyxy[0].tolist(), detection_time

app = Flask(__name__)

@app.route('/upload', methods=['POST'])
def upload_file():
    file = request.files['file']
    start_total = time.time()

    start = time.time()
    image = Image.open(io.BytesIO(file.read()))
    #print(image.size)
    decode_time = time.time() - start
    #print(f"Read Time: {decode_time}")

    detections, detect_time = detect_image(image)

    total_time = time.time() - start_total
    #print(f"Total Time: {total_time}")

    return jsonify(detections)

if __name__ == '__main__':
    app.run(debug=False)
