FROM ubuntu:20.04
RUN apt-get update
RUN apt-get upgrade -y
RUN apt install -y wget tar
RUN wget https://dl.google.com/go/go1.17.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.17.linux-amd64.tar.gz
RUN groupadd -g 1000 group
RUN useradd -m -u 1000 -g 1000 user
USER user
WORKDIR /home/user/ws
ENV PATH $PATH:/usr/local/go/bin
#RUN echo "export PATH=$PATH:/usr/local/go/bin" >> /home/user/.bashrc
