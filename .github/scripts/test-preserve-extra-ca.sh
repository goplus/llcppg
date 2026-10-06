#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
HELPER="$SCRIPT_DIR/preserve-extra-ca.sh"
INTEGRATION=false
if [ "$#" -gt 0 ]; then
  if [ "$#" -ne 1 ] || [ "$1" != "--integration" ]; then
    echo "Usage: $0 [--integration]" >&2
    exit 2
  fi
  INTEGRATION=true
fi

WORK_DIR="$(mktemp -d "${TMPDIR:-/tmp}/llcppg-ca-test.XXXXXX")"
SERVER_PID=""
PASSED=0

cleanup() {
  if [ -n "$SERVER_PID" ]; then
    kill "$SERVER_PID" 2>/dev/null || true
    wait "$SERVER_PID" 2>/dev/null || true
  fi
  rm -rf "$WORK_DIR"
}
trap cleanup EXIT

fail() {
  echo "FAIL: $*" >&2
  exit 1
}

pass() {
  PASSED=$((PASSED + 1))
  printf 'PASS: %s\n' "$1"
}

make_certificate() {
  local certificate_name="$1"
  local ca_constraint="$2"
  local common_name="${3:-$certificate_name}"
  cat > "$WORK_DIR/$certificate_name.conf" <<EOF
[req]
distinguished_name = subject
x509_extensions = extensions
prompt = no
[subject]
CN = $common_name
[extensions]
basicConstraints = critical,CA:$ca_constraint
EOF
  openssl req -x509 -newkey rsa:2048 -nodes -sha256 -days 1 \
    -config "$WORK_DIR/$certificate_name.conf" \
    -keyout "$WORK_DIR/$certificate_name.key" \
    -out "$WORK_DIR/$certificate_name.crt" > /dev/null 2>&1
}

make_certificate distribution-root TRUE
make_certificate extra-root-one TRUE
make_certificate extra-root-two TRUE
make_certificate leaf-certificate FALSE 'CA:TRUE'

source "$HELPER"

new_case() {
  CASE_DIR="$WORK_DIR/$1"
  SYSTEM_DIR="$CASE_DIR/system sources"
  LOCAL_DIR="$CASE_DIR/local sources"
  BUNDLE="$CASE_DIR/trust bundle.crt"
  mkdir -p "$SYSTEM_DIR" "$LOCAL_DIR"
  cp "$WORK_DIR/distribution-root.crt" "$SYSTEM_DIR/distribution-root.crt"
  cp "$WORK_DIR/distribution-root.crt" "$BUNDLE"
}

run_preservation() {
  bash -c 'source "$1"; shift; preserve_extra_ca "$@"' \
    bash "$HELPER" "${1:-}" "$BUNDLE" "$SYSTEM_DIR" "$LOCAL_DIR"
}

local_certificate_count() {
  find "$LOCAL_DIR" -type f -name '*.crt' | wc -l | tr -d ' '
}

assert_preserved() {
  local fingerprint
  fingerprint="$(ca_certificate_fingerprint "$WORK_DIR/$1.crt")"
  [ -f "$LOCAL_DIR/llcppg-extra-ca-$fingerprint.crt" ] || fail "$1 was not preserved"
}

new_case no-extra
run_preservation
[ "$(local_certificate_count)" -eq 0 ] || fail "Managed certificates were duplicated"
pass "No extra CA leaves local sources unchanged"

new_case missing-bundle
rm "$BUNDLE"
run_preservation
touch "$BUNDLE"
run_preservation
[ "$(local_certificate_count)" -eq 0 ] || fail "An absent or empty bundle created certificates"
pass "Missing and empty bundles do not block initial package installation"

