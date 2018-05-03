FROM ubuntu:latest

RUN apt-get update
RUN apt-get install -y nginx

RUN rm -v /etc/nginx/nginx.conf
ADD docker/nginx.conf /etc/nginx/

COPY build/ /var/www/

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]