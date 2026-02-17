# build
FROM node:18.19.1 AS build
WORKDIR /app
COPY package*.json ./
RUN npm config set strict-ssl false
RUN npm install
COPY . .
RUN npm run build

# runtime
FROM nginx:alpine
# install envsubst
RUN apk add --no-cache gettext
COPY --from=build /app/build /usr/share/nginx/html
COPY public/config.js /usr/share/nginx/html/config.js
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
EXPOSE 80
CMD ["/entrypoint.sh"]