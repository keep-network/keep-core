# Build environment
FROM node:11.15.0-alpine as build

WORKDIR /app

ENV PATH /app/node_modules/.bin:$PATH

RUN apk add --update --no-cache \
  git \
  build-base \
  python

COPY package.json /app/package.json
COPY package-lock.json /app/package-lock.json
COPY .env* /app/
COPY craco.config.js /app/

# Use `https://` instead of unauthenticated `git://` protocol
RUN git config --global url."https://".insteadOf git://

# Install from lockfile.
RUN npm ci --ignore-scripts

COPY src /app/src
COPY public /app/public

RUN NODE_OPTIONS=--max_old_space_size=4096 npm run build

# Production environment
FROM nginx:1.17-alpine

COPY --from=build /app/build /usr/share/nginx/html

RUN rm /etc/nginx/conf.d/default.conf
COPY nginx.conf /etc/nginx/conf.d

EXPOSE 80

## Set `daemon off` so the nginx is run in the foreground.
CMD ["nginx", "-g", "daemon off;"]
