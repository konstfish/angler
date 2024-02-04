FROM nginx:latest

RUN apt update && apt install -y curl

RUN rm -v /etc/nginx/nginx.conf

ADD nginx.conf /etc/nginx/

CMD ["nginx", "-g", "daemon off;"]