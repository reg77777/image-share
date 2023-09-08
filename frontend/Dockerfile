FROM node:16.15.0-alpine
RUN mkdir -p /home/node/workspace
RUN chown node:node /home/node/workspace
WORKDIR /home/node/workspace
COPY workspace/package*.json ./
RUN chown node:node /home/node/workspace/package*.json
RUN chown -R $(whoami) $(npm config get prefix)/lib/node_modules
USER node
RUN npm install
RUN npm install npm@9.6.5
