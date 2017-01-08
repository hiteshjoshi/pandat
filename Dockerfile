FROM smebberson/alpine-redis
MAINTAINER Hitesh Joshi <hi@hitesh.io>

# Add the files
ADD dist/pandat-linux-amd64 /app/pandat

ENV REDIS=0.0.0.0:6379

# Expose the ports for redis
EXPOSE 9090

CMD ["/app/pandat"]