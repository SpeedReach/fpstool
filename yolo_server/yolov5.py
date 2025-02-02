import torch
from flask import Flask, request, jsonify
import io
import time
import cv2
from PIL import Image
import numpy as np
from torch.quantization import quantize_dynamic


# Model
global device, model
max_det = 1000

device = torch.device('cuda' if torch.cuda.is_available() else 'cpu')
model = torch.hub.load("ultralytics/yolov5", "custom", path="best.engine", device='cuda:0')
warmup_image = np.random.randint(0, 255, (416, 416, 3), dtype=np.uint8)
model(warmup_image, size=416)
print("warmed up")



def detect_image(image):
    start = time.time()
    results = model(image, size=416)
    detection_time = time.time() - start
    print(f"Detect Time: {detection_time}")
    print(results)
    return results.xyxy[0].tolist(), detection_time

app = Flask(__name__)

@app.route('/upload', methods=['POST'])
def upload_file():
    file = request.files['file']
    start_total = time.time()

    start = time.time()
    image = Image.open(io.BytesIO(file.read()))
    print(image.size)
    decode_time = time.time() - start
    print(f"Read Time: {decode_time}")

    detections, detect_time = detect_image(image)

    total_time = time.time() - start_total
    print(f"Total Time: {total_time}")

    return jsonify(detections)

if __name__ == '__main__':
    app.run(debug=False)
