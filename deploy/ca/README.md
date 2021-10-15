# Commands to generate certificate

## Generate a CA certificate
Edit the ca.conf for change in Root CA details.
```
openssl genrsa -out ca.key 2048
openssl req -x509 -new -nodes -key ca.key -days 100000 -out ca.crt -extensions v3_req  -extensions v3_ca -config ca.conf
```

## Create a server certificate.
Edit the server.conf for domain(sidecar-injector.default.svc) name change. If you are hosting your app in a domain change the `commonName` and the `DNS.1` under `[alt_names]`. If you are hosting your app in an IP change the `commonName` and the `IP.1` under `[alt_names]`.
```
openssl genrsa -out server.key 2048
openssl req -new -key server.key -out server.csr -config server.conf
```
## Sign the server certificate with the above CA
```
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 100000 -extensions v3_req -extfile server.conf
```

## Create configmap/secret with server.key and server.crt
```
kubectl create configmap auth-tls --from-file=/path/to/server/cert
```
OR
```
kubectl create secret tls auth-tls --cert=path/to/cert/server.crt  --key=path/to/key/server.key
```

