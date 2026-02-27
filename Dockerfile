FROM node:20-alpine

WORKDIR /app

# Install only production deps
COPY server/package.json ./package.json
RUN npm install --omit=dev

# Copy application source
COPY server/. .

EXPOSE 8080
CMD ["npm", "start"]
