FROM node:latest

WORKDIR /app

COPY package*.json ./

RUN npm install

COPY . .

ARG VITE_API_BASE_URL

ENV VITE_API_BASE_URL=${VITE_API_BASE_URL}

RUN echo "VITE_API_BASE_URL=${VITE_API_BASE_URL}" > .env

RUN npm run build

ENV HOST=0.0.0.0

EXPOSE 3000

CMD ["sh", "-c", "VITE_API_BASE_URL=${VITE_API_BASE_URL} npm start"]
