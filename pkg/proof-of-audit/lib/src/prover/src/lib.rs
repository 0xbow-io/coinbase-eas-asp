use sp1_core::{SP1Verifier, SP1Prover, SP1Stdin};
use serde_json;
use libc;
use std::ffi::{CStr, CString};

const ELF: &[u8] = include_bytes!("../../../elf/riscv32im-succinct-zkvm-elf");

fn get_string_c_char(s: *const libc::c_char) -> String {
    let c_str = unsafe {
        assert!(!s.is_null());
        CStr::from_ptr(s)
    };

    return c_str.to_str().unwrap().to_string();
}


#[no_mangle]
pub extern "C" fn generate_sp1_proof_ffi(
    secret: *const libc::c_char,
    url_path: *const libc::c_char,
    nonce: *const libc::c_char,
    timestamp: *const libc::c_char,
    payload: *const libc::c_char,
    signature: *const libc::c_char,
    commitment_id: *const libc::c_char,
    public_id: *const libc::c_char,
) -> *mut libc::c_char {

    let mut stdin = SP1Stdin::new();
    
    stdin.write(&get_string_c_char(secret));

    stdin.write(&get_string_c_char(url_path));

    stdin.write(&get_string_c_char(nonce));

    stdin.write(&get_string_c_char(timestamp));

    stdin.write(&get_string_c_char(payload));

    stdin.write(&get_string_c_char(signature));

    stdin.write(&get_string_c_char(commitment_id));
    stdin.write(&get_string_c_char(public_id));

    let mut proof = SP1Prover::prove(ELF, stdin).expect("proving failed");
    if !SP1Verifier::verify(ELF, &proof).is_ok() {
    return CString::new("invalid proof").unwrap().into_raw();
    }

    let commitment_ok = proof.stdout.read::<bool>();
    if !commitment_ok {
        return CString::new("invalid commitment").unwrap().into_raw();
    }

    // remove inputs as it contains private inputs
    proof.stdin = SP1Stdin::new();

    // serialize proof
    let proof_serialized = serde_json::to_string(&proof).unwrap();

    // convert to CString
    return CString::new(proof_serialized).unwrap().into_raw();
}


