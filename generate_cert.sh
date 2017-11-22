rm -rf cert/*

openssl genrsa -out ./cert/private.key 2048
openssl req -x509 -days 365 -key cert/private.key -out cert/domain.crt -new -subj "/CN=localhost"
openssl x509 -text -noout -in ./cert/domain.crt