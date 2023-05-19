FROM node:19.7-alpine as BUILD_IMAGE

ENV NODE_ENV development

WORKDIR /app

COPY . .

RUN npm install

RUN npm run build

FROM node:19.7-alpine 

WORKDIR /app

# copy from build image
COPY --from=BUILD_IMAGE app/build/ ./build
COPY --from=BUILD_IMAGE app/package.json .

EXPOSE 3000

CMD ["node", "build/index.js"]



