#!/bin/bash
# Makes the downloaded Revolut Business spec loadable by oapi-codegen.
#
# Some schemas put prose into numeric keywords, such as `maximum: P7D` on a
# duration or a quoted sentence as the `minimum` of a date-time parameter.
# kin-openapi expects a number there and refuses the whole document, so every
# numeric keyword whose value is not a number (or a 3.0 boolean exclusive
# flag) is dropped. The prose stays in the schema descriptions.
#
# Some path templates name a placeholder in camelCase while the operation
# declares the parameter in snake_case, which every other path uses. The
# generator cannot match them, so camelCase placeholders in path keys are
# rewritten to snake_case.
#
# An inline property schema, such as the `code` enum of ExchangeReason, gets
# the Go type name <Schema><Property>, which can collide with a component
# schema of the same name (ExchangeReasonCode). Such inline schemas get an
# x-go-type-name of <Schema><Property>Value instead.
set -euo pipefail
trap 'on_error "$?" "$LINENO"' ERR

LOG_FILE="${LOG_FILE:-/tmp/$(basename "$0" .sh).log}"
exec > >(tee -a "$LOG_FILE") 2>&1

SPEC_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/api"
readonly SPEC_DIR
readonly SOURCE_SPEC="${SPEC_DIR}/business.yaml"
readonly OUTPUT_SPEC="${SPEC_DIR}/business.patched.yaml"
readonly NUMERIC_KEYWORD='^[[:space:]]*(minimum|maximum|exclusiveMinimum|exclusiveMaximum|multipleOf):'
readonly VALID_VALUE=':[[:space:]]*(-?[0-9.]+|true|false)[[:space:]]*$'

on_error() {
	local exit_code="$1"
	local line="$2"

	printf '{"level":"ERROR","file":"%s","line":%d,"msg":"command failed","exit_code":%d}\n' \
		"${BASH_SOURCE[1]}" "$line" "$exit_code" >&2
}

log() {
	local level="$1"
	shift

	jq -nc \
		--arg time "$(date -u '+%Y-%m-%dT%H:%M:%SZ')" \
		--arg level "$level" \
		--arg file "${BASH_SOURCE[1]##*/}" \
		--arg line "${BASH_LINENO[0]}" \
		--arg func "${FUNCNAME[1]:-main}" \
		--arg msg "$*" \
		'{time:$time,level:$level,file:$file,line:($line|tonumber),func:$func,msg:$msg}' >&2
}

main() {
	if [[ ! -f "$SOURCE_SPEC" ]]; then
		log ERROR "source spec missing path=${SOURCE_SPEC}"
		exit 1
	fi

	local tmp
	tmp="$(mktemp "${SPEC_DIR}/.business.patched.yaml.XXXXXX")"
	# shellcheck disable=SC2064 # expand tmp now, the local is gone when EXIT fires
	trap "rm -f -- '${tmp}'" EXIT

	awk -v key="$NUMERIC_KEYWORD" -v ok="$VALID_VALUE" '
		function snake(s, out, i, c) {
			out = ""
			for (i = 1; i <= length(s); i++) {
				c = substr(s, i, 1)
				if (c ~ /[A-Z]/) {
					out = out "_" tolower(c)
				} else {
					out = out c
				}
			}
			return out
		}
		function pascal(s, out, parts, n, i) {
			out = ""
			n = split(s, parts, "_")
			for (i = 1; i <= n; i++) {
				out = out toupper(substr(parts[i], 1, 1)) substr(parts[i], 2)
			}
			return out
		}
		/^components:/ { in_components = 1 }
		in_components && /^  schemas:/ { in_schemas = 1 }
		in_schemas && /^  [A-Za-z]/ && !/^  schemas:/ { in_schemas = 0 }
		NR == FNR {
			if (in_schemas && match($0, /^    [A-Za-z0-9_]+:[[:space:]]*$/)) {
				name = $0
				gsub(/[[:space:]:]/, "", name)
				schemas[name] = 1
			}
			next
		}
		FNR == 1 { in_components = 0; in_schemas = 0 }
		pending != "" {
			if ($0 !~ /\$ref:/) {
				print "          x-go-type-name: " pending
				collisions++
			}
			pending = ""
		}
		in_schemas && /^    [A-Za-z0-9_]+:[[:space:]]*$/ {
			schema = $0
			gsub(/[[:space:]:]/, "", schema)
		}
		in_schemas && /^        [a-z0-9_]+:[[:space:]]*$/ {
			prop = $0
			gsub(/[[:space:]:]/, "", prop)
			if ((schema pascal(prop)) in schemas) {
				pending = schema pascal(prop) "Value"
			}
		}
		$0 ~ key && $0 !~ ok { dropped++; next }
		/^  \/.*:[[:space:]]*$/ && /\{[a-z]+[A-Z][A-Za-z]*\}/ {
			line = $0
			result = ""
			while (match(line, /\{[a-z]+[A-Z][A-Za-z]*\}/)) {
				result = result substr(line, 1, RSTART - 1) snake(substr(line, RSTART, RLENGTH))
				line = substr(line, RSTART + RLENGTH)
			}
			$0 = result line
			renamed++
		}
		{ print }
		END { print dropped + 0, renamed + 0, collisions + 0 > "/dev/stderr" }
	' "$SOURCE_SPEC" "$SOURCE_SPEC" >"$tmp" 2>"${tmp}.count"

	local dropped renamed collisions
	read -r dropped renamed collisions <"${tmp}.count"
	rm -f -- "${tmp}.count"

	chmod 0644 "$tmp"
	mv -- "$tmp" "$OUTPUT_SPEC"
	log INFO "patched spec dropped_numeric_keywords=${dropped} renamed_path_keys=${renamed} renamed_inline_types=${collisions}"
}

main "$@"
