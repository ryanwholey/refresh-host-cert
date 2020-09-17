docker run --rm -it --network=host -v $(pwd):/app -w /app refresh-host-cert \
  -approle-id=ff3525cc-cdb8-bda0-5f20-367b2000c4f3 \
  -approle-secret=86c7c633-9576-3e74-4dc4-135f44856c34 \
  -host-signer-path=/ssh-host-signer/sign/host \
  -vault-addr=http://localhost:8200 \
  -public-key-path=/app/id_rsa.pub \
  -cert-path=/app/id_rsa-cert.pub
