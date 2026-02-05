(** cn: Main CLI entrypoint.
    
    Usage: cn <command> [subcommand] [options]
    
    Everything runs through cn. *)

open Cn_command
open Cn_config
open Cn_output

let version = "1.0.0"

(* Help text *)
let help_text = {|cn - Coherent Network agent CLI

Usage: cn <command> [options]

Commands:
  init [name]         Create new hub
  status              Show hub state
  inbox               Manage inbound messages
    check             List inbound branches
    process           Materialize as threads
    flush             Execute decisions
  peer                Manage peers
    list              List peers
    add <name> <url>  Add peer
    remove <name>     Remove peer
    sync              Fetch all peers
  send <peer> <branch>  Push branch to peer
  thread              Manage threads
    new <topic>       Create adhoc thread
    daily             Today's daily thread
    list              List recent threads
  doctor              Check system health
  config              Show/set configuration

Aliases:
  i = inbox, s = status, p = peer, t = thread, d = doctor

Flags:
  --json              Machine-readable output
  --quiet, -q         Minimal output
  --verbose, -v       Detailed output
  --dry-run           Show what would happen
  --help, -h          Show help
  --version, -V       Show version

Examples:
  cn init sigma       Create hub named 'sigma'
  cn inbox check      List inbound branches
  cn peer add pi git@github.com:user/cn-pi.git
  cn doctor           Check system health
|}

(* Load or create config *)
let load_config () =
  match detect_hub_path () with
  | Some hub_path ->
      let config_file = config_path hub_path in
      if Sys.file_exists config_file then
        let ic = open_in config_file in
        let content = really_input_string ic (in_channel_length ic) in
        close_in ic;
        Some (config_of_json ~hub_path content)
      else
        let name = derive_name hub_path in
        Some (default_config ~name ~hub_path)
  | None -> None

(* Run a command *)
let run_command flags cmd config =
  match cmd with
  | Help ->
      print_endline help_text;
      0
  | Version ->
      println (Printf.sprintf "cn %s" version);
      0
  | Status ->
      Status.run ~flags ~config
  | Doctor ->
      Doctor.run ~flags ~config
  | Inbox Inbox_check ->
      println (info "Running inbox check...");
      (* TODO: wire to inbox tool *)
      println (warn "Not yet implemented - use: node tools/dist/inbox.js check");
      1
  | Inbox Inbox_process ->
      println (info "Running inbox process...");
      println (warn "Not yet implemented - use: node tools/dist/inbox.js process");
      1
  | Inbox Inbox_flush ->
      println (info "Running inbox flush...");
      println (warn "Not yet implemented - use: node tools/dist/inbox.js flush");
      1
  | Init name ->
      let hub_name = match name with
        | Some n -> n
        | None -> derive_name (Sys.getcwd ())
      in
      println (info (Printf.sprintf "Initializing hub: %s" hub_name));
      println (warn "Not yet implemented");
      1
  | _ ->
      println (warn (Printf.sprintf "Command not yet implemented: %s" (string_of_command cmd)));
      1

(* Main *)
let () =
  let args = Array.to_list Sys.argv |> List.tl in
  let flags, cmd_args = parse_flags args in
  
  match parse_command cmd_args with
  | None ->
      if List.length cmd_args > 0 then
        println (fail (Printf.sprintf "Unknown command: %s" (List.hd cmd_args)))
      else
        print_endline help_text;
      exit 1
  | Some Help ->
      print_endline help_text;
      exit 0
  | Some Version ->
      println (Printf.sprintf "cn %s" version);
      exit 0
  | Some cmd ->
      (* Commands that don't need config *)
      let needs_config = match cmd with
        | Help | Version | Init _ -> false
        | _ -> true
      in
      let config = 
        if needs_config then
          match load_config () with
          | Some c -> c
          | None ->
              println (fail "Not in a cn hub. Run 'cn init' first.");
              exit 1
        else
          default_config ~name:"" ~hub_path:(Sys.getcwd ())
      in
      let exit_code = run_command flags cmd config in
      exit exit_code
