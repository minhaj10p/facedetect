FROM ubuntu
  
#RUN wget https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
WORKDIR /go/src/github.com/minhaj10p/facerecog
COPY . .
#RUN dep ensure -v
#RUN go install github.com/minhaj10p/facerecog/...

RUN apt-get update -y
RUN apt-get -y install python3-pip
RUN pip3 install face_recognition -y

CMD ["go", "run", "main.go"]
