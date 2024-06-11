import { describe, test, expect, beforeAll, beforeEach } from '@jest/globals';
import { algorandFixture } from '@algorandfoundation/algokit-utils/testing';
import { WritingFirstContractClient } from '../contracts/clients/WritingFirstContractClient';

const fixture = algorandFixture();

let appClient: WritingFirstContractClient;

describe('Runner', () => {
  beforeEach(fixture.beforeEach);

  beforeAll(async () => {
    await fixture.beforeEach();
    const { algod, testAccount } = fixture.context;

    appClient = new WritingFirstContractClient(
      {
        sender: testAccount,
        resolveBy: 'id',
        id: 0,
      },
      algod
    );

    await appClient.create.createApplication({});
  });

  test('should get proposal', async () => {
    const proposal = await appClient.getProposal({});
    expect(proposal.return?.valueOf()).toBe('This is a proposal.');
  });
});
