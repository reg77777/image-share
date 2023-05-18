import json
from pathlib import Path

import numpy as np
import torch
import torchvision
from PIL import Image
from torch.nn import functional as F
from torch.utils.data import DataLoader, Dataset
from torchvision import transforms
from torchvision.datasets.utils import download_url


def get_device(use_gpu):
    if use_gpu and torch.cuda.is_available():
        # これを有効にしないと、計算した勾配が毎回異なり、再現性が担保できない。
        torch.backends.cudnn.deterministic = True
        return torch.device("cuda")
    else:
        return torch.device("cpu")


# デバイスを選択する。
device = get_device(use_gpu=True)

model = torchvision.models.resnet50(pretrained=True).to(device)

transform = transforms.Compose(
    [
        transforms.Resize(256),  # (256, 256) で切り抜く。
        transforms.CenterCrop(224),  # 画像の中心に合わせて、(224, 224) で切り抜く
        transforms.ToTensor(),  # テンソルにする。
        transforms.Normalize(
            mean=[0.485, 0.456, 0.406], std=[0.229, 0.224, 0.225]
        ),  # 標準化する。
    ]
)
def get_classes():
    if not Path("data/imagenet_class_index.json").exists():
        # ファイルが存在しない場合はダウンロードする。
        download_url("https://git.io/JebAs", "data", "imagenet_class_index.json")

    # クラス一覧を読み込む。
    with open("data/imagenet_class_index.json") as f:
        data = json.load(f)
        class_names = [x["ja"] for x in data]

    return class_names


from pydantic import BaseModel

class_names = get_classes()

import base64

from flask import Flask, g, request, jsonify
from flask_cors import CORS
app=Flask(__name__)
CORS(app)

@app.route("/",methods=["POST"])
def post_item():
    print("posted")
    json_data=request.get_json()
    data=base64.b64decode(json_data["Data"])
    filename="./tmp."+json_data["Extension"]
    print(filename)
    with open(filename,mode="wb") as f:
        f.truncate(0)
        f.seek(0)
        f.write(data)
    img = Image.open(filename).convert("RGB")
    inputs = transform(img)
    inputs = inputs.unsqueeze(0).to(device)
    model.eval()
    outputs = model(inputs)
    batch_probs = F.softmax(outputs, dim=1)
    batch_probs, batch_indices = batch_probs.sort(dim=1, descending=True)
    ans=""
    for probs, indices in zip(batch_probs, batch_indices):
        ans=get_classes()[indices[0]]
    print(ans)
    return ans

@app.route("/test",methods=["GET"])
def test():
    return "test"

if __name__ == "__main__":
    app.run(debug=False,host="0.0.0.0",port=80)
