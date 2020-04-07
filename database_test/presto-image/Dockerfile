FROM ubuntu:18.04
RUN apt-get update -y && apt-get upgrade -y
RUN apt-get install wget -y
RUN apt-get install python -y
RUN apt-get install openjdk-8-jre-headless -y
RUN cd /opt && \
wget https://repo1.maven.org/maven2/com/facebook/presto/presto-server/0.232/presto-server-0.232.tar.gz && \
tar -xzvf presto-server-0.232.tar.gz && \
mkdir /opt/presto-server-0.232/etc && \
mkdir /opt/presto-server-0.232/etc/catalog
RUN mkdir /opt/presto && mkdir /opt/presto/data

RUN cd /opt/presto-server-0.232/etc && \
echo 'coordinator=true' > config.properties && \
echo 'node-scheduler.include-coordinator=true' >> config.properties && \
echo 'http-server.http.port=8080' >> config.properties && \
echo 'query.max-memory=5GB' >> config.properties && \
echo 'query.max-memory-per-node=1GB' >> config.properties && \
echo 'query.max-total-memory-per-node=2GB' >> config.properties && \
echo 'discovery-server.enabled=true' >> config.properties && \
echo 'discovery.uri=http://localhost:8080' >> config.properties

RUN cd /opt/presto-server-0.232/etc && \
echo '-server' > jvm.config && \
echo '-Xmx16G' >> jvm.config && \
echo '-XX:+UseG1GC' >> jvm.config && \
echo '-XX:G1HeapRegionSize=32M' >> jvm.config && \
echo '-XX:+UseGCOverheadLimit' >> jvm.config && \
echo '-XX:+ExplicitGCInvokesConcurrent' >> jvm.config && \
echo '-XX:+HeapDumpOnOutOfMemoryError' >> jvm.config && \
echo '-XX:+ExitOnOutOfMemoryError' >> jvm.config

RUN cd /opt/presto-server-0.232/etc && \
echo 'node.environment=production' > node.properties && \
echo 'node.id=ffffffff-ffff-ffff-ffff-ffffffffffff' >> node.properties && \
echo 'node.data-dir=/opt/presto/data' >> node.properties

COPY catalog /opt/presto-server-0.232/etc/catalog

WORKDIR /opt/presto-server-0.232
CMD ["./bin/launcher", "run"]