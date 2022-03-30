IMAGE_NAME="api-proxy"
docker build . -t $IMAGE_NAME
docker run -i -t --rm --env-file ./.env -p 9527:9527 $IMAGE_NAME
