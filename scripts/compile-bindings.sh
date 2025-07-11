#!/bin/bash

BINDING_DIR=./pkg/bindings
JSON_DIR=.devkit/contracts/out

function create_binding {
    contract_name=$1
    echo "Generating binding for $contract_name"

    mkdir -p $BINDING_DIR/${contract_name}

    contract_json_path="${JSON_DIR}/${contract_name}.sol/${contract_name}.json"
    echo "Contract JSON path: $contract_json_path"

    binding_out_dir="${BINDING_DIR}/${contract_name}"
    mkdir -p $binding_out_dir || true

    cat $contract_json_path | jq -r '.abi' > $binding_out_dir/tmp.abi
    cat $contract_json_path | jq -r '.bytecode.object' > $binding_out_dir/tmp.bin

    abigen \
        --bin=$binding_out_dir/tmp.bin \
        --abi=$binding_out_dir/tmp.abi \
        --pkg="${contract_name}" \
        --out=$BINDING_DIR/$contract_name/binding.go \
        > /dev/null 2>&1

    if [[ $? == "1" ]];
    then
        echo "Failed to generate binding for $contract_json_path"
    fi
    echo "Binding for $contract_name generated at $BINDING_DIR/$contract_name/binding.go"
    rm $binding_out_dir/tmp.abi
    rm $binding_out_dir/tmp.bin
}

# forge b

contracts=$(find contracts/src -type f -name "*.sol" )
IFS=$'\n'

mkdir -p $JSON_DIR || true

create_binding "OmniVRF"
create_binding "ITaskMailbox"


# for contract_name in $contracts; do
# 	contract_name=$(basename $contract_name .sol)
# 	create_binding $contract_name
# done
