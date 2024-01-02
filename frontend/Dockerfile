FROM oven/bun:latest as builder

WORKDIR /opt/web
COPY package.json ./
#COPY yarn.lock ./

ENV PATH /opt/web/node_modules/.bin:$PATH

RUN bun install

COPY . /opt/web

RUN bun run build

EXPOSE 3000

FROM nginx:latest
COPY --from=builder /opt/web/dist /usr/share/nginx/html

RUN rm -v /etc/nginx/nginx.conf
ADD nginx.conf /etc/nginx/

CMD ["nginx", "-g", "daemon off;"]