#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn verify_sp1_proof_works() {

        // prepare the inputs as c_chars
        let secret_c_str = CString::new("qnsec_dFzHeJ5iQbefXDH1akAKow==").unwrap();
        let secret_c_char: *const libc::c_char = secret_c_str.as_ptr() as *const libc::c_char;

        let url_path_c_str = CString::new("/webhook/2057b8e5-de11-4c65-8e39-92507d20de80").unwrap();
        let url_path_c_char: *const libc::c_char = url_path_c_str.as_ptr() as *const libc::c_char;

        let nonce_c_str = CString::new("632e2d63-d253-4a06-ab77-d565e806e5e1").unwrap();
        let nonce_c_char: *const libc::c_char = nonce_c_str.as_ptr() as *const libc::c_char;
        
        let timestamp_c_str = CString::new("2024-03-12 04:04:14.11330824 +0000 UTC m=+17208.257632042").unwrap();
        let timestamp_c_char: *const libc::c_char = timestamp_c_str.as_ptr() as *const libc::c_char;

        let payload_str = r#"{"matchedReceipts":[{"blockHash":"0x6644e74c3ff45ac8d514cdfc7393bcd193067c47c138842427a46a49d528573a","blockNumber":"0xb2bbad","contractAddress":"","cumulativeGasUsed":"0x70ef7","effectiveGasPrice":"0x187d3","from":"0x8844591d47f17bca6f5df8f6b64f4a739f1c0080","gasUsed":"0x450b5","logs":[{"address":"0x4200000000000000000000000000000000000021","blockHash":"0x6644e74c3ff45ac8d514cdfc7393bcd193067c47c138842427a46a49d528573a","blockNumber":"0xb2bbad","data":"0x8133f214f7bdaf516f03655db7406ba3c7945e5e4849238e43fdea6ef7a25cd4","logIndex":"0x1","removed":false,"topics":["0x8bf46bf4cfd674fa735a3d63ec1c9ad4153f033c290341f3a588b75685141b35","0x000000000000000000000000ff9418c67d18c8e067141bd77be43e32c4c3abe7","0x000000000000000000000000357458739f90461b99789350868cd7cf330dd7ee","0xf8b05c79f090979bf4a80270aba232dff11a10d9ca55c4f88de95317970f0de9"],"transactionHash":"0x63a3aef220d84947b12b0d771fa0df510eb753cab0c218efa1861b0d7c3d4567","transactionIndex":"0x4"},{"address":"0x2c7ee1e5f416dff40054c27a62f7b357c4e8619c","blockHash":"0x6644e74c3ff45ac8d514cdfc7393bcd193067c47c138842427a46a49d528573a","blockNumber":"0xb2bbad","data":"0x000000000000000000000000d867cbed445c37b0f95cc956fe6b539bdef7f32f","logIndex":"0x2","removed":false,"topics":["0x7fd54fcc14543b4db08cef4cd9fb23a6670c072d8a44cb0f1817d35b474176ca","0x000000000000000000000000ff9418c67d18c8e067141bd77be43e32c4c3abe7","0xf8b05c79f090979bf4a80270aba232dff11a10d9ca55c4f88de95317970f0de9","0x8133f214f7bdaf516f03655db7406ba3c7945e5e4849238e43fdea6ef7a25cd4"],"transactionHash":"0x63a3aef220d84947b12b0d771fa0df510eb753cab0c218efa1861b0d7c3d4567","transactionIndex":"0x4"}],"logsBloom":"0x00000000000000000000000040000000100000000000000000000000000000000001000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000008000000000000000000020040000000000000000000000000000000000000002000000000800000000000000000000000000000000400001000000000000000000000000010000000000000000000008010800000000020000000000010000000000000002000000000000000800000000000000000000000000000000004000000000000000000000000000800010000200000000000000000000000000000000000000000080000","status":"0x1","to":"0x357458739f90461b99789350868cd7cf330dd7ee","transactionHash":"0x63a3aef220d84947b12b0d771fa0df510eb753cab0c218efa1861b0d7c3d4567","transactionIndex":"0x4","type":"0x2"}],"matchedTransactions":[{"accessList":[],"blockHash":"0x6644e74c3ff45ac8d514cdfc7393bcd193067c47c138842427a46a49d528573a","blockNumber":"0xb2bbad","chainId":"0x2105","from":"0x8844591d47f17bca6f5df8f6b64f4a739f1c0080","gas":"0x927c0","gasPrice":"0x187d3","hash":"0x63a3aef220d84947b12b0d771fa0df510eb753cab0c218efa1861b0d7c3d4567","input":"0x56feed5e000000000000000000000000ff9418c67d18c8e067141bd77be43e32c4c3abe7","maxFeePerGas":"0x31214","maxPriorityFeePerGas":"0x186a0","nonce":"0x11211","r":"0xf015d09c9d8b274b703fad93020d38300aba7a9a0cbc786f05dac65dd0e0795b","s":"0x2d77121dc4ed30832b960e025b2b4b961e5103c8ad7e8f629a1eb381977c556c","to":"0x357458739f90461b99789350868cd7cf330dd7ee","transactionIndex":"0x4","type":"0x2","v":"0x0","value":"0x0"}]}"#
        .to_string();

        let payload_c_str = CString::new(payload_str).unwrap();
        let payload_c_str_char: *const libc::c_char = payload_c_str.as_ptr() as *const libc::c_char;

        let signature_c_str = CString::new("tsPXXWc53kXWob6ZbrHCxDSUrljtmk40d5vGGtFXvbs=").unwrap();
        let signature_c_str_char: *const libc::c_char = signature_c_str.as_ptr() as *const libc::c_char;

        let commitment_id_c_str = CString::new("0x8133f214f7bdaf516f03655db7406ba3c7945e5e4849238e43fdea6ef7a25cd4").unwrap();
        let commitment_id_c_str_char: *const libc::c_char = commitment_id_c_str.as_ptr() as *const libc::c_char;

        let public_id_c_str = CString::new("0x000000000000000000000000ff9418c67d18c8e067141bd77be43e32c4c3abe7").unwrap();
        let public_id_c_str_char: *const libc::c_char = public_id_c_str.as_ptr() as *const libc::c_char;

        let proof = generate_sp1_proof_ffi(
            secret_c_char, 
            url_path_c_char,
            nonce_c_char,
            timestamp_c_char,
            payload_c_str_char,
            signature_c_str_char,
            commitment_id_c_str_char,
            public_id_c_str_char
            );

        let proof_str = get_string_c_char(proof);
        assert!(proof_str.contains("proof"));
    }
}

