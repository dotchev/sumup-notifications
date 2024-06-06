provider "tls" {}

// CA
resource "tls_private_key" "ca" {
  algorithm = "RSA"
}

resource "local_file" "ca_private_key" {
  content  = tls_private_key.ca.private_key_pem
  filename = "${path.module}/ca.key.pem"
}

resource "tls_self_signed_cert" "ca" {
  private_key_pem = tls_private_key.ca.private_key_pem

  subject {
    common_name  = "CA"
    organization = "Notifications Ltd"
  }

  is_ca_certificate = true

  validity_period_hours = 8760

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "cert_signing",
  ]
}

resource "local_file" "ca_cert" {
  content  = tls_self_signed_cert.ca.cert_pem
  filename = "${path.module}/ca.crt.pem"
}

// Server
resource "tls_private_key" "gateway" {
  algorithm = "RSA"
}

resource "local_file" "server_private_key" {
  content  = tls_private_key.gateway.private_key_pem
  filename = "${path.module}/gateway.key.pem"
}

resource "tls_cert_request" "gateway" {
  private_key_pem = tls_private_key.gateway.private_key_pem

  subject {
    common_name  = "gateway"
    organization = "Notifications Ltd"
  }

  dns_names = ["gateway", "localhost"]
}

resource "tls_locally_signed_cert" "gateway" {
  cert_request_pem = tls_cert_request.gateway.cert_request_pem

  ca_private_key_pem = tls_private_key.ca.private_key_pem
  ca_cert_pem        = tls_self_signed_cert.ca.cert_pem

  validity_period_hours = 8760

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "server_auth",
  ]
}

resource "local_file" "server_cert" {
  content  = tls_locally_signed_cert.gateway.cert_pem
  filename = "${path.module}/gateway.crt.pem"
}

// Client
resource "tls_private_key" "client" {
  algorithm = "RSA"
}

resource "local_file" "client_private_key" {
  content  = tls_private_key.client.private_key_pem
  filename = "${path.module}/client.key.pem"
}

resource "tls_cert_request" "client" {
  private_key_pem = tls_private_key.client.private_key_pem

  subject {
    common_name  = "client"
    organization = "Notifications Ltd"
  }
}

resource "tls_locally_signed_cert" "client" {
  cert_request_pem = tls_cert_request.client.cert_request_pem

  ca_private_key_pem = tls_private_key.ca.private_key_pem
  ca_cert_pem        = tls_self_signed_cert.ca.cert_pem

  validity_period_hours = 8760

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "client_auth",
  ]
}

resource "local_file" "client_cert" {
  content  = tls_locally_signed_cert.client.cert_pem
  filename = "${path.module}/client.crt.pem"
}
