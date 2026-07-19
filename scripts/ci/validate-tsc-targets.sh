#!/usr/bin/env bash
# Validate a TSC 0.1 registry and its declared targets using the parsing and
# prefix/suffix glob semantics of usurobor/tsc@26aab502. This is intentionally
# a bounded compatibility projection, not a second TSC executor.

set -euo pipefail

usage() {
  echo "usage: $0 --registry PATH --target-root DIR [--target-revision SHA] --target NAME..." >&2
  exit 2
}

registry=""
target_root=""
target_revision=""
targets=()
while (($#)); do
  case "$1" in
    --registry) registry="$2"; shift 2 ;;
    --target-root) target_root="$2"; shift 2 ;;
    --target-revision) target_revision="$2"; shift 2 ;;
    --target) targets+=("$2"); shift 2 ;;
    *) usage ;;
  esac
done

[[ -n "$registry" && -n "$target_root" && ${#targets[@]} -gt 0 ]] || usage
[[ -f "$registry" ]] || { echo "registry missing or not a regular file: $registry" >&2; exit 1; }

tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT
registry_rows="$tmpdir/registry.tsv"

# This mirrors target_registry.ml: trim lines, ignore blank/comment lines,
# split key/value at the first '=', strip one surrounding quote pair, retain
# current [target.NAME] state, and bind subsequent manifest entries. We add
# fail-closed format/header/duplicate checks around the accepted 0.1 shape.
if ! awk -v out="$registry_rows" '
  function trim(s) { sub(/^[[:space:]]+/, "", s); sub(/[[:space:]]+$/, "", s); return s }
  function value_of(line, p, v) {
    p=index(line, "="); v=trim(substr(line, p+1));
    if (length(v)>=2 && substr(v,1,1)=="\"" && substr(v,length(v),1)=="\"")
      v=substr(v,2,length(v)-2);
    return v
  }
  {
    line=trim($0); if (line=="" || substr(line,1,1)=="#") next;
    p=index(line,"=");
    if (p>0) {
      key=trim(substr(line,1,p-1)); val=value_of(line);
      if (key=="format") format=val;
      else if (key=="default_target") default_target=val;
      else if (key=="manifest" && current!="") {
        if (seen[current]++) { print "duplicate target registration: " current > "/dev/stderr"; bad=1 }
        print current "\t" val > out;
      }
      next;
    }
    if (substr(line,1,8)=="[target.") {
      if (substr(line,length(line),1)!="]" || length(line)<=9) {
        print "malformed target section: " line > "/dev/stderr"; bad=1; next
      }
      current=substr(line,9,length(line)-9);
    }
  }
  END {
    if (format=="") { print "missing format field in registry" > "/dev/stderr"; bad=1 }
    else if (format!="tsc-target-registry/0.1") {
      print "unsupported registry format: " format > "/dev/stderr"; bad=1
    }
    if (bad) exit 1
  }
' "$registry"; then
  exit 1
fi

if [[ -n "$target_revision" ]]; then
  git -C "$target_root" cat-file -e "${target_revision}^{commit}" 2>/dev/null || {
    echo "target revision is unavailable: $target_revision" >&2
    exit 1
  }
  git -C "$target_root" ls-tree -r --name-only "$target_revision" >"$tmpdir/all-paths"
else
  (cd "$target_root" && find . \( -type f -o -type l \) -print | sed 's#^\./##' | sort) >"$tmpdir/all-paths"
fi

parse_manifest() {
  local manifest="$1" parsed="$2"
  awk -v out="$parsed" '
    function trim(s) { sub(/^[[:space:]]+/, "", s); sub(/[[:space:]]+$/, "", s); return s }
    function scalar(line, p, v) {
      p=index(line,"="); v=trim(substr(line,p+1));
      if (length(v)>=2 && substr(v,1,1)=="\"" && substr(v,length(v),1)=="\"")
        v=substr(v,2,length(v)-2);
      return v
    }
    {
      line=trim($0);
      if (in_array) {
        if (line=="]") { in_array=0; next }
        n=split(line,q,"\""); if (n>=3) print array_key "\t" q[2] > out;
        next
      }
      if (line=="" || substr(line,1,1)=="#") next;
      p=index(line,"="); if (!p) next;
      key=trim(substr(line,1,p-1)); val=scalar(line);
      if (key=="format" || key=="name" || key=="kind" || key=="description")
        print "META:" key "\t" val > out;
      if (index(line,"[") && (key=="include" || key=="exclude" || key=="optional" || key=="include_targets")) {
        if (val!="[]") { in_array=1; array_key=toupper(key) }
      }
    }
  ' "$manifest"
}

manifest_for() {
  local name="$1"
  awk -F '\t' -v n="$name" '$1==n { print $2; found=1; exit } END { if (!found) exit 1 }' "$registry_rows"
}

declare -A checked counts active
check_target() {
  local name="$1"
  [[ -z "${checked[$name]:-}" ]] || return 0
  [[ -z "${active[$name]:-}" ]] || { echo "include_targets cycle at target '$name'" >&2; return 1; }
  active[$name]=1

  local rel manifest parsed manifest_name kind format
  if ! rel="$(manifest_for "$name")"; then
    echo "target '$name' is not exactly registered in $registry" >&2
    return 1
  fi
  [[ -n "$rel" && "$rel" != /* && "$rel" != *".."* ]] || {
    echo "target '$name' has unsafe manifest path: ${rel:-<empty>}" >&2; return 1;
  }
  manifest="$target_root/$rel"
  [[ -f "$manifest" ]] || { echo "target '$name' manifest missing: $rel" >&2; return 1; }
  parsed="$tmpdir/manifest-${name//[^A-Za-z0-9_.-]/_}.tsv"
  : >"$parsed"
  parse_manifest "$manifest" "$parsed" || { echo "target '$name' manifest malformed: $rel" >&2; return 1; }
  format="$(awk -F '\t' '$1=="META:format" {print $2; exit}' "$parsed")"
  manifest_name="$(awk -F '\t' '$1=="META:name" {print $2; exit}' "$parsed")"
  kind="$(awk -F '\t' '$1=="META:kind" {print $2; exit}' "$parsed")"
  [[ "$format" == "tsc-target/0.1" ]] || { echo "target '$name' manifest has unsupported format: ${format:-<empty>}" >&2; return 1; }
  [[ "$manifest_name" == "$name" ]] || { echo "target '$name' manifest name mismatch: ${manifest_name:-<empty>}" >&2; return 1; }
  case "$kind" in theory|implementation|aggregate) ;; *) echo "target '$name' manifest has unknown kind: ${kind:-<empty>}" >&2; return 1 ;; esac

  local own="$tmpdir/files-${name//[^A-Za-z0-9_.-]/_}"
  : >"$own"
  while IFS=$'\t' read -r tag pattern; do
    [[ "$tag" == "INCLUDE" ]] || continue
    if [[ "$pattern" == *'*'* ]]; then
      local prefix="${pattern%%\**}" suffix="${pattern##*\*}"
      while IFS= read -r path; do
        [[ "$path" == "$prefix"* && "$path" == *"$suffix" ]] && printf '%s\n' "$path" >>"$own"
      done <"$tmpdir/all-paths"
    elif grep -Fxq -- "$pattern" "$tmpdir/all-paths"; then
      printf '%s\n' "$pattern" >>"$own"
    fi
  done <"$parsed"
  while IFS=$'\t' read -r tag pattern; do
    [[ "$tag" == "EXCLUDE" ]] || continue
    local filtered="$own.filtered"
    if [[ "$pattern" == *'*'* ]]; then
      local prefix="${pattern%%\**}" suffix="${pattern##*\*}"
      awk -v p="$prefix" -v s="$suffix" '!(index($0,p)==1 && (s=="" || substr($0,length($0)-length(s)+1)==s))' "$own" >"$filtered"
    else
      grep -Fvx -- "$pattern" "$own" >"$filtered" || true
    fi
    mv "$filtered" "$own"
  done <"$parsed"
  sort -u "$own" -o "$own"

  local total
  total="$(wc -l <"$own" | tr -d ' ')"
  while IFS=$'\t' read -r tag nested; do
    [[ "$tag" == "INCLUDE_TARGETS" ]] || continue
    check_target "$nested" || return 1
    total=$((total + counts[$nested]))
  done <"$parsed"
  ((total > 0)) || { echo "target '$name' resolves to an empty bundle" >&2; return 1; }
  counts[$name]="$total"
  checked[$name]=1
  unset 'active[$name]'
  printf 'PASS target %s: manifest=%s files=%s\n' "$name" "$rel" "$total"
}

for target in "${targets[@]}"; do
  check_target "$target"
done
