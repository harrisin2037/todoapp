FROM node:latest

WORKDIR /app

COPY package*.json ./

RUN npm install

COPY . .

ARG API_BASE_URL

ENV API_BASE_URL=${API_BASE_URL}

ARG VITE_API_BASE_URL

ENV VITE_API_BASE_URL=${VITE_API_BASE_URL}

RUN npm run build

ENV HOST=0.0.0.0

EXPOSE 3000

CMD ["npm", "start"]
