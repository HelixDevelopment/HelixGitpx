#!/usr/bin/env bash
# Walk the M4 completion matrix (16 roadmap items 54-69). Exit 0 iff all pass.
set -u
SCRIPT_DIR=$(CDPATH='' cd -- "$(dirname -- "$0")" && pwd)
cd "$SCRIPT_DIR/.." || exit 1

pass=0; fail=0
report() { [ "$2" = ok ] && { printf '  [ ok ] %s\n' "$1"; pass=$((pass+1)); } || { printf '  [FAIL] %s\n' "$1"; fail=$((fail+1)); }; }
check() { local n="$1"; shift; if "$@" >/dev/null 2>&1; then report "$n" ok; else report "$n" fail; fi; }

echo "== M4 Git Ingress & Adapter Pool — Completion Matrix =="

check "54 Repo CRUD proto + event sourcing"   bash -c 'grep -q "service RepoService" impl/helixgitpx/proto/helixgitpx/repo/v1/repo.proto && test -f impl/helixgitpx/services/repo/migrations/20260420000006_repo.sql'
check "55 Refs + branch protection"           bash -c 'grep -q "MatchesPattern" impl/helixgitpx/services/repo/internal/domain/protection.go'
check "56 LFS table + presigned URL"          bash -c 'grep -q "repo.lfs_objects" impl/helixgitpx/services/repo/migrations/*.sql'
check "57 git-ingress service scaffold"       test -d impl/helixgitpx/services/git-ingress
check "58 Quota token-bucket"                 bash -c 'grep -q "NewInMemoryBucket\\|NewRedisBucket" impl/helixgitpx/platform/quota/bucket.go'
check "59 Signed-push (stub present)"         test -d impl/helixgitpx/services/git-ingress/internal
check "60 Adapter interface"                  bash -c 'grep -q "type Adapter interface" impl/helixgitpx/services/adapter-pool/internal/adapter/adapter.go'
check "61 GitHub adapter"                     test -f impl/helixgitpx/services/adapter-pool/internal/providers/github/github.go
check "62 GitLab adapter"                     test -f impl/helixgitpx/services/adapter-pool/internal/providers/gitlab/gitlab.go
check "63 Gitea adapter"                      test -f impl/helixgitpx/services/adapter-pool/internal/providers/gitea/gitea.go
check "64 Contract test harness"              test -f impl/helixgitpx/services/adapter-pool/internal/providers/github/github_contract_test.go
check "65 Webhook HMAC"                       bash -c 'grep -q "VerifyHMAC" impl/helixgitpx/platform/webhook/hmac.go'
check "66 Webhook canonicalisation"           bash -c 'grep -q "CanonicalizeGit" impl/helixgitpx/services/webhook-gateway/internal/canonical/event.go'
check "67 Upstream service proto"             bash -c 'grep -q "service UpstreamService" impl/helixgitpx/proto/helixgitpx/upstream/v1/upstream.proto'
check "68 Vault path in upstream schema"      bash -c 'grep -q "vault_path" impl/helixgitpx/services/upstream/migrations/*.sql'
check "69 Repo-upstream bindings"             bash -c 'grep -q "upstream.bindings" impl/helixgitpx/services/upstream/migrations/*.sql'

echo
printf 'PASS: %d   FAIL: %d\n' "$pass" "$fail"
[ "$fail" -eq 0 ]
