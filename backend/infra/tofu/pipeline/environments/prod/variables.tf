variable "app_name" {
  type    = string
  default = "ns-scoring-backend"
}

variable "region" {
  type = string
}

variable "environment" {
  type = string
}

variable "profile" {
  type    = string
  default = ""
}

variable "github_repo" {
  type    = string
  default = "d-ginty4/ns-scoring-backend"
}
