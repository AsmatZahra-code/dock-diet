# ==========================================
# Stage 1: Dependencies & Build Stage
# ==========================================
FROM node:22-alpine AS builder

WORKDIR /usr/app

# Pehle sirf package files copy karte hain taake caching achi ho
COPY package*.json ./

# Saari dependencies (devDependencies samait) install karo build ke liye
RUN npm ci

# Baaki ka source code copy karo aur build/compile karo agar zaroorat ho
COPY . .
# RUN npm run build  # (Agar aap TypeScript ya koi bundler use kar rahi hain)


# ==========================================
# Stage 2: Production Runtime Stage
# ==========================================
FROM node:22-alpine AS runner

WORKDIR /usr/app

# Security best practice: Non-root user par app chalana
USER node

# Sirf production dependencies install karo (dev dependencies chhor do)
COPY --chown=node:node package*.json ./
RUN npm ci --omit=dev && npm cache clean --force

# Builder stage se sirf compiled/zaroori files uthao (choti image ke liye)
COPY --chown=node:node --from=builder /usr/app/src ./src
COPY --chown=node:node --from=builder /usr/app/app.js ./app.js

EXPOSE 4000

ENV NODE_ENV=production
ENV PORT=4000

CMD ["node", "app.js"]