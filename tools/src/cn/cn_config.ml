(** cn_config: Hub configuration management.
    
    Config stored in .cn/config.json *)

type peer = {
  name: string;
  url: string;
}

type config = {
  name: string;
  version: string;
  hub_path: string;
  peers: peer list;
  default_branch: string;
  default_remote: string;
}

let default_config ~name ~hub_path = {
  name;
  version = "1.0.0";
  hub_path;
  peers = [];
  default_branch = "main";
  default_remote = "origin";
}

(* Config file path *)
let config_path hub_path = 
  Filename.concat hub_path ".cn/config.json"

(* JSON serialization - manual for simplicity *)
let peer_to_json (p : peer) =
  Printf.sprintf {|{"name":"%s","url":"%s"}|} p.name p.url

let config_to_json c =
  let peers_json = c.peers |> List.map peer_to_json |> String.concat "," in
  Printf.sprintf {|{
  "name": "%s",
  "version": "%s",
  "hub_path": "%s",
  "peers": [%s],
  "defaults": {
    "branch": "%s",
    "remote": "%s"
  }
}|}
    c.name c.version c.hub_path peers_json c.default_branch c.default_remote

(* Simple JSON parsing helpers *)
let extract_string_field json field =
  let pattern = Printf.sprintf {|"%s":\s*"([^"]*)"|}  field in
  let re = Str.regexp pattern in
  try
    let _ = Str.search_forward re json 0 in
    Some (Str.matched_group 1 json)
  with Not_found -> None

(* Parse config from JSON string *)
let config_of_json ~hub_path json =
  let name = extract_string_field json "name" |> Option.value ~default:"agent" in
  let version = extract_string_field json "version" |> Option.value ~default:"1.0.0" in
  let default_branch = extract_string_field json "branch" |> Option.value ~default:"main" in
  let default_remote = extract_string_field json "remote" |> Option.value ~default:"origin" in
  (* TODO: parse peers properly *)
  {
    name;
    version;
    hub_path;
    peers = [];
    default_branch;
    default_remote;
  }

(* Detect hub path from current directory *)
let detect_hub_path () =
  let cwd = Sys.getcwd () in
  let config = Filename.concat cwd ".cn/config.json" in
  if Sys.file_exists config then Some cwd
  else
    (* Check parent for cn-* pattern *)
    let parent = Filename.dirname cwd in
    let parent_config = Filename.concat parent ".cn/config.json" in
    if Sys.file_exists parent_config then Some parent
    else None

(* Derive name from hub path: /path/to/cn-sigma -> sigma *)
let derive_name hub_path =
  hub_path
  |> Filename.basename
  |> fun s ->
    if String.length s > 3 && String.sub s 0 3 = "cn-" then
      String.sub s 3 (String.length s - 3)
    else s
