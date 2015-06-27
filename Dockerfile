FROM scratch
MAINTAINER Atheatos <atheatos.engr@gmail.com>
COPY turmoil /
COPY params.ini /
ENTRYPOINT ["/turmoil","-config=params.ini","-logtostderr=true"]
