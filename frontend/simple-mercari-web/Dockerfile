FROM node:16-alpine

WORKDIR /app
COPY package.json /app
COPY package-lock.json /app

RUN npm install

COPY . /app

RUN addgroup -g 1001 build && adduser -D -G build -u 1001 build
RUN chown -R build:build /app

ENV HOST 0.0.0.0
ENV PORT 3000
ARG API_URL

USER 1001

CMD ["npm", "start"]
