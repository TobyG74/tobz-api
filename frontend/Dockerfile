FROM node:20-alpine AS build
WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

ARG VITE_API_BASE=http://localhost:8080/api/v1
ARG VITE_CAPTCHA_ENABLED=false
ARG VITE_TURNSTILE_SITEKEY=1x00000000000000000000AA
ENV VITE_API_BASE=$VITE_API_BASE \
    VITE_CAPTCHA_ENABLED=$VITE_CAPTCHA_ENABLED \
    VITE_TURNSTILE_SITEKEY=$VITE_TURNSTILE_SITEKEY

COPY . .
RUN npm run build

FROM nginx:alpine AS serve
COPY --from=build /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
