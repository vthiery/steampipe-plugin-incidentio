#!/usr/bin/env bash
# Run a smoke-test query against every incidentio table and report results.
# Exit code is the number of failed tables (0 = all passed).

set -euo pipefail

# Colour codes (disabled if not a terminal)
if [ -t 1 ]; then
  GREEN="\033[0;32m"
  RED="\033[0;31m"
  YELLOW="\033[0;33m"
  RESET="\033[0m"
else
  GREEN="" RED="" YELLOW="" RESET=""
fi

PASS=0
FAIL=0
SKIP=0

run_test() {
  local table="$1"
  local query="$2"
  printf "  %-40s" "$table"
  local output
  if output=$(steampipe query "$query" 2>&1); then
    printf "${GREEN}PASS${RESET}\n"
    ((PASS++)) || true
  else
    # A 403 / missing scope is treated as a skip, not a hard failure
    if echo "$output" | grep -q "missing_required_scope"; then
      printf "${YELLOW}SKIP${RESET} (missing API scope)\n"
      ((SKIP++)) || true
    else
      printf "${RED}FAIL${RESET}\n"
      echo "$output" | sed 's/^/    /'
      ((FAIL++)) || true
    fi
  fi
}

echo ""
echo "incident.io Steampipe plugin — table smoke tests"
echo "================================================="
echo ""

run_test "incidentio_incident"         "select id, reference, name from incidentio_incident limit 1"
run_test "incidentio_action"           "select id, description, status from incidentio_action limit 1"
run_test "incidentio_severity"         "select id, name, rank from incidentio_severity limit 1"
run_test "incidentio_incident_type"    "select id, name, is_default from incidentio_incident_type limit 1"
run_test "incidentio_followups"        "select id, title, status from incidentio_followups limit 1"
run_test "incidentio_incident_updates" "select id, incident_id, created_at from incidentio_incident_updates limit 1"
run_test "incidentio_users"            "select id, name, email from incidentio_users limit 1"
run_test "incidentio_alerts"           "select id, title, status from incidentio_alerts limit 1"
run_test "incidentio_incident_roles"   "select id, name, role_type from incidentio_incident_roles limit 1"
run_test "incidentio_incident_statuses" "select id, name, category from incidentio_incident_statuses limit 1"
run_test "incidentio_custom_fields"    "select id, name, field_type from incidentio_custom_fields limit 1"
run_test "incidentio_escalations"      "select id, title, status from incidentio_escalations limit 1"

echo ""
echo "-------------------------------------------------"
printf "Results: ${GREEN}%d passed${RESET}  ${RED}%d failed${RESET}  ${YELLOW}%d skipped${RESET}\n" "$PASS" "$FAIL" "$SKIP"
echo ""

exit "$FAIL"
