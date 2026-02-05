(** cn_command: Command type definitions and routing.
    
    Each command is a variant. Subcommands are nested. *)

(* Global flags *)
type flags = {
  json: bool;
  quiet: bool;
  verbose: bool;
  dry_run: bool;
}

let default_flags = {
  json = false;
  quiet = false;
  verbose = false;
  dry_run = false;
}

(* Inbox subcommands *)
type inbox_cmd =
  | Inbox_check
  | Inbox_process
  | Inbox_flush
  | Inbox_log

(* Peer subcommands *)
type peer_cmd =
  | Peer_list
  | Peer_add of string * string  (* name, url *)
  | Peer_remove of string
  | Peer_sync

(* Thread subcommands *)
type thread_cmd =
  | Thread_new of string
  | Thread_daily
  | Thread_list

(* Config subcommands *)
type config_cmd =
  | Config_show
  | Config_set of string * string
  | Config_get of string

(* Top-level commands *)
type command =
  | Init of string option           (* optional name *)
  | Status
  | Inbox of inbox_cmd
  | Peer of peer_cmd
  | Send of string * string         (* peer, branch *)
  | Run of string list              (* action args *)
  | Thread of thread_cmd
  | Config of config_cmd
  | Doctor
  | Update of bool                  (* check_only *)
  | Help
  | Version

(* Aliases *)
let expand_alias = function
  | "i" -> "inbox"
  | "s" -> "status"
  | "p" -> "peer"
  | "t" -> "thread"
  | "c" -> "config"
  | "d" -> "doctor"
  | other -> other

(* Parse inbox subcommand *)
let parse_inbox_cmd = function
  | [] | ["check"] -> Some Inbox_check
  | ["process"] -> Some Inbox_process
  | ["flush"] -> Some Inbox_flush
  | ["log"] -> Some Inbox_log
  | _ -> None

(* Parse peer subcommand *)
let parse_peer_cmd = function
  | [] | ["list"] -> Some Peer_list
  | ["add"; name; url] -> Some (Peer_add (name, url))
  | ["remove"; name] -> Some (Peer_remove name)
  | ["sync"] -> Some Peer_sync
  | _ -> None

(* Parse thread subcommand *)
let parse_thread_cmd = function
  | [] | ["list"] -> Some Thread_list
  | ["new"; topic] -> Some (Thread_new topic)
  | ["daily"] -> Some Thread_daily
  | _ -> None

(* Parse config subcommand *)
let parse_config_cmd = function
  | [] | ["show"] -> Some Config_show
  | ["set"; key; value] -> Some (Config_set (key, value))
  | ["get"; key] -> Some (Config_get key)
  | _ -> None

(* Parse command from args *)
let rec parse_command args =
  match args with
  | [] | ["help"] | ["-h"] | ["--help"] -> Some Help
  | ["--version"] | ["-V"] -> Some Version
  | "init" :: rest -> 
      let name = match rest with [n] -> Some n | _ -> None in
      Some (Init name)
  | ["status"] -> Some Status
  | ["doctor"] -> Some Doctor
  | "update" :: rest ->
      let check_only = List.mem "--check" rest in
      Some (Update check_only)
  | "inbox" :: rest -> 
      parse_inbox_cmd rest |> Option.map (fun c -> Inbox c)
  | "peer" :: rest ->
      parse_peer_cmd rest |> Option.map (fun c -> Peer c)
  | "thread" :: rest ->
      parse_thread_cmd rest |> Option.map (fun c -> Thread c)
  | "config" :: rest ->
      parse_config_cmd rest |> Option.map (fun c -> Config c)
  | "send" :: peer :: branch :: _ -> Some (Send (peer, branch))
  | "run" :: action_args -> Some (Run action_args)
  | [alias] -> 
      let expanded = expand_alias alias in
      if expanded <> alias then parse_command [expanded]
      else None
  | _ -> None

(* Parse flags from args, return (flags, remaining_args) *)
let parse_flags args =
  let rec go flags remaining = function
    | [] -> (flags, List.rev remaining)
    | "--json" :: rest -> go { flags with json = true } remaining rest
    | "-q" :: rest | "--quiet" :: rest -> go { flags with quiet = true } remaining rest
    | "-v" :: rest | "--verbose" :: rest -> go { flags with verbose = true } remaining rest
    | "--dry-run" :: rest -> go { flags with dry_run = true } remaining rest
    | arg :: rest -> go flags (arg :: remaining) rest
  in
  go default_flags [] args

(* String representation for logging *)
let string_of_command = function
  | Init name -> "init" ^ (match name with Some n -> " " ^ n | None -> "")
  | Status -> "status"
  | Inbox Inbox_check -> "inbox check"
  | Inbox Inbox_process -> "inbox process"
  | Inbox Inbox_flush -> "inbox flush"
  | Inbox Inbox_log -> "inbox log"
  | Peer Peer_list -> "peer list"
  | Peer (Peer_add (n, _)) -> "peer add " ^ n
  | Peer (Peer_remove n) -> "peer remove " ^ n
  | Peer Peer_sync -> "peer sync"
  | Send (p, b) -> "send " ^ p ^ " " ^ b
  | Run args -> "run " ^ String.concat " " args
  | Thread Thread_list -> "thread list"
  | Thread (Thread_new t) -> "thread new " ^ t
  | Thread Thread_daily -> "thread daily"
  | Config Config_show -> "config"
  | Config (Config_set (k, _)) -> "config set " ^ k
  | Config (Config_get k) -> "config get " ^ k
  | Doctor -> "doctor"
  | Update check -> "update" ^ (if check then " --check" else "")
  | Help -> "help"
  | Version -> "version"
