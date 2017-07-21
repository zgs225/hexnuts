FROM golang:1.9

RUN go get github.com/Sirupsen/logrus

EXPOSE 5678
EXPOSE 5679

ENV HEXNUTS_PATH $GOPATH/src/github.com/zgs225/hexnuts
COPY . $HEXNUTS_PATH
RUN cd $HEXNUTS_PATH && go install -ldflags '-w'
ENTRYPOINT ['hexnuts', 'server']
