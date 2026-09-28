# frontend/Dockerfile

# ==========================================
# BUILD
# ==========================================

FROM node:24-alpine AS builder

WORKDIR /app

RUN corepack enable

COPY package.json pnpm-lock.yaml ./

RUN pnpm install --frozen-lockfile

COPY . .

ARG VITE_WS_URL

ENV VITE_WS_URL=$VITE_WS_URL

RUN pnpm build


# ==========================================
# RUNTIME
# ==========================================

FROM nginx:alpine

COPY --from=builder \
    /app/dist \
    /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]