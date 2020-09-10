output "host_approle" {
  value = {
    role   = vault_approle_auth_backend_role.host.role_id
    secret = vault_approle_auth_backend_role_secret_id.host.secret_id
  }
}
