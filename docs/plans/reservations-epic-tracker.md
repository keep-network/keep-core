# UTXO Reservations (M1) - epic tracker

This is the live tracker for every milestone and open item across the reservation PRs. Update it in the same change that completes an item.

- History up to 2026-09-28 is in `docs/plans/m1-delivery.md` on branch `docs/reservations-spec`. This file replaces its status sections from that date on.
- Scope: M1 only (create-only vault). M2 work (reserved redemption, multi-term) is listed under "Related, out of scope".
- Exposure rule: fixes to code that already runs on mainnet are tracked privately. They ship in a coordinated client release, and this file does not list them.

Last updated: 2026-10-10.

## Where we are

All work that can run without the deployer key is done. The next step is the Sepolia deploy, which the key holder runs (milestone M4). After that come the testnet drill (M5) and publication of the review fixes (M6).

```mermaid
flowchart LR
  M0[M0 review fixes merged] --> M1[M1 audit and drill blockers]
  M1 --> M2[M2 tests and fork dry-run]
  M2 --> M3[M3 Sepolia deploy script]
  M3 --> M4[M4 Sepolia deploy]
  M4 --> M5[M5 testnet drill]
  M5 --> M6[M6 publish fixes, reconcile with dev]
  M6 --> M7[M7 mainnet readiness]
```

## Related PRs

| Repo | PR | Role | State (2026-10-10) |
|---|---|---|---|
| keep-core | #4282 `reservations-epic` -> `dev` | client tracker | open, draft, conflicts with `dev` |
| keep-core | #4324, #4343, #4346 | review-fix rounds | merged into `reservations-epic` |
| keep-core | #4238 | earlier attempt | superseded by #4282, kept for reference |
| tbtc-v2 | #1116 `reservations-upgrade` -> `dev` | contracts tracker | open, draft, conflicts with `dev` |
| tbtc-v2 | #1161 | cross-repo review fixes | merged into `reservations-upgrade` |
| keep-core | branch `docs/reservations-spec` | M1 spec | no PR yet |

Related, out of scope for M1: tbtc-v2 #1129 (M2 reserved redemption), #1162 (multi-term), #1088 (superseded pre-M1 stack root).

## Milestones

### M0 - review-fix rounds merged - DONE

- [x] keep-core #4324, #4343, #4346 merged into `reservations-epic`.
- [x] tbtc-v2 #1161 merged into `reservations-upgrade`.

### M1 - security audit and drill blockers - DONE

- [x] Full audit of #4282 and the matching tbtc-v2 contracts.
- [x] Each finding routed by exposure: reservation-only findings are fixed for this epic; findings in code already on mainnet are tracked privately.
- [x] 7 drill-blocking reservation fixes made, each with a test that fails before the fix.
- [x] Third-party re-anchor hold cap per target wallet (tbtc-v2), with the client sending the member proof it needs.
- [x] Bridge guard: every reservation and deposit parameter setter refuses a value set that makes reserved deposits impossible (`now + revealAhead <= refundDeadline <= now + term + 24h`, and amount at least dust).
- [x] Client gauge `reservation_params_infeasible` warns when the live parameters fail the same check.
- [x] Spec updated for the hold cap and the feasibility guard (`docs/reservations-spec`).

Fixes are held in a review line that is not public. Publication is M6.

### M2 - tests and fork dry-runs - DONE

- [x] Chained test: acceptance followed by re-anchor.
- [x] tbtc-v2 #1116 line merged with `dev`. All fixture suites pass locally (about 540 tests in 16 suites).
- [x] Mainnet-fork activation dry-run: 22/22.
- [x] Mainnet-fork lifecycle (reveal, acceptance request, acceptance proof): 25/25.
- [x] Lint rule that bans the self-catching `expect.fail` inside `try` pattern, with a policy test.

### M3 - Sepolia deploy script - DONE

- [x] Runbook reordered so that the client with the activation block ships before `setVaultStatus`.
- [x] One script with three modes: fork rehearsal, `PHASE=prepare` (steps 1-6 and 10, resumable), and `PHASE=activate ACTIVATION_BLOCK=<B>` (steps 11-12).
- [x] Live mode refuses a signer other than the deployer EOA, and refuses to activate before block B.
- [x] Rehearsed on a Sepolia fork (23 pass, 3 off-chain skips). Rehearsed live against a local Sepolia-fork node, with a killed and resumed prepare and real-time governance waits.

### M4 - Sepolia deploy - NEXT (waits on the key holder)

- [ ] Recommended, waiting for approval: before resuming, validate the saved prepare file. Tie it to the chain ID and the Bridge address, and require chain ID 11155111 on `--network sepolia`.
- [ ] Key holder runs `PHASE=prepare` on Sepolia and shares the printed addresses.
- [ ] Publish npm `@keep-network/tbtc-v2@sepolia` with the reservation artifacts.
- [ ] Set the Sepolia activation block B in `pkg/tbtc/coordination.go`, regenerate the bindings, build the testnet client.
- [ ] Drill operators (existing Sepolia operators) run that build.
- [ ] Key holder runs `PHASE=activate ACTIVATION_BLOCK=<B>` at or after block B.
- [ ] Read-only check of the deployed state against the runbook values.

### M5 - testnet drill - TODO

- [ ] Fill in the remaining TBD items in the drill runbook (addresses, block B, operator list).
- [ ] Run the drill: reserved deposit, acceptance, re-anchor, timeout and strand paths, plus alarms.
- [ ] Record the results and any new findings. Route each finding by the exposure rule.

### M6 - publish and reconcile with dev - TODO

- [ ] Open the reservation-only fixes as public PRs into `reservations-epic` and `reservations-upgrade` (only after the drill).
- [ ] Resolve the conflicts of #4282 and #1116 with `dev`.
- [ ] Open a PR for the spec branch `docs/reservations-spec`.
- [ ] Ship the mainnet-code fixes in a coordinated client release (tracked privately).

### M7 - mainnet readiness - TODO

- [ ] Decide the fix for the known limit on third-party re-anchor holds (accepted for testnet only).
- [ ] Choose launch values that pass the feasibility guard. Amount caps must be at least dust (0.01 BTC); a 12-month term is fine.
- [ ] Keep a Bridge bytecode budget: 1,009 bytes are left under EIP-170.
- [ ] Mainnet runbook through Timelock and Council Safe, rehearsed on a mainnet fork.

## Open items not tied to a milestone

- [ ] Fork lifecycle step L3: also assert the depositor's TBTC balance change (and the fee recipient's change, if a fee is charged), not only `mintedAmount`.

## Decisions

| Date | Decision |
|---|---|
| 2026-10-08 | Split fixes by exposure: reservation-only fixes go public after approval; mainnet-code fixes stay private until a release ships them. |
| 2026-10-10 | The fork lifecycle covers reveal and acceptance only. Unit tests and the drill cover re-anchor, timeout and strand. |
| 2026-10-10 | Accept the known third-party hold limit for testnet; decide a fix before mainnet. |
| 2026-10-10 | Script and rehearse the Sepolia runbook. The key holder runs it; the fee reserve target is 0. |
| 2026-10-10 | The drill uses the existing Sepolia operators. |
| 2026-10-10 | Publish the fixes after the drill. |
| 2026-10-10 | Approved: Bridge feasibility guard, lint ban, client feasibility gauge. Declined: a nightly fork dry-run in CI. |
