# Base Sepolia EAS schema registration

This is the one manual chain setup required by the attestation writer. Full
local setup and Docker documentation belongs to PR4.

1. Create a dedicated test wallet. Never reuse a wallet that holds real funds.
2. Add Base Sepolia (chain ID `84532`) and obtain test ETH from a Base-supported
   faucet.
3. Open the Base Sepolia EAS explorer at
   <https://base-sepolia.easscan.org/schema/create> and connect the test wallet.
4. Register this exact schema:

   ```text
   bytes32 dayDigest,bytes32 subjectDayCommitment
   ```

   Use the zero resolver (`0x0000000000000000000000000000000000000000`)
   and make the schema revocable.
5. Copy the resulting schema UID into `EAS_SCHEMA_UID`.

The writer targets the canonical Base Sepolia EAS contract
`0x4200000000000000000000000000000000000021`.

Set these environment variables at runtime:

```text
BASE_SEPOLIA_RPC_URL=<Base Sepolia HTTPS or WSS RPC URL>
EAS_PRIVATE_KEY=<dedicated test wallet private key, hex>
EAS_SCHEMA_UID=<registered bytes32 schema UID, 0x-prefixed>
EAS_SUBJECT_HMAC_SECRET=<independent high-entropy server secret, at least 32 bytes>
```

Do not put real values in source control, command arguments, logs, or error
messages. The HMAC secret must not be the wallet private key.
