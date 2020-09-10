resource "vault_mount" "ssh_host" {
  path = "ssh-host-signer"
  type = "ssh"
}

resource "vault_ssh_secret_backend_ca" "host_ca" {
  backend              = vault_mount.ssh_host.path
  generate_signing_key = true
}

resource "vault_auth_backend" "approle" {
  type = "approle"
}

resource "vault_approle_auth_backend_role" "host" {
  backend        = vault_auth_backend.approle.path
  role_name      = "host"
  token_policies = ["default", vault_policy.host.name]
}

resource "vault_approle_auth_backend_role_secret_id" "host" {
  backend   = vault_auth_backend.approle.path
  role_name = vault_approle_auth_backend_role.host.role_name
}

data "vault_policy_document" "host" {
  rule {
    path         = "ssh-host-signer/sign/host"
    capabilities = ["update"]
    description  = "Allow hosts to sign their own certs"
  }
  rule {
    path         = "ssh-host-signer/config/ca"
    capabilities = ["read"]
    description  = "Allow hosts read ca cert"
  }
}

resource "vault_policy" "host" {
  name   = "host"
  policy = data.vault_policy_document.host.hcl
}

resource "vault_ssh_secret_backend_role" "host" {
  name     = "host"
  backend  = vault_mount.ssh_host.path
  key_type = "ca"

  allow_host_certificates = true
  allowed_domains         = var.hosted_zone
  allow_subdomains        = true
}
