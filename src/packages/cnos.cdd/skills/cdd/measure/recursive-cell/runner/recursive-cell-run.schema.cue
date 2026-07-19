package recursivecell

#Axis: "alpha" | "beta" | "gamma"
#Level: "L0" | "L1" | "L2" | "L3" | "L4"
#Score: number & >=0 & <=1

#TargetRecord: {
	target: string
	prompt_sha256: string & =~"^[0-9a-f]{64}$"
	response_sha256: string & =~"^[0-9a-f]{64}$"
	report_sha256: string & =~"^[0-9a-f]{64}$"
	report_path: !=""
}

#LevelRecord: {
	level: #Level
	target: string
	alpha: #Score
	beta: #Score
	gamma: #Score
	C_sigma_math: #Score
	C_sigma_num: #Score
}

#Invariant: {
	id: =~"^H(0[1-9]|1[0-3])$"
	status: "pass" | "fail" | "unknown"
	evidence: !=""
	level: #Level | "system"
	primary_axis: #Axis
	next_mca: string | null
}

#RecursiveCellRun: {
	schema: "cnos-recursive-cell-run/0.1"
	target: "cc662-r7-recursive"
	status: "calibration-source-zero-standing"
	targets: [
		#TargetRecord & {target: "cc662-system"},
		#TargetRecord & {target: "cc662-l0"},
		#TargetRecord & {target: "cc662-l1"},
		#TargetRecord & {target: "cc662-l2"},
		#TargetRecord & {target: "cc662-l3"},
		#TargetRecord & {target: "cc662-l4"},
	]
	cross_level: {
		levels: [
			#LevelRecord & {level: "L0", target: "cc662-l0"},
			#LevelRecord & {level: "L1", target: "cc662-l1"},
			#LevelRecord & {level: "L2", target: "cc662-l2"},
			#LevelRecord & {level: "L3", target: "cc662-l3"},
			#LevelRecord & {level: "L4", target: "cc662-l4"},
		]
		C_sigma_cross_math: #Score
		C_sigma_cross_num: #Score
		semantics: "unweighted-geometric-mean-of-L0-L4-canonical-TSC-C_sigma"
	}
	hard_invariant_gate: {
		passed: bool
		items: [
			#Invariant & {id: "H01"}, #Invariant & {id: "H02"},
			#Invariant & {id: "H03"}, #Invariant & {id: "H04"},
			#Invariant & {id: "H05"}, #Invariant & {id: "H06"},
			#Invariant & {id: "H07"}, #Invariant & {id: "H08"},
			#Invariant & {id: "H09"}, #Invariant & {id: "H10"},
			#Invariant & {id: "H11"}, #Invariant & {id: "H12"},
			#Invariant & {id: "H13"},
		]
	}
	bottleneck: {
		level: #Level
		axis: #Axis
		tie_break: "lowest-level-index-then-alpha-beta-gamma"
	}
	disposition: {
		value: "complete" | "complete_with_residuals" | "hold" | "request_planning" | "request_working" | "request_human" | "block"
		reason: !=""
		standing: "none"
	}
	provenance: {
		declared_cm_revision: string & =~"^[0-9a-f]{40}$"
		target_revision: string & =~"^[0-9a-f]{40}$"
		cm_authority_bundle_sha256: string & =~"^[0-9a-f]{64}$"
		registry_sha256: string & =~"^[0-9a-f]{64}$"
		instruction_sha256: string & =~"^[0-9a-f]{64}$"
		runner_sha256: string & =~"^[0-9a-f]{64}$"
		output_schema_sha256: string & =~"^[0-9a-f]{64}$"
		invariant_assessment_template_sha256: string & =~"^[0-9a-f]{64}$"
		invariant_assessment_prompt_sha256: string & =~"^[0-9a-f]{64}$"
		invariant_assessment_sha256: string & =~"^[0-9a-f]{64}$"
		declared_engine_revision: string & =~"^[0-9a-f]{40}$"
		coh_executable_sha256: string & =~"^[0-9a-f]{64}$"
		engine_version: !=""
		run_time_utc: string & =~"^[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}(\\.[0-9]+)?Z$"
		mode: "State-A external-response hybrid ingestion"
	}
}
