# deploy/terraform/helixgitpx-aws/variables.tf
variable "prefix" {
  description = "Name prefix for all resources (e.g. 'helixgitpx-prod-eu')."
  type        = string
  validation {
    condition     = can(regex("^[a-z][a-z0-9-]{2,30}[a-z0-9]$", var.prefix))
    error_message = "prefix must be lowercase alphanumeric / hyphen, 4-32 chars."
  }
}

variable "region" {
  description = "AWS region (e.g. eu-west-1)."
  type        = string
}

variable "environment" {
  description = "Environment label (staging / prod-eu / prod-us / ...)."
  type        = string
}

variable "owner_email" {
  description = "Contact email used in resource tags."
  type        = string
}

variable "cost_center" {
  description = "Cost attribution tag."
  type        = string
  default     = "engineering"
}

variable "cidr" {
  description = "VPC CIDR range. /16 recommended."
  type        = string
  default     = "10.42.0.0/16"
  validation {
    condition     = can(cidrhost(var.cidr, 0))
    error_message = "cidr must be a valid IPv4 CIDR block."
  }
}

variable "eks_public_endpoint" {
  description = "Whether EKS API server is reachable over the public internet. Keep false for production unless you have a strong reason."
  type        = bool
  default     = false
}

variable "use_cnpg" {
  description = "If true, run Postgres in-cluster via CloudNativePG. If false, use RDS Aurora."
  type        = bool
  default     = true
}

variable "use_strimzi" {
  description = "If true, run Kafka in-cluster via Strimzi. If false, use MSK."
  type        = bool
  default     = true
}

variable "enable_fips" {
  description = "Enable FIPS-validated endpoints and enforce FIPS mode in nodes."
  type        = bool
  default     = false
}

variable "enable_wafv2" {
  description = "Attach an AWS WAFv2 ruleset in front of the Envoy Gateway ELB."
  type        = bool
  default     = true
}

# ----- Example tfvars -------------------------------------
#
# prefix       = "helixgitpx-prod-eu"
# region       = "eu-west-1"
# environment  = "prod-eu"
# owner_email  = "platform@helixgitpx.example.com"
# cost_center  = "platform-eu"
# cidr         = "10.42.0.0/16"
# use_cnpg     = true
# use_strimzi  = true
# enable_wafv2 = true
