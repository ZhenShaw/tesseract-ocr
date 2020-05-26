#FROM golang:latest as builder
## 安装依赖
#RUN apt-get update -qq && apt-get install -y -qq libtesseract-dev libleptonica-dev

FROM zhenshaw/tesseract:golang as builder
WORKDIR /home
ADD . /home
RUN go build -mod=vendor -o app main.go

#=============================

# 运行容器
#FROM ubuntu:20.04
#RUN apt-get update -qq && apt-get install -y -qq libtesseract-dev libleptonica-dev

# 下载数据文件，并设置数据文件目录
#RUN apt-get install -y -qq  tesseract-ocr-eng
#ENV TESSDATA_PREFIX=/usr/share/tesseract-ocr/4.00/tessdata

FROM zhenshaw/tesseract
COPY --from=builder /home/app /bin/
CMD ["app"]