new_case generic-extras
cat "$WORK_DIR/extra-root-one.crt" "$WORK_DIR/extra-root-two.crt" >> "$BUNDLE"
cp "$BUNDLE" "$CASE_DIR/original.crt"
run_preservation
[ "$(local_certificate_count)" -eq 2 ] || fail "Expected two generic extra CA certificates"
assert_preserved extra-root-one
assert_preserved extra-root-two
cmp "$BUNDLE" "$CASE_DIR/original.crt" || fail "The original bundle was modified"
run_preservation > "$CASE_DIR/repeated-output"
[ "$(local_certificate_count)" -eq 2 ] || fail "Repeated preservation created duplicates"
[ ! -s "$CASE_DIR/repeated-output" ] || fail "Registered certificates were preserved again"
pass "Provider-independent preservation is idempotent and leaves the bundle intact"

new_case duplicate-extra
cat "$WORK_DIR/extra-root-one.crt" "$WORK_DIR/extra-root-one.crt" >> "$BUNDLE"
run_preservation
[ "$(local_certificate_count)" -eq 1 ] || fail "Duplicate bundle entries created multiple sources"
pass "Repeated PEM entries are deduplicated by certificate fingerprint"

new_case registered-local
cp "$WORK_DIR/extra-root-one.crt" "$LOCAL_DIR/company-root.crt"
cat "$WORK_DIR/extra-root-one.crt" >> "$BUNDLE"
run_preservation
[ "$(local_certificate_count)" -eq 1 ] || fail "An existing local CA was duplicated"
cmp "$LOCAL_DIR/company-root.crt" "$WORK_DIR/extra-root-one.crt" || fail "An existing local CA was overwritten"
pass "Existing local CA sources are retained without duplication"

new_case managed-multiple
cat "$WORK_DIR/extra-root-one.crt" "$WORK_DIR/extra-root-two.crt" > "$SYSTEM_DIR/multiple.crt"
cat "$WORK_DIR/extra-root-one.crt" "$WORK_DIR/extra-root-two.crt" >> "$BUNDLE"
run_preservation
[ "$(local_certificate_count)" -eq 0 ] || fail "Multi-certificate managed sources were not recognized"
pass "All certificates in a managed source are recognized"

new_case managed-symlink
ln -s "$WORK_DIR/extra-root-one.crt" "$SYSTEM_DIR/linked-root.crt"
cat "$WORK_DIR/extra-root-one.crt" >> "$BUNDLE"
run_preservation
[ "$(local_certificate_count)" -eq 0 ] || fail "A managed symlink was not recognized"
pass "Managed certificate symlinks are recognized"

new_case disabled-system
cp "$WORK_DIR/extra-root-one.crt" "$SYSTEM_DIR/disabled-root.crt"
printf '!disabled-root.crt\n' > "$CASE_DIR/ca-certificates.conf"
cat "$WORK_DIR/extra-root-one.crt" >> "$BUNDLE"
run_preservation
[ "$(local_certificate_count)" -eq 0 ] || fail "A disabled system root was promoted to local trust"
pass "Disabled distribution certificates are not promoted to local sources"

new_case removed-system
cp "$WORK_DIR/extra-root-one.crt" "$SYSTEM_DIR/obsolete-root.crt"
cat "$WORK_DIR/extra-root-one.crt" "$WORK_DIR/extra-root-two.crt" >> "$BUNDLE"
run_preservation
rm "$SYSTEM_DIR/obsolete-root.crt"
[ "$(local_certificate_count)" -eq 1 ] || fail "A distribution certificate was preserved as an extra"
assert_preserved extra-root-two
pass "Distribution roots are not frozen when preserving extras"

new_case missing-source-directories
rm -rf "$SYSTEM_DIR" "$LOCAL_DIR"
cat "$WORK_DIR/extra-root-one.crt" > "$BUNDLE"
run_preservation
[ "$(local_certificate_count)" -eq 1 ] || fail "Missing source directories prevented preservation"
pass "An injected bundle can be preserved before CA package installation"

new_case incomplete-pem
printf '\n-----BEGIN CERTIFICATE-----\ninvalid\n' >> "$BUNDLE"
if run_preservation > "$CASE_DIR/output" 2>&1; then
  fail "An incomplete PEM certificate was accepted"
