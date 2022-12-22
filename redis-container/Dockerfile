FROM redis:latest
RUN mkdir /var/lib/redis/
RUN mkdir /var/log/redis/
RUN chmod 777 /var/lib/redis
RUN chmod 777 /var/log/redis/
RUN touch /var/log/redis/redis-server.log
COPY redis.conf /etc/redis/
EXPOSE 6379
CMD ["redis-server", "/etc/redis/redis.conf"]
