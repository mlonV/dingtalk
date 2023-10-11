FROM centos:latest
COPY config.yaml.sample /data/config.yaml
COPY dingtalk /data/dingtalk
WORKDIR /data
RUN mkdir /data/logs
CMD [ "./dingtalk" ]
