#!/bin/bash
# validate.sh - Comprehensive validation script for dtiam CLI
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Counters
PASS=0
FAIL=0
WARN=0

# Helper functions
pass() {
    echo -e "${GREEN}[PASS]${NC} $1"
    ((PASS++))
}

fail() {
    echo -e "${RED}[FAIL]${NC} $1"
    ((FAIL++))
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
    ((WARN++))
}

info() {
    echo -e "      $1"
}

section() {
    echo ""
    echo "=========================================="
    echo "$1"
    echo "=========================================="
}

# Check if binary exists
check_binary() {
    if [ -f "./bin/dtiam" ]; then
        pass "Binary exists at ./bin/dtiam"
        return 0
    else
        warn "Binary not found, building..."
        make build
        if [ -f "./bin/dtiam" ]; then
            pass "Binary built successfully"
            return 0
        else
            fail "Failed to build binary"
            return 1
        fi
    fi
}

# Run Go tests
run_tests() {
    section "Running Unit Tests"
    if go test -v -race ./... 2>&1; then
        pass "All unit tests passed"
    else
        fail "Unit tests failed"
    fi
}

# Run Go vet
run_vet() {
    section "Running Go Vet"
    if go vet ./... 2>&1; then
        pass "Go vet passed"
    else
        fail "Go vet found issues"
    fi
}

# Run Go fmt check
run_fmt_check() {
    section "Checking Code Formatting"
    UNFORMATTED=$(gofmt -l . 2>&1 | grep -v vendor || true)
    if [ -z "$UNFORMATTED" ]; then
        pass "All code is formatted"
    else
        fail "Code needs formatting:"
        echo "$UNFORMATTED"
    fi
}

# Validate CLI help
validate_cli_help() {
    section "Validating CLI Help Output"

    if ./bin/dtiam --help > /dev/null 2>&1; then
        pass "Main help works"
    else
        fail "Main help failed"
    fi

    # Check major command groups
    COMMANDS=("get" "describe" "create" "delete" "config" "user" "service-user" "group" "bulk" "export" "analyze" "boundary" "account" "cache")
    for cmd in "${COMMANDS[@]}"; do
        if ./bin/dtiam "$cmd" --help > /dev/null 2>&1; then
            pass "Help for '$cmd' works"
        else
            fail "Help for '$cmd' failed"
        fi
    done
}

# Validate command subcommands
validate_subcommands() {
    section "Validating Subcommands"

    # Get subcommands
    GET_SUBCMDS=("groups" "users" "policies" "bindings" "environments" "boundaries")
    for subcmd in "${GET_SUBCMDS[@]}"; do
        if ./bin/dtiam get "$subcmd" --help > /dev/null 2>&1; then
            pass "get $subcmd help works"
        else
            fail "get $subcmd help failed"
        fi
    done

    # Service-user subcommands (accessed via service-user, not get)
    SERVICE_USER_SUBCMDS=("list" "get" "create" "update" "delete" "add-to-group" "remove-from-group" "list-groups")
    for subcmd in "${SERVICE_USER_SUBCMDS[@]}"; do
        if ./bin/dtiam service-user "$subcmd" --help > /dev/null 2>&1; then
            pass "service-user $subcmd help works"
        else
            fail "service-user $subcmd help failed"
        fi
    done

    # Bulk subcommands
    BULK_SUBCMDS=("add-users-to-group" "remove-users-from-group" "create-groups" "create-bindings" "export-group-members")
    for subcmd in "${BULK_SUBCMDS[@]}"; do
        if ./bin/dtiam bulk "$subcmd" --help > /dev/null 2>&1; then
            pass "bulk $subcmd help works"
        else
            fail "bulk $subcmd help failed"
        fi
    done

    # Export subcommands
    EXPORT_SUBCMDS=("all" "group" "policy")
    for subcmd in "${EXPORT_SUBCMDS[@]}"; do
        if ./bin/dtiam export "$subcmd" --help > /dev/null 2>&1; then
            pass "export $subcmd help works"
        else
            fail "export $subcmd help failed"
        fi
    done

    # Analyze subcommands
    ANALYZE_SUBCMDS=("user-permissions" "group-permissions" "permissions-matrix" "policy" "least-privilege" "effective-user" "effective-group")
    for subcmd in "${ANALYZE_SUBCMDS[@]}"; do
        if ./bin/dtiam analyze "$subcmd" --help > /dev/null 2>&1; then
            pass "analyze $subcmd help works"
        else
            fail "analyze $subcmd help failed"
        fi
    done

    # Account subcommands
    ACCOUNT_SUBCMDS=("limits" "check-capacity" "subscriptions" "forecast")
    for subcmd in "${ACCOUNT_SUBCMDS[@]}"; do
        if ./bin/dtiam account "$subcmd" --help > /dev/null 2>&1; then
            pass "account $subcmd help works"
        else
            fail "account $subcmd help failed"
        fi
    done
}

# Validate output formats
validate_output_formats() {
    section "Validating Output Format Flags"

    FORMATS=("table" "json" "yaml" "csv" "wide" "plain")
    for fmt in "${FORMATS[@]}"; do
        # Test that -o flag is accepted (will fail without config, but should parse)
        if ./bin/dtiam get groups -o "$fmt" --help > /dev/null 2>&1; then
            pass "Output format '$fmt' is recognized"
        else
            fail "Output format '$fmt' is not recognized"
        fi
    done
}

# Validate version output
validate_version() {
    section "Validating Version Command"

    if ./bin/dtiam version 2>&1 | grep -qi "version"; then
        pass "Version command works"
        info "$(./bin/dtiam version 2>&1 | head -1)"
    else
        warn "Version command may not show version info"
    fi
}

# Check dependencies
check_dependencies() {
    section "Checking Dependencies"

    if go mod verify 2>&1; then
        pass "Go module dependencies verified"
    else
        fail "Go module verification failed"
    fi

    # Check for any missing dependencies
    if go mod tidy -v 2>&1 | grep -q "unused"; then
        warn "Some dependencies may be unused"
    else
        pass "Dependencies are clean"
    fi
}

# Summary
summary() {
    section "Validation Summary"
    echo -e "${GREEN}Passed:${NC} $PASS"
    echo -e "${YELLOW}Warnings:${NC} $WARN"
    echo -e "${RED}Failed:${NC} $FAIL"

    if [ $FAIL -gt 0 ]; then
        echo ""
        echo -e "${RED}Validation FAILED${NC}"
        exit 1
    else
        echo ""
        echo -e "${GREEN}Validation PASSED${NC}"
        exit 0
    fi
}

# Main
main() {
    echo "dtiam CLI Validation Script"
    echo "==========================="

    check_binary || exit 1
    check_dependencies
    run_fmt_check
    run_vet
    run_tests
    validate_version
    validate_cli_help
    validate_subcommands
    validate_output_formats
    summary
}

main "$@"
