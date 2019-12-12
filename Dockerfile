FROM node:10-alpine

COPY app.js package.json /app/
COPY public /app/public/

WORKDIR /app

RUN npm i

CMD ["npm","start"]
