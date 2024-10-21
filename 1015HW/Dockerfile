FROM golang:1.23

WORKDIR /TODO

COPY go.mod go.mod
COPY go.sum go.sum
COPY main.go main.go
COPY class02_HW.go class2_HW.go
COPY error.go error.go


RUN echo "Before build "
RUN ls 

RUN go build -v
RUN echo "After build " 
RUN ls -all

ENTRYPOINT [ "./gogo" ]