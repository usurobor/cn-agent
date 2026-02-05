(** cn status: Show hub state at a glance. *)

open Cn_output
open Cn_config

let run ~flags:_ ~config =
  println (info (Printf.sprintf "cn-agent hub: %s" config.name));
  println "";
  
  (* Hub info *)
  println (check_item "hub" true);
  println (check_item_val "name" config.name true);
  println (check_item_val "path" config.hub_path true);
  
  (* Peers *)
  println "";
  let peer_count = List.length config.peers in
  println (check_item_val "peers" (string_of_int peer_count) (peer_count > 0));
  config.peers |> List.iter (fun (p : Cn_config.peer) ->
    println (Printf.sprintf "  %s %s" sym_bullet p.name)
  );
  
  (* TODO: inbox status, recent threads *)
  println "";
  println (inactive "[status] ok version=1.0.0");
  0
