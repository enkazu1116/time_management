# Time Management API

Go製のtime-management APIです。現在のapplication dataとattestation UIDは
process内memoryに保存されます。Base Sepolia EASへのattestationだけがchain上に
残ります。

## Base Sepolia tamper-detection demo

このdemoはtestnet専用です。mainnet、送金、実資金を扱いません。
`ATTEST_ADMIN_TOKEN`はlocal demo用のstatic write guardであり、一般的な認証機構
ではありません。

### 1. Test walletとtest ETH

1. MetaMaskやRabbyで、このdemo専用の新しいaccountを作成します。実資金を持つ
   walletを再利用しないでください。
2. Base Sepolia（chain ID `84532`）をwalletへ追加します。RPC URLは利用する
   providerから取得してください。
3. [Baseのfaucet一覧](https://docs.base.org/base-chain/tools/network-faucets)
   から専用accountへBase Sepolia ETHを取得します。
4. 専用accountのprivate keyをwalletからexportします。画面共有、shell history、
   source code、chat、logへ貼らないでください。

### 2. EAS schemaを一度だけ登録

1. [Base Sepolia EAS schema作成画面](https://base-sepolia.easscan.org/schema/create)
   を開き、専用test walletを接続します。
2. 次のschemaを正確に登録します。

   ```text
   bytes32 dayDigest,bytes32 subjectDayCommitment
   ```

3. Resolverはzero address
   `0x0000000000000000000000000000000000000000`、Revocableは`true`にします。
4. transaction確定後、表示された32-byte schema UIDを控えます。

APIはBase Sepoliaのcanonical EAS contract
`0x4200000000000000000000000000000000000021`を使用します。詳細は
[`docs/base-sepolia-eas-schema.md`](docs/base-sepolia-eas-schema.md)にもあります。

### 3. Local environment

`.env.example`は変数名を確認するためのdummy templateです。private keyを含む実値は
`.env`を含め、どのfileにも保存しません。起動するterminalのenvironmentへだけ設定
します。`read -s`で入力したprivate keyは画面にもshell historyにも残りません。

```sh
read -r -p 'Base Sepolia RPC URL: ' BASE_SEPOLIA_RPC_URL
export BASE_SEPOLIA_RPC_URL
read -r -p 'EAS schema UID: ' EAS_SCHEMA_UID
export EAS_SCHEMA_UID
read -r -s -p 'Dedicated test wallet private key: ' EAS_PRIVATE_KEY
printf '\n'
export EAS_PRIVATE_KEY
export EAS_SUBJECT_HMAC_SECRET="$(openssl rand -hex 32)"
export ATTEST_ADMIN_TOKEN="$(openssl rand -hex 32)"
```

HMAC secretとadmin tokenは毎回変更すると、既存attestationの検索・POST authorization
が変わります。demo session中は同じterminalでComposeを操作してください。HMAC secret
をwallet private keyと共用しないでください。`.env`は念のためGitとDocker build
contextから除外済みで、`.env.example`にはdummy値以外を入れません。

### 4. 起動

```sh
docker compose up --build -d
```

同じterminalでhealth checkを確認します。これにより、前のstepでexportしたsecretはfileへ
保存せず、そのままcurl walkthroughにも利用できます。

```sh
curl --fail http://localhost:8080/healthz
```

以下のwalkthrough用に固定IDをshell変数へ設定します。環境変数は起動時と同じterminal
に既に設定されています。

```sh
API=http://localhost:8080
USER_ID=00000000-0000-4000-8000-000000000101
WORK_LOG_ID=00000000-0000-4000-8000-000000000201
WORK_TYPE_ID=00000000-0000-4000-8000-000000000301
BREAK_LOG_ID=00000000-0000-4000-8000-000000000401
WEEKLY_TASK_ID=00000000-0000-4000-8000-000000000501
DAILY_TASK_ID=00000000-0000-4000-8000-000000000601
ACHIEVEMENT_ID=00000000-0000-4000-8000-000000000701
```

### 5. Recordsを作成

Work log:

```sh
curl --fail-with-body -X POST "$API/work-logs" \
  -H 'Content-Type: application/json' \
  -d "{
    \"work_log_id\":\"$WORK_LOG_ID\",
    \"user_id\":\"$USER_ID\",
    \"work_type_id\":\"$WORK_TYPE_ID\",
    \"work_type_name\":\"demo work\",
    \"work_date\":\"2026-11-10\",
    \"started_at\":\"2026-11-10T09:00:00+09:00\",
    \"ended_at\":\"2026-11-10T12:00:00+09:00\",
    \"work_duration_seconds\":10800
  }"
```

Break log:

```sh
curl --fail-with-body -X POST "$API/break-logs" \
  -H 'Content-Type: application/json' \
  -d "{
    \"break_log_id\":\"$BREAK_LOG_ID\",
    \"work_log_id\":\"$WORK_LOG_ID\",
    \"user_id\":\"$USER_ID\",
    \"started_at\":\"2026-11-10T10:30:00+09:00\",
    \"ended_at\":\"2026-11-10T10:45:00+09:00\",
    \"break_duration_seconds\":900
  }"
```

Weekly taskとdaily task:

```sh
curl --fail-with-body -X POST "$API/weekly-tasks" \
  -H 'Content-Type: application/json' \
  -d "{
    \"weekly_task_id\":\"$WEEKLY_TASK_ID\",
    \"user_id\":\"$USER_ID\",
    \"work_type_id\":\"$WORK_TYPE_ID\",
    \"week_start_date\":\"2026-11-09\",
    \"target_duration_seconds\":36000,
    \"goal_description\":\"finish the demo\"
  }"

curl --fail-with-body -X POST "$API/daily-tasks" \
  -H 'Content-Type: application/json' \
  -d "{
    \"daily_task_id\":\"$DAILY_TASK_ID\",
    \"weekly_task_id\":\"$WEEKLY_TASK_ID\",
    \"user_id\":\"$USER_ID\",
    \"work_date\":\"2026-11-10\",
    \"progress_description\":\"demo complete\",
    \"progress_rate\":0.75
  }"
```

Daily achievement:

```sh
curl --fail-with-body -X POST "$API/summaries/daily" \
  -H 'Content-Type: application/json' \
  -d "{
    \"daily_achievement_id\":\"$ACHIEVEMENT_ID\",
    \"user_id\":\"$USER_ID\",
    \"work_date\":\"2026-11-10\"
  }"
```

### 6. Attestしてmatchを確認

POSTだけにdemo bearer tokenが必要です。

```sh
curl --fail-with-body -X POST "$API/attestations" \
  -H "Authorization: Bearer $ATTEST_ADMIN_TOKEN" \
  -H 'Content-Type: application/json' \
  -d "{\"user_id\":\"$USER_ID\",\"date\":\"2026-11-10\"}"

curl --fail-with-body \
  "$API/attestations/verify?user_id=$USER_ID&date=2026-11-10"
```

2つ目のresponseは`"status":"match"`になります。Attestation transactionの確定まで
POSTは待機するため、RPCやnetworkによって時間がかかる場合があります。

### 7. Recordを編集してmismatchを確認

現在のmemory adapterでは同じIDのPOSTがrecordを置き換えます。work type nameだけを
変更して再送します。

```sh
curl --fail-with-body -X POST "$API/work-logs" \
  -H 'Content-Type: application/json' \
  -d "{
    \"work_log_id\":\"$WORK_LOG_ID\",
    \"user_id\":\"$USER_ID\",
    \"work_type_id\":\"$WORK_TYPE_ID\",
    \"work_type_name\":\"edited after attestation\",
    \"work_date\":\"2026-11-10\",
    \"started_at\":\"2026-11-10T09:00:00+09:00\",
    \"ended_at\":\"2026-11-10T12:00:00+09:00\",
    \"work_duration_seconds\":10800
  }"

curl --fail-with-body \
  "$API/attestations/verify?user_id=$USER_ID&date=2026-11-10"
```

responseは`"status":"mismatch"`になります。再度attestすると新しいattestationが
作られ、直前のUIDが`refUID`になります。

### Demoの制約

- APIをrestartするとmemory上のrecordsとUID locatorは消えます。その場合verifyは
  `not-attested`になりますが、作成済みEAS attestationはchain上に残ります。
- RPC failure、wrong chain、wrong attester/schema、revoked/malformed attestationは
  `not-attested`ではなくerrorです。
- Static bearer tokenはlocal demo用です。public deploymentには通常のauthentication、
  authorization、rate limiting、secret managementが必要です。

終了時:

```sh
docker compose down
unset EAS_PRIVATE_KEY EAS_SUBJECT_HMAC_SECRET ATTEST_ADMIN_TOKEN
```
