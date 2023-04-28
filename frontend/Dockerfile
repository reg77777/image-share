FROM node:16.15.0-alpine
WORKDIR /home/node/ws
COPY ws/package*.json .
RUN chown -R node:node package*.json
RUN npm install
RUN mkdir -p node_modules/.cache
RUN chown -R node node_modules/.cache
USER node
