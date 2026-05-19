FROM oven/bun:1

WORKDIR /app

COPY package.json bun.lockb* ./

RUN bun install --production

COPY . .

RUN bun run prisma:generate

EXPOSE 3000

CMD ["sh", "-c", "bun run prisma db push --skip-generate && bun run dev"]
