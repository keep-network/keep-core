import { smock } from "@defi-wonderland/smock"
import type { FakeContract } from "@defi-wonderland/smock"
import type { IRandomBeaconStaking, RandomBeacon } from "../../typechain"

// eslint-disable-next-line import/prefer-default-export
export async function fakeTokenStaking(
  randomBeacon: RandomBeacon
): Promise<FakeContract<IRandomBeaconStaking>> {
  const tokenStaking = await smock.fake<IRandomBeaconStaking>(
    "IRandomBeaconStaking",
    {
      address: await randomBeacon.callStatic.staking(),
    }
  )

  return tokenStaking
}
