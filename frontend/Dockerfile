FROM node:latest

WORKDIR /app

COPY package*.json ./

RUN npm install

COPY . .

COPY .env* ./

RUN npm run build

ENV HOST=0.0.0.0

EXPOSE 3000

CMD ["npm", "start"]
