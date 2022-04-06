import { smock } from "@defi-wonderland/smock"

import type { FakeContract } from "@defi-wonderland/smock"
import type { RandomBeacon, TokenStaking } from "../../typechain"

// eslint-disable-next-line import/prefer-default-export
export async function fakeTokenStaking(
  randomBeacon: RandomBeacon
): Promise<FakeContract<TokenStaking>> {
  const tokenStaking = await smock.fake<TokenStaking>("TokenStaking", {
    address: await randomBeacon.callStatic.staking(),
  })

  return tokenStaking
}
