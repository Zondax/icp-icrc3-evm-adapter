# Based on https://nextjs.org/docs/deployment
# Install dependencies only when needed
FROM node:18-alpine AS base
RUN apk --no-cache --update add libc6-compat git npm python3 make g++ vips-dev

WORKDIR /app

COPY package.json yarn.lock ./
RUN yarn install --frozen-lockfile

COPY . .
RUN chmod u+x scripts/*

RUN yarn codegen && yarn build

# Production image, copy all the files and run next
FROM zondax/subql-node-ethereum:v3.5.3 AS runner

ENV NODE_ENV production

RUN addgroup -g 1001 -S nodejs
RUN adduser -S service -u 1001

USER service

WORKDIR /app

# You only need to copy next.config.js if you are NOT using the default configuration
# COPY --from=builder /app/next.config.js ./
COPY --from=base /app/dist ./dist
COPY --chown=service:nodejs --from=base /app/project.yaml ./project.yaml
COPY --from=base /app/node_modules ./node_modules
COPY --from=base /app/package.json ./package.json
COPY --from=base /app/schema.graphql ./schema.graphql
COPY --from=base /app/abis ./abis
COPY --from=base /app/scripts ./scripts


# Next.js collects completely anonymous telemetry data about general usage.
# Learn more here: https://nextjs.org/telemetry
# Uncomment the following line in case you want to disable telemetry.
# RUN npx next telemetry disable

# Set Entry point and command
ENTRYPOINT ["/sbin/tini", "--", "/usr/local/lib/node_modules/@subql/node-ethereum/bin/run"]
CMD ["-f","/app"]
