package recursivecell

#Axis:   "alpha" | "beta" | "gamma"
#Level:  "L0" | "L1" | "L2" | "L3" | "L4"
#Score:  number & >=0 & <=1
#Digest: string & =~"^[0-9a-f]{64}$"

#ResponseDigests: {
	"cc662-system": #Digest
	"cc662-l0":     #Digest
	"cc662-l1":     #Digest
	"cc662-l2":     #Digest
	"cc662-l3":     #Digest
	"cc662-l4":     #Digest
}

#TargetRecord: {
	target:          string
	prompt_sha256:   #Digest
	response_sha256: #Digest
	report_sha256:   #Digest
	report_path:     !=""
}

#LevelRecord: {
	level:           #Level
	target:          string
	alpha:           #Score
	beta:            #Score
	gamma:           #Score
	bottleneck_axis: #Axis
	C_sigma_math:    #Score
	C_sigma_num:     #Score
}

#Invariant: {
	id:           =~"^H(0[1-9]|1[0-3])$"
	status:       "pass" | "fail" | "unknown"
	evidence:     !=""
	level:        #Level | "system"
	primary_axis: #Axis
	next_mca:     string | null
}

#RecursiveCellRun: {
	schema: "cnos-recursive-cell-run/0.2"
	target: "cc662-r7-recursive"
	status: "calibration-source-zero-standing"
	targets: [
		#TargetRecord & {target: "cc662-system", report_path: "reports/cc662-system.json"},
		#TargetRecord & {target: "cc662-l0", report_path: "reports/cc662-l0.json"},
		#TargetRecord & {target: "cc662-l1", report_path: "reports/cc662-l1.json"},
		#TargetRecord & {target: "cc662-l2", report_path: "reports/cc662-l2.json"},
		#TargetRecord & {target: "cc662-l3", report_path: "reports/cc662-l3.json"},
		#TargetRecord & {target: "cc662-l4", report_path: "reports/cc662-l4.json"},
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
		C_sigma_cross_num:  #Score
		semantics:          "unweighted-geometric-mean-of-L0-L4-canonical-TSC-C_sigma"
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
		level:     #Level
		axis:      #Axis
		tie_break: "lowest-level-index-then-TSC-bottleneck-axis"
	}
	disposition: {
		value:    "complete" | "complete_with_residuals" | "hold" | "request_planning" | "request_working" | "request_human" | "block"
		reason:   !=""
		standing: "none"
	}
	provenance: {
		declared_cm_revision:                 string & =~"^[0-9a-f]{40}$"
		target_revision:                      string & =~"^[0-9a-f]{40}$"
		cm_authority_bundle_sha256:           #Digest
		skill_sha256:                         #Digest
		registry_sha256:                      #Digest
		instruction_sha256:                   #Digest
		runner_sha256:                        #Digest
		output_schema_sha256:                 #Digest
		invariant_assessment_template_sha256: #Digest
		invariant_assessment_prompt_sha256:   #Digest
		invariant_assessment_sha256:          #Digest
		response_sha256:                      #ResponseDigests
		declared_engine_revision:             string & =~"^[0-9a-f]{40}$"
		coh_executable_sha256:                #Digest
		engine_version:                       !=""
		run_time_utc:                         string & =~"^[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}(\\.[0-9]+)?Z$"
		mode:                                 "State-A external-response hybrid ingestion"
	}
}

#PublicationReport: {
	target: string
	path:   string
	sha256: #Digest
}

#PublicationArtifact: {
	path:   string
	sha256: #Digest
}

#PublicationSuccess: {
	schema: "cnos-recursive-cell-publication/0.2"
	target: "cc662-r7-recursive"
	status: "complete"
	result: {
		path:   "recursive-cell-run.json"
		sha256: #Digest
	}
	reports: [
		#PublicationReport & {target: "cc662-system", path: "reports/cc662-system.json"},
		#PublicationReport & {target: "cc662-l0", path: "reports/cc662-l0.json"},
		#PublicationReport & {target: "cc662-l1", path: "reports/cc662-l1.json"},
		#PublicationReport & {target: "cc662-l2", path: "reports/cc662-l2.json"},
		#PublicationReport & {target: "cc662-l3", path: "reports/cc662-l3.json"},
		#PublicationReport & {target: "cc662-l4", path: "reports/cc662-l4.json"},
	]
	emission: {
		prompt_manifest: #PublicationArtifact & {path: "emission/prompt-digests.json"}
		invariant_assessment_prompt: #PublicationArtifact & {path: "emission/invariant-assessment-prompt.md"}
		prompts: [
			#PublicationReport & {target: "cc662-system", path: "emission/prompts/cc662-system.md"},
			#PublicationReport & {target: "cc662-l0", path: "emission/prompts/cc662-l0.md"},
			#PublicationReport & {target: "cc662-l1", path: "emission/prompts/cc662-l1.md"},
			#PublicationReport & {target: "cc662-l2", path: "emission/prompts/cc662-l2.md"},
			#PublicationReport & {target: "cc662-l3", path: "emission/prompts/cc662-l3.md"},
			#PublicationReport & {target: "cc662-l4", path: "emission/prompts/cc662-l4.md"},
		]
	}
	inputs: {
		responses: [
			#PublicationReport & {target: "cc662-system", path: "inputs/responses/cc662-system.json"},
			#PublicationReport & {target: "cc662-l0", path: "inputs/responses/cc662-l0.json"},
			#PublicationReport & {target: "cc662-l1", path: "inputs/responses/cc662-l1.json"},
			#PublicationReport & {target: "cc662-l2", path: "inputs/responses/cc662-l2.json"},
			#PublicationReport & {target: "cc662-l3", path: "inputs/responses/cc662-l3.json"},
			#PublicationReport & {target: "cc662-l4", path: "inputs/responses/cc662-l4.json"},
		]
		invariant_assessment: #PublicationArtifact & {path: "inputs/invariant-assessment.json"}
	}
}
