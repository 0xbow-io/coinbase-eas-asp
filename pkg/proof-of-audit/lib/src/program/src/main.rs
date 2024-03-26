
#![no_main]
sp1_zkvm::entrypoint!(main);

use hmac::{Hmac, Mac};
use base64::{Engine as _, engine::general_purpose};
use hex;
use sha2::{Sha256, Digest};
use serde_json::Value; 

const BASE_CHAIN_ID: &str = "0x2105";
const COINBASE_EAS_TOPIC: &str = "0x8bf46bf4cfd674fa735a3d63ec1c9ad4153f033c290341f3a588b75685141b35";
const BASE_EAS_ADDR: &str = "0x4200000000000000000000000000000000000021";
const COINBASE_EAS_HASH: &str = "0x000000000000000000000000357458739f90461b99789350868cd7cf330dd7ee";
const COINBASE_EAS_SCHEMA_ID: &str = "0xf8b05c79f090979bf4a80270aba232dff11a10d9ca55c4f88de95317970f0de9";



/*
    Proof of payload origin: 
        - payload is hashed with url_path
        - A signature is generated using the Hmac hash (secret, paylod hash, nonce, timestamp)
        - Prove that secret & url_path is known to generate the signature
        - Prove that the generated signature matches the expected signature
    
    Once origin of payload is verified, review the payload to verify decision on target commitmentID (aka Attestation UID)
        - iterate through payload and find Tx Receipt log of EAS Attestation: 
            - commitment_id (EAS UUID)
            - public_id (Attested Wallet Address)
*/

fn verify_sig(secret: &str, nonce: &str, timestamp: &str, body_hash: &str, signature: &str) -> bool {
    type HmacSha256 = Hmac<Sha256>;
    let mut h = HmacSha256::new_from_slice(secret.as_bytes())
    .expect("HMAC can take key of any size");

    // generate hash of nonce + computedHashStr + timestamp
    h.update(format!("{}{}{}", nonce, body_hash, timestamp).as_bytes());

    let result = h.finalize().into_bytes();

    // std encoding to base64
    let expected_sig = general_purpose::STANDARD.encode(result);

    return expected_sig == signature;
}

fn verify_commitment(body: &str, commitment_id: &str, public_id: &str) -> bool  {
    // Generic JSON parsing
    // Seek for evidence that commitmentID (Attestation UID) should be member of Inclusion Set.
    let json: Value = serde_json::from_str(&body).unwrap();

    let matched_txs = &json["matchedTransactions"];

    if let serde_json::Value::Array(matched_receipts) = &json["matchedReceipts"] {
       for (index, receipt) in matched_receipts.iter().enumerate() {
            if let serde_json::Value::Array(logs) = &receipt["logs"] {
                for log in logs {
                    // Detect EAS Attestation ADDR
                    if  log["address"].to_string().trim_matches('"').to_lowercase().eq(&BASE_EAS_ADDR) 
                        && log["data"].to_string().trim_matches('"').to_lowercase().eq(&commitment_id)  // EAS UUID
                    {
                        if let serde_json::Value::Array(topics) = &log["topics"] {
                            // Verfy Attestation 
                            if  topics[0].to_string().trim_matches('"').to_lowercase().eq(COINBASE_EAS_TOPIC) 
                                && topics[1].to_string().trim_matches('"').to_lowercase().eq(&public_id)  // Attested Wallet Address
                                && topics[2].to_string().trim_matches('"').to_lowercase().eq(COINBASE_EAS_HASH)   
                                && topics[3].to_string().trim_matches('"').to_lowercase().eq(COINBASE_EAS_SCHEMA_ID)   
                            {
                                // check the chainID is correct 
                                if let Value::Array(transactions) = matched_txs {
                                   if transactions[index] ["chainId"].to_string().trim_matches('"').to_lowercase().eq(BASE_CHAIN_ID) {
                                    return true;
                                   }
                                }
                            }
                        }
                    }
                }
            }
        }
    }
    return false
}



pub fn main() { 
   
    let mut hasher = Sha256::new();
    let mut include = false;

    // private inputs
    let secret = sp1_zkvm::io::read::<String>();
    let url_path = sp1_zkvm::io::read::<String>();

    // read the nonce from webhook response
    // public input
    let nonce = sp1_zkvm::io::read::<String>();
 
    // read the timestamp from webhook response
    // public input
    let timestamp: String = sp1_zkvm::io::read::<String>();
    
    // read the payload body from webhook response
    // public input
    let body = sp1_zkvm::io::read::<String>();

    // payload signature 
    let signature = sp1_zkvm::io::read::<String>();

    // EAS UUID
    let commitment_id = sp1_zkvm::io::read::<String>();

    // Attested Wallet Address
    let public_id = sp1_zkvm::io::read::<String>();


    if secret.len() > 0 && url_path.len() > 0{
        // compute expected hash
        hasher.update(format!("{}{}", url_path, body));
        let computed_body_hash = hex::encode(hasher.finalize());

        // verify signature
        // proof of payload origin
        if verify_sig(&secret, &nonce, &timestamp, &computed_body_hash, &signature)  {
            include = verify_commitment(&body, &commitment_id, &public_id);
        }
    }

    sp1_zkvm::io::write(&include);
}
