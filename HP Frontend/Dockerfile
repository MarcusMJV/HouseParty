# 1. Build the Vue app
FROM node:20-alpine AS builder

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .
RUN npm run build

# 2. Serve the built files using nginx
FROM nginx:stable-alpine

# Copy built files to nginx's public folder
COPY --from=builder /app/dist /usr/share/nginx/html

# Optional: custom nginx config if you're using Vue Router history mode
# COPY nginx.conf /etc/nginx/nginx.conf

# Expose port 80 (matches what nginx uses)
EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
