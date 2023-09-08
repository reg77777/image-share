FROM nvidia/cuda:11.6.0-devel-ubuntu20.04
ENV LIBRARY_PATH /usr/local/cuda/lib64/stubs
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update
RUN apt-get upgrade -y
RUN apt-get install -y tzdata
ENV TZ=Asia/Tokyo
RUN apt-get install -y python3 python3-pip
RUN pip3 install torch torchvision
RUN pip3 install flask flask_cors pydantic
RUN apt install -y wget tar
RUN groupadd -g 1000 group
RUN useradd -m -u 1000 -g 1000 user
USER user
WORKDIR /home/user/workspace
