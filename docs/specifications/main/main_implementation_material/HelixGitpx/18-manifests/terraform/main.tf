# deploy/terraform/helixgitpx-aws/main.tf
# Terraform module that provisions the AWS infrastructure required to
# run HelixGitpx (managed-dedicated or customer-managed deployment).
#
# Uses:
#   - VPC with public + private + database subnets across 3 AZs
#   - EKS (K8s 1.31) with Karpenter-managed node pools
#   - RDS for PostgreSQL (for customers who don't want CNPG on EKS)
#   - MSK for Kafka (alt: run Strimzi on EKS)
#   - S3 buckets for object store + backups
#   - KMS keys + IAM roles for IRSA (service accounts)
#
# Usage:
#   module "helixgitpx" {
#     source  = "github.com/vasic-digital/helixgitpx-terraform//aws?ref=v1"
#     prefix  = "helixgitpx-prod-eu"
#     region  = "eu-west-1"
#     cidr    = "10.42.0.0/16"
#     # ... see variables.tf
#   }

terraform {
  required_version = ">= 1.9"
  required_providers {
    aws        = { source = "hashicorp/aws",        version = "~> 5.60" }
    kubernetes = { source = "hashicorp/kubernetes", version = "~> 2.32" }
    helm       = { source = "hashicorp/helm",       version = "~> 2.14" }
  }
}

provider "aws" {
  region = var.region
  default_tags {
    tags = {
      Project     = "helixgitpx"
      Environment = var.environment
      Owner       = var.owner_email
      ManagedBy   = "terraform"
      CostCenter  = var.cost_center
    }
  }
}

# --------------------------------------------------------------
# Locals
# --------------------------------------------------------------
locals {
  azs = slice(data.aws_availability_zones.available.names, 0, 3)
  tags = {
    Environment = var.environment
    Project     = "helixgitpx"
  }
}

data "aws_availability_zones" "available" {
  state = "available"
}

# --------------------------------------------------------------
# KMS — one key per data class
# --------------------------------------------------------------
module "kms_app" {
  source  = "terraform-aws-modules/kms/aws"
  version = "~> 3.0"

  aliases                = ["alias/${var.prefix}-app"]
  description            = "HelixGitpx application data"
  enable_key_rotation    = true
  deletion_window_in_days = 30
  tags                   = local.tags
}

module "kms_backups" {
  source  = "terraform-aws-modules/kms/aws"
  version = "~> 3.0"

  aliases                = ["alias/${var.prefix}-backups"]
  description            = "HelixGitpx backups"
  enable_key_rotation    = true
  multi_region           = true
  deletion_window_in_days = 30
  tags                   = local.tags
}

# --------------------------------------------------------------
# VPC
# --------------------------------------------------------------
module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "~> 5.13"

  name = var.prefix
  cidr = var.cidr
  azs  = local.azs

  public_subnets   = [for i, az in local.azs : cidrsubnet(var.cidr, 4, i)]
  private_subnets  = [for i, az in local.azs : cidrsubnet(var.cidr, 4, i + 4)]
  database_subnets = [for i, az in local.azs : cidrsubnet(var.cidr, 4, i + 8)]

  enable_nat_gateway     = true
  single_nat_gateway     = false         # one per AZ
  one_nat_gateway_per_az = true
  enable_dns_hostnames   = true
  enable_dns_support     = true

  # Flow logs → CloudWatch + S3
  enable_flow_log                      = true
  create_flow_log_cloudwatch_iam_role  = true
  create_flow_log_cloudwatch_log_group = true

  public_subnet_tags  = { "kubernetes.io/role/elb"            = 1 }
  private_subnet_tags = { "kubernetes.io/role/internal-elb"   = 1 }

  tags = local.tags
}

