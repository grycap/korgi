# Dockerfile to try thing with Tensorflow in a Jupyter Notebook
#
# 20220119 - Estibaliz Parcero

# https://hub.docker.com/r/tensorflow/tensorflow
FROM tensorflow/tensorflow:2.6.0-gpu

RUN python3 -m pip install --upgrade pip

RUN pip install matplotlib
RUN pip install jupyter

#sudo docker run -it --net=host -p 8888:8888 --gpus all -v /hdd500:/data tf-notebooks jupyter notebook --ip 0.0.0.0 --port 8888 --allow-root
