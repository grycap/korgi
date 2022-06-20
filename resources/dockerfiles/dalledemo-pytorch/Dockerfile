# Dockerfile: Test DALL-E in a Jupyter Notebook
# https://github.com/openai/dall-e
#
# 20220408 - Estibaliz Parcero

# https://hub.docker.com/r/pytorch/pytorch/
FROM pytorch/pytorch:1.10.0-cuda11.3-cudnn8-runtime

RUN apt-get -y update
RUN apt-get -y install git

RUN python3 -m pip install --upgrade pip

RUN pip install matplotlib
RUN pip install jupyter
RUN pip install DALL-E

# cd $PROJECT/recources/dockerfiles/dalledemo
# sudo docker build . -t pytorch-notebook
# sudo docker run -it -p 8888:8888 --gpus all -v /data/notebooks:/data pytorch-notebook jupyter notebook /data --ip 0.0.0.0 --allow-root
