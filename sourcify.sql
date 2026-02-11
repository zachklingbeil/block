
Enum "signature_type_enum" {
  "function"
  "event"
  "error"
}

Table "code" {
  "code_hash" bytea [not null]
  "created_at" timestamp [not null, default: `now()`]
  "updated_at" timestamp [not null, default: `now()`]
  "created_by" "character varying" [not null, default: `CURRENT_USER`]
  "updated_by" "character varying" [not null, default: `CURRENT_USER`]
  "code_hash_keccak" bytea [not null]
  "code" bytea

  Indexes {
    code_hash_keccak [type: btree, name: "code_code_hash_keccak"]
  }
}

Table "compiled_contracts" {
  "id" uuid [not null, default: `gen_random_uuid()`]
  "created_at" timestamp [not null, default: `now()`]
  "updated_at" timestamp [not null, default: `now()`]
  "created_by" "character varying" [not null, default: `CURRENT_USER`]
  "updated_by" "character varying" [not null, default: `CURRENT_USER`]
  "compiler" "character varying" [not null]
  "version" "character varying" [not null]
  "language" "character varying" [not null]
  "name" "character varying" [not null]
  "fully_qualified_name" "character varying" [not null]
  "compiler_settings" jsonb [not null]
  "compilation_artifacts" jsonb [not null]
  "creation_code_hash" bytea
  "creation_code_artifacts" jsonb
  "runtime_code_hash" bytea [not null]
  "runtime_code_artifacts" jsonb [not null]

  Indexes {
    creation_code_hash [type: btree, name: "compiled_contracts_creation_code_hash"]
    runtime_code_hash [type: btree, name: "compiled_contracts_runtime_code_hash"]
  }
}

Table "compiled_contracts_signatures" {
  "id" uuid [not null, default: `gen_random_uuid()`]
  "compilation_id" uuid [not null]
  "signature_hash_32" bytea [not null]
  "signature_type" signature_type_enum [not null]
  "created_at" timestamp [not null, default: `now()`]

  Indexes {
    signature_hash_32 [type: btree, name: "compiled_contracts_signatures_signature_idx"]
    (signature_type, signature_hash_32) [type: btree, name: "compiled_contracts_signatures_type_signature_idx"]
  }
}

Table "compiled_contracts_sources" {
  "id" uuid [not null, default: `gen_random_uuid()`]
  "compilation_id" uuid [not null]
  "source_hash" bytea [not null]
  "path" "character varying" [not null]

  Indexes {
    compilation_id [type: btree, name: "compiled_contracts_sources_compilation_id"]
    source_hash [type: btree, name: "compiled_contracts_sources_source_hash"]
  }
}

Table "contract_deployments" {
  "id" uuid [not null, default: `gen_random_uuid()`]
  "created_at" timestamp [not null, default: `now()`]
  "updated_at" timestamp [not null, default: `now()`]
  "created_by" "character varying" [not null, default: `CURRENT_USER`]
  "updated_by" "character varying" [not null, default: `CURRENT_USER`]
  "chain_id" bigint [not null]
  "address" bytea [not null]
  "transaction_hash" bytea
  "block_number" numeric
  "transaction_index" numeric
  "deployer" bytea
  "contract_id" uuid [not null]

  Indexes {
    address [type: btree, name: "contract_deployments_address"]
    contract_id [type: btree, name: "contract_deployments_contract_id"]
  }
}

Table "contracts" {
  "id" uuid [not null, default: `gen_random_uuid()`]
  "created_at" timestamp [not null, default: `now()`]
  "updated_at" timestamp [not null, default: `now()`]
  "created_by" "character varying" [not null, default: `CURRENT_USER`]
  "updated_by" "character varying" [not null, default: `CURRENT_USER`]
  "creation_code_hash" bytea
  "runtime_code_hash" bytea [not null]

  Indexes {
    creation_code_hash [type: btree, name: "contracts_creation_code_hash"]
    (creation_code_hash, runtime_code_hash) [type: btree, name: "contracts_creation_code_hash_runtime_code_hash"]
    runtime_code_hash [type: btree, name: "contracts_runtime_code_hash"]
  }
}

Table "schema_migrations" {
  "version" "character varying" [not null]
}

Table "session" {
  "sid" "character varying" [not null]
  "sess" json [not null]
  "expire" timestamp(6) [not null]

  Indexes {
    expire [type: btree, name: "IDX_session_expire"]
  }
}

Table "signatures" {
  "signature_hash_32" bytea [not null]
  "signature_hash_4" bytea
  "signature" "character varying" [not null]
  "created_at" timestamp [not null, default: `now()`]

  Indexes {
    signature_hash_4 [type: btree, name: "signatures_hash_4_idx"]
  }
}

