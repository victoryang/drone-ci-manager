FROM {{ .FromImage }}
MAINTAINER maintainer@xueqiu.com

ADD *.zip /data/deploy/
ADD deploy_1_stop.sh /data/deploy/deploy_1_stop.sh
ADD deploy_2_replace.sh /data/deploy/deploy_2_replace.sh
ADD deploy_3_start.sh /data/deploy/deploy_3_start.sh

{{ .BuildDependency }}

RUN /data/deploy/deploy_2_replace.sh
RUN chown -R app:app /data/deploy