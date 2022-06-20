# Dockerfile: Test DALL-E in a Jupyter Notebook
# https://github.com/borisdayma/dalle-mini
#
# 20220513 - Estibaliz Parcero
# FROM nvidia/cuda:11.6.0-devel-ubuntu20.04
FROM tensorflow/tensorflow:2.7.1-gpu

# declare the image name
# ENV IMG_NAME=11.6.0-runtime-ubuntu20.04 \
    # declare what jaxlib tag to use
    # if a CI/CD system is expected to pass in these arguments
    # the dockerfile should be modified accordingly
    # JAXLIB_VERSION=0.3.0

# install pip
RUN python3 -m pip install --upgrade pip
# install git
RUN apt-get -y update
RUN apt-get -y install git

# install dependencies via pip
RUN pip install numpy scipy six wheel jaxlib
RUN pip install jupyter

# cd $PROJECT/recources/dockerfiles/dalledemo-jax
# sudo docker build . -t dalledemo-jax
# sudo docker run -it -p 8888:8888 --gpus all -v /data/notebooks:/data dalledemo-jax jupyter notebook /data --ip 0.0.0.0 --allow-root 