fi
[ "$(local_certificate_count)" -eq 0 ] || fail "Malformed input partially modified local trust"
grep -q 'Failed to parse CA certificates' "$CASE_DIR/output" || fail "Missing PEM error message"
pass "Incomplete PEM input fails before changing trust"

new_case nested-pem
printf '\n-----BEGIN CERTIFICATE-----\n-----BEGIN CERTIFICATE-----\n' >> "$BUNDLE"
if run_preservation > "$CASE_DIR/output" 2>&1; then
  fail "Nested PEM boundaries were accepted"
fi
[ "$(local_certificate_count)" -eq 0 ] || fail "Nested PEM input modified trust"
pass "Nested PEM boundaries fail before changing trust"

new_case unexpected-pem-end
printf '\n-----END CERTIFICATE-----\n' >> "$BUNDLE"
if run_preservation > "$CASE_DIR/output" 2>&1; then
  fail "An unmatched PEM end boundary was accepted"
fi
[ "$(local_certificate_count)" -eq 0 ] || fail "An unmatched PEM boundary modified trust"
pass "Unmatched PEM end boundaries fail before changing trust"

new_case invalid-x509
printf '\n-----BEGIN CERTIFICATE-----\ninvalid\n-----END CERTIFICATE-----\n' >> "$BUNDLE"
if run_preservation > "$CASE_DIR/output" 2>&1; then
  fail "An invalid X.509 certificate was accepted"
fi
[ "$(local_certificate_count)" -eq 0 ] || fail "Invalid X.509 input modified trust"
grep -q 'Invalid CA certificate' "$CASE_DIR/output" || fail "Missing X.509 error message"
pass "Invalid X.509 input fails closed"

new_case invalid-source
printf 'not a certificate\n' > "$SYSTEM_DIR/broken.crt"
cat "$WORK_DIR/extra-root-one.crt" >> "$BUNDLE"
if run_preservation > "$CASE_DIR/output" 2>&1; then
  fail "An invalid managed source was ignored"
fi
[ "$(local_certificate_count)" -eq 0 ] || fail "An invalid source caused a partial trust update"
pass "Invalid managed sources cannot cause false extra-CA detection"

new_case non-ca-extra
cat "$WORK_DIR/extra-root-one.crt" "$WORK_DIR/leaf-certificate.crt" >> "$BUNDLE"
if run_preservation > "$CASE_DIR/output" 2>&1; then
  fail "A leaf certificate was promoted to a local CA"
fi
[ "$(local_certificate_count)" -eq 0 ] || fail "Validation failure partially installed an extra CA"
grep -q 'without CA basic constraints' "$CASE_DIR/output" || fail "Missing CA constraint error message"
pass "All extras are validated before any certificate is installed"

new_case no-openssl
mkdir "$CASE_DIR/empty-bin"
if PATH="$CASE_DIR/empty-bin" /bin/bash -c 'source "$1"; shift; preserve_extra_ca "$@"' \
  bash "$HELPER" "" "$BUNDLE" "$SYSTEM_DIR" "$LOCAL_DIR" > "$CASE_DIR/output" 2>&1; then
  fail "A nonempty bundle was silently ignored without OpenSSL"
fi
grep -q 'OpenSSL is required' "$CASE_DIR/output" || fail "Missing OpenSSL error message"
pass "Missing OpenSSL fails before dependency installation"

new_case privilege-wrapper
mkdir "$CASE_DIR/bin"
cat > "$CASE_DIR/bin/sudo" <<'EOF'
#!/usr/bin/env bash
printf '%s\n' "$1" >> "$PRESERVE_CA_SUDO_LOG"
exec "$@"
EOF
chmod +x "$CASE_DIR/bin/sudo"
cat "$WORK_DIR/extra-root-one.crt" >> "$BUNDLE"
PATH="$CASE_DIR/bin:$PATH" PRESERVE_CA_SUDO_LOG="$CASE_DIR/sudo-log" run_preservation sudo
printf 'mkdir\ninstall\n' > "$CASE_DIR/expected-sudo-log"
cmp "$CASE_DIR/sudo-log" "$CASE_DIR/expected-sudo-log" || fail "Privileged writes did not use the wrapper"
assert_preserved extra-root-one
pass "Root-required writes use the installer's privilege wrapper"

