FROM golang:1.11-stretch

# ENV SETUP  
RUN apt-get update -y
RUN apt-get -y install python3-pip
RUN apt-get -y install cmake
RUN pip3 install imutils
RUN pip3 install opencv-python
RUN pip3 install face_recognition 
RUN apt install -y libsm6 
RUN apt install -y libxext6
RUN apt install -y libxrender1


# APP SETUP
WORKDIR /go/src/github.com/minhaj10p/facedetect
COPY . .
EXPOSE 8080
CMD ["go", "run", "main.go"]