Table "sources" {
  "source_hash" bytea [not null]
  "source_hash_keccak" bytea [not null]
  "content" "character varying" [not null]
  "created_at" timestamp [not null, default: `now()`]
  "updated_at" timestamp [not null, default: `now()`]
  "created_by" "character varying" [not null, default: `CURRENT_USER`]
  "updated_by" "character varying" [not null, default: `CURRENT_USER`]
}

Table "sourcify_matches" {
  "id" bigint [not null]
  "verified_contract_id" bigint [not null]
  "creation_match" "character varying"
  "runtime_match" "character varying"
  "created_at" timestamp [not null, default: `now()`]
  "updated_at" timestamp [not null, default: `now()`]
  "metadata" json [not null]

  Indexes {
    verified_contract_id [type: btree, name: "sourcify_matches_verified_contract_id_idx"]
  }
}

Table "sourcify_sync" {
  "id" bigint [not null]
  "chain_id" numeric [not null]
  "address" bytea [not null]
  "match_type" "character varying" [not null]
  "synced" boolean [not null, default: false]
  "created_at" timestamp [not null, default: `now()`]
}

Table "verification_jobs" {
  "id" uuid [not null, default: `gen_random_uuid()`]
  "started_at" timestamp [not null, default: `now()`]
  "completed_at" timestamp
  "chain_id" bigint [not null]
  "contract_address" bytea [not null]
  "verified_contract_id" bigint
  "error_code" "character varying"
  "error_id" uuid
  "error_data" json
  "verification_endpoint" "character varying" [not null]
  "hardware" "character varying"
  "compilation_time" bigint

  Indexes {
    (chain_id, contract_address) [type: btree, name: "verification_jobs_chain_id_address_idx"]
  }
}

Table "verification_jobs_ephemeral" {
  "id" uuid [not null, default: `gen_random_uuid()`]
  "recompiled_creation_code" bytea
  "recompiled_runtime_code" bytea
  "onchain_creation_code" bytea
  "onchain_runtime_code" bytea
  "creation_transaction_hash" bytea
}

Table "verified_contracts" {
  "id" bigint [not null]
  "created_at" timestamp [not null, default: `now()`]
  "updated_at" timestamp [not null, default: `now()`]
  "created_by" "character varying" [not null, default: `CURRENT_USER`]
  "updated_by" "character varying" [not null, default: `CURRENT_USER`]
  "deployment_id" uuid [not null]
  "compilation_id" uuid [not null]
  "creation_match" boolean [not null]
  "creation_values" jsonb
  "creation_transformations" jsonb
  "creation_metadata_match" boolean
  "runtime_match" boolean [not null]
  "runtime_values" jsonb
  "runtime_transformations" jsonb
  "runtime_metadata_match" boolean

  Indexes {
    compilation_id [type: btree, name: "verified_contracts_compilation_id"]
    deployment_id [type: btree, name: "verified_contracts_deployment_id"]
  }
}

Ref "compiled_contracts_creation_code_hash_fkey":"code"."code_hash" < "compiled_contracts"."creation_code_hash"

Ref "compiled_contracts_runtime_code_hash_fkey":"code"."code_hash" < "compiled_contracts"."runtime_code_hash"

Ref "compiled_contracts_signatures_compilation_id_fkey":"compiled_contracts"."id" < "compiled_contracts_signatures"."compilation_id"

Ref "compiled_contracts_signatures_signature_hash_32_fkey":"signatures"."signature_hash_32" < "compiled_contracts_signatures"."signature_hash_32"

Ref "compiled_contracts_sources_compilation_id_fkey":"compiled_contracts"."id" < "compiled_contracts_sources"."compilation_id"

Ref "compiled_contracts_sources_source_hash_fkey":"sources"."source_hash" < "compiled_contracts_sources"."source_hash"

Ref "contract_deployments_contract_id_fkey":"contracts"."id" < "contract_deployments"."contract_id"

Ref "contracts_creation_code_hash_fkey":"code"."code_hash" < "contracts"."creation_code_hash"

Ref "contracts_runtime_code_hash_fkey":"code"."code_hash" < "contracts"."runtime_code_hash"

Ref "sourcify_matches_verified_contract_id_fk":"verified_contracts"."id" < "sourcify_matches"."verified_contract_id" [update: restrict, delete: restrict]

Ref "verification_jobs_ephemeral_id_fk":"verification_jobs"."id" < "verification_jobs_ephemeral"."id" [update: cascade, delete: cascade]

Ref "verification_jobs_verified_contract_id_fk":"verified_contracts"."id" < "verification_jobs"."verified_contract_id" [update: restrict, delete: restrict]

Ref "verified_contracts_compilation_id_fkey":"compiled_contracts"."id" < "verified_contracts"."compilation_id"

Ref "verified_contracts_deployment_id_fkey":"contract_deployments"."id" < "verified_contracts"."deployment_id"