new_case installer-order
mkdir -p "$CASE_DIR/repo/.github/scripts" "$CASE_DIR/bin"
cp "$SCRIPT_DIR/install-llgo.sh" "$CASE_DIR/repo/.github/scripts/install-llgo.sh"
cat > "$CASE_DIR/repo/.github/scripts/preserve-extra-ca.sh" <<'EOF'
preserve_extra_ca() {
  printf 'preserve\n' >> "$INSTALLER_CALL_LOG"
  exit 77
}
EOF
cat > "$CASE_DIR/bin/uname" <<'EOF'
#!/usr/bin/env bash
printf '%s\n' "$INSTALLER_TEST_OS"
EOF
cat > "$CASE_DIR/bin/id" <<'EOF'
#!/usr/bin/env bash
printf '0\n'
EOF
for command_name in apt-get brew; do
  cat > "$CASE_DIR/bin/$command_name" <<'EOF'
#!/usr/bin/env bash
printf 'package-manager\n' >> "$INSTALLER_CALL_LOG"
exit 78
EOF
done
chmod +x "$CASE_DIR/bin/"*
installer_status=0
PATH="$CASE_DIR/bin:$PATH" INSTALLER_TEST_OS=Linux INSTALLER_CALL_LOG="$CASE_DIR/linux-log" \
  LLVM_VERSION=999999 bash "$CASE_DIR/repo/.github/scripts/install-llgo.sh" \
  > "$CASE_DIR/linux-output" 2>&1 || installer_status=$?
[ "$installer_status" -eq 77 ] || fail "CA preservation did not execute before Linux package installation"
printf 'preserve\n' > "$CASE_DIR/expected-linux-log"
cmp "$CASE_DIR/linux-log" "$CASE_DIR/expected-linux-log" || fail "A package operation preceded CA preservation"
installer_status=0
PATH="$CASE_DIR/bin:$PATH" INSTALLER_TEST_OS=Darwin INSTALLER_CALL_LOG="$CASE_DIR/macos-log" \
  bash "$CASE_DIR/repo/.github/scripts/install-llgo.sh" > "$CASE_DIR/macos-output" 2>&1 \
  || installer_status=$?
[ "$installer_status" -eq 78 ] || fail "The macOS installation path unexpectedly changed"
printf 'package-manager\n' > "$CASE_DIR/expected-macos-log"
cmp "$CASE_DIR/macos-log" "$CASE_DIR/expected-macos-log" || fail "CA preservation executed on macOS"
pass "Linux preserves CA before apt; macOS remains unchanged"

if "$INTEGRATION"; then
  [ "$(uname -s)" = Linux ] || fail "The integration test requires Linux"
  for command_name in update-ca-certificates python3 curl git; do
    command -v "$command_name" >/dev/null || fail "The integration test requires $command_name"
  done
  new_case https-integration
  mkdir "$CASE_DIR/etc" "$CASE_DIR/hooks" "$CASE_DIR/server"
  BUNDLE="$CASE_DIR/etc/ca-certificates.crt"
  cat "$WORK_DIR/distribution-root.crt" "$WORK_DIR/extra-root-one.crt" > "$BUNDLE"
  cp "$BUNDLE" "$CASE_DIR/original.crt"
  printf 'distribution-root.crt\n' > "$CASE_DIR/ca-certificates.conf"
  openssl req -new -newkey rsa:2048 -nodes -subj /CN=127.0.0.1 \
    -keyout "$CASE_DIR/server.key" -out "$CASE_DIR/server.csr" > /dev/null 2>&1
  printf 'subjectAltName=IP:127.0.0.1\nbasicConstraints=CA:FALSE\n' > "$CASE_DIR/server.ext"
  openssl x509 -req -sha256 -days 1 -in "$CASE_DIR/server.csr" \
    -CA "$WORK_DIR/extra-root-one.crt" -CAkey "$WORK_DIR/extra-root-one.key" \
    -set_serial 1 -extfile "$CASE_DIR/server.ext" -out "$CASE_DIR/server.crt" > /dev/null 2>&1
  git -c init.defaultBranch=main init --bare "$CASE_DIR/server/repo.git" > /dev/null 2>&1
  git --git-dir="$CASE_DIR/server/repo.git" update-server-info
  python3 - "$CASE_DIR/server" "$CASE_DIR/server.crt" "$CASE_DIR/server.key" "$CASE_DIR/port" \
    > "$CASE_DIR/server.log" 2>&1 <<'PY' &
