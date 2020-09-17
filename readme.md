# refresh-host-cert

A simple app that requests a signed host certificate from Vault.

docker run --rm -it --network=host -v /etc/ssh/:/etc/ssh/ -w /app refresh-host-cert \
  -approle-id=123 \
  -approle-secret=456 \
  -host-signer-path=/ssh-host-signer/sign/host \
  -vault-addr=https://vault.com \
  -public-key-path=/etc/ssh/ssh_host_rsa_key.pub \
  -cert-path=/etc/ssh/ssh_host_rsa_key-cert.pub
