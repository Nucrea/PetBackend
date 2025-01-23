FROM node:20.18.2-alpine

WORKDIR /app
COPY . .
RUN npm install
RUN npm run build

EXPOSE 3000

CMD ["npm", "start"]