import functools
import http.server
import pathlib
import ssl
import sys

directory, certificate, private_key, port_file = sys.argv[1:]
handler = functools.partial(http.server.SimpleHTTPRequestHandler, directory=directory)
server = http.server.ThreadingHTTPServer(("127.0.0.1", 0), handler)
context = ssl.SSLContext(ssl.PROTOCOL_TLS_SERVER)
context.load_cert_chain(certificate, private_key)
server.socket = context.wrap_socket(server.socket, server_side=True)
pathlib.Path(port_file).write_text(str(server.server_port))
server.serve_forever()
PY
  SERVER_PID=$!
  attempts=0
  while [ ! -s "$CASE_DIR/port" ]; do
    kill -0 "$SERVER_PID" 2>/dev/null || fail "The HTTPS fixture server exited"
    attempts=$((attempts + 1))
    [ "$attempts" -lt 50 ] || fail "The HTTPS fixture server did not start"
    sleep 0.1
  done
  HTTPS_URL="https://127.0.0.1:$(cat "$CASE_DIR/port")/repo.git"
  curl --disable --noproxy '*' --max-time 10 --fail --silent --show-error \
    --cacert "$BUNDLE" "$HTTPS_URL/HEAD" > "$CASE_DIR/https-output"

  check_git_https() {
    env -u GIT_SSL_NO_VERIFY GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null \
      GIT_SSL_CAINFO="$BUNDLE" git -c http.sslVerify=true -c http.proxy= \
      -c http.lowSpeedTime=5 -c http.lowSpeedLimit=1 ls-remote "$HTTPS_URL"
  }

  check_git_https > /dev/null

  rebuild_fixture_bundle() {
    update-ca-certificates --certsdir "$SYSTEM_DIR" --localcertsdir "$LOCAL_DIR" \
      --etccertsdir "$CASE_DIR/etc" --certsconf "$CASE_DIR/ca-certificates.conf" \
      --hooksdir "$CASE_DIR/hooks" > "$CASE_DIR/update-log" 2>&1
  }

  rebuild_fixture_bundle
  curl_status=0
  curl --disable --noproxy '*' --max-time 10 --fail --silent --show-error \
    --cacert "$BUNDLE" "$HTTPS_URL/HEAD" > "$CASE_DIR/control-output" 2>&1 || curl_status=$?
  [ "$curl_status" -eq 60 ] || fail "The control did not reproduce TLS certificate verification failure"
  if check_git_https > "$CASE_DIR/control-git-output" 2>&1; then
    fail "Git unexpectedly trusted the control after bundle regeneration"
  fi
  grep -qi 'certificate' "$CASE_DIR/control-git-output" || fail "Git failed for a reason other than certificate verification"
  pass "Without preservation, real bundle regeneration breaks HTTPS and Git TLS"

  cp "$CASE_DIR/original.crt" "$BUNDLE"
  run_preservation
  rebuild_fixture_bundle
  curl --disable --noproxy '*' --max-time 10 --fail --silent --show-error \
    --cacert "$BUNDLE" "$HTTPS_URL/HEAD" > "$CASE_DIR/https-output"
  grep -q 'ref: refs/heads/main' "$CASE_DIR/https-output" || fail "The HTTPS fixture response was incorrect"
  check_git_https > /dev/null
  pass "With preservation, HTTPS and Git TLS survive real bundle regeneration"
fi

printf '\nAll %s CA preservation checks passed.\n' "$PASSED"
