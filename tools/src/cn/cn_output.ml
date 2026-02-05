(** cn_output: UX-CLI compliant terminal output.
    
    Colors encode semantics, never decoration.
    Symbols required alongside color for accessibility. *)

(* ANSI color codes *)
let reset = "\027[0m"
let green = "\027[32m"
let red = "\027[31m"
let yellow = "\027[33m"
let cyan = "\027[36m"
let magenta = "\027[35m"
let dim = "\027[2m"

(* Check NO_COLOR environment variable *)
let no_color () =
  match Sys.getenv_opt "NO_COLOR" with
  | Some _ -> true
  | None -> false

let color code text =
  if no_color () then text
  else code ^ text ^ reset

(* Semantic output functions *)
let ok text = color green ("✓ " ^ text)
let fail text = color red ("✗ " ^ text)
let warn text = color yellow ("⚠ " ^ text)
let info text = color cyan text
let cmd text = color magenta text
let inactive text = color dim text

(* Symbols for non-color contexts *)
let sym_ok = "✓"
let sym_fail = "✗"
let sym_warn = "⚠"
let sym_blocked = "⏸"
let sym_bullet = "•"
let sym_arrow = "→"

(* Status line with dot leaders *)
let status_line label width status =
  let dots = String.make (width - String.length label) '.' in
  Printf.sprintf "%s%s %s" label dots status

(* Checklist item *)
let check_item label ok =
  let width = 25 in
  let status = if ok then color green sym_ok else color red sym_fail in
  status_line label width status

(* Checklist item with value *)
let check_item_val label value ok =
  let width = 25 in
  let status = 
    if ok then color green (sym_ok ^ " " ^ value) 
    else color red (sym_fail ^ " " ^ value) 
  in
  status_line label width status

(* Section header *)
let section title =
  "\n" ^ color cyan ("=== " ^ String.uppercase_ascii title ^ " ===")

(* Error block *)
let error_block cause items fix_cmds rerun =
  let cause_line = color red (sym_fail ^ " Cannot continue — " ^ cause) in
  let items_str = items |> List.map (fun i -> "  " ^ sym_bullet ^ " " ^ i) |> String.concat "\n" in
  let fix_header = "\nFix by running:" in
  let fix_str = fix_cmds |> List.mapi (fun i c -> 
    Printf.sprintf "  %d) %s" (i + 1) (cmd c)) |> String.concat "\n" in
  let rerun_str = Printf.sprintf "\nAfter completing the steps above, run:\n\n  %s" (cmd rerun) in
  String.concat "\n\n" [cause_line; items_str; fix_header; fix_str; rerun_str]

(* Machine-readable footer *)
let machine_footer pairs version =
  let kvs = pairs |> List.map (fun (k, v) -> k ^ "=" ^ v) |> String.concat " " in
  color dim (Printf.sprintf "\n[status] %s version=%s" kvs version)

(* Progress indicator *)
let progress phase =
  info (phase ^ "...")

(* Print helpers *)
let println s = print_endline s
let print_ok s = println (ok s)
let print_fail s = println (fail s)
let print_warn s = println (warn s)
let print_info s = println (info s)
let print_cmd s = println ("  " ^ cmd s)
