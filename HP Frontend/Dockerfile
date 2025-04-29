# 1. Build the Vue app
FROM node:20-alpine AS builder

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .

# Double-check: ensure build completes
RUN npm run build

# 2. Serve the built files using nginx
FROM nginx:stable-alpine

# Copy built static files
COPY --from=builder /app/dist /usr/share/nginx/html

# Default nginx.conf is fine unless you need custom history mode
EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