# --------------------------------------------------------------
# EKS
# --------------------------------------------------------------
module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "~> 20.24"

  cluster_name    = var.prefix
  cluster_version = "1.31"
  subnet_ids      = module.vpc.private_subnets
  vpc_id          = module.vpc.vpc_id

  cluster_endpoint_public_access  = var.eks_public_endpoint
  cluster_endpoint_private_access = true

  cluster_addons = {
    coredns    = { most_recent = true }
    kube-proxy = { most_recent = true }
    vpc-cni    = { most_recent = true, before_compute = true }
    aws-ebs-csi-driver = { most_recent = true }
  }

  # Encryption of etcd
  cluster_encryption_config = {
    provider_key_arn = module.kms_app.key_arn
    resources        = ["secrets"]
  }

  # OIDC → IRSA
  enable_irsa = true

  # Small system pool; apps scale via Karpenter (installed later).
  eks_managed_node_groups = {
    system = {
      name            = "system"
      instance_types  = ["m7i.large"]
      min_size        = 3
      max_size        = 6
      desired_size    = 3
      capacity_type   = "ON_DEMAND"
      ami_type        = "BOTTLEROCKET_x86_64"
      labels          = { "workload-class" = "system" }
      taints          = [{ key = "dedicated", value = "system", effect = "NO_SCHEDULE" }]
    }
  }

  tags = local.tags
}

# --------------------------------------------------------------
# Karpenter for app / data / gpu node pools
# --------------------------------------------------------------
module "karpenter" {
  source  = "terraform-aws-modules/eks/aws//modules/karpenter"
  version = "~> 20.24"

  cluster_name          = module.eks.cluster_name
  enable_v1_permissions = true
  irsa_use_name_prefix  = false
  tags                  = local.tags
}

# NodeClasses + NodePools are deployed via the helixgitpx-platform chart;
# see deploy/helm/helixgitpx/values-*.yaml.

# --------------------------------------------------------------
# RDS — PostgreSQL (used if `var.use_cnpg == false`)
# --------------------------------------------------------------
module "rds" {
  count = var.use_cnpg ? 0 : 1

  source  = "terraform-aws-modules/rds-aurora/aws"
  version = "~> 9.10"

  name           = "${var.prefix}-pg"
  engine         = "aurora-postgresql"
  engine_version = "16.3"

  database_name    = "helixgitpx"
  master_username  = "helixgitpx_dba"
  manage_master_user_password = true

  instance_class   = "db.r7g.2xlarge"
  instances        = { 1 = {}, 2 = {}, 3 = {} }

  vpc_id               = module.vpc.vpc_id
  db_subnet_group_name = module.vpc.database_subnet_group_name
  security_group_rules = {
    from_eks = {
      source_security_group_id = module.eks.node_security_group_id
    }
  }

  storage_encrypted   = true
  kms_key_id          = module.kms_app.key_arn
  backup_retention_period = 35
  preferred_backup_window = "02:00-03:00"
  deletion_protection     = true

  monitoring_interval = 60
  performance_insights_enabled = true
  performance_insights_retention_period = 31

  tags = local.tags
}

# --------------------------------------------------------------
# MSK — Kafka (used if `var.use_strimzi == false`)
# --------------------------------------------------------------
module "msk" {
  count = var.use_strimzi ? 0 : 1

  source  = "terraform-aws-modules/msk-kafka-cluster/aws"
  version = "~> 2.7"

  name                   = "${var.prefix}-kafka"
  kafka_version          = "3.7.x"
  number_of_broker_nodes = 6

  broker_node_client_subnets  = module.vpc.private_subnets
  broker_node_instance_type   = "kafka.m7g.xlarge"
  broker_node_security_groups = [module.eks.node_security_group_id]
  broker_node_storage_info = {
    ebs_storage_info = { volume_size = 1000 }
  }

  encryption_in_transit_client_broker = "TLS"
  encryption_in_transit_in_cluster    = true
  encryption_at_rest_kms_key_arn      = module.kms_app.key_arn

  enhanced_monitoring    = "PER_TOPIC_PER_BROKER"
  open_monitoring_prometheus_jmx_exporter = true
  open_monitoring_prometheus_node_exporter = true

  configuration_info = {
    arn      = aws_msk_configuration.this[0].arn
    revision = aws_msk_configuration.this[0].latest_revision
  }

  tags = local.tags
}

