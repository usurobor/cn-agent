(** cn_workflow.ml — Orchestrator IR runtime (`cn.orchestrator.v1`).

    Implements the mechanical workflow engine specified in
    `docs/reference/runtime/ORCHESTRATORS.md` §7-8.

    Responsibilities:
    - Parse an orchestrator manifest (`cn.orchestrator.v1` JSON)
      into typed IR.
    - Validate it structurally: no duplicate step ids, all step
      refs resolve, every dispatched op is in `permissions.ops`,
      every `llm` step requires `permissions.llm = true`.
    - Discover installed orchestrators under
      `.cn/vendor/packages/*/orchestrators/<id>/orchestrator.json`
      per the package-system content-class rule (§10.3).
    - Execute a loaded IR step by step: deterministic control
      flow (`if`, `match`, `return`, `fail`) plus typed `op`
      dispatch via [Cn_executor.execute_op]. Receipts/events are
      emitted to the global trace session (no-op when no session).

    Deferred to a later cycle (documented in docs/reference/runtime/ORCHESTRATORS.md §8.2):
    - [parallel] step kind (cnos has no async model)
    - [llm] step execution (prompt + context injection mechanism
      not yet designed)

    Naming: this module is distinct from the pre-existing
    [Cn_orchestrator] which implements the N-pass LLM bind loop.
    Two orthogonal orchestration concerns, two modules.

    v3.40.0 (#182 Move 2 slice 3) note: the 6 pure IR types and 10
    pure functions (parsers, validator, helpers, result combinator)
    that used to live in this module were extracted into
    [src/lib/cn_workflow_ir.ml]. [Cn_workflow_ir] is now the
    canonical authority for [trigger], [permissions], [step],
    [orchestrator], [issue_kind], [issue] and their parse/validate
    surface. This module re-exports each type below via OCaml
    type-equality syntax and delegates each pure function as a
    one-line let-binding so existing callers ([cn_workflow_test.ml]
    plus the IO functions below that pattern-match on step/issue
    constructors) compile unchanged. [Cn_workflow] retains every
    IO-side function ([parse_file], [discover], [doctor_issues],
    [execute_step], [execute], [typed_op_of_op_step], [trace_event],
    [find_step], [env_get], [as_bool], [as_string]) and the three
    IO-transit types ([load_outcome], [installed], [outcome]) —
    these last are pure by shape but only consumed by IO code in
    this module and in [cn_runtime_contract.ml::build_orchestrator_registry],
    so per issue #196 option (b) they stay here with their
    consumers. *)

(* === Types (re-exported from Cn_workflow_ir for caller compatibility) === *)

type trigger = Cn_workflow_ir.trigger = {
  trigger_kind : string;       (** "command" | "schedule" | "event" *)
  trigger_name : string;
}

type permissions = Cn_workflow_ir.permissions = {
  llm : bool;
  ops : string list;           (** allowed op kind names *)
  external_effects : bool;
}

type step = Cn_workflow_ir.step =
  | Op_step of {
      id : string;
      op : string;             (** op kind name e.g. "fs_read" *)
      args : (string * Cn_json.t) list;
      bind : string option;
    }
  | Llm_step of {
      id : string;
      prompt : string;
      bind : string option;
    }
  | If_step of {
      id : string;
      cond : string;           (** name of bound boolean *)
      then_ref : string;       (** step id to jump to on true *)
      else_ref : string;       (** step id to jump to on false *)
    }
  | Match_step of {
      id : string;
      input : string;          (** bound value name to match on *)
      cases : (string * string) list;  (** literal-value -> step id *)
      default : string option; (** step id when no case matches *)
    }
  | Return_step of {
      id : string;
      value : Cn_json.t;
    }
  | Fail_step of {
      id : string;
      message : string;
    }

type orchestrator = Cn_workflow_ir.orchestrator = {
  kind : string;
  name : string;
  trigger : trigger;
  permissions : permissions;
  steps : step list;
}

type issue_kind = Cn_workflow_ir.issue_kind =
  | Missing_field of string
  | Unknown_step_kind of string
  | Duplicate_step_id of string
  | Invalid_step_ref of string
  | Permission_gap of string
  | Llm_without_permission
  | Empty_steps

type issue = Cn_workflow_ir.issue = {
  issue_kind : issue_kind;
  message : string;
}

(* === Pure helpers delegated to Cn_workflow_ir ===

    Note: [let ( let* ) ...] is NOT re-exported here because no IO
    function in this module uses the result-bind operator after the
    extraction — all the parsers that used [let*] were moved to
    [Cn_workflow_ir]. Re-exporting would create a dead binding. *)

let require_string = Cn_workflow_ir.require_string
let parse_string_list = Cn_workflow_ir.parse_string_list
let parse_trigger = Cn_workflow_ir.parse_trigger
let parse_permissions = Cn_workflow_ir.parse_permissions
let parse_step = Cn_workflow_ir.parse_step
let parse = Cn_workflow_ir.parse
let step_id = Cn_workflow_ir.step_id
let validate = Cn_workflow_ir.validate
let manifest_orchestrator_ids = Cn_workflow_ir.manifest_orchestrator_ids

(* === IO-side parse_file === *)

let parse_file path =
  if not (Cn_ffi.Fs.exists path) then
    Error (Printf.sprintf "missing orchestrator file: %s" path)
  else
    match Cn_json.parse (Cn_ffi.Fs.read path) with
    | Error msg -> Error (Printf.sprintf "invalid JSON: %s" msg)
    | Ok json -> parse json

(* === Discovery — IO-transit types (kept in cn_workflow.ml per
   issue #196 option (b): pure by shape but only consumed by IO
   functions in this module and by cn_runtime_contract.ml, which
   pattern-matches on Loaded/Load_error — moving them would force
   caller churn without any pure-consumer benefit) === *)

type load_outcome =
  | Loaded of orchestrator
  | Load_error of string

type installed = {
  package : string;
  id : string;
  path : string;
  outcome : load_outcome;
}

let discover ~hub_path =
  Cn_assets.list_installed_packages hub_path
  |> List.concat_map (fun (pkg_name, pkg_dir) ->
    let manifest_path = Cn_ffi.Path.join pkg_dir "cn.package.json" in
    if not (Cn_ffi.Fs.exists manifest_path) then []
    else
      match Cn_json.parse (Cn_ffi.Fs.read manifest_path) with
      | Error msg ->
          Printf.eprintf
            "cn: workflow: package %s: cannot parse %s: %s\n"
            pkg_name manifest_path msg;
          []
      | Ok json ->
          manifest_orchestrator_ids ~pkg_name json
          |> List.map (fun orch_id ->
            let path = Cn_ffi.Path.join pkg_dir
              (Printf.sprintf "orchestrators/%s/orchestrator.json" orch_id) in
            let outcome = match parse_file path with
              | Ok o -> Loaded o
              | Error msg -> Load_error msg
            in
            { package = pkg_name; id = orch_id; path; outcome }))

(* === Doctor validation === *)

(** Collect a doctor-facing summary of every installed orchestrator:
    load failures and schema-level issues each become a line. *)
let doctor_issues ~hub_path =
  discover ~hub_path
  |> List.concat_map (fun i ->
    match i.outcome with
    | Load_error msg ->
        [ Printf.sprintf "%s/%s: %s" i.package i.id msg ]
    | Loaded o ->
        validate o
        |> List.map (fun (iss : issue) ->
          Printf.sprintf "%s/%s: %s" i.package i.id iss.message))

(* === Execution === *)

type outcome =
  | Completed of Cn_json.t
  | Failed of string

(** Build a [Cn_shell.typed_op] from an [Op_step] invocation. *)
let typed_op_of_op_step ~op_kind (op_name : string) (args : (string * Cn_json.t) list) =
  { Cn_shell.kind = op_kind;
    op_id = Some (Printf.sprintf "orch-%s" op_name);
    fields = args }

let trace_event ~component ~event ~status ~reason ~details =
  let reason_opt = if reason = "" then None else Some reason in
  Cn_trace.gemit
    ~component
    ~layer:Cn_trace.Body
    ~event
    ~severity:Cn_trace.Info
    ~status
    ?reason:reason_opt
    ~details
    ()

(** Look up a step by id. Returns [None] if not found. *)
let find_step orch id =
  List.find_opt (fun s -> step_id s = id) orch.steps

(** Render a bound-value reference in an [if]/[match] step. The
    environment holds JSON values, and the condition/input names refer
    to either a previously bound name or the literal string in the IR.
    v1 keeps this simple: look up by name, return [None] if not bound. *)
let env_get (env : (string * Cn_json.t) list) name =
  List.assoc_opt name env

(** Convert a bound value to a boolean for [if]-step cond. *)
let as_bool = function
  | Cn_json.Bool b -> b
  | Cn_json.String "true" -> true
  | Cn_json.String "false" -> false
  | Cn_json.Int 0 -> false
  | Cn_json.Int _ -> true
  | Cn_json.Null -> false
  | _ -> true

(** Convert a bound value to a string for [match]-step input. *)
let as_string = function
  | Cn_json.String s -> s
  | Cn_json.Bool true -> "true"
  | Cn_json.Bool false -> "false"
  | Cn_json.Int i -> string_of_int i
  | Cn_json.Null -> "null"
  | v -> Cn_json.to_string v

(** Execute one step, producing either a new environment + next step
    id, or a terminal outcome. *)
let rec execute_step ~hub_path ~shell_config (o : orchestrator) env step =
  let name = o.name in
  let id = step_id step in
  trace_event
    ~component:"workflow"
    ~event:"workflow.step.start"
    ~status:Cn_trace.Ok_
    ~reason:""
    ~details:[
      "workflow", Cn_json.String name;
      "step_id", Cn_json.String id;
    ];
  match step with
  | Return_step s ->
      trace_event
        ~component:"workflow"
        ~event:"workflow.step.complete"
        ~status:Cn_trace.Ok_
        ~reason:""
        ~details:[
          "workflow", Cn_json.String name;
          "step_id", Cn_json.String s.id;
          "terminal", Cn_json.String "return";
        ];
      `Terminal (Completed s.value)

  | Fail_step s ->
      trace_event
        ~component:"workflow"
        ~event:"workflow.step.complete"
        ~status:Cn_trace.Error_status
        ~reason:s.message
        ~details:[
          "workflow", Cn_json.String name;
          "step_id", Cn_json.String s.id;
          "terminal", Cn_json.String "fail";
        ];
      `Terminal (Failed (Printf.sprintf "workflow fail at step '%s': %s"
                           s.id s.message))

  | Op_step s ->
      (* Permission gate — defence in depth; validator also flags. *)
      if not (List.mem s.op o.permissions.ops) then begin
        trace_event
          ~component:"workflow"
          ~event:"workflow.step.complete"
          ~status:Cn_trace.Blocked
          ~reason:"permission_gap"
          ~details:[
            "workflow", Cn_json.String name;
            "step_id", Cn_json.String s.id;
            "op", Cn_json.String s.op;
          ];
        `Terminal (Failed (Printf.sprintf
          "workflow '%s' step '%s': op '%s' not permitted by permissions.ops"
          name s.id s.op))
      end else
        (match Cn_shell.op_kind_of_string s.op with
         | None ->
             `Terminal (Failed (Printf.sprintf
               "workflow '%s' step '%s': unknown op '%s'"
               name s.id s.op))
         | Some op_kind ->
             let top = typed_op_of_op_step ~op_kind s.op s.args in
             let trigger_id = Printf.sprintf "workflow-%s-%s" name s.id in
             let receipt =
               Cn_executor.execute_op
                 ~hub_path
                 ~trigger_id
                 ~config:shell_config
                 top
             in
             (match receipt.status with
              | Cn_shell.Ok_status ->
                  let artifacts_json =
                    receipt.artifacts
                    |> List.map (fun (a : Cn_shell.artifact) ->
                      Cn_json.Object [
                        "path", Cn_json.String a.path;
                        "hash", Cn_json.String a.hash;
                        "size", Cn_json.Int a.size;
                      ])
                  in
                  let bound =
                    Cn_json.Object [
                      "status", Cn_json.String "ok";
                      "kind", Cn_json.String receipt.kind;
                      "artifacts", Cn_json.Array artifacts_json;
                    ]
                  in
                  let env' = match s.bind with
                    | Some key -> (key, bound) :: env
                    | None -> env
                  in
                  trace_event
                    ~component:"workflow"
                    ~event:"workflow.step.complete"
                    ~status:Cn_trace.Ok_
                    ~reason:""
                    ~details:[
                      "workflow", Cn_json.String name;
                      "step_id", Cn_json.String s.id;
                      "op", Cn_json.String s.op;
                      "bind", Cn_json.String (Option.value s.bind ~default:"");
                    ];
                  `Continue (env', next_step_after o s.id)
              | _ ->
                  trace_event
                    ~component:"workflow"
                    ~event:"workflow.step.complete"
                    ~status:Cn_trace.Error_status
                    ~reason:receipt.reason
                    ~details:[
                      "workflow", Cn_json.String name;
                      "step_id", Cn_json.String s.id;
                      "op", Cn_json.String s.op;
                    ];
                  `Terminal (Failed (Printf.sprintf
                    "workflow '%s' step '%s': op '%s' failed: %s"
                    name s.id s.op receipt.reason))))

  | Llm_step s ->
      (* v1 stub: parsed + validated + permission-checked, but not
         executed. The prompt + context injection mechanism is still
         being designed (PLAN-174 §Risk). *)
      trace_event
        ~component:"workflow"
        ~event:"workflow.step.complete"
        ~status:Cn_trace.Error_status
        ~reason:"llm_not_implemented"
        ~details:[
          "workflow", Cn_json.String name;
          "step_id", Cn_json.String s.id;
        ];
      `Terminal (Failed (Printf.sprintf
        "workflow '%s' step '%s': llm step not yet implemented in this release"
        name s.id))

  | If_step s ->
      let cond_val = env_get env s.cond in
      let branch = match cond_val with
        | Some v -> if as_bool v then s.then_ref else s.else_ref
        | None ->
            (* Unbound cond → treat as false, take else branch. This
               is explicit default behaviour, not a silent skip. *)
            s.else_ref
      in
      (match find_step o branch with
       | Some target ->
           trace_event
             ~component:"workflow"
             ~event:"workflow.step.complete"
             ~status:Cn_trace.Ok_
             ~reason:""
             ~details:[
               "workflow", Cn_json.String name;
               "step_id", Cn_json.String s.id;
               "branch", Cn_json.String branch;
             ];
           `Jump (env, target)
       | None ->
           `Terminal (Failed (Printf.sprintf
             "workflow '%s' step '%s': branch target '%s' not found"
             name s.id branch)))

  | Match_step s ->
      let input_val = env_get env s.input in
      let key = match input_val with
        | Some v -> as_string v
        | None -> ""
      in
      let branch = match List.assoc_opt key s.cases, s.default with
        | Some target, _ -> Some target
        | None, Some d -> Some d
        | None, None -> None
      in
      (match branch with
       | None ->
           `Terminal (Failed (Printf.sprintf
             "workflow '%s' step '%s': no case for input '%s' and no default"
             name s.id key))
       | Some target ->
           match find_step o target with
           | Some tgt ->
               trace_event
                 ~component:"workflow"
                 ~event:"workflow.step.complete"
                 ~status:Cn_trace.Ok_
                 ~reason:""
                 ~details:[
                   "workflow", Cn_json.String name;
                   "step_id", Cn_json.String s.id;
                   "branch", Cn_json.String target;
                 ];
               `Jump (env, tgt)
           | None ->
               `Terminal (Failed (Printf.sprintf
                 "workflow '%s' step '%s': match target '%s' not found"
                 name s.id target)))

(** Return the step immediately following the given step id in the
    manifest's step list, or [None] if this was the last step. *)
and next_step_after o id =
  let rec walk = function
    | [] -> None
    | s :: next :: _ when step_id s = id -> Some next
    | _ :: rest -> walk rest
  in
  walk o.steps

(** Execute the orchestrator from its first step.

    [env] supplies a pre-seeded binding environment; the default
    [[]] starts clean. Tests use this to pre-bind scalars so the
    [if] and [match] branches can be exercised directly without
    needing a working [llm] step (which is deferred in v1). *)
let execute ~hub_path ?(shell_config = Cn_shell.default_shell_config)
            ?(env = []) o =
  match o.steps with
  | [] -> Failed (Printf.sprintf "workflow '%s' has no steps" o.name)
  | first :: _ ->
      let rec loop env step =
        match execute_step ~hub_path ~shell_config o env step with
        | `Terminal outcome -> outcome
        | `Continue (env', Some next) -> loop env' next
        | `Continue (_, None) ->
            Failed (Printf.sprintf
              "workflow '%s' step '%s' fell off the end without return/fail"
              o.name (step_id step))
        | `Jump (env', next) -> loop env' next
      in
      loop env first
