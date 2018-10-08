FROM node:10-alpine

COPY app.js package.json /app/
COPY node_modules /app/node_modules/
COPY public /app/public/

WORKDIR /app

CMD ["npm","start"]