resource "aws_msk_configuration" "this" {
  count = var.use_strimzi ? 0 : 1

  name              = "${var.prefix}-msk"
  kafka_versions    = ["3.7.x"]
  server_properties = <<-EOT
    auto.create.topics.enable=false
    default.replication.factor=3
    min.insync.replicas=2
    num.partitions=12
    log.retention.hours=168
    delete.topic.enable=true
    unclean.leader.election.enable=false
  EOT
}

# --------------------------------------------------------------
# S3 buckets
# --------------------------------------------------------------
resource "aws_s3_bucket" "object_store" {
  bucket        = "${var.prefix}-objects"
  force_destroy = false
  tags          = merge(local.tags, { DataClass = "confidential" })
}

resource "aws_s3_bucket_server_side_encryption_configuration" "object_store" {
  bucket = aws_s3_bucket.object_store.id
  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm     = "aws:kms"
      kms_master_key_id = module.kms_app.key_arn
    }
    bucket_key_enabled = true
  }
}

resource "aws_s3_bucket_versioning" "object_store" {
  bucket = aws_s3_bucket.object_store.id
  versioning_configuration { status = "Enabled" }
}

resource "aws_s3_bucket_object_lock_configuration" "object_store" {
  bucket              = aws_s3_bucket.object_store.id
  object_lock_enabled = "Enabled"

  rule {
    default_retention {
      mode = "COMPLIANCE"
      days = 90
    }
  }
}

resource "aws_s3_bucket" "backups" {
  bucket        = "${var.prefix}-backups"
  force_destroy = false
  tags          = merge(local.tags, { DataClass = "confidential" })
}

resource "aws_s3_bucket_server_side_encryption_configuration" "backups" {
  bucket = aws_s3_bucket.backups.id
  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm     = "aws:kms"
      kms_master_key_id = module.kms_backups.key_arn
    }
    bucket_key_enabled = true
  }
}

resource "aws_s3_bucket_versioning" "backups" {
  bucket = aws_s3_bucket.backups.id
  versioning_configuration { status = "Enabled" }
}

# --------------------------------------------------------------
# IRSA roles for HelixGitpx workloads
# --------------------------------------------------------------
module "irsa_object_store" {
  source  = "terraform-aws-modules/iam/aws//modules/iam-role-for-service-accounts-eks"
  version = "~> 5.44"

  role_name = "${var.prefix}-s3-access"

  role_policy_arns = {
    S3ReadWrite = aws_iam_policy.s3_rw.arn
  }

  oidc_providers = {
    eks = {
      provider_arn               = module.eks.oidc_provider_arn
      namespace_service_accounts = [
        "helixgitpx:repo-service",
        "helixgitpx:git-ingress",
        "helixgitpx:release-service",
      ]
    }
  }
}

resource "aws_iam_policy" "s3_rw" {
  name   = "${var.prefix}-s3-rw"
  policy = data.aws_iam_policy_document.s3_rw.json
}

data "aws_iam_policy_document" "s3_rw" {
  statement {
    actions   = ["s3:GetObject", "s3:PutObject", "s3:DeleteObject", "s3:ListBucket"]
    resources = [
      aws_s3_bucket.object_store.arn,
      "${aws_s3_bucket.object_store.arn}/*",
    ]
  }
  statement {
    actions   = ["kms:Decrypt", "kms:GenerateDataKey"]
    resources = [module.kms_app.key_arn]
  }
}

# --------------------------------------------------------------
# Outputs
# --------------------------------------------------------------
output "eks_cluster_endpoint" { value = module.eks.cluster_endpoint }
output "eks_cluster_name"     { value = module.eks.cluster_name     }
output "vpc_id"                { value = module.vpc.vpc_id           }
output "private_subnet_ids"    { value = module.vpc.private_subnets  }
output "object_store_bucket"   { value = aws_s3_bucket.object_store.bucket }
output "backups_bucket"        { value = aws_s3_bucket.backups.bucket }
output "kms_app_arn"           { value = module.kms_app.key_arn      }
output "kms_backups_arn"       { value = module.kms_backups.key_arn  }
output "msk_bootstrap_brokers_tls" {
  value = var.use_strimzi ? null : module.msk[0].bootstrap_brokers_tls
}
output "rds_writer_endpoint" {
  value = var.use_cnpg ? null : module.rds[0].cluster_endpoint
}
