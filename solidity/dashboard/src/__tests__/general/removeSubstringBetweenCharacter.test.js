import { removeSubstringBetweenCharacter } from "../../utils/general.utils"

describe("Test removeSubstringBetweenCharacter function", () => {
  it("should return proper values", () => {
    expect(
      removeSubstringBetweenCharacter(
        "0x1E47B39F8909Edcb5cC71ec7Bd0117c48B88cd26/dashboard",
        "/",
        0
      )
    ).toBe("dashboard")
    expect(
      removeSubstringBetweenCharacter(
        "/0x1E47B39F8909Edcb5cC71ec7Bd0117c48B88cd26/dashboard",
        "/",
        0
      )
    ).toBe("/dashboard")
    expect(
      removeSubstringBetweenCharacter("test1/test2/test3/test4", "/", 2)
    ).toBe("test1/test2/test4")
    expect(
      removeSubstringBetweenCharacter("/test1/test2/test3/test4", "/", 2)
    ).toBe("/test1/test2/test4")
    expect(
      removeSubstringBetweenCharacter(" /test1/test2/test3/test4", "/", 2)
    ).toBe(" /test1/test3/test4")
  })

  it("should return same input if the character is not in the input", () => {
    expect(
      removeSubstringBetweenCharacter("test1/test2/test3/test4", "!", 2)
    ).toBe("test1/test2/test3/test4")
  })

  it("should return same input if there is no word in a give occurrenceIndex", () => {
    expect(
      removeSubstringBetweenCharacter("test1/test2/test3/test4", "/", -10)
    ).toBe("test1/test2/test3/test4")

    expect(
      removeSubstringBetweenCharacter("test1/test2/test3/test4", "/", 10)
    ).toBe("test1/test2/test3/test4")
  })
})
