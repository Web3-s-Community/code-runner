package main

import (
	"fmt"

	"code-runner/module/algorand/biz"
	"code-runner/module/code/models"
)

func main() {
	code := "import { Contract } from '@algorandfoundation/tealscript';\n\n// eslint-disable-next-line no-unused-vars\nclass AlgoContracts extends Contract {\n  /**\n   * Calculates the sum of two numbers\n   *\n   * @param a\n   * @param b\n   * @returns The sum of a and b\n   */\n  private getSum(a: number, b: number): number {\n    return a + b;\n  }\n\n  /**\n   * Calculates the difference between two numbers\n   *\n   * @param a\n   * @param b\n   * @returns The difference between a and b.\n   */\n  private getDifference(a: number, b: number): number {\n    return a >= b ? a - b : b - a;\n  }\n\n  /**\n   * A method that takes two numbers and does either addition or subtraction\n   *\n   * @param a The first number\n   * @param b The second number\n   * @param operation The operation to perform. Can be either 'sum' or 'difference'\n   *\n   * @returns The result of the operation\n   */\n  doMath(a: number, b: number, operation: string): number {\n    let result: number;\n\n    if (operation === 'sum') {\n      result = this.getSum(a, b);\n    } else if (operation === 'difference') {\n      result = this.getDifference(a, b);\n    } else throw Error('Invalid operation');\n\n    return result;\n  }\n}\n"
	CodeTest := "import { describe, test, expect, beforeAll, beforeEach } from '@jest/globals';\nimport { algorandFixture } from '@algorandfoundation/algokit-utils/testing';\nimport { AlgoContractsClient } from '../contracts/clients/AlgoContractsClient';\nimport * as algokit from '@algorandfoundation/algokit-utils';\n\nconst fixture = algorandFixture();\nalgokit.Config.configure({ populateAppCallResources: true });\n\nlet appClient: AlgoContractsClient;\n\ndescribe('AlgoContracts', () => {\n  beforeEach(fixture.beforeEach);\n\n  beforeAll(async () => {\n    await fixture.beforeEach();\n    const { algod, testAccount } = fixture.context;\n\n    appClient = new AlgoContractsClient(\n      {\n        sender: testAccount,\n        resolveBy: 'id',\n        id: 0,\n      },\n      algod\n    );\n\n    await appClient.create.createApplication({});\n  });\n\n  test('sum', async () => {\n    const a = 13;\n    const b = 37;\n    const sum = await appClient.doMath({ a, b, operation: 'sum' });\n    expect(sum.return?.valueOf()).toBe(BigInt(a + b));\n  });\n\n  test('difference', async () => {\n    const a = 13;\n    const b = 37;\n    const diff = await appClient.doMath({ a, b, operation: 'difference' });\n    expect(diff.return?.valueOf()).toBe(BigInt(a >= b ? a - b : b - a));\n  });\n});\n"

	_submission := models.SubmissionPlayground{
		Code:     code,
		CodeTest: CodeTest,
	}

	output, err := biz.ExecuteCodeTest(_submission)

	fmt.Println(fmt.Sprintf("message: %s \n err: %s", output, err))
}
