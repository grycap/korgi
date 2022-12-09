# Dockerfile to create a container for training Tensorflow Models
#
# 20221209 - Estibaliz Parcero

# https://hub.docker.com/r/tensorflow/tensorflow
FROM tensorflow/tensorflow:2.11.0-gpu

RUN apt-get update && apt-get install -y git

RUN python3 -m pip install --upgrade pip

RUN git clone https://github.com/tensorflow/models.git

RUN export PYTHONPATH=$PYTHONPATH:/models

RUN pip3 install --user -r models/official/requirements.txt
