(** cn doctor: Self-diagnosis and health checks. *)

open Cn_output
open Cn_config

(* Run shell command, return Some output or None *)
let run_cmd cmd =
  try
    let ic = Unix.open_process_in cmd in
    let output = input_line ic in
    let _ = Unix.close_process_in ic in
    Some (String.trim output)
  with _ -> None

(* Check if command exists *)
let command_exists cmd =
  run_cmd (Printf.sprintf "which %s 2>/dev/null" cmd) |> Option.is_some

(* Get git config value *)
let git_config key =
  run_cmd (Printf.sprintf "git config --get %s 2>/dev/null" key)

(* Check if file exists *)
let file_exists path = Sys.file_exists path

(* Check if remote is reachable *)
let remote_reachable remote =
  run_cmd (Printf.sprintf "git ls-remote --exit-code %s HEAD 2>/dev/null" remote) |> Option.is_some

type check_result = Ok of string | Fail of string | Warn of string

let run ~flags:_ ~config =
  println (info (Printf.sprintf "Checking %s health..." config.name));
  println "";
  
  let checks = ref [] in
  let add_check name result = checks := (name, result) :: !checks in
  
  (* Git installed *)
  (match run_cmd "git --version" with
   | Some v -> 
       let version = String.sub v 12 (min 6 (String.length v - 12)) in
       add_check "git" (Ok version)
   | None -> add_check "git" (Fail "not installed"));
  
  (* Git user.name *)
  (match git_config "user.name" with
   | Some name -> add_check "git config user.name" (Ok name)
   | None -> add_check "git config user.name" (Fail "not set"));
  
  (* Git user.email *)
  (match git_config "user.email" with
   | Some email -> add_check "git config user.email" (Ok email)
   | None -> add_check "git config user.email" (Fail "not set"));
  
  (* Hub directory *)
  if file_exists config.hub_path then
    add_check "hub directory" (Ok "exists")
  else
    add_check "hub directory" (Fail "not found");
  
  (* Config file *)
  let config_file = config_path config.hub_path in
  if file_exists config_file then
    add_check ".cn/config.json" (Ok "valid")
  else
    add_check ".cn/config.json" (Fail "missing");
  
  (* SOUL.md *)
  let soul_file = Filename.concat config.hub_path "spec/SOUL.md" in
  if file_exists soul_file then
    add_check "spec/SOUL.md" (Ok "exists")
  else
    add_check "spec/SOUL.md" (Warn "missing (optional)");
  
  (* Peers *)
  let peer_count = List.length config.peers in
  add_check "peers" (Ok (Printf.sprintf "%d configured" peer_count));
  
  (* Print results *)
  let oks = ref 0 in
  let warns = ref 0 in
  let fails = ref 0 in
  
  !checks |> List.rev |> List.iter (fun (name, result) ->
    match result with
    | Ok v -> 
        println (check_item_val name v true);
        incr oks
    | Warn v ->
        println (status_line name 25 (color yellow (sym_warn ^ " " ^ v)));
        incr warns
    | Fail v ->
        println (check_item_val name v false);
        incr fails
  );
  
  println "";
  
  (* Summary *)
  if !fails = 0 then begin
    print_ok "All critical checks passed.";
    println (machine_footer [
      ("ok", string_of_int !oks);
      ("warn", string_of_int !warns);
      ("fail", "0")
    ] "1.0.0");
    0
  end else begin
    print_fail (Printf.sprintf "%d issue(s) found" !fails);
    println "";
    println "Fix by running:";
    println "";
    
    (* Suggest fixes based on failures *)
    let fix_num = ref 1 in
    !checks |> List.rev |> List.iter (fun (name, result) ->
      match result with
      | Fail _ when name = "git config user.name" ->
          println (Printf.sprintf "  %d) %s" !fix_num (cmd "git config --global user.name \"Your Name\""));
          incr fix_num
      | Fail _ when name = "git config user.email" ->
          println (Printf.sprintf "  %d) %s" !fix_num (cmd "git config --global user.email \"you@example.com\""));
          incr fix_num
      | Fail _ when name = "git" ->
          println (Printf.sprintf "  %d) Install git: %s" !fix_num (cmd "apt install git"));
          incr fix_num
      | _ -> ()
    );
    
    println "";
    println (Printf.sprintf "Then run: %s" (cmd "cn doctor"));
    println (machine_footer [
      ("ok", string_of_int !oks);
      ("warn", string_of_int !warns);
      ("fail", string_of_int !fails)
    ] "1.0.0");
    1
  end
