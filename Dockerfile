FROM node:10-alpine

COPY . /app
WORKDIR /app

CMD ["npm","start